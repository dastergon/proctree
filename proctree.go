// build +linux +darwin

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/xlab/treeprint"
)

const version = "1.0"

// Process contains system process fields and information.
type Process struct {
	Pid     int
	ppid    int
	pgid    int
	user    string
	command string
}

var (
	psCmd        = "ps"
	psCmdArgs    string
	startPID     = 1
	startProcess Process
	ppidtable    map[int][]Process
)

var (
	flagTreeDepth = flag.Int("l", -1, "Print tree to n level deep")
	flagOnlyUser  = flag.String("u", "", "Show only branches containing process of <user>")
	flagNotRoot   = flag.Bool("U", false, "Do not show branches containing only root processes")
	flagPgid      = flag.Bool("g", false, "Show process group ids")
	flagOnlyPid   = flag.Int("p", -1, "Show only branches containing process <pid>")
	flagVersion   = flag.Bool("version", false, "Outputs the version of proctree.")
)

func init() {
	switch runtime.GOOS {
	case "linux":
		psCmdArgs = "-eo uid,pid,ppid,pgid,args"
	case "darwin":
		psCmdArgs = "-axwwo user,pid,ppid,pgid,command"
	}
}

func main() {
	// Modfify the default usage output.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options...] [<pid>] (defaults to PID 1)\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *flagVersion {
		fmt.Println("proctree:", version)
		os.Exit(0)
	}

	flags := flag.Args()
	// Change the starting PID if another PID is given (defaults to PID 1) .
	if len(flags) != 0 {
		var err error
		// Process ID to start from.
		startPID, err = strconv.Atoi(flags[0])
		if err != nil {
			fmt.Println("proctree:", "Please provide a unique Process ID (integer)")
			os.Exit(1)
		}
	}
	// Execute the ps command with the proper arguments and get the stdout and stderr.
	out, err := exec.Command(psCmd, psCmdArgs).CombinedOutput()
	if err != nil {
		fmt.Println("proctree:", err)
		os.Exit(1)
	}

	// Split the output in multiple lines.
	parts := strings.Split(string(out), "\n")
	ppidtable = map[int][]Process{}
	for lineIdx := range parts {

		// Skip column titles.
		if lineIdx == 0 {
			continue
		}

		// Parse information for a process.
		psInfoLine := strings.Fields(parts[lineIdx])
		if len(psInfoLine) == 0 {
			continue
		}

		pid, _ := strconv.Atoi(psInfoLine[1])
		ppid, _ := strconv.Atoi(psInfoLine[2])
		pgid, _ := strconv.Atoi(psInfoLine[3])
		processInfo := Process{
			user:    psInfoLine[0], // username on darwin & uid on linux
			Pid:     pid,
			ppid:    ppid,
			pgid:    pgid,
			command: psInfoLine[4],
		}
		if startPID == processInfo.Pid {
			startProcess = processInfo
		}
		ppidtable[ppid] = append(ppidtable[ppid], processInfo)
	}

	if startProcess == (Process{}) {
		fmt.Println("proctree:", "The given process ID does not exist.")
		os.Exit(0)
	}

	// Build the process tree.
	tree := treeprint.New()
	output := constructOutput(startProcess)
	treeprint.EdgeTypeStart = treeprint.EdgeType(output) // Set the name of the root of the tree.
	buildProcessTree(startProcess, true, tree, 0)
	fmt.Println(tree.String())
}

// buildProcessTree iterates the available processes in a depth-first fashion and
// constructs the final process tree.
func buildProcessTree(parentProcess Process, root bool, tree treeprint.Tree, depth int) {
	if depth == *flagTreeDepth {
		return
	}
	if !root {
		output := constructOutput(parentProcess)
		tree = tree.AddBranch(output)
	}
	ppid := parentProcess.Pid
	for index := range ppidtable[ppid] {
		process := ppidtable[ppid][index]
		if *flagOnlyUser != "" && process.user != *flagOnlyUser {
			continue
		}
		if *flagNotRoot && process.user == "root" {
			continue
		}

		if len(ppidtable[process.ppid]) > 0 {
			buildProcessTree(process, false, tree, depth+1)
		} else {
			output := constructOutput(process)
			tree.AddNode(output)
		}
	}
}

// constructOutput constructs the name of the branch and node
// according to the given flags.
func constructOutput(process Process) string {
	output := strconv.Itoa(process.Pid) + " " + process.user
	if *flagPgid {
		output = output + " " + strconv.Itoa(process.pgid)
	}
	output = output + " (" + process.command + ")"
	return output
}
