package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gocolly/colly/v2"
)

func newAuthenticatedCollector(username, password string) *colly.Collector {
	c := colly.NewCollector()
	err := c.SetCookies(geniosDomain, getCookies(username, password))
	if err != nil {
		panic(err)
	}
	return c
}

func getCookies(username, password string) []*http.Cookie {
	v := make(url.Values)
	v.Set("bibLoginLayer.number", username)
	v.Set("bibLoginLayer.password", password)
	v.Set("bibLoginLayer.terms_cb", "1")
	v.Set("bibLoginLayer.terms", "1")
	v.Set("bibLoginLayer.gdpr_cb", "1")
	v.Set("bibLoginLayer.gdpr", "1")
	v.Set("EVT.srcId", "bibLoginLayer_c0")
	v.Set("EVT.scrollTop", "0")
	v.Set("eventHandler", "loginClicked")
	v.Set("state", getCurrentState())
	r, err := http.PostForm(fmt.Sprintf("%s/formEngine/doAction", geniosDomain), v)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	return r.Cookies()
}

func getCurrentState() (state string) {
	c := colly.NewCollector()
	c.OnHTML("#layer_overlay + script + script", func(script *colly.HTMLElement) {
		state = script.Text[3847 : 3847+720]
	})
	c.Visit(geniosDomain)
	return
}
