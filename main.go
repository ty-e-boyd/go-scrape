package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
)

type Article struct {
	Url, Name string
}

func main() {
	// creating our industries slice, and setting up the new Colly Collector
	var articles []Article
	c := colly.NewCollector()

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"

	// our first (and only right now) OnHtml handler that will run on any successfully hit html files
	c.OnHTML(".crayons-story", func(e *colly.HTMLElement) {
		url := e.ChildAttr(".crayons-story__hidden-navigation-link", "href")
		name := e.ChildText(".crayons-story__hidden-navigation-link")

		fmt.Printf("found url: %v, and name: %v\n\n", url, name)

		if url != "" && name != "" {
			article := Article{
				Url:  url,
				Name: name,
			}

			fmt.Printf("article made::%v\n\n", article)
			articles = append(articles, article)
		} else {
			fmt.Println("found record, missing url or name")
		}
	})

	// going and attempting the scrape
	err := c.Visit("https://dev.to/")
	if err != nil {
		log.Fatal("error on visit --")
	}

	// some results of our scraping
	fmt.Printf("articles :::: %v\n\n", articles)
	fmt.Printf("number scraped : %v\n", len(articles))

	// going to write the results to a csv -> found in the root of the project rn
	file, err := os.Create("dev-to_articles.csv")
	if err != nil {
		log.Fatalln("error creating the csv", err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln("unable to close file stream", err.Error())
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"url", "name"}
	err = writer.Write(headers)
	if err != nil {
		log.Fatalln("unable to write headers to file", err.Error())
		return
	}

	for _, article := range articles {
		record := []string{article.Url, article.Name}

		err := writer.Write(record)
		if err != nil {
			log.Fatalln("unable to write record to file", err.Error(), record)
		}
	}
}
