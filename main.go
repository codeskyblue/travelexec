package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func ignore(info os.FileInfo) bool {
	if info.IsDir() {
		if info.Name() != "." && info.Name() != ".." &&
			strings.HasPrefix(info.Name(), ".") { // ignore hidden dir
			return true
		}
	} else {
		return strings.HasPrefix(info.Name(), ".")
	}
	return false
}

func pathWalk(path string, depth int) (files []string, err error) {
	files = make([]string, 0)
	baseNumSeps := strings.Count(path, string(os.PathSeparator))
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			pathDepth := strings.Count(path, string(os.PathSeparator)) - baseNumSeps
			if pathDepth > depth {
				return filepath.SkipDir
			}
			if ignore(info) {
				return filepath.SkipDir
			}
		} else if info.Mode().IsRegular() && !ignore(info) {
			files = append(files, path)
		}
		return nil
	})
	return
}

var (
	depth    = flag.Int("depth", 0, "depth of dir search")
	path     = flag.String("path", "./", "root path")
	replacer = flag.String("i", "{}", "filename replacer")
	command  = flag.String("c", "./{}", "spec how to run command")
)

func work(files []string) (err error) {
	tmpl, err := template.ParseFiles("rep.tmpl")
	if err != nil {
		return
	}
	startTime := time.Now()
	data := map[string]interface{}{
		"StartTime": startTime.Format("2006-01-02 15:04:05"),
	}
	for i, file := range files {
		c := strings.Replace(*command, *replacer, file, -1)
		fmt.Printf("%d\t%s\t%s\n", i, file, c)
	}
	endTime := time.Now()
	data["TimeCost"] = endTime.Sub(startTime)

	err = renderTemplate("test.html", tmpl, data)
	if err != nil {
		return
	}
	return nil
}

func init() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	var err error
	files, err := pathWalk(*path, *depth)
	if err != nil {
		log.Fatal(err)
	}
	err = work(files)
	if err != nil {
		log.Fatal(err)
	}
}
