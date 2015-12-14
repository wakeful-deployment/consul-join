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

type Resolver interface {
	LookupHost(host string) ([]string, error)
}

type resolver struct{}

func (r resolver) LookupHost(host string) ([]string, error) {
	return net.LookupHost(host)
}

func resolveARecords(res Resolver, domain string) []string {
	ips, err := res.LookupHost(domain)

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

func fullCommandArgs(joinArgs []string, node string, advertise string, bootstrapExpect int64, osArgs []string) []string {
	commandArgs := []string{consul, "agent"}

	if node != "" {
		commandArgs = append(commandArgs, fmt.Sprintf("-node=%s", node))
	}

	if advertise != "" {
		commandArgs = append(commandArgs, fmt.Sprintf("-advertise=%s", advertise))
	}

	for _, joinArg := range joinArgs {
		commandArgs = append(commandArgs, joinArg)
	}

	for index, arg := range osArgs {
		if index != 0 {
			commandArgs = append(commandArgs, arg)
		}

		if arg == "-server" {
			commandArgs = append(commandArgs, fmt.Sprintf("-bootstrap-expect=%d", bootstrapExpect))
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
	advertise, _ := os.LookupEnv("ADVERTISE")
	node, _ := os.LookupEnv("NODE")
	_, debug := os.LookupEnv("DEBUG")
	res := resolver{}

	var args []string

	if joinIPExists && joinDNSExists {
		fmt.Fprint(os.Stderr, "cannot have both JOINIP and JOINDNS together")
		syscall.Exit(1)
	}

	if joinIPExists {
		args = []string{joinIPArg(joinIP)}
	} else if joinDNSExists {
		args = joinIPArgs(resolveARecords(res, joinDNS))
	} else {
		args = []string{}
	}

	runArgs := fullCommandArgs(args, node, advertise, bootstrapExpect, os.Args)
	runCommand(runArgs, debug)
}
