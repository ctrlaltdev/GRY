package main

import (
	"fmt"
	"net/url"
)

func LogErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// ValidateURL checks if the provided URL string is valid
func ValidateURL(urlStr string) error {
	parsedURL, err := url.Parse(urlStr)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("invalid URL format")
	}
	return nil
}
