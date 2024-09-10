package main

import (
    "net/http"
    "io"
    "errors"
    "fmt"
)


var CodeHTTPError   = errors.New("Error: HTTP Response Returned Error-Level Code!")
var BadContentTypeHTTPError   = errors.New("Error: Content-Type mismatch!")

func getHTML(rawURL string) (string, error) {
    res, err := http.Get(rawURL)
    if err != nil {
        return "", err
    }
    defer res.Body.Close()

    if res.StatusCode >= http.StatusBadRequest {
        return "", fmt.Errorf(": %v | %v | Status Code: %v", rawURL, CodeHTTPError, res.StatusCode)
    }
    if res.Header.Get("Content-Type")[0:9] != "text/html" {
        return "", fmt.Errorf(": %v | %v | Content-Type: %v | Body: %v", rawURL, BadContentTypeHTTPError, res.Header.Get("Content-Type"), res.Body)
    }

    data, err := io.ReadAll(res.Body)
    if err != nil {
        return "", err
    }

    return string(data), nil
}
