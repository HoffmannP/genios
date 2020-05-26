package main

import (
	"log"

	"github.com/gocolly/colly/v2"
)

type grabber struct {
	id     string
	uri    string
	domain string
	auth   bool
	c      *colly.Collector
	t      chan todo
}

// Todo is an Interface for all the things that need to be done and can
type todo interface {
	URL() string
	Name() (string, map[string]string)
	OnHTML() (h HTMLHandler)
	OnResponse(*grabber, *colly.Response)
	OnScraped(*grabber, *colly.Response)
	Post() map[string]string
}

// HTMLHandler is a collection of selectors and HTMLHandlers
type HTMLHandler map[string]func(*grabber, *colly.HTMLElement)

// ResponseHandler is a function to be called to handle a response
type ResponseHandler func(*grabber, colly.Response)

func newGrabber(domain string) *grabber {
	return &grabber{
		domain: domain,
		t:      make(chan todo, 100),
	}
}

func (g *grabber) Authenticate(username, password string) {
	g.c = colly.NewCollector()
	state := make(chan string)
	go func() {
		g.c.OnHTML(
			"#layer_overlay + script + script",
			func(script *colly.HTMLElement) {
				state <- script.Text[3847 : 3847+720]
			},
		)
		g.c.Visit(g.domain)
	}()

	g.c.Post(
		g.domain+"/stream/downloadConsole",
		map[string]string{
			"bibLoginLayer.number":   username,
			"bibLoginLayer.password": password,
			"bibLoginLayer.terms_cb": "1",
			"bibLoginLayer.terms":    "1",
			"bibLoginLayer.gdpr_cb":  "1",
			"bibLoginLayer.gdpr":     "1",
			"EVT.srcId":              "bibLoginLayer_c0",
			"EVT.scrollTop":          "0",
			"eventHandler":           "loginClicked",
			"state":                  <-state,
		},
	)
}

func (g *grabber) AddTodo(t todo) {
	typename, parameter := t.Name()
	log.Printf("Queueing %s to todo [%+v]", typename, parameter)
	g.t <- t
}

func (g *grabber) Run() {
	if !g.auth {
		panic("Not authenticated yet!") // TODO
	}
	active := true
	for active {
		select {
		case t := <-g.t:
			g.do(t)
		default:
			active = false
		}
	}
}

func (g *grabber) do(t todo) {
	c := g.c.Clone()
	for s, h := range t.OnHTML() {
		c.OnHTML(s, func(e *colly.HTMLElement) { h(g, e) })
	}
	c.OnResponse(func(r *colly.Response) { t.OnResponse(g, r) })
	c.OnScraped(func(r *colly.Response) { t.OnResponse(g, r) })
	p := t.Post()
	typename, parameter := t.Name()
	log.Printf("Running %s to todo [%+v]", typename, parameter)
	if len(p) > 0 {
		c.Post(g.domain+t.URL(), p)
	} else {
		c.Visit(g.domain + t.URL())
	}
}
