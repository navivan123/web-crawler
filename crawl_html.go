package main

import (
    "fmt"
    "net/url"
    "sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
    maxPages           int 
}


func (cfg *config) crawlPage(rawCurrentURL string) {
    cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()
    
    if cfg.pagesLen() >= cfg.maxPages {
		return
	}
	
    currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - crawlPage: couldn't parse URL '%s': %v\n", currentURL, err)
		return
	}

	// skip other websites
	if currentURL.Hostname() != cfg.baseURL.Hostname() {
		return
	}

    normCurrentURL, err := normalizeURL(rawCurrentURL)
    if err != nil {
        return 
    }

    isFirst := cfg.addPageVisit(normCurrentURL)
    if !isFirst {
        return
    }

    html, err := getHTML(rawCurrentURL)
    if err != nil {
        fmt.Printf("Error while getting HTML from URL", err)
        return 
    }

    urls, err := getURLsFromHTML(html, cfg.baseURL)
    if err != nil {
        fmt.Printf("Error while getting URLs from HTML", err)
        return 
    }
    //fmt.Printf("\nURLs: %v\n", urls)
    for _, url := range urls {
        cfg.wg.Add(1)
        go cfg.crawlPage(url)
    }
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, visited := cfg.pages[normalizedURL]; visited {
		cfg.pages[normalizedURL]++
		return false
	}

	cfg.pages[normalizedURL] = 1
	return true
}

func configure(rawBaseURL string, maxConcurrency, maxPages int) (*config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	return &config{
		pages:              make(map[string]int),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
        maxPages:           maxPages,
	}, nil
}


func (cfg *config) pagesLen() int {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	return len(cfg.pages)
}
