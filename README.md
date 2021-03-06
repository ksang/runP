# runP
[![Build Status](https://travis-ci.org/ksang/runP.svg?branch=master)](https://travis-ci.org/ksang/runP) [![Go Report Card](https://goreportcard.com/badge/github.com/ksang/runP)](https://goreportcard.com/report/github.com/ksang/runP)

Run any number of processes simultaneously for performance testing

### Features

- Precise wall clock time gathering
- Independent environment variable injection

### Usage

    Usage of ./runp:
      -c string
        	full command with arguments, e.g "ifconfig -a"
      -e string
        	env to pass to sub-processes, semi-column to divide env entries,
            pipe to divide different env of sub-processes,
            e.g: "PATH=/usr/local;OS=Linux|OS=Darwin"
      -n int
        	the number of processes to run (default 2)
      -s	suppress outputs from process

### Example

    $./build/runp -c "sleep 5" -n 8 -s
    Process 3:      Elapsed Time: 5.001953797s
    Process 7:      Elapsed Time: 5.002043916s
    Process 1:      Elapsed Time: 5.002514304s
    Process 5:      Elapsed Time: 5.00284948s
    Process 4:      Elapsed Time: 5.003651915s
    Process 2:      Elapsed Time: 5.004092578s
    Process 0:      Elapsed Time: 5.004677919s
    Process 6:      Elapsed Time: 5.00522753s
