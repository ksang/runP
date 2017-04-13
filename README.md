# runP
[![Build Status](https://travis-ci.org/ksang/runP.svg?branch=master)](https://travis-ci.org/ksang/runP) [![Go Report Card](https://goreportcard.com/badge/github.com/ksang/runP)](https://goreportcard.com/report/github.com/ksang/runP)

Run any number of processes simultaneously for performance testing

    runp -h
    -c string
          full command with arguments, e.g "ifconfig -a"
    -n int
          the number of processes to run, default is 2 (default 2)
    -s	suppress the stdout/stderr from process
