package traceroute

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type TraceRoute struct {
	RemoteAddr string
}

type Result struct {
	Output string
	Hops   []*Hop
}

type Hop struct {
	Id   int
	T1   string
	T2   string
	T3   string
	Dest string
}

func Run(remoteAddr string) (*Result, error) {
	c, err := exec.Command("tracert", remoteAddr).CombinedOutput()
	if err != nil {
		log.Printf("error performing tracert %s", err.Error())
		return nil, err
	}

	res := &Result{
		Output: string(c),
	}

	log.Printf("parsing output %s", string(c))
	hops, err := parseOutput(string(c))
	if err != nil {
		return res, err
	}

	res.Hops = hops

	return res, nil
}

func parseOutput(out string) ([]*Hop, error) {
	lines := strings.Split(out, "\n")

	if !strings.Contains(lines[1], "Tracing route to") {
		return nil, errors.New(fmt.Sprintf("bad tracert, expected Tracing route to line got %s", lines[1]))
	}

	var firstLine int
	for k,v := range lines {
		if strings.Contains(v, "over a maximum of") {
			firstLine = k
			break
		}
	}

	var lastLine int
	for k, v := range lines {
		if strings.Contains(v, "Trace complete.") {
			lastLine = k
			break
		}
	}

	hops := make([]*Hop, 0)
	for _, v := range lines[firstLine+2 : lastLine-1] {
		hop, err := strconv.Atoi(string(v[2]))
		if err != nil {
			log.Printf(v)
			return nil, err
		}
		t1 := v[5:12]
		t2 := v[14:21]
		t3 := v[22:30]
		addr := v[32:]
		hops = append(hops, &Hop{
			Id:   hop,
			T1:   t1,
			T2:   t2,
			T3:   t3,
			Dest: addr,
		})
	}

	return hops, nil
}
