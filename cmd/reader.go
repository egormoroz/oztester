package main

import (
	"cmp"
	"os"
	"strconv"
	"strings"
)

func caseCmp(x, y string) int {
	ix, xErr := strconv.Atoi(x)
	iy, yErr := strconv.Atoi(y)
	if xErr == nil && yErr == nil {
		return cmp.Compare(ix, iy)
	}
	return cmp.Compare(x, y)
}

func findCases() ([]string, error) {
	dirEntries, err := os.ReadDir(*testsDir)
	if err != nil {
		return nil, err
	}

	type pair struct {
		input, output bool
	}
	pairs := make(map[string]pair, len(dirEntries))
	cases := make([]string, 0, len(dirEntries)/2)
	for _, e := range dirEntries {
		cut, ok := strings.CutSuffix(e.Name(), ".a")
		p := pairs[cut]
		if ok {
			p.output = true
		} else {
			p.input = true
		}
		pairs[cut] = p

		if p.input && p.output {
			cases = append(cases, cut)
		}
	}

	return cases, nil
}
