package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func main() {

	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL := args[0]

	fmt.Printf("starting crawl of: %s\n", baseURL)

	pages := make(map[string]int)

	crawlPage(baseURL, baseURL, pages)
}

func getHTML(rawURL string) (string, error) {
	res, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if !strings.Contains(res.Header.Get("content-type"), "text/html") {
		return "", errors.New("content-type")
	}
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return
	}

	if currentURL.Host != baseURL.Host {
		return
	}

	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	_, exists := pages[normalizedCurrentURL]

	if exists {
		pages[normalizedCurrentURL] += 1
		return
	} else {
		pages[normalizedCurrentURL] = 1
	}

	fmt.Printf("Getting HTML for %s\n", rawCurrentURL)

	pageHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}

	fmt.Printf("Got HTML for %s\n", rawCurrentURL)

	pageURLs, err := getURLsFromHTML(pageHTML, rawBaseURL)
	if err != nil {
		return
	}

	for _, link := range pageURLs {
		crawlPage(rawBaseURL, link, pages)
	}
}
