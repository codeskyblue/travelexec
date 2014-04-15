## travelexec
This is command tool for program test. He will help you travel the directories and find all files and execute specified command.

So I called it `travelexec`.

This program is write by golang. To got the command tool, need to install golang development before. (**Only tested in linux**)

### How to install
1. install golang(skip it if already exists), see how to install: <http://golang.org/doc/install>
2. run `go get github.com/codeskyblue/travelexec`

### QuickStart
for example, there are 2 test file in the current directory.

 	test_a.py
	test_b.py
  
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
