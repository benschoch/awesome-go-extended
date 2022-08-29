package awesomego

import (
	"sorted-awesome-go/pkg/hash"
	"time"
)

type Package struct {
	Name        string
	URL         string
	Stars       uint
	Forks       uint
	UpdatedAt   time.Time
	Description string
}

func NewPackage(packageName string, url string, stars uint, forks uint, updatedAt time.Time, description string) Package {
	return Package{
		Name:        packageName,
		URL:         url,
		Stars:       stars,
		Forks:       forks,
		UpdatedAt:   updatedAt,
		Description: description,
	}
}

type Section struct {
	Title          string
	Description    string
	Anchor         string
	GithubPackages []Package
	OtherPackages  []Package
}

func NewSection(title, description string, githubPackages, otherPackages []Package) Section {
	return Section{
		Title:          title,
		Description:    description,
		Anchor:         hash.Sha1(title),
		GithubPackages: githubPackages,
		OtherPackages:  otherPackages,
	}
}
