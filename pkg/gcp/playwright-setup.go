package gcp

import (
	"net"
	"os/exec"
	"runtime"
	"time"

	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
)

func init() {
	err := playwright.Install(&playwright.RunOptions{
		Stdout: nil,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// return playwright, browser, and page.
func getPlaywright() (*playwright.Playwright, playwright.Browser, playwright.Page) {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatal(err)
	}

	openChrome(false)

	browser, err := pw.Chromium.ConnectOverCDP("http://localhost:9222")
	if err != nil {
		log.Fatal(err)
	}

	defaultContext := browser.Contexts()[0]
	page := defaultContext.Pages()[0]

	return pw, browser, page
}

func openChrome(headless bool) {
	// find default chrome executable
	var chromePath string
	switch runtime.GOOS {
	case "linux":
		chromePath = "/usr/bin/google-chrome"
	case "windows":
		chromePath = "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe"
	case "darwin":
		chromePath = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	default:
		log.Fatal("unsupported platform")
	}

	args := []string{"--remote-debugging-port=9222"}
	if headless {
		args = append(args, "--headless")
	}

	cmd := exec.Command(chromePath, args...)
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start chrome: %v", err)
	}

	// wait for the browser debugging port to be available
	if err := waitForPort("localhost:9222"); err != nil {
		log.Fatalf("Failed to wait for the debugging port: %v", err)
	}
}

func waitForPort(port string) error {
	for {
		conn, err := net.Dial("tcp", port)
		if err == nil {
			conn.Close()
			return nil
		}

		time.Sleep(1 * time.Second)
	}
}
