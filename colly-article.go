/*
g {
	colly
	pubId
	domain

}
*/

func (a *article)New(id string) {
	a.id = id
	a.uri = strings.Sprintf("%%s/document/%%s__%s", id)
	a.images = make(string, 0)
}

func (c *article) Run(g* grabber) []handler {
	return [
		[".pdfAttachments div a", a.pdfHandler],
		[".divDocument > .moduleDocumentGraphic", a.imageHandler],
		[".divDocument", a.textHandler]
	]
}

func (a *article) pdfHandler(e *colly.HTMLElement) {
	id := e.Attr("id")
	g.AddTodo(WisoPdf.New(id))
	a.pdf = id
}

func (a *article) imageHandler(e *colly.HTMLElement) {
	uri := e.ChildAttr("a", "href")
	if uri == "" {
		println("Keine URI gefunden (" + name + ")")
		return
	}
	g.AddTodo(WisoImage.New(uri))
	a.images = append(a.images, uri)
}

func (a *article) textHandler(e *colly.HTMLElement) {
	a.title = e.ChildText(".boldLarge")
	a.text = e.ChildText(".text")
	a.author = e.ChildText(".italc")
	a.quelle = e.ChildText("table tr:nth-child(1) pre")
	a.resort = e.ChildText("table tr:nth-child(2) pre")
	a.edition = e.ChildText("table tr:nth-child(3) pre")
}