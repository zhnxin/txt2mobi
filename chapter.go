package main

import (
	"bytes"
	"html"

	"github.com/zhnxin/mobi"
)

var (
	isEscape    = true
	isParagraph = false
	LineFeed    = []byte("<br/>")
	PStart      = []byte("<p>")
	PEnd        = []byte("</p>")
	Blank       = []byte{}
)

func SetIsEscape(flag bool) {
	isEscape = flag
}
func SetIsParagraph(flag bool) {
	isParagraph = flag
}

type chapterContent struct {
	Title   string
	Content []byte
}

func (c *chapterContent) Append(content []byte) {
	if !isEscape {
		content = []byte(html.EscapeString(string(content)))
	}
	if isParagraph {
		if len(content) > 1 {
			c.Content = append(c.Content, bytes.Join([][]byte{PStart, content, PEnd}, Blank)...)
		} else {
			c.Content = append(c.Content, LineFeed...)
		}
	} else {
		c.Content = append(c.Content, LineFeed...)
		if len(content) > 1 {
			c.Content = append(c.Content, content...)
		}
	}

}

func (c *chapterContent) SetTitle(title string) {
	c.Title = title
}

func (c *chapterContent) Restore(title string) {
	c.Title = title
	c.Content = make([]byte, 0)
}

type Chapter struct {
	content        chapterContent
	subChapterList []chapterContent
	subLength      uint
}

func NewChapter(title string) Chapter {
	return Chapter{
		content: chapterContent{
			Title:   title,
			Content: []byte{},
		},
		subChapterList: make([]chapterContent, 0),
	}
}

func (c *Chapter) Restore(title string) {
	c.subLength = 0
	c.content.Title = title
	c.content.Content = make([]byte, 0)
	c.subChapterList = make([]chapterContent, 0)
}

func (c *Chapter) AddSubChapter(title string) {
	subChapter := chapterContent{
		Title:   title,
		Content: make([]byte, 0),
	}
	c.subChapterList = append(c.subChapterList, subChapter)
	c.subLength += 1
}

func (c *Chapter) Append(content []byte) {
	if c.subLength < 1 {
		c.content.Append(content)
	} else {
		c.subChapterList[c.subLength-1].Append(content)
	}
}

func (c *Chapter) Flush(mobiWriter mobi.Builder) {
	chapter := mobiWriter.NewChapter(c.content.Title, c.content.Content)
	for _, content := range c.subChapterList {
		chapter.AddSubChapter(content.Title, content.Content)
	}
}
