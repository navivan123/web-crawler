package main

import (
    "fmt"
    "os"
    "strconv"
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

    for url, count := range cfg.pages {
        fmt.Printf("\nPage! => URL: %v | # Visited: %v\n", url, count)
    }
}

func printReport(pages map[string]int, baseURL string) {
    fmt.Printf("=============================\nREPORT for %v\n=============================", baseURL)
}
