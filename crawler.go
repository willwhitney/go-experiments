package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"math/rand"
)

// if the given error exists, print it and panic
func errorify(err error) {
	if err != nil {
		fmt.Printf("%s", err)
		panic(err)
	}
}

// not used
func fetch(address string) (results string) {
	resp, err := http.Get(address)
	errorify(err)
	contents, err := ioutil.ReadAll(resp.Body)
	errorify(err)
	results = string(contents)
	
	return
}

func Crawl(url string, urls chan []string) (links []string) {
	doc, err := goquery.NewDocument(url)
	errorify(err)
	
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		value := s.Text()
		fmt.Println(value)
	})
	
	links = make([]string, 0, 100)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		
		// only take wikipedia links to real pages
		if exists {
			switch {
			case strings.Contains(link, "edit") || strings.Contains(link, "disambiguation"):
				return
			case strings.HasPrefix(link, "//"):
				return
			case strings.HasPrefix(link, "#"):
				return
			case strings.HasPrefix(link, "/") && strings.Contains(link, "wiki"):
				links = append(links, "http://en.wikipedia.org" + link)
			default:
				return
			}
		}
	})
	
	urls <- links
	return
	
}

func main() {
	visited := map[string] int {
		"http://en.wikipedia.org/wiki/Outer_space": 0,
	}
	urls := make(chan []string, 1000)
	liveThreads := 1
	maxThreads := 200
	go Crawl("http://en.wikipedia.org/wiki/Outer_space", urls)

	// listen to the channel `urls`, which pipes in lists of urls
	for options := range urls {
		liveThreads--
		
		// same as while(true) in other languages
		for {
			// pick a random option from the list of urls, and remove it from the list
			index := rand.Intn(len(options))
			option := options[index]
			options = append(options[:index], options[index:]...)
			
			// check it out in a new thread, if it's new to us
			_, seen := visited[option]
			if !seen {
				go Crawl(option, urls)
					liveThreads++
			}
			
			// only stop when we've filled our 100 threads or when the list is empty
			if liveThreads >= maxThreads || len(options) < 1 {
				break
			}
		}
	}
	
}
