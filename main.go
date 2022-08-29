package main

import (
	"flag"
	"fmt"
	"sorted-awesome-go/internal/awesomego"
	"sorted-awesome-go/pkg/cache"
	"sorted-awesome-go/pkg/github"
	"sorted-awesome-go/pkg/writer"
)

const (
	templateFile   = "./template/index.html"
	htmlOutputFile = "./index.html"
)

var flushCache bool

func main() {
	flag.BoolVar(&flushCache, "flush", false, "flush cache")
	flag.Parse()

	localCache, err := cache.NewLocalCache()
	if err != nil {
		panic(err)
	}

	if flushCache {
		err := localCache.EraseAll()
		if err != nil {
			panic(err)
		}
	}

	githubClient := github.NewClient(localCache)

	awesomeGoReader := awesomego.NewReader(localCache, githubClient)

	result, err := awesomeGoReader.Read()
	if err != nil {
		panic(err)
	}

	err = writer.WriteHTMLFile(htmlOutputFile, templateFile, result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("🦄 Wrote output to %q\n", htmlOutputFile)
}
