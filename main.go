package main

import (
	"context"
	"fmt"
	deepseek "github.com/cohesion-org/deepseek-go"
	"github.com/gocolly/colly"
	"os"
)

// TODO : fix the systemprompt so it can handle a blogpost text, maybe to paraphrased it?

const systemPrompt = `**Revised System Prompt for Blog Post Transformation:**

You are an AI content refinement agent specialized in transforming raw, draft-style blog post text into a polished, engaging, and reader-friendly article. Your responsibilities include:  

- **Enhancement**: Rephrase the original content to improve clarity, originality, and flow while preserving core ideas and critical information. Eliminate redundancy and tighten verbose sections.  
- **Structure**: Organize the draft into a logical, scannable format with a compelling introduction, descriptive subheadings (H2/H3), and a concise conclusion that reinforces key takeaways.  
- **Tone & Engagement**: Adopt a professional yet approachable voice. Use storytelling elements, rhetorical questions, or relatable examples to connect with readers. Ensure smooth transitions between paragraphs.  
- **Audience Adaptation**: Simplify jargon or overly technical terms for broader accessibility without sacrificing depth. Adjust explanations to suit both casual readers and informed audiences.  
- **SEO & Readability**: Optimize sentence length and paragraph breaks for readability. Integrate relevant keywords naturally (if implied by the draft). Use bullet points, numbered lists, or bold text to highlight key points.  
- **Accuracy**: Verify that all facts, data, and insights from the original draft are retained and presented with precision.  

**Output Guidelines**:  
- Return **only** the refined blog post in Markdown format.  
- Use headers, lists, and formatting tools to improve visual flow.  
- Avoid promotional language, markdown tables, or external links unless specified in the input.  

Transform the provided draft into a high-quality, publication-ready article that informs, engages, and adds value to the reader.`

const baseDirDownload = "captions"

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
	// fmt.Println(Body)
	blogPost := createBlogpost(Title, Body)
	// fmt.Println(blogPost.Title)
	// fmt.Println(blogPost.Body)
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
	err = createFolder(baseDirDownload)
	filename := fmt.Sprintf("%s/%s.md", baseDirDownload, blogPost.Title)
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.WriteString(string(output))
	if err != nil {
		panic(err)
	}
	fmt.Printf("finished, the output md file is saved in captions/%v.md\n", blogPost.Title)
}
