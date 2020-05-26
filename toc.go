package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type toc struct {
	generic
	date  time.Time
	count int
}

func newTOC(name string, date time.Time, count int) (t *toc) {
	t.id = name
	t.date = date
	t.count = count
	return
}

func (t *toc) URL() string {
	return fmt.Sprintf(
		"/toc_list/%s?issueName=%s&offset=%d&max=%d",
		t.id,
		t.date.Format("02.01.2006"),
		0,
		t.count,
	)
}

func (t *toc) OnHTML() (h HTMLHandler) {
	h["a.boxDocumentTitle"] = t.handleTOC
	return
}

func (t *toc) handleTOC(g *grabber, e *colly.HTMLElement) {
	href, exists := e.DOM.Attr("href")
	if !exists {
		return
	}
	g.AddTodo(newArticle(strings.Split(href, "/")[2]))
}
