package main

import "github.com/gocolly/colly/v2"

type generic struct {
	id string
}

func (x *generic) URL() string {
	return "/" + x.id
}

func (x *generic) OnHTML() (h HTMLHandler) {
	return
}

func (x *generic) OnResponse(g *grabber, r *colly.Response) {
	return
}

func (x *generic) OnScraped(g *grabber, r *colly.Response) {
	return
}

func (x *generic) Post() (p map[string]string) {
	return
}
