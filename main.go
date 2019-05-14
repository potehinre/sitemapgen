package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/potehinre/sitemapgen/sitemap"
	"github.com/potehinre/sitemapgen/spider"
)

func main() {
	host := flag.String("url", "http://golang.org", "url to generate sitemap")
	maxdepth := flag.Int("maxdepth", 1, "max crawl depth")
	workerCount := flag.Int("parallel", 8, "crawler worker count")
	outputFile := flag.String("output", "sitemap.xml", "name of result sitemap file")
	flag.Parse()
	client := &http.Client{}
	urls := spider.Run(client, *host, *maxdepth, *workerCount)
	output := sitemap.Generate(urls)
	if err := ioutil.WriteFile(*outputFile, output, 0644); err != nil {
		log.Fatalf("Error writing to a file %s: %s", *outputFile, err)
	}
}
