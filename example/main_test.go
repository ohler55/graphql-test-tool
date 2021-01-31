package main_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/ohler55/graphql-test-tool/gtt"
)

var (
	testPort int
	ruby     bool
)

func TestMain(m *testing.M) {
	flag.BoolVar(&ruby, "ruby", ruby, "run the ruby server instead of go server")
	flag.Parse()

	code, err := run(m)
	if err != nil {
		fmt.Println(err.Error())
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	code = 1
	var addr *net.TCPAddr
	if addr, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var ln *net.TCPListener
		if ln, err = net.ListenTCP("tcp", addr); err == nil {
			testPort = ln.Addr().(*net.TCPAddr).Port
			ln.Close()
		}
		if err != nil {
			return
		}
	}
	var cmd *exec.Cmd
	if ruby {
		cmd = exec.Command("ruby", "song.rb", "-p", strconv.Itoa(testPort))
	} else {
		cmd = exec.Command("go", "run", "main.go", "-p", strconv.Itoa(testPort))
	}
	stdout, _ := cmd.StdoutPipe()
	if err = cmd.Start(); err != nil {
		return
	}
	var out []byte
	done := make(chan bool)
	go func() {
		out, _ = ioutil.ReadAll(stdout)
		done <- true
	}()
	defer func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		// Should be redundant but for some reason stdout is not always
		// closed when the process is killed.
		stdout.Close()
		<-done
		if testing.Verbose() {
			fmt.Println(string(out))
		}
	}()

	for i := 0; i < 10; i++ {
		u := fmt.Sprintf("http://localhost:%d", testPort)
		var r *http.Response
		if r, err = http.Get(u); err == nil {
			r.Body.Close()
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
	if err == nil {
		code = m.Run()
	}
	return
}

func gttTest(t *testing.T, filepath string) {
	uc, err := gtt.NewUseCase(filepath)
	if err != nil {
		t.Fatal(err.Error())
	}
	r := gtt.Runner{
		Server:   fmt.Sprintf("http://localhost:%d", testPort),
		Base:     "/graphql",
		Indent:   2,
		UseCases: []*gtt.UseCase{uc},
	}
	if testing.Verbose() {
		r.ShowComments = true
		r.ShowResponses = true
		r.ShowRequests = true
	}

	if err = r.Run(); err != nil {
		t.Fatal(err.Error())
	}
}
