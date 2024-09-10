package main

import (
    "fmt"
    "os"
    "strconv"
    "sort"
)

func main() {
    if len(os.Args) < 4 {
        fmt.Println("no args provided")
        fmt.Println("usage: crawler <baseURL> <maxConcurrency> <maxPages>")
        os.Exit(1)
    }
    if len(os.Args) > 4 {
        fmt.Println("too many arguments provided")
        fmt.Println("usage: crawler <baseURL> <maxConcurrency> <maxPages>")
        os.Exit(1)
    }

    baseURL        := os.Args[1]
    arg2           := os.Args[2]
    arg3           := os.Args[3]

    maxConcurrency, err := strconv.Atoi(arg2)
    if err != nil { 
        fmt.Println("Error parsing max concurrency")
        os.Exit(1)
    }
    
    maxPages, err := strconv.Atoi(arg3)
    if err != nil { 
        fmt.Println("Error parsing max pages")
        os.Exit(1)
    }



    fmt.Printf("starting crawl of: %v\n", baseURL)
    
	cfg, err := configure(baseURL, maxConcurrency, maxPages)
	if err != nil {
		fmt.Printf("Error - configure: %v", err)
		return
	}

    cfg.wg.Add(1)
    go cfg.crawlPage(baseURL)
    cfg.wg.Wait()


    printReport(cfg.pages, baseURL)
}

type PageResults struct {
    site     string
    numLinks int
}

func printReport(pages map[string]int, baseURL string) {
    links := make([]PageResults, 0)
    for url, count := range pages {
        links = append(links, PageResults{site: url, numLinks: count})
    }
    
    sort.Slice(links, func(i, j int) bool  {
        return links[i].numLinks > links[j].numLinks 
    })

    sort.Slice(links, func(i, j int) bool  {
        return links[i].site < links[j].site && links[i].numLinks == links[j].numLinks
    })

    fmt.Printf("=============================\nREPORT for %v\n=============================\n\n", baseURL)

    for _, link := range links {
        fmt.Printf("Found %v internal links to %v\n", link.numLinks, link.site)
    }
}
