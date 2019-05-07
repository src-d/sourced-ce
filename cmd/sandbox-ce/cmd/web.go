package cmd

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/smacker/superset-compose/cmd/sandbox-ce/compose"
)

type webCmd struct {
	Command `name:"web" short-description:"Web"`
}

func (c *webCmd) Execute(args []string) error {
	if err := <-OpenUI(10 * time.Second); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(&webCmd{})
}

func openUI() error {
	var stdout, stderr bytes.Buffer
	for {
		if err := compose.RunWithIO(context.Background(),
			os.Stdin, &stdout, &stderr, "port", "superset", "8088"); err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	address := strings.TrimSpace(stdout.String())
	if address == "" {
		return fmt.Errorf("no address found")
	}

	for {
		if _, err := net.Dial("tcp", address); err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	if err := browser.OpenURL(fmt.Sprintf("http://%s", address)); err != nil {
		errors.Wrap(err, "cannot open browser")
	}

	return nil
}

// OpenUI opens the browser with the UI.
//
// If opening the UI raises an error or took more time then `timeout`, then the
// error sent to the returned channel.
func OpenUI(timeout time.Duration) <-chan error {
	done := make(chan error)

	go func() {
		ch := make(chan error)
		go func() {
			ch <- openUI()
		}()

		select {
		case err := <-ch:
			done <- errors.Wrap(err, "an error occured while opening ui")
		case <-time.After(timeout):
			done <- fmt.Errorf("opening the ui took more than %v", timeout)
		}
	}()

	return done
}
