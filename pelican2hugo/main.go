package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func convert2Hugo(articles []string) {
	for _, article := range articles {
		log.Printf("Converting %v to hugo\n", article)
		file, err := os.Open(article)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		f, err := os.Create(article + ".hugo")
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
		// begin frontmatter for hugo
		f.WriteString("---\n")
		metadataEnd := false

		scanner := bufio.NewScanner(file)
		var line string
		for scanner.Scan() {
			line = string(scanner.Text())

			if strings.HasPrefix(line, "Title:") {
				title := strings.Split(line, ":")
				line = "title: " + title[1]
			}
			if strings.HasPrefix(line, ":Title:") {
				title := strings.Split(line, ":")
				line = "title: " + title[2]
			}

			if strings.HasPrefix(line, "Date:") {
				date := strings.Split(line, ":")
				line = "date: " + date[1]
			}
			if strings.HasPrefix(line, ":Date:") {
				date := strings.Split(line, ":")
				line = "date: " + date[2]
			}

			if strings.HasPrefix(line, "Category:") {
				category := strings.Split(line, ":")
				line = "categories: " + category[1]
				metadataEnd = true
			}

			if strings.HasPrefix(line, ":Category:") {
				category := strings.Split(line, ":")
				line = "categories: " + category[2]
				metadataEnd = true
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

			_, err := f.WriteString(line + "\n")
			if err != nil {
				log.Fatal(err)
			}

			if metadataEnd {
				f.WriteString("---\n")
				metadataEnd = false
			}

			// Issue a `Sync` to flush writes to stable storage.
			f.Sync()
		}
	}

}

func main() {
	var articles []string
	err := filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".rst" || filepath.Ext(path) == ".md" {
			articles = append(articles, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	convert2Hugo(articles)

}
