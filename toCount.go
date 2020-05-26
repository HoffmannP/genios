package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type toCount struct {
	generic
	date time.Time
}

func newToCount(name string, date time.Time) (c *toCount) {
	c.id = name
	c.date = date
	return
}

func (c *toCount) URL() string {
	return fmt.Sprintf(
		"/toc_list/%s?issueName=%s&max=%d",
		c.id,
		c.date.Format("02.01.2006"),
		1,
	)
}

func (c *toCount) OnHTML() (h HTMLHandler) {
	h["#content > div > div.gridDetailLeft > div.moduleDirectoryHeader.clearfix > div.floatLeft > strong"] = c.handleCount
	return
}

func (c *toCount) handleCount(g *grabber, e *colly.HTMLElement) {
	count, err := strconv.Atoi(strings.Split(e.DOM.Text(), " ")[2])
	if err != nil {
		panic(err)
	}
	g.AddTodo(newTOC(c.id, c.date, count))
}
