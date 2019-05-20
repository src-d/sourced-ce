package cmd

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/src-d/superset-compose/cmd/sandbox-ce/compose"
)

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
			os.Stdin, &stdout, nil, "port", "superset", "8088"); err == nil {
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

	select {
	case err := <-ch:
		return errors.Wrap(err, "an error occurred while opening the UI")
	case <-time.After(timeout):
		return fmt.Errorf("opening the UI took more than %v", timeout)
	}
}
