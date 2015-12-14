package main

import (
	"testing"
)

func TestSingleIP(t *testing.T) {
	expectation := "-join=127.0.0.1"

	str := joinIPArg("127.0.0.1")

	if str != expectation {
		t.Errorf("%s doesn't match %s", expectation, str)
	}
}

func TestMultipleIPs(t *testing.T) {
	expectation := []string{"-join=10.0.0.1", "-join=10.0.0.2"}

	ips := []string{"10.0.0.1", "10.0.0.2"}
	strs := joinIPArgs(ips)

	matches := true
	for index, arg := range strs {
		if expectation[index] != arg {
			matches = false
			break
		}
	}

	if !matches {
		t.Errorf("%v doesn't match %v", strs, expectation)
	}
}

type fakeResolver struct {
	servers []string
}

func (f fakeResolver) LookupHost(host string) ([]string, error) {
	return f.servers, nil
}

func TestDNSResolution(t *testing.T) {
	expectation := []string{"10.0.0.1", "10.0.0.2"}

	res := fakeResolver{servers: []string{"10.0.0.1", "10.0.0.2"}}
	ips := resolveARecords(res, "example.com")

	matches := true
	for index, ip := range ips {
		if expectation[index] != ip {
			matches = false
			break
		}
	}

	if !matches {
		t.Errorf("%v doesn't match %v", ips, expectation)
	}
}

func TestBootstrapExpectForOne(t *testing.T) {
	expectation := int64(1)

	result := bootstrapExpectFromLookup("1", true)

	if result != expectation {
		t.Errorf("%v doesn't match %d", result, expectation)
	}
}

func TestBootstrapExpectForLetter(t *testing.T) {
	expectation := int64(3)

	result := bootstrapExpectFromLookup("A", true)

	if result != expectation {
		t.Errorf("%v doesn't match %d", result, expectation)
	}
}

func TestFullArgsForServer(t *testing.T) {
	expectation := []string{"consul", "agent", "-join=127.0.0.1", "-bootstrap-expect=3", "-server"}
	joinArgs := []string{"-join=127.0.0.1"}
	bootstrap := int64(3)
	osArgs := []string{"-server"}

	result := fullCommandArgs(joinArgs, "", "", bootstrap, osArgs)

	matches := true
	for index, arg := range result {
		if expectation[index] != arg {
			matches = false
			break
		}
	}

	if !matches {
		t.Errorf("%v doesn't match %v", result, expectation)
	}
}

func TestFullArgsForServerIncludingNodeAndAdvertise(t *testing.T) {
	expectation := []string{"consul", "agent", "-node=foo", "-advertise=10.1.1.1", "-join=127.0.0.1", "-bootstrap-expect=3", "-server"}
	joinArgs := []string{"-join=127.0.0.1"}
	node := "foo"
	advertise := "10.1.1.1"
	bootstrap := int64(3)
	osArgs := []string{"-server"}

	result := fullCommandArgs(joinArgs, node, advertise, bootstrap, osArgs)

	matches := true
	for index, arg := range result {
		if expectation[index] != arg {
			matches = false
			break
		}
	}

	if !matches {
		t.Errorf("%v doesn't match %v", result, expectation)
	}
}

func TestFullArgs(t *testing.T) {
	expectation := []string{"consul", "agent", "-join=127.0.0.1"}
	joinArgs := []string{"-join=127.0.0.1"}
	bootstrap := int64(3)
	osArgs := []string{}

	result := fullCommandArgs(joinArgs, "", "", bootstrap, osArgs)

	matches := true
	for index, arg := range result {
		if expectation[index] != arg {
			matches = false
			break
		}
	}

	if !matches {
		t.Errorf("%v doesn't match %v", result, expectation)
	}
}

func TestFullArgsIncludingNodeAndAdvertise(t *testing.T) {
	expectation := []string{"consul", "agent", "-node=foo", "-advertise=10.1.1.1", "-join=127.0.0.1"}
	joinArgs := []string{"-join=127.0.0.1"}
	node := "foo"
	advertise := "10.1.1.1"
	bootstrap := int64(3)
	osArgs := []string{}

	result := fullCommandArgs(joinArgs, node, advertise, bootstrap, osArgs)

	matches := true
	for index, arg := range result {
		if expectation[index] != arg {
			matches = false
			break
		}
	}

	if !matches {
		t.Errorf("%v doesn't match %v", result, expectation)
	}
}
