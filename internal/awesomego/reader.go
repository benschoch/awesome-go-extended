package awesomego

import (
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"sorted-awesome-go/pkg/cache"
	"sorted-awesome-go/pkg/github"
	"sorted-awesome-go/pkg/hash"
	"strings"
	"sync"
)

const (
	// all credits go to https://github.com/avelino/awesome-go
	awesomeGoURL  = "https://raw.githubusercontent.com/avelino/awesome-go/main/README.md"
	githubBaseURL = "https://github.com/"
)

type Reader struct {
	cache        cache.Handler
	githubClient github.ClientInterface
}

func NewReader(cache cache.Handler, githubClient github.ClientInterface) *Reader {
	return &Reader{cache: cache, githubClient: githubClient}
}

func (r *Reader) Read() (map[string]Section, error) {
	readmeContent, err := r.fetchReadMe()
	if err != nil {
		return nil, err
	}
	if readmeContent == "" {
		return nil, errors.New("got empty readme content")
	}

	return r.parseReadme(readmeContent)
}

func (r *Reader) parseReadme(readme string) (map[string]Section, error) {
	rawSections := strings.Split(readme, "\n## ")
	searchedSections := make(map[int]Section)
	sectionMutex := sync.RWMutex{}
	var wg sync.WaitGroup
	reachedPackages := false
	for i, rawSection := range rawSections {
		// continue to next headline until we passed the Contents headline
		hasContentsPrefix := strings.HasPrefix(rawSection, "Contents\n")
		if !reachedPackages && !hasContentsPrefix {
			continue
		} else if hasContentsPrefix {
			reachedPackages = true
			continue
		}

		// parse sections in parallel
		wg.Add(1)
		go func(i int, rawSection string) {
			defer wg.Done()
			section := r.parseSection(rawSection)
			if section != nil {
				sectionMutex.Lock()
				searchedSections[i] = *section
				sectionMutex.Unlock()
			}
		}(i, rawSection)
	}
	wg.Wait()

	foundSections := make(map[string]Section, len(searchedSections))
	for _, section := range searchedSections {
		foundSections[section.Title] = section
	}

	return foundSections, nil
}

func (r *Reader) parseSection(rawSection string) *Section {
	if !strings.Contains(rawSection, githubBaseURL) {
		return nil
	}
	sectionLines := strings.Split(rawSection, "\n")
	if len(sectionLines) == 0 {
		log.Println("skipping unexpectedly formatted section")
		return nil
	}

	title, description, packageList := r.readSectionLines(sectionLines)
	log.Printf("parsing packages in section %q", title)
	githubPackages, otherPackages := r.parsePackageList(packageList)
	if len(githubPackages) == 0 && len(otherPackages) == 0 {
		log.Printf("skipping empty section %q", title)
		return nil
	}

	section := NewSection(title, description, githubPackages, otherPackages)

	return &section
}

func (r *Reader) fetchReadMe() (string, error) {
	key := hash.Sha1(awesomeGoURL)
	if r.cache.Has(key) {
		log.Println("reading readme from cache")
		return r.cache.ReadString(key), nil
	}

	response, err := http.Get(awesomeGoURL)
	if err != nil {
		return "", err
	}

	readmeContent, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	err = r.cache.Write(key, readmeContent)
	if err != nil {
		return "", err
	}

	return string(readmeContent), nil
}

func (r *Reader) readSectionLines(lines []string) (title string, description string, packageList string) {
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "- [") {
			packageList = packageList + "\n" + line
		} else if strings.HasPrefix(line, "_") {
			description = strings.Trim(strings.TrimSpace(line), "_")
		} else if strings.HasPrefix(line, "**") {
			continue
		} else if !strings.HasPrefix(line, "#") && title == "" {
			title = strings.Trim(line, " #")
		}
	}

	return title, description, packageList
}

func (r *Reader) parsePackageList(packageListContent string) (
	githubPackages []Package,
	otherPackages []Package,
) {
	listParts := strings.Split(packageListContent, "\n")
	for _, listPartLine := range listParts {
		listPartLine = strings.TrimSpace(listPartLine)
		if listPartLine == "" {
			continue
		}
		if !strings.HasPrefix(listPartLine, "- [") {
			log.Printf("skip unexpected list items: %s", listPartLine)
			continue
		}

		packageName := r.parsePackageName(listPartLine)
		packageURL := r.parseURL(listPartLine)
		description := r.parseDescription(listPartLine)

		stars, forks, updatedAt, fromGithub := r.githubClient.GetDetails(packageURL)

		p := NewPackage(packageName, packageURL, stars, forks, updatedAt, description)
		if fromGithub {
			githubPackages = append(githubPackages, p)
		} else {
			otherPackages = append(otherPackages, p)
		}
	}
	return
}

func (r *Reader) parsePackageName(line string) string {
	return regexp.MustCompile(`^\s*- \[([^]]+).*$`).ReplaceAllString(line, "$1")
}

func (r *Reader) parseURL(line string) (packageURL string) {
	urlPattern := regexp.MustCompile(`^[^]]+]\(([^)]+)\).*$`)
	if urlPattern.MatchString(line) {
		packageURL = urlPattern.ReplaceAllString(line, "$1")
	}
	return packageURL
}

func (r *Reader) parseDescription(line string) (description string) {
	// trim everything before description
	return regexp.MustCompile(`^\s*- \[[^]]+](\([^)]+\))?\s+(-\s+)?(.+)$`).ReplaceAllString(line, "$3")
}
