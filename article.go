package main

import (
	"github.com/gocolly/colly/v2"
)

type article struct {
	generic
	images  []string
	pdf     string
	title   string
	text    string
	author  string
	quelle  string
	resort  string
	edition string
	page    int
}

func newArticle(id string) *article {
	return &article{
		generic: generic{
			id: id,
		},
		images: make([]string, 1),
	}
}

func (a *article) URL() string {
	return "/document/" + a.id
}

func (a *article) OnHTML() HTMLHandler {
	return HTMLHandler{
		".pdfAttachments div a":                 a.pdfHandler,
		".divDocument > .moduleDocumentGraphic": a.imageHandler,
		"divDocument":                           a.textHandler,
	}
}

func (a *article) pdfHandler(g *grabber, e *colly.HTMLElement) {
	id := e.Attr("id")
	g.AddTodo(newPdf(id, a.id))
	a.pdf = id
}

func (a *article) imageHandler(g *grabber, e *colly.HTMLElement) {
	uri := e.ChildAttr("a", "href")
	if uri == "" {
		println("Keine URI gefunden (" + a.id + ")")
		return
	}
	g.AddTodo(newDownload(uri))
	a.images = append(a.images, uri)
}

func (a *article) textHandler(g *grabber, e *colly.HTMLElement) {
	a.title = e.ChildText(".boldLarge")
	a.text = e.ChildText(".text")
	a.author = e.ChildText(".italc")
	a.quelle = e.ChildText("table tr:nth-child(1) pre")
	a.resort = e.ChildText("table tr:nth-child(2) pre")
	a.edition = e.ChildText("table tr:nth-child(3) pre")
}

func (a *article) OnScrap() {
	// TODO: Write Article to File
}
