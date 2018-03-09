package client

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
	c := &AtcoderClient{}
	c.Login()
}

func Test_func(t *testing.T, fn func()) {
	stdout, stderr := stubIO("", fn)
	t.Log(stdout)
	t.Log(stderr)
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
