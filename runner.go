package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/fatih/color"
)

type result int

const (
	resultOK  result = iota
	resultWA         // wrong answer
	resultTL         // time limit
	resultCC         // context cancelled
	resultErr        // runtime error
)

type resFmt struct {
	name  string
	color color.Attribute
}

var result2Fmt = map[result]resFmt{
	resultOK:  {"OK", color.FgGreen},
	resultWA:  {"WA", color.FgYellow},
	resultTL:  {"TL", color.FgMagenta},
	resultCC:  {"CC", color.FgCyan},
	resultErr: {"ERR", color.FgRed},
}

func (r result) String() string {
	if f, ok := result2Fmt[r]; ok {
		return color.New(f.color).Sprint(f.name)
	}
	return "UNK"
}

type caseReport struct {
	name string
	res  result
	err  error
}

func (r caseReport) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s\t%v", r.name, r.res)
	if *verbose {
		fmt.Fprintf(&sb, "\t%v", r.err)
	}
	return sb.String()
}

func launch(ctx context.Context, cases []string) []caseReport {
	var wg sync.WaitGroup
	inCh := make(chan string, 1)
	outCh := make(chan caseReport, 1)

	wg.Add(*concurrency)
	defer wg.Wait()
	for range *concurrency {
		go runner(ctx, inCh, outCh, &wg)
	}

	go func() {
		defer close(inCh)
		for _, testCase := range cases {
			select {
			case inCh <- testCase:
			case <-ctx.Done():
				return
			}
		}
	}()

	res := make([]caseReport, 0, len(cases))
	for range len(cases) {
		select {
		case r := <-outCh:
			res = append(res, r)
			if *asyncPrint {
				fmt.Println(r)
			}
		case <-ctx.Done():
			return res
		}
	}
	return res
}

func runner(
	ctx context.Context,
	inCh <-chan string,
	outCh chan<- caseReport,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for testCase := range inCh {
		res, err := runCase(ctx, testCase)
		select {
		case outCh <- caseReport{testCase, res, err}:
		case <-ctx.Done():
			return
		}
	}
}

func runCase(ctx context.Context, testCase string) (result, error) {
	inp, ans, err := readCase(fmt.Sprint(*testsDir, "/", testCase))
	if err != nil {
		return resultErr, err
	}

	out, err := run(ctx, inp)
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return resultTL, err
	case errors.Is(err, context.Canceled):
		return resultCC, err
	case err != nil:
		return resultErr, err
	}

	if *normalizeWS {
		ans = bytes.TrimSpace(bytes.ReplaceAll(ans, []byte("\r\n"), []byte("\n")))
		out = bytes.TrimSpace(bytes.ReplaceAll(out, []byte("\r\n"), []byte("\n")))
	}

	if bytes.Equal(out, ans) {
		return resultOK, nil
	}
	return resultWA, nil
}

func readCase(casePath string) (inp []byte, out []byte, err error) {
	if inp, err = os.ReadFile(casePath); err != nil {
		err = fmt.Errorf("failed to read file %s: %w", casePath, err)
		return
	}
	if out, err = os.ReadFile(casePath + ".a"); err != nil {
		err = fmt.Errorf("failed to read file %s.a: %w", casePath, err)
		return
	}
	return
}

func run(ctx context.Context, inp []byte) (out []byte, err error) {
	ctx, cancel := context.WithTimeout(ctx, *timeLimit)
	defer cancel()

	cmd := exec.CommandContext(ctx, *executable)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	_, _ = stdin.Write(inp)
	_ = stdin.Close()

	if err := cmd.Wait(); err != nil {
		if cerr := ctx.Err(); cerr != nil {
			return nil, cerr
		}
		return nil, fmt.Errorf("command failed: %w", err)
	}

	return stdout.Bytes(), nil
}
