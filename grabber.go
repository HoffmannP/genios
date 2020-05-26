package main

import (
	"github.com/gocolly/colly/v2"
)

type grabber struct {
	id     string
	uri    string
	domain string
	c      *colly.Collector
	t      chan Todo
}

// Todo is an Interface for all the things that need to be done and can
type Todo interface {
	URL() string
	OnHTML() (h HTMLHandler)
	OnResponse(*grabber, *colly.Response)
	OnScraped(*grabber, *colly.Response)
	Post() map[string]string
}

// HTMLHandler is a collection of selectors and HTMLHandlers
type HTMLHandler map[string]func(*grabber, *colly.HTMLElement)

// ResponseHandler is a function to be called to handle a response
type ResponseHandler func(*grabber, colly.Response)

func newGrabber(domain string) (g *grabber) {
	g.domain = domain
	g.t = make(chan Todo, 100)
	return
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

func (g *grabber) AddTodo(todo Todo) {
	g.t <- todo
}

func (g *grabber) Run() {
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

func (g *grabber) do(t Todo) {
	c := g.c.Clone()
	for s, h := range t.OnHTML() {
		c.OnHTML(s, func(e *colly.HTMLElement) { h(g, e) })
	}
	c.OnResponse(func(r *colly.Response) { t.OnResponse(g, r) })
	c.OnScraped(func(r *colly.Response) { t.OnResponse(g, r) })
	p := t.Post()
	if len(p) > 0 {
		c.Post(g.domain+t.URL(), p)
	} else {
		c.Visit(g.domain + t.URL())
	}
}
