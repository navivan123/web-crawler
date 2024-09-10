package main

import (
        "net/url"
        "errors"
        "fmt"
        "strings"
        "golang.org/x/net/html"
)

var ParseURLError  = errors.New("Error: Could Not Parse URLs!")
var SchemeURLError = errors.New("Error: Not a valid scheme!")
var PathURLError   = errors.New("Error: URL is Path!")
var NoURLsError    = errors.New("Error: No URLs detected!")

func normalizeURL(rawURL string) (string, error) {
    parsedURL, err := parseURL(rawURL)
    if err != nil {
        return "", err
    }

    normURL := parsedURL.Host + parsedURL.Path
    normURL  = strings.ToLower(normURL)
    normURL  = strings.TrimSuffix(normURL,"/")

    return normURL, nil
}

func parseURL(rawURL string) (*url.URL, error) {
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        return nil, ParseURLError
    }

    if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" && parsedURL.Scheme != "" {
        return nil, SchemeURLError
    }
    if parsedURL.Scheme == "" {
        return nil, PathURLError
    }

    return parsedURL, nil
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
    readHTML  := strings.NewReader(htmlBody) 
    nodesHTML, err := html.Parse(readHTML)
    if err != nil {
        return []string{}, err
    }
    
    links := recurseGetURLsFromHTML(nodesHTML, baseURL)
    if links == nil {
        return nil, NoURLsError
    }

    return links, nil
}

func recurseGetURLsFromHTML(node *html.Node, baseURL *url.URL) []string {
    links := []string{}

    //fmt.Printf("\nNode : %v\n\nType : %v\n\nData : %v\n\nNamespace : %v\n\nAttributes : %v\n", node, node.Type, node.Data, node.Namespace, node.Attr)
    if node.Data == "a" {
        for _, attr := range node.Attr {
            if attr.Key == "href" {
                _, err := parseURL(attr.Val)
                fmt.Printf("\n\n\nThis is url: %v\n\n\n", attr.Val)

                if err == nil || errors.Is(err, SchemeURLError){
                    links = append(links, attr.Val)
                }
                if errors.Is(err, PathURLError) {
                    links = append(links, baseURL.String() + attr.Val)
                }
            }
        }
    }

    if node.FirstChild == nil {
        return links
    }

    nodeChild := node.FirstChild
    for {
        if node.LastChild == nodeChild {
            links = append(links, recurseGetURLsFromHTML(nodeChild, baseURL)...)
            break
        }
        links = append(links, recurseGetURLsFromHTML(nodeChild, baseURL)...)
        nodeChild = nodeChild.NextSibling
    }

    return links
}
