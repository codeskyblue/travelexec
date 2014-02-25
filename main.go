package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
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

func dumpFile(filename string, data interface{}) (err error) {
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	err = json.NewEncoder(fd).Encode(data)
	return
}

func loadFile(filename string, data interface{}) (err error) {
	fd, err := os.Open(filename)
	if err != nil {
		return
	}
	defer fd.Close()
	err = json.NewDecoder(fd).Decode(data)
	return
}

func renderJson(filename string, data interface{}) (err error) {
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	err = json.NewEncoder(fd).Encode(data)
	return
}

//func saveText(filename string,

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
			if matched, _ := regexp.Match(*include, []byte(info.Name())); matched {
				files = append(files, path)
			}
		}
		return nil
	})
	return
}

var (
	depth      = flag.Int("depth", 0, "depth of dir search")
	path       = flag.String("path", "./", "root path")
	replacer   = flag.String("r", "{}", "filename replacer")
	command    = flag.String("c", "./{}", "spec how to run command")
	resultHtml = flag.String("html", "test.html", "output html")
	verbose    = flag.Bool("v", false, "show verbose info")
	resultJson = flag.String("json", ".out.json", "output json")

	include = flag.String("I", ".*", "regex to match include file")
	//exclude = flag.String("x", "\\.*", "regex to exclude file")
	reload = flag.Bool("reload", false, "reload failed file, and run again")
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
	Filename  string
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

		prefix := fmt.Sprintf("\r\033[36m>>>\033[0m %-5d", i)
		format := fmt.Sprintf(prefix+"%%-30v    %-24s\n", c) //file) //file)
		// show current exec file
		fmt.Printf(prefix+"exec %s ...", strconv.Quote(file))

		output := bytes.NewBuffer(nil)
		cmd := exec.Command("/bin/bash", "-c", c)
		if *verbose {
			cmd.Stdout = io.MultiWriter(os.Stdout, output)
			cmd.Stderr = io.MultiWriter(os.Stderr, output)
		} else {
			cmd.Stdout = output
			cmd.Stderr = output
		}
		err = cmd.Run()
		if err != nil {
			r.Error = err
			fmt.Printf(format, fmt.Sprintf("\033[33m"+"err: %s"+"\033[0m", err))
		} else {
			fmt.Printf(format, "\033[32m"+"success"+"\033[0m")
		}
		r.Command = c
		r.Filename = file
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

func fileExists(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.Mode().IsRegular()
}

func main() {
	var err error
	var taskcfg = &TaskConfig{}
	if *reload && fileExists(*resultJson) {
		err = loadFile(*resultJson, taskcfg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		files, err := pathWalk(*path, *depth)
		if err != nil {
			log.Fatal(err)
		}
		taskcfg.Command = *command
		taskcfg.Replacer = *replacer
		taskcfg.Files = files
	}

	results := work(taskcfg)
	errfiles := []string{}
	errCnt := 0
	for _, r := range results {
		if r.Error != nil {
			errCnt++
			errfiles = append(errfiles, r.Filename)
		}
	}
	// save to restart again
	taskcfg.Files = errfiles
	err = dumpFile(*resultJson, taskcfg)
	if err != nil {
		log.Fatal(err)
	}

	// reder html
	startTime := time.Now()
	data := map[string]interface{}{}
	data["StartTime"] = startTime.Format("2006-01-02 15:04:05")
	endTime := time.Now()
	data["TimeCost"] = endTime.Sub(startTime)
	data["Tasks"] = results
	hostname, _ := os.Hostname()
	data["Host"] = hostname
	data["Total"] = len(taskcfg.Files)
	data["FailCount"] = errCnt

	err = renderTemplate(*resultHtml, "rep.tmpl", data)
	if err != nil {
		log.Fatal(err)
	}
}
