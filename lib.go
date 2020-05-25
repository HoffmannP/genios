package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const geniosDomain = "https://bib-jena.genios.de"
const location = "/usr/local/share/wiso"

func main() {
	c := newAuthenticatedCollector("L0075062", "14092010")
	newspaper := "OTZ"
	today := time.Now()
	base := location + "/" + newspaper + "/" + today.Format("2006-01-02")
	if err := os.Mkdir(base, 0755); err != nil && !os.IsExist(err) {
		panic(err)
	}
	/*
		count := getCount(c, newspaper, today)
		entries := getTOC(c, newspaper, today, count, 0)
		createArticles(entries, base)
	*/
	if err := loadArticles(c, base); err != nil {
		panic(err)
	}

}

func createArticles(entries []string, base string) {
	os.Mkdir(base+"/new", 0755)
	for _, entry := range entries {
		createArticle(entry, base+"/new")
	}
}

func createArticle(uri, target string) {
	// /document/OTZ__doc7anuoz9qdxj19tbki17ao/toc/0?all -> doc7anuoz9qdxj19tbki17ao
	name := strings.Split(uri, "/")[2][5:]
	handler, err := os.OpenFile(target+"/"+name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	handler.Close()
}

func loadArticles(c *colly.Collector, base string) error {
	dir, err := os.Open(base + "/new")
	if err != nil {
		return err
	}
	files, err := dir.Readdirnames(1000)
	if err != nil {
		return err
	}
	for i, name := range files {
		println(fmt.Sprintf(
			"%s/document/OTZ__%s",
			geniosDomain,
			name,
		))
		downloadFile(
			c,
			base+"/"+name+".html",
			fmt.Sprintf(
				"%s/document/OTZ__%s",
				geniosDomain,
				name,
			),
		)
		// loadArticle(c, base, name)
		if i > 4 {
			break
		}
	}
	return os.Remove(base + "/new")
}

func loadArticle(c *colly.Collector, base, name string) {
	var article map[string]string
	images := 1
	c.OnHTML(".pdfAttachments div a", handlePdf(c.Clone(), base, name))
	c.OnHTML(".divDocument > .moduleDocumentGraphic", func(e *colly.HTMLElement) {
		uri := e.ChildAttr("a", "href")
		if uri == "" {
			println("Keine URI gefunden (" + name + ")")
			return
		}
		err := downloadFile(c.Clone(), fmt.Sprintf("%s/%s_%02d.jpg", base, name, images), uri)
		if err != nil {
			panic(err)
		}
		images++
	})
	c.OnHTML(".divDocument", func(e *colly.HTMLElement) {
		article = map[string]string{
			"title":   e.ChildText(".boldLarge"),
			"text":    e.ChildText(".text"),
			"author":  e.ChildText(".italc"),
			"quelle":  e.ChildText("table tr:nth-child(1) pre"),
			"resort":  e.ChildText("table tr:nth-child(2) pre"),
			"edition": e.ChildText("table tr:nth-child(3) pre"),
		}
	})
	c.OnScraped(writeArticle(base, name, article, images))
	c.Visit(fmt.Sprintf(
		"%s/document/OTZ__%s",
		geniosDomain,
		name,
	))
}

func handlePdf(c *colly.Collector, base, name string) colly.HTMLCallback {
	return func(e *colly.HTMLElement) {
		uri := downloadPdf(c.Clone(), name, e.Attr("id"))
		err := downloadFile(c.Clone(), fmt.Sprintf("%s/%s.pdf", base, name), uri)
		if err != nil {
			panic(err)
		}
	}
}

func writeArticle(base, name string, article map[string]string, images int) colly.ScrapedCallback {
	return func(r *colly.Response) {
		err := writeHugoFile(base, name, article, images)
		if err != nil {
			panic(err)
		}
		os.Remove(base + "/new/" + name)
	}
}

func downloadPdf(c *colly.Collector, name, id string) (uri string) {
	c.OnResponse(func(r *colly.Response) {
		uri = geniosDomain + string(bytes.Split(r.Body, []byte("\""))[11])
	})
	c.Post(
		fmt.Sprintf("%s/stream/downloadConsole", geniosDomain),
		map[string]string{
			"srcId":    id,
			"id":       "OTZ__" + name,
			"type":     "-7",
			"sourceId": id,
		},
	)
	return
}

func downloadFile(c *colly.Collector, filepath string, uri string) (err error) {
	c.OnResponse(func(r *colly.Response) {
		err = r.Save(filepath)
		if err != nil {
			return
		}
		err = os.Chmod(filepath, 0644)
	})
	err = c.Visit(uri)
	return
}

func writeHugoFile(base, name string, article map[string]string, images int) error {
	filename := fmt.Sprintf("%s/%s.md", base, name)
	f, err := os.Create(filename)
	os.Chmod(filename, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	imageFiles := make([]string, images-1)
	for i := 1; i < images; i++ {
		imageFiles[i-1] = fmt.Sprintf("\"%s_%02d.jpg\", ", name, i)
	}
	_, err = fmt.Fprintf(f,
		"+++\n"+
			"title: \"%s\"\n"+
			"author: \"%s\"\n"+
			"images: %s\n"+
			"resources: \"%s\"\n"+
			"categories: %s\n"+
			"tags: %s\n"+
			"+++\n"+
			"\n"+
			"%s\n"+
			"%s",
		article["title"],
		article["author"],
		printList(imageFiles),
		name+".pdf",
		article["resort"],
		printList(strings.Split(article["edition"], "; ")),
		article["text"],
		article["quelle"],
	)
	return err
}

func printList(l []string) (s string) {
	s = "["
	for _, e := range l {
		s += "\"" + e + "\", "
	}
	if len(l) > 0 {
		s = s[:len(s)-2]
	}
	s += "]"
	return
}

func getCount(c *colly.Collector, name string, date time.Time) (count int) {
	selector := "#content > div > div.gridDetailLeft > div.moduleDirectoryHeader.clearfix > div.floatLeft > strong"
	c.OnHTML(selector, func(e *colly.HTMLElement) {
		innerText := e.DOM.First().Text()
		count, _ = strconv.Atoi(strings.Split(innerText, " ")[2])
	})
	c.Visit(fmt.Sprintf(
		"%s/toc_list/%s?issueName=%s&max=1",
		geniosDomain,
		name,
		date.Format("02.01.2006"),
	))
	return
}

func getTOC(c *colly.Collector, name string, date time.Time, count, offset int) (entries []string) {
	selector := "a.boxDocumentTitle"
	c.OnHTML(selector, func(e *colly.HTMLElement) {
		href, exists := e.DOM.First().Attr("href")
		if exists {
			entries = append(entries, href)
		}
	})
	c.Visit(fmt.Sprintf(
		"%s/toc_list/%s?issueName=%s&offset=%d&max=%d",
		geniosDomain,
		name,
		date.Format("02.01.2006"),
		offset,
		count,
	))
	return
}
