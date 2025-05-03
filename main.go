package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

var maxConcurrentConnections = 10
var maxPages = 100
var usageString = "usage: ./crawler URL [maxConcurrency] [maxPages]"

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {

	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		fmt.Println(usageString)
		os.Exit(1)
	}

	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		fmt.Println(usageString)
		os.Exit(1)
	}

	rawBaseURL := strings.Trim(args[0], "/")
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Println("cannot parse base url")
		os.Exit(1)
	}

	if len(args) > 1 {
		maxConcurrentConnections, err = strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("invalid maxConcurrency integer")
			fmt.Println(usageString)
			os.Exit(1)
		}
	}

	if len(args) > 2 {
		maxPages, err = strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("invalid maxPages integer")
			fmt.Println(usageString)
			os.Exit(1)
		}
	}

	cfg := config{
		pages:              make(map[string]int),
		baseURL:            baseURL,
		concurrencyControl: make(chan struct{}, maxConcurrentConnections),
		wg:                 &sync.WaitGroup{},
		mu:                 &sync.Mutex{},
		maxPages:           maxPages,
	}

	fmt.Printf("starting crawl of: %s with maxConcurrency: %v and maxPages: %v\n", baseURL, maxConcurrentConnections, cfg.maxPages)

	cfg.wg.Add(1)
	go cfg.crawlPage(rawBaseURL)

	cfg.wg.Wait()
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

func (cfg *config) crawlPage(rawCurrentURL string) {

	defer cfg.wg.Done()
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return
	}

	if currentURL.Host != cfg.baseURL.Host {
		return
	}

	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	cfg.mu.Lock()
	_, exists := cfg.pages[normalizedCurrentURL]

	if exists {
		cfg.pages[normalizedCurrentURL] += 1
	} else {
		cfg.pages[normalizedCurrentURL] = 1
	}

	cfg.mu.Unlock()

	if exists {
		return
	}

	fmt.Printf("Getting HTML for %s\n", rawCurrentURL)

	cfg.concurrencyControl <- struct{}{}
	defer func() { <-cfg.concurrencyControl }()

	pageHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}

	fmt.Printf("Got HTML for %s\n", rawCurrentURL)

	pageURLs, err := getURLsFromHTML(pageHTML, cfg.baseURL.String())
	if err != nil {
		return
	}

	for _, link := range pageURLs {
		cfg.mu.Lock()
		if len(cfg.pages) >= cfg.maxPages {
			cfg.mu.Unlock()
			return
		}
		cfg.mu.Unlock()

		cfg.wg.Add(1)
		go cfg.crawlPage(link)
	}
}
