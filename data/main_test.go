package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func Test_solve(t *testing.T) {
	sampleInput := []string{"sampleInput-placeholder"}
	sampleOutput := []string{"sampleOutput-placeholder"}
	for i, input := range sampleInput {
		stdout, _ := stubIO(input, solve)
		if got, want := stdout, sampleOutput[i]; got != want {
			t.Fatalf("wrong answer: got %s, want %s", got, want)
		}
	}

}

func stubIO(inbuf string, fn func()) (string, string) {
	inr, inw, _ := os.Pipe()
	outr, outw, _ := os.Pipe()
	errr, errw, _ := os.Pipe()

	orgStdin := os.Stdin
	orgStdout := os.Stdout
	orgStderr := os.Stderr

	inw.Write([]byte(inbuf))
	inw.Close()
	os.Stdin = inr
	os.Stdout = outw
	os.Stderr = errw
	fn()
	os.Stdin = orgStdin
	os.Stdout = orgStdout
	os.Stderr = orgStderr
	outw.Close()
	outbuf, _ := ioutil.ReadAll(outr)
	errw.Close()
	errbuf, _ := ioutil.ReadAll(errr)

	return string(outbuf), string(errbuf)
}
