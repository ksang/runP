package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/ksang/runP"
)

var (
	// Command line with arguments for target process
	command string
	// Process number to run
	procNum int
	// Suppress stdout/stderr from processes
	suppress bool
)

func init() {
	flag.StringVar(&command, "c", "", "full command with arguments, e.g \"ifconfig -a\"")
	flag.IntVar(&procNum, "n", 2, "the number of processes to run")
	flag.BoolVar(&suppress, "s", false, "suppress outputs from process")
}

func main() {
	flag.Parse()
	if len(command) == 0 {
		fmt.Println("ERROR: must provide a command")
		flag.PrintDefaults()
		return
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)

	arg := runP.Arg{
		Command:  command,
		ProcNum:  procNum,
		Suppress: suppress,
	}
	manager := runP.New(arg)

	go func() {
		<-sc
		fmt.Println("Interrupted")
		manager.Quit()
	}()

	if err := manager.Start(); err != nil {
		log.Fatal(err)
	}
	manager.PrintResult()
}
