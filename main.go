package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func renderTemplate(filename string, tmplFile string, data interface{}) (err error) {
	tmpl, err := template.ParseFiles("rep.tmpl")
	if err != nil {
		return
	}
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
	outfile  = flag.String("o", "test.html", "output file")
)

type TaskConfig struct {
	Replacer string
	Command  string
	Files    []string
}

type TaskResult struct {
	StartTime string
	TimeCost  time.Duration
	Command   string
	Output    string // console out
	Source    string // source code
	Error     error
}

func work(cfg *TaskConfig) (results []TaskResult) {
	var err error
	results = []TaskResult{}
	for i, file := range cfg.Files {
		start := time.Now()
		r := TaskResult{
			StartTime: start.Format("15:04:05"),
		}
		c := strings.Replace(cfg.Command, cfg.Replacer, file, -1)
		fmt.Printf(">>> %d\t%s\t%s\n", i, file, c)
		output := bytes.NewBuffer(nil)
		cmd := exec.Command("/bin/bash", "-c", c)
		cmd.Stdout = io.MultiWriter(os.Stdout, output)
		cmd.Stderr = io.MultiWriter(os.Stderr, output)
		err = cmd.Run()
		if err != nil {
			r.Error = err
		}
		r.Command = c
		r.TimeCost = time.Now().Sub(start)
		r.Output = string(output.Bytes())
		r.Source = "unfinished(todo)"
		results = append(results, r)
	}
	return
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
	taskcfg := &TaskConfig{}
	taskcfg.Command = *command
	taskcfg.Replacer = *replacer
	taskcfg.Files = files
	results := work(taskcfg)
	if err != nil {
		log.Fatal(err)
	}
	errCnt := 0
	for _, r := range results {
		if r.Error != nil {
			errCnt++
		}
	}
	startTime := time.Now()
	data := map[string]interface{}{}
	data["StartTime"] = startTime.Format("2006-01-02 15:04:05")
	endTime := time.Now()
	data["TimeCost"] = endTime.Sub(startTime)
	data["Tasks"] = results
	hostname, _ := os.Hostname()
	data["Host"] = hostname
	data["Total"] = len(files)
	data["FailCount"] = errCnt

	err = renderTemplate(*outfile, "rep.tmpl", data)
	if err != nil {
		log.Fatal(err)
	}
}
