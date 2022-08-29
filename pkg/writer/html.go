package writer

import (
	"html/template"
	"os"
	"sort"
	"sorted-awesome-go/internal/awesomego"
)

type Output struct {
	Sections []awesomego.Section
}

func WriteHTMLFile(outputFile, templateFile string, sections map[string]awesomego.Section) error {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	fileHandler, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer fileHandler.Close()

	keys := make([]string, 0, len(sections))
	for k := range sections {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sortedSections []awesomego.Section
	for _, k := range keys {
		sortedSections = append(sortedSections, sections[k])
	}

	return tmpl.Execute(fileHandler, Output{Sections: sortedSections})
}
