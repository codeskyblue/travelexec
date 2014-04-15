## travelexec
This is command tool for program test. It will help you **travel** the directories and find all files and **exec**ute specified command.

So I call it `travelexec`.

This program is write by golang(Require go version>=1.2). To get the command tool, you need to install golang development before. (**Only tested in linux**)

Current Version: v0.1.0415

see changelog [HERE](CHANGELOG.md)

FEATURES:

1. HTML report output
2. command timeout
3. group kill when timeout
4. use regex to filter files
5. configuration file support

Already used in two project testing for a month. But the usage of travelexec maybe changed in the future.
### How to install
1. install golang*(skip it if already exists)*, see how to install: <http://golang.org/doc/install>
2. run `go get github.com/codeskyblue/travelexec`

### QuickStart
for example, there are three files in the current directory.

 	test_a.py
	test_b.py
        lib.py
  
run through

        travelexec -I '^test_.*\.py$' -c 'python {}'

use regex `^test_.*\.py$` to find files. `{}` will be replaced as filename. The result is same as

        python test_a.py
        python test_b.py

### config file support
with config file, you don't need to prepare parameters for command.

first, generete a sample config file. default config file is `.travel.yml`

        travelexec --init
  
### how to run last failed files.
        travelexec --reload

### integerate with jenkins
add such command into **Execute Shell**

	mkdir -p ${WORKSPACE}/travelrep
	travelexec --html ${WORKSPACE}/travelrep/index.html
  
Achieve HTML report. set base dir (travelrep), index file (index.html)

if test failed, travelexec exitcode will be not 0.
### there are still a lot this README not metion about.
use `travelexec -h` for more help.

report issues here: <https://github.com/codeskyblue/travelexec/issues> or send mail to me through ssx205@gmail.com

### LICENSE
Apache License 2.0
