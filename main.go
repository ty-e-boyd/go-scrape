package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
)

type Industry struct {
	Url, Image, Name string
}

func main() {
	// creating our industries slice, and setting up the new Colly Collector
	var industries []Industry
	c := colly.NewCollector()

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"

	// our first (and only right now) OnHtml handler that will run on any successfully hit html files
	c.OnHTML(".e-con-inner .section_cases_mobile_item", func(e *colly.HTMLElement) {
		url := e.ChildAttr(".elementor-image-box-img a", "href")
		image := e.ChildAttr(".elementor-image-box-img img", "data-lazy-src")
		name := e.ChildText(".elementor-image-box-content .elementor-image-box-title")

		fmt.Printf("found url: %v, image: %v, and name: %v\n\n", url, image, name)

		if url != "" && image != "" && name != "" {
			industry := Industry{
				Url:   url,
				Image: image,
				Name:  name,
			}

			fmt.Printf("industry made::%v\n\n", industry)
			industries = append(industries, industry)
		} else {
			fmt.Println("found record, missing url, image, or name")
		}
	})

	// going and attempting the scrape
	err := c.Visit("https://brightdata.com/")
	if err != nil {
		log.Fatal("error on visit --")
	}

	// some results of our scraping
	fmt.Printf("industries :::: %v\n\n", industries)
	fmt.Printf("number scraped : %v\n", len(industries))

	// going to write the results to a csv -> found in the root of the project rn
	file, err := os.Create("industries.csv")
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

	headers := []string{"url", "image", "name"}
	err = writer.Write(headers)
	if err != nil {
		log.Fatalln("unable to write headers to file", err.Error())
		return
	}

	for _, industry := range industries {
		record := []string{industry.Url, industry.Image, industry.Name}

		err := writer.Write(record)
		if err != nil {
			log.Fatalln("unable to write record to file", err.Error(), record)
		}
	}
}
