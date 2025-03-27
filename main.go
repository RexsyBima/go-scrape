package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

type Paragraph string
type Blogpost struct {
	Title string
	Body  Paragraph
}

func createBlogpost(title string, body Paragraph) Blogpost {
	return Blogpost{
		Title: title,
		Body:  Paragraph(body),
	}
}

func main() {
	url := os.Args[1]
	fmt.Println(url)
	panic("exit here")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.StatusCode)
	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("expected 200 but got %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	fmt.Println(buf.String())

}
