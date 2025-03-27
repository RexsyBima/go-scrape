package main

import (
	"fmt"

	"github.com/gocolly/colly"
	// "strings"
	"testing"
)

func Test_main(t *testing.T) {
	tests := []struct {
		url      string // description of this test case
		expected string // description of this test case
	}{
		{"https://maximilian-schwarzmueller.com/articles/whats-the-mcp-model-context-protocol-hype-all-about/", " Looking Beyond the Hype: What's MCP All About? "},
		{"https://maximilian-schwarzmueller.com/articles/what-are-machine-and-byte-code-anyways/", " What Are \"Machine Code\" & \"Byte Code\" Anyways? "},
	}
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			c := colly.NewCollector()
			var Title string
			var Body Paragraph
			c.OnHTML("h1.text-4xl", func(e *colly.HTMLElement) {
				Title = e.Text
			})
			c.OnHTML("#content", func(e *colly.HTMLElement) {
				Body = Paragraph(e.Text)
			})
			_ = c.Visit(tt.url)
			if Title != tt.expected {
				t.Errorf("expected %s but got %s", tt.expected, Title)
			}
			fmt.Println(Body)

		})
	}
}
