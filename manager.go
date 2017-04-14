package runP

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Manager is runP controller for all the runP tasks
type Manager struct {
	arg     Arg
	readyCh chan struct{}
	kickCh  chan struct{}
	quitCh  chan struct{}
	wg      sync.WaitGroup

	rMu     sync.RWMutex
	results []*Result
}

// Result is the metadata for each sub process
type Result struct {
	ProcessID   int
	Stdout      string
	Stderr      string
	Error       error
	ElapsedTime time.Duration
}

// New a runP manager with arguments provided
func New(arg Arg) *Manager {
	return &Manager{
		arg: arg,
	}
}

// Start the manager and run sub processses
func (m *Manager) Start() error {
	rc := make(chan struct{}, m.arg.ProcNum)
	kc := make(chan struct{}, m.arg.ProcNum)
	qc := make(chan struct{}, m.arg.ProcNum)
	m.readyCh = rc
	m.kickCh = kc
	m.quitCh = qc
	m.wg.Add(m.arg.ProcNum)
	for i := 0; i < m.arg.ProcNum; i++ {
		n := i
		go m.RunProc(n)
	}
	count := 0
	for {
		select {
		case <-rc:
			count++
		}
		if count == m.arg.ProcNum {
			break
		}
	}
	for i := 0; i < m.arg.ProcNum; i++ {
		kc <- struct{}{}
	}
	m.wg.Wait()
	return nil
}

// RunProc run one process according to Manager config
func (m *Manager) RunProc(n int) {
	defer m.wg.Done()
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
		start  time.Time
	)

	res := &Result{
		ProcessID: n,
	}

	// append result first, in case sub process hang, still can get it's output
	// by intrrupting the runP main process
	m.rMu.Lock()
	m.results = append(m.results, res)
	m.rMu.Unlock()

	args := strings.Split(m.arg.Command, " ")
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// process is ready to run
	m.readyCh <- struct{}{}

	<-m.kickCh
	start = time.Now()

	go func() {
		// received quit signal
		<-m.quitCh
		// if the process is already exited
		if cmd == nil {
			return
		}
		//res.ElapsedTime = time.Since(start)
		res.Stdout = stdout.String()
		res.Stderr = stderr.String()
		// kill the sub process
		if err := cmd.Process.Kill(); err != nil {
			log.Println(err)
		}
	}()

	err := cmd.Run()
	res.ElapsedTime = time.Since(start)
	if err == nil {
		res.Stdout = stdout.String()
		res.Stderr = stderr.String()
	}
	res.Error = err
}

// PrintResult prints process outputs and timing information
func (m *Manager) PrintResult() {
	if !m.arg.Suppress {
		for i, r := range m.results {
			m.printSepLine()
			fmt.Printf("Process %d:\n", r.ProcessID)
			if len(r.Stdout) > 0 {
				fmt.Printf("Stdout:\n%s", r.Stdout)
			}
			if len(r.Stderr) > 0 {
				fmt.Printf("Stderr:\n%s", r.Stderr)
			}
			if r.Error != nil {
				fmt.Printf("Error:\n%s\n", r.Error.Error())
			}
			if i == len(m.results)-1 {
				m.printSepLine()
			}
		}
	}
	for _, r := range m.results {
		fmt.Printf("Process %d:\tElapsed Time: %s\n", r.ProcessID, r.ElapsedTime.String())
	}
}

func (m *Manager) printSepLine() {
	fmt.Println("==========================================================")
}

// Quit signal subprocess goroutinue to quit
func (m *Manager) Quit() {
	for i := 0; i < m.arg.ProcNum; i++ {
		m.quitCh <- struct{}{}
	}
}
