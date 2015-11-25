package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

const consul = "consul"

func joinIPArg(ip string) string {
	return fmt.Sprintf("-join=%s", ip)
}

func joinIPArgs(ips []string) []string {
	acc := []string{}

	for _, ip := range ips {
		acc = append(acc, joinIPArg(ip))
	}

	return acc
}

func resolveARecords(domain string) []string {
	ips, err := net.LookupHost(domain)

	if err == nil {
		return ips
	} else {
		errorStr := fmt.Sprintf("error resolving %s\n%s", domain, err.Error())
		fmt.Fprint(os.Stderr, errorStr)
		return []string{}
	}
}

func bootstrapExpectFromLookup(s string, exists bool) int64 {
	if exists {
		parsed, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			return parsed
		}
	}

	return 3
}

func fullCommandArgs(joinArgs []string, bootstrapExpect int64) []string {
	commandArgs := []string{consul, "agent"}

	var isServer = false
	for _, arg := range os.Args {
		if arg == "-server" {
			isServer = true
		}
	}

	if isServer {
		commandArgs = append(commandArgs, fmt.Sprintf("-bootstrap-expect=%d", bootstrapExpect))
	}

	for _, joinArg := range joinArgs {
		commandArgs = append(commandArgs, joinArg)
	}

	for index, arg := range os.Args {
		if index != 0 {
			commandArgs = append(commandArgs, arg)
		}
	}

	return commandArgs
}

func runCommand(argv []string, debug bool) {
	fullCommandPath, pathErr := exec.LookPath(consul)

	if pathErr != nil {
		fmt.Fprint(os.Stderr, pathErr.Error())
		syscall.Exit(1)
	}

	if debug {
		fmt.Println(fullCommandPath, argv)
		syscall.Exit(0)
	} else {
		envv := os.Environ()
		execErr := syscall.Exec(fullCommandPath, argv, envv)

		if execErr != nil {
			fmt.Fprint(os.Stderr, execErr.Error())
			syscall.Exit(1)
		}
	}
}

func main() {
	joinIP, joinIPExists := os.LookupEnv("JOINIP")
	joinDNS, joinDNSExists := os.LookupEnv("JOINDNS")
	bootstrapExpect := bootstrapExpectFromLookup(os.LookupEnv("BOOTSTRAP_EXPECT"))
	_, debug := os.LookupEnv("DEBUG")

	var args []string

	if joinIPExists && joinDNSExists {
		fmt.Fprint(os.Stderr, "cannot have both JOINIP and JOINDNS together")
		syscall.Exit(1)
	}

	if joinIPExists {
		args = []string{joinIPArg(joinIP)}
	} else if joinDNSExists {
		args = joinIPArgs(resolveARecords(joinDNS))
	} else {
		args = []string{}
	}

	runArgs := fullCommandArgs(args, bootstrapExpect)
	runCommand(runArgs, debug)
}
