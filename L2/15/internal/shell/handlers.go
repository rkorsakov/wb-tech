package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func (shell *Shell) handleCdCommand(args []string) {
	var path string
	if len(args) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "cd: %v\n", err)
			return
		}
		path = home
	} else {
		path = args[0]
	}

	err := os.Chdir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cd: %v\n", err)
	}
}

func (shell *Shell) handlePwdCommand() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pwd: %v\n", err)
		return
	}
	fmt.Println(wd)
}

func (shell *Shell) handleEchoCommand(args []string) {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg)
	}
	fmt.Println()
}

func (shell *Shell) handlePsCommand() {
	cmd := exec.Command("ps", "aux")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ps: %v\n", err)
	}
}

func (shell *Shell) handleKillCommand(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "kill: usage: kill <pid>\n")
		return
	}

	pidStr := args[0]
	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid <= 0 {
		fmt.Fprintf(os.Stderr, "kill: invalid PID: %s\n", pidStr)
		return
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "kill: cannot find process %d: %v\n", pid, err)
		return
	}

	err = process.Kill()
	if err != nil {
		fmt.Fprintf(os.Stderr, "kill: failed to kill process %d: %v\n", pid, err)
		return
	}

	fmt.Printf("Killed process %d\n", pid)
}
