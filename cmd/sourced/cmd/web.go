package cmd

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/src-d/sourced-ce/cmd/sourced/compose"
)

// The service name used in docker-compose.yml for the srcd/sourced-ui image
const containerName = "sourced-ui"

type webCmd struct {
	Command `name:"web" short-description:"Open the web interface in your browser" long-description:"Open the web interface in your browser"`
}

func (c *webCmd) Execute(args []string) error {
	return OpenUI(2 * time.Second)
}

func init() {
	rootCmd.AddCommand(&webCmd{})
}

func openUI() error {
	var stdout bytes.Buffer
	// wait for the container to start, it can take a while in some cases
	for {
		if err := compose.RunWithIO(context.Background(),
			os.Stdin, &stdout, nil, "port", containerName, "8088"); err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	address := strings.TrimSpace(stdout.String())
	if address == "" {
		return fmt.Errorf("no address found")
	}

	// docker-compose returns 0.0.0.0 which is correct for the bind address
	// but incorrect as connect address
	url := fmt.Sprintf("http://%s", strings.Replace(address, "0.0.0.0", "127.0.0.1", 1))

	for {
		client := http.Client{Timeout: time.Second}
		if _, err := client.Get(url); err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	if err := browser.OpenURL(url); err != nil {
		errors.Wrap(err, "cannot open browser")
	}

	return nil
}

// OpenUI opens the browser with the UI.
func OpenUI(timeout time.Duration) error {
	ch := make(chan error)
	go func() {
		ch <- openUI()
	}()

	if timeout > 5*time.Second {
		stopSpinner := startSpinner("Initializing source{d}...")
		defer stopSpinner()
	}

	select {
	case err := <-ch:
		return errors.Wrap(err, "an error occurred while opening the UI")
	case <-time.After(timeout):
		return fmt.Errorf("error opening the UI, the container is not running after %v", timeout)
	}
}

type spinner struct {
	msg      string
	charset  []int
	interval time.Duration

	stop chan bool
}

func startSpinner(msg string) func() {
	charset := []int{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
	if runtime.GOOS == "windows" {
		charset = []int{'|', '/', '-', '\\'}
	}

	s := &spinner{
		msg:      msg,
		charset:  charset,
		interval: 200 * time.Millisecond,
		stop:     make(chan bool),
	}
	s.Start()

	return s.Stop
}

func (s *spinner) Start() {
	go s.printLoop()
}

func (s *spinner) Stop() {
	s.stop <- true
}

func (s *spinner) printLoop() {
	i := 0
	for {
		select {
		case <-s.stop:
			fmt.Println(s.msg)
			return
		default:
			char := string(s.charset[i%len(s.charset)])
			if runtime.GOOS == "windows" {
				fmt.Printf("\r%s %s", s.msg, char)
			} else {
				fmt.Printf("%s %s\033[A", s.msg, char)
			}

			time.Sleep(s.interval)
		}

		i++
		if len(s.charset) == i {
			i = 0
		}
	}
}
