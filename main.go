package main

import (
	"context"
	"fmt"
	deepseek "github.com/cohesion-org/deepseek-go"
	"github.com/gocolly/colly"
	"os"
)

// TODO : fix the systemprompt so it can handle a blogpost text, maybe to paraphrased it?

const systemPrompt = `You are an AI transformation agent tasked with converting raw YouTube caption texts about knowledge into a polished, engaging, and readable blog post. Your responsibilities include:
- **Paraphrasing**: Transform the original caption text into fresh, original content while preserving the key information and insights.
- **Structure**: Organize the content into a well-defined structure featuring a captivating introduction, clearly delineated subheadings in the body, and a strong conclusion.
- **Engagement**: Ensure the blog post is outstanding by using a professional yet conversational tone, creating smooth transitions, and emphasizing clarity and readability.
- **Retention of Key Elements**: Maintain all essential elements and core ideas from the original text, while enhancing the narrative to captivate the reader.
- **Adaptation**: Simplify technical details if necessary, ensuring that the transformed content is accessible to a broad audience without losing depth or accuracy.
- **Quality**: Aim for a high-quality article that is both informative and engaging, ready for publication.

Follow these guidelines to generate a comprehensive, coherent, and outstanding blog post from the provided YouTube captions text.
Your final output should be **only** the paraphrased text, styled in Markdown format.`

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
	c := colly.NewCollector()
	var Title string
	var Body Paragraph
	c.OnHTML("h1.text-4xl", func(e *colly.HTMLElement) {
		Title = e.Text
	})
	c.OnHTML("#content", func(e *colly.HTMLElement) {
		Body = Paragraph(e.Text)
	})
	_ = c.Visit(url)
	fmt.Println(Body)
	blogPost := createBlogpost(Title, Body)
	fmt.Println(blogPost.Title)
	fmt.Println(blogPost.Body)
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	// Create a chat completion request
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: deepseek.ChatMessageRoleUser, Content: blogPost.Title},
		},
	}
	// Send the request and handle the response
	deepseek_ctx := context.Background()
	response, err := client.CreateChatCompletion(deepseek_ctx, request)
	if err != nil {
		panic(err)
	}
	// Print the response
	output := response.Choices[0].Message.Content
	fmt.Println("Response:", output)
}
