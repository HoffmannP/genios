package main

import (
	"bytes"
	"reflect"

	"github.com/gocolly/colly/v2"
)

type pdf struct {
	generic
	parentID string
}

func newPdf(id, parentID string) (p *pdf) {
	return &pdf{
		generic: generic{
			id: id,
		},
		parentID: parentID,
	}
}

func (p *pdf) Name() (string, map[string]string) {
	return reflect.TypeOf(p).Name(), map[string]string{"id": p.id, "parentId": p.parentID}
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
