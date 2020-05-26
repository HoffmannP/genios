package main

import (
	"os"

	"github.com/gocolly/colly/v2"
)

type download struct {
	generic
}

func newDownload(uri string) (d *download) {
	d.id = uri
	return
}

func (d *download) URL() string {
	return d.id
}

func (d *download) OnScraped(g *grabber, r *colly.Response) {
	filepath := "" // TODO
	err := r.Save(filepath)
	if err != nil {
		return
	}
	os.Chmod(filepath, 0644)
}
