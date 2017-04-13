package runP

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Manager struct {
	arg     Arg
	readyCh chan struct{}
	kickCh  chan struct{}
	wg      sync.WaitGroup

	rMu     sync.RWMutex
	results []Result
}

type Result struct {
	ProcessID   int
	Stdout      string
	Stderr      string
	Error       error
	ElapsedTime time.Duration
}

func New(arg Arg) *Manager {
	return &Manager{
		arg: arg,
	}
}

func (m *Manager) Start() error {
	rc := make(chan struct{}, m.arg.ProcNum)
	kc := make(chan struct{}, m.arg.ProcNum)
	m.readyCh = rc
	m.kickCh = kc
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

func (m *Manager) RunProc(n int) {
	defer m.wg.Done()
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
		start  time.Time
	)

	res := Result{
		ProcessID: n,
	}

	args := strings.Split(m.arg.Command, " ")
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// process is ready to run
	m.readyCh <- struct{}{}

	<-m.kickCh
	start = time.Now()
	err := cmd.Run()
	res.ElapsedTime = time.Since(start)

	if err == nil {
		res.Stdout = stdout.String()
		res.Stderr = stderr.String()
	}
	res.Error = err

	m.rMu.Lock()
	m.results = append(m.results, res)
	m.rMu.Unlock()
}

func (m *Manager) PrintResult() {
	if !m.arg.Suppress {
		for i, r := range m.results {
			m.PrintSepLine()
			fmt.Printf("Process %d:\nStdout:\n%s\nStderr:\n%s\nError:\n%v\n",
				r.ProcessID, r.Stdout, r.Stderr, r.Error)
			if i == len(m.results)-1 {
				m.PrintSepLine()
			}
		}
	}
	for _, r := range m.results {
		fmt.Printf("Process %d:\tElapsed Time: %s\n", r.ProcessID, r.ElapsedTime.String())
	}
}

func (m *Manager) PrintSepLine() {
	fmt.Println("==========================================================")
}
