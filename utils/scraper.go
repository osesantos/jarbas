package utils

import (
	"errors"
	"strings"

	"github.com/gocolly/colly"
	"github.com/osesantos/resulto"
)

func ScrapeText(url string) resulto.Result[string] {
	var extractedText strings.Builder

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.5")
		//		r.Headers.Set("Referer", "https://www.example.com")
		r.Headers.Set("DNT", "1")
		r.Headers.Set("Connection", "keep-alive")
	})

	// Extract text from paragraphs
	c.OnHTML("p", func(e *colly.HTMLElement) {
		extractedText.WriteString(e.Text + "\n")
	})

	err := c.Visit(url)
	if err != nil {
		return resulto.Failure[string](errors.New("failed to scrape text: " + err.Error()))
	}

	return resulto.Success(extractedText.String())
}
