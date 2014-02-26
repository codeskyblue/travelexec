package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/shxsun/flags"
	"github.com/shxsun/goyaml"
)

func renderTemplate(filename string, tmplFile string, data interface{}) (err error) {
	tmpl, err := template.ParseFiles(tmplFile)
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
	out, err := goyaml.Marshal(data)
	if err != nil {
		return
	}
	return ioutil.WriteFile(filename, out, 0644)
}

func loadFile(filename string, data interface{}) (err error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	return goyaml.Unmarshal(raw, data)
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
			if matched, _ := regexp.Match(mycnf.Include, []byte(info.Name())); matched {
				files = append(files, path)
			}
		}
		return nil
	})
	return
}

type GlobalConfig struct {
	Replacer string `short:"i" description:"replacer"`
	Command  string `short:"c" description:"specify how to process each file"`
	Include  string `short:"I" long:"include-regex" description:"regex set which file can be included"`
	Path     string `short:"p" long:"path" description:"path for search"`
	Depth    int    `short:"d" long:"depth" description:"depth to travel directory tree"`
	Verbose  bool   `short:"v" long:"verbose" description:"show verbose output"`
	Result   string `yaml:"report-html" short:"H" long:"out-html" description:"output result as html"`
	Timeout  string `short:"t" long:"timeout" description:"timeout for each exec"`
	Reload   bool   `short:"r" long:"reload" description:"reload all failed cmd, run again:`
	Exclude  string `yaml:"-"`

	Version  bool `yaml:"-"`
	InitYaml bool `yaml:"-" long:"init" description:"create a sample .trival.yml and exit"`
}

var mycnf = &GlobalConfig{
	Depth:    0,
	Path:     "./",
	Replacer: "{}",
	Version:  false,
	Result:   "test.html",
	Timeout:  "30m",
	Command:  "./{}",
	Include:  ".*",
	Verbose:  false,
}

func init() {
	if fileExists(".travel.yml") {
		loadFile(".travel.yml", mycnf)
	}
	args, err := flags.Parse(mycnf)
	_, _ = args, err
	if err != nil {
		/*if err == flags.ErrorType(flags.ErrHelp) {
			os.Exit(0)
		} */
		os.Exit(1)
	}
	if mycnf.InitYaml {
		dumpFile(".travel.yml", mycnf)
		os.Exit(0)
	}
	if mycnf.Version {
		fmt.Println(VERSION)
		os.Exit(0)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type TaskConfig struct {
	Replacer string
	Command  string
	Timeout  time.Duration
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

func Go(f func() error) chan error {
	ch := make(chan error)
	go func() {
		err := f()
		ch <- err
	}()
	return ch
}

func groupKill(cmd *exec.Cmd) (err error) {
	var pid, pgid int
	if cmd.Process != nil {
		pid = cmd.Process.Pid
		c := exec.Command("/bin/ps", "-o", "pgid", "-p", strconv.Itoa(pid), "--no-header")
		var out []byte
		out, err = c.Output()
		if err != nil {
			return
		}
		_, err = fmt.Sscanf(string(out), "%d", &pgid)
		if err != nil {
			return
		}
		err = exec.Command("/bin/kill", "-TERM", "-"+strconv.Itoa(pgid)).Run()
	}
	return
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
		format := fmt.Sprintf(prefix+"%%-14v %-24s\n", c) //file) //file)
		// show current exec file
		fmt.Printf(prefix+"exec %s ...", strconv.Quote(file))
		if mycnf.Verbose {
			fmt.Printf("\n")
		}

		output := bytes.NewBuffer(nil)
		cmd := exec.Command("/bin/bash", "-c", c)
		if mycnf.Verbose {
			cmd.Stdout = io.MultiWriter(os.Stdout, output)
			cmd.Stderr = io.MultiWriter(os.Stderr, output)
		} else {
			cmd.Stdout = output
			cmd.Stderr = output
		}
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.SysProcAttr.Setpgid = true

		// handle timeout
		if err = cmd.Start(); err == nil {
			select {
			case <-time.After(cfg.Timeout):
				if cmd.Process != nil {
					//cmd.Process.Kill()
					groupKill(cmd)
					err = errors.New("signal: terminated")
				}
			case err = <-Go(cmd.Wait):
			}
		}

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

func fileExists(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.Mode().IsRegular()
}

func selfPath() string {
	return filepath.Dir(os.Args[0])
}

const STATE_FILE = ".out.json"

func main() {
	var err error
	var taskcfg = &TaskConfig{}
	var startTime = time.Now()
	if mycnf.Reload && fileExists(STATE_FILE) {
		err = loadFile(STATE_FILE, taskcfg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		files, err := pathWalk(mycnf.Path, mycnf.Depth)
		if err != nil {
			log.Fatal(err)
		}
		taskcfg.Command = mycnf.Command
		taskcfg.Replacer = mycnf.Replacer
		timeout, err := time.ParseDuration(mycnf.Timeout)
		if err != nil {
			log.Fatal(err)
		}
		taskcfg.Timeout = timeout
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

	// reder html
	data := map[string]interface{}{}
	data["StartTime"] = startTime.Format("2006-01-02 15:04:05")
	endTime := time.Now()
	data["TimeCost"] = endTime.Sub(startTime)
	data["Tasks"] = results
	hostname, _ := os.Hostname()
	data["Host"] = hostname
	data["Total"] = len(taskcfg.Files)
	data["FailCount"] = errCnt

	tmplPath := filepath.Join(selfPath(), ".rep.tmpl")
	if !fileExists(tmplPath) {
		log.Println("create .rep.tmpl")
		ioutil.WriteFile(tmplPath, []byte(defaultTemplate), 0644)
	}
	err = renderTemplate(mycnf.Result, tmplPath, data)
	if err != nil {
		log.Fatal(err)
	}

	// save to restart again
	taskcfg.Files = errfiles
	err = dumpFile(STATE_FILE, taskcfg)
	if err != nil {
		log.Fatal(err)
	}
}
