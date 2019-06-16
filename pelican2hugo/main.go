package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func convert2Hugo(articles []string, targetDir string) {
	for _, article := range articles {
		log.Printf("Converting %v to hugo\n", article)
		file, err := os.Open(article)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		fName := filepath.Base(article)
		fileExtension := filepath.Ext(article)
		filename := strings.TrimSuffix(fName, fileExtension)

		f, err := os.Create(filepath.Join(targetDir, filename+fileExtension))
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
		// begin frontmatter for hugo
		f.WriteString("---\n")
		metadataEnd := false

		scanner := bufio.NewScanner(file)
		var line string
		var postTitle string

		for scanner.Scan() {
			line = string(scanner.Text())

			if strings.HasPrefix(line, "Title:") {
				title := strings.Split(line, ":")
				line = "title: "
				postTitle = ""
				for _, t := range title[1:len(title)] {
					postTitle += t
				}
				line += postTitle
			}
			if strings.HasPrefix(line, ":Title:") {
				title := strings.Split(line, ":")
				line = "title: "
				postTitle = ""
				for _, t := range title[2:len(title)] {
					postTitle += t
				}
				line += postTitle
			}

			if strings.HasPrefix(line, "Date:") {
				date := strings.Split(line, ":")
				date[1] = strings.TrimPrefix(date[1], " ")
				// date[1] can either be 2018-11-09 hh:mm or 2018-11-09
				dateSplit := strings.Split(date[1], " ")
				if len(dateSplit) == 1 {
					line = "date: " + date[1]
				}
				if len(dateSplit) == 2 {
					line = "date: " + dateSplit[0]
				}
			}
			if strings.HasPrefix(line, ":Date:") {
				date := strings.Split(line, ":")
				date[2] = strings.TrimPrefix(date[2], " ")
				// date[2] can either be 2018-11-09 hh:mm or 2018-11-09
				dateSplit := strings.Split(date[2], " ")
				if len(dateSplit) == 1 {
					line = "date: " + date[2]
				}
				if len(dateSplit) == 2 {
					line = "date: " + dateSplit[0]
				}
			}

			if strings.HasPrefix(line, "Category:") {
				category := strings.Split(line, ":")
				line = "categories:\n- " + category[1]
				metadataEnd = true
			}

			if strings.HasPrefix(line, ":Category:") {
				category := strings.Split(line, ":")
				line = "categories:\n- " + category[2]
				metadataEnd = true
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

			if strings.HasPrefix(line, ".. code-include::") {
				f.WriteString(".. code::\n\n")
				filePath := strings.Split(line, "::")[1]
				filePath = strings.TrimPrefix(filePath, " ")
				filePath = strings.TrimRight(filePath, "\n")

				f1, err := os.Open(filePath)
				if err != nil {
					log.Fatal(err)
				}
				defer f1.Close()

				scanner := bufio.NewScanner(f1)

				for scanner.Scan() {
					l := "   " + string(scanner.Text()) + "\n"
					_, err := f.WriteString(l)
					if err != nil {
						log.Fatal(err)
					}
				}
			} else if strings.HasPrefix(line, ":lexer:") {
				continue
			} else {

				_, err := f.WriteString(line + "\n")
				if err != nil {
					log.Fatal(err)
				}
			}

			if metadataEnd {
				_, err := f.WriteString("aliases:\n")
				if err != nil {
					log.Fatal(err)
				}

				postTitle = strings.TrimPrefix(postTitle, " ")

				dashed := strings.ReplaceAll(postTitle, " ", "-")
				dashedLower := strings.ToLower(dashed)

				specialChars := []string{".", "'", "+", "`", ":", "/", "!", ",", "(", ")", "?"}
				transformedString := dashedLower
				for _, c := range specialChars {
					transformedString = strings.ReplaceAll(transformedString, c, "")
				}

				_, err = f.WriteString(fmt.Sprintf("- /%s\n", transformedString+".html"))
				if err != nil {
					log.Fatal(err)
				}
				f.WriteString("---\n")
				metadataEnd = false
			}

			// Issue a `Sync` to flush writes to stable storage.
			f.Sync()
		}
	}

}

func main() {

	postsDir := os.Args[1]

	// todo: pages, notes, etc
	var articles []string
	files, err := ioutil.ReadDir(postsDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".rst" || filepath.Ext(file.Name()) == ".md" {
			articles = append(articles, filepath.Join(postsDir, file.Name()))
		}
	}
	convert2Hugo(articles, os.Args[2])
}
