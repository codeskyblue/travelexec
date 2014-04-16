package main

/*
	2014/02/24 first release(contains all basic functions)
	2014/02/26 add timeout
	2014/02/26 add group kill
	2014/02/26 use 3rd lib(go-flags,goyaml). add .travel.yml for spec setting + bug fix
	2014/02/28 exit_status not eq 0 when some failed
	2014/02/28 send html to platform support
	2014/03/01 catch sigint (Ctrl+C)
	2014/03/03 fix sigint bug
	2014/03/04 fix dirWalk depth bug
	2014/04/15 clean temporary file
*/
const VERSION = "0.1.0415"
