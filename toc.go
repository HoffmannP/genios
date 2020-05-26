package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type toc struct {
	generic
	date  time.Time
	count int
}

func newTOC(name string, date time.Time, count int) *toc {
	return &toc{
		generic: generic{
			id: name,
		},
		date:  date,
		count: count,
	}
}

func (t *toc) Name() (string, map[string]string) {
	return reflect.TypeOf(t).Name(), map[string]string{"id": t.id, "date": t.date.Format("02.01.2006"), "count": strconv.Itoa(t.count)}
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
