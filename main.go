package main

import (
	"time"
)

// const location = "/usr/local/share/wiso"

func main() {
	grabber := newGrabber("https://bib-jena.genios.de")
	grabber.AddTodo(newToCount("OTZ", time.Now()))
	grabber.Run()
}

/*
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
*/
