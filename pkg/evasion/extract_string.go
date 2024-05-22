package evasion

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func ExtractMatchedStringFromURL(url, pattern string) (string, error) {
	// Perform the GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Compile the regex pattern
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	// Find the first match in the text
	match := re.FindString(string(body))
	if match == "" {
		return "", fmt.Errorf("[!] No match found")
	}

	// Return the matched string
	return match, nil
}
