package main

import (
	"flag"
	"log"
	"os"
	"text/template"
	"time"
)

func renderTemplate(filename string, tmpl *template.Template, data interface{}) (err error) {
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	err = tmpl.Execute(fd, data)
	return
}

func init() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	tmpl, err := template.ParseFiles("rep.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	data := map[string]interface{}{
		"StartTime": time.Now().Format("2006-01-02 15:04:05"),
	}
	err = renderTemplate("test.html", tmpl, data)
	if err != nil {
		log.Fatal(err)
	}
}
