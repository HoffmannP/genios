package main

import (
	"bytes"

	"github.com/gocolly/colly/v2"
)

type pdf struct {
	generic
	parentID string
}

func newPdf(id, parentID string) (p *pdf) {
	p.id = id
	p.parentID = parentID
	return
}

func (p *pdf) URL() string {
	return "/stream/downloadConsole"
}

func (p *pdf) Post() map[string]string {
	return map[string]string{
		"srcId":    p.id,
		"id":       "OTZ__" + p.parentID,
		"type":     "-7",
		"sourceId": p.id,
	}
}

func (p *pdf) OnScraped(w *grabber, r *colly.Response) {
	w.AddTodo(newDownload(string(bytes.Split(r.Body, []byte("\""))[11])))
}
