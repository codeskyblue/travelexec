## travelexec

This is tool is for qa-test. It find all test files and execute it, Judge a case is ok depends on if the exitcode eq 0.

(golang binary bool) travel file tree and call exec (like find + xargs, but offers more functions)

Only tested on linux.

## How to use
first download binary file(through gobuild.io) to you system, and add it to your $PATH

for example, there are 2 test file need to run.

 	test_a.py
	test_b.py
  
There are two thing we need to case about.

1. find out which file to run
2. specify how to run this test file

travelexec use regext to match files.

	travelexec -I '^test_.*\.py$' -c 'python {}'
  
travelexec will filter out files which basename match `^test_.*\.py$`, before call shell, filename will replace `{}`

so, the following command will be executed.

	python test_a.py
	python test_b.py

result will be saved into test.html.

### config file support
  travelexec -I '^test_.*\.py$' -c 'python {}' --init
  
after you run this command, a conf file `.travel.yml` will be generated. Next time, you just need to run `travelexec`, that is very helpful to save old settings.

	travelexec --reload
  
this command will only run last failed test cases.

### integerate with jenkins
add such command into **Execute Shell**

	mkdir -p ${WORKSPACE}/travelrep
	travelexec --html ${WORKSPACE}/travelrep/index.html
  
Achieve HTML report. set base dir (travelrep), index file (index.html)

if test failed, travelexec exitcode will be not 0.
### there are still a lot this README not metion about.
use `travelexec -h` for more help.
