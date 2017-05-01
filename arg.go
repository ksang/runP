package runP

// Arg is the arguments parsed from command line
type Arg struct {
	// Command line with arguments for target process
	Command string
	// Process number to run
	ProcNum int
	// Suppress stdout/stderr from processes
	Suppress bool
	// environment varibales
	Env [][]string
}
