package logic

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	charset      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	minUrlLength = 5
	maxUrlLength = 10
)

func IsRealURL(u string) (bool, error) {
	if u == "" {
		return false, nil
	}

	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = "http://" + u
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Head(u)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 300, nil
}

func ValidateUrl(u string) bool {
	if u == "" {
		return false
	}
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = "http://" + u
	}
	parsedURL, err := url.ParseRequestURI(u)
	if err != nil {
		return false
	}
	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

func generateRandomString(length int) string {
	shortURL := make([]byte, length)
	for i := range shortURL {
		shortURL[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortURL)
}

func generateUrl(db *sql.DB) string {
	var shortenedUrl string
	for {
		shortenedUrl = generateRandomString(rand.Int()%minUrlLength + (maxUrlLength - minUrlLength))

		var exists int
		err := db.QueryRow("SELECT COUNT(1) FROM urls WHERE short = ?", shortenedUrl).Scan(&exists)
		if err != nil {
			log.Fatal(err)
		}

		if exists == 0 {
			break
		}
	}

	return shortenedUrl
}

func AddUrlToDb(url string, db *sql.DB) string {
	var existingShort string
	err := db.QueryRow("SELECT short FROM urls WHERE url = ?", url).Scan(&existingShort)
	if err == nil {
		fmt.Printf("Shortened URL: %s\n", existingShort)
		return existingShort
	}

	if err != sql.ErrNoRows {
		log.Fatal(err)
	}

	shortenedUrl := generateUrl(db)

	_, err = db.Exec("INSERT INTO urls (url, short) VALUES (?, ?)", url, shortenedUrl)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Shortened URL: %s\n", shortenedUrl)

	return shortenedUrl
}
