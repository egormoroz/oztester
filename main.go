package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"time"
)

type diffResult string

var (
	executable  = flag.String("x", "./bin/sol1", "path to executable")
	testsDir    = flag.String("i", "./tests", "path to tests directory")
	concurrency = flag.Int("n", 1, "concurrency")
	verbose     = flag.Bool("v", false, "verbose")
	normalizeWS = flag.Bool("w", false, "normalize whitespace")
	timeLimit   = flag.Duration("t", time.Second, "time limit")
	asyncPrint  = flag.Bool("a", false, "print testcase result immediatelly in async mode")
)

func main() {
	flag.Parse()

	cases, err := findCases()
	if err != nil {
		log.Fatal(err)
	}
	slices.SortFunc(cases, caseCmp)

	fmt.Printf("%d cases found: %v\n\n", len(cases), cases)

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	res := launch(ctx, cases)
	slices.SortFunc(res, func(u, v caseReport) int {
		return caseCmp(u.name, v.name)
	})

	m := make(map[result]int)
	for _, r := range res {
		if !*asyncPrint {
			fmt.Println(r)
		}

		m[r.res]++
	}

	fmt.Println()
	fmt.Printf("[%d/%d] OK %d WA %d TL %d CC %d ERR %d\n",
		m[resultOK], len(cases), m[resultOK], m[resultWA],
		m[resultTL], m[resultCC], m[resultErr],
	)
}
