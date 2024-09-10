package main

import (
        "testing"
        "errors"
        "reflect"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		expected      string
                errURL        error
	}{
		{
			name:     "remove scheme https no slash",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove scheme https slash",
			inputURL: "https://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove scheme http no slash",
			inputURL: "http://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove scheme http slash",
			inputURL: "http://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		}, 
		{
			name:     "invalid url",
			inputURL: "http//blog.boot.dev/path/",
			expected: "",
                        errURL: PathURLError,
		}, 
		{
			name:     "longer path",
                        inputURL: "http://blog.boot.dev/path/to/stuff/",
			expected: "blog.boot.dev/path/to/stuff",
		}, 

	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
                                if errors.Is(err, tc.errURL) && tc.expected == actual{
                                        return
                                }
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetURLs(t *testing.T) {
	tests := []struct {
                name        string
                inputURL    string
                inputBody   string
                expected  []string
                errURL      error       
        }{
                {
                        name:     "absolute and relative URLs",
                        inputURL: "https://blog.boot.dev",
                        inputBody: `
                        <html>
	                        <body>
		                        <a href="/path/one">
		                        	<span>Boot.dev</span>
	                        	</a>
	                        	<a href="https://other.com/path/one">
	                        		<span>Boot.dev</span>
	                        	</a>
                        	</body>
                        </html>
                        `,
                        expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
                },
                {
                        name:     "No URLs",
                        inputURL: "https://blog.boot.dev",
                        inputBody: `
                        <html>
	                        <body>
		                        <span>Boot.dev</span>
	                        	<span>Shoot.dev</span>
                        	</body>
                        </html>
                        `,
                        expected: []string{},
                        errURL: nil,
                },
                {
                        name:     "Nested URLs",
                        inputURL: "https://blog.boot.dev",
                        inputBody: `
                        <html>
	                        <body>
                                        <ul>
                                                <li>Check this first  link <a href="http://damn.com/s/">link</a></li>
                                                <li>Check this second link <a href="http://d.com/">link2</a></li>
                                                <li>Check this last link <a href="/path/shoot"></a></li>
                                        </ul>
		                        <span>Boot.dev</span>
	                        	<span>Shoot.dev</span>
                                        <h2>
                                                <p>Here is a quote from WWF's website (<a href="https://google.com/pig"></a>):</p>
                                                <p>lol jk</p>
                                        </h2>
                        	</body>
                        </html>
                        `,
                        expected: []string{"http://damn.com/s/", "http://d.com/", "https://blog.boot.dev/path/shoot", "https://google.com/pig"},
                },
        }

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
                                if errors.Is(err, tc.errURL) && reflect.DeepEqual(actual, tc.expected) {
                                        return
                                }
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
			}
                        if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
