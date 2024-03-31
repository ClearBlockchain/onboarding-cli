package gcp

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/ClearBlockchain/onboarding-cli/pkg/utils"
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
)

var (
	ErrUnsupportedPlatform = fmt.Errorf("unsupported platform")
	cpWg sync.WaitGroup
	profileDataDir string
	profileDataDirErr error
)

type PlaywrightInfo struct {
	Playwright *playwright.Playwright
	Context	playwright.BrowserContext
	Browser    playwright.Browser
	Page       playwright.Page
}

func init() {
	err := playwright.Install(&playwright.RunOptions{
		Stdout: nil,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Playwright installed")

	// deep copy and mark as done
	cpWg.Add(1)
	go func() {
		profileDataDir, profileDataDirErr = deepCopyDefaultChromeProfile()
		if profileDataDirErr != nil {
			log.Errorf("Failed to copy default chrome profile: %v", profileDataDirErr)
		}
	}()
}

// return playwright, browser, and page.
func getPlaywright() (pwi *PlaywrightInfo, err error) {
	log.Debug("Getting playwright")
	pw, err := playwright.Run()
	if err != nil {
		log.Fatal(err)
	}

	var browser playwright.Browser

	log.Debug("Connecting to the browser")
	browser, err = pw.Chromium.ConnectOverCDP("http://localhost:9222")
	if err != nil {
		log.Debug("Couldn't find browser with port 9222, opening a new one")
		_, err := openChrome()
		if err != nil {
			log.Errorf("Failed to open chrome: %v", err)
			return nil, err
		}

		log.Debug("Retrying to connect to the browser")
		browser, err = pw.Chromium.ConnectOverCDP("http://localhost:9222")
		if err != nil {
			log.Errorf("Failed to connect to the browser: %v", err)
			return nil, err
		}

		log.Debug("Connected to the browser")
	}

	log.Debug("Getting context and page")
	defaultContext, page, err := getContextAndPage(browser)
	if err != nil {
		log.Errorf("Failed to get context and page: %v", err)
		return nil, err
	}

	return &PlaywrightInfo{
		Playwright: pw,
		Context: defaultContext,
		Browser:    browser,
		Page:       page,
	}, nil
}

func getContextAndPage(browser playwright.Browser) (playwright.BrowserContext, playwright.Page, error) {
	var defaultContext playwright.BrowserContext
	var page playwright.Page

	if len(browser.Contexts()) == 0 {
		dc, err := browser.NewContext()
		if err != nil {
			log.Errorf("Failed to create a new context: %v", err)
			return nil, nil, err
		}

		defaultContext = dc
	} else {
		defaultContext = browser.Contexts()[0]
	}

	if len(defaultContext.Pages()) == 0 {
		p, err := defaultContext.NewPage()
		if err != nil {
			log.Errorf("Failed to create a new page: %v", err)
			return nil, nil, err
		}

		page = p
	} else {
		page = defaultContext.Pages()[0]
	}

	return defaultContext, page, nil
}

func (p *PlaywrightInfo) Close() {
	// close each of the open tabs
	for _, page := range p.Browser.Contexts()[0].Pages() {
		page.Close()
	}

	// close the browser
	p.Browser.Close()
	p.Playwright.Stop()
}

func getExecutablePath() (string, error) {
	switch runtime.GOOS {
	case "linux":
		return "/usr/bin/google-chrome", nil
	case "windows":
		return "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe", nil
	case "darwin":
		return "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", nil
	default:
		return "", ErrUnsupportedPlatform
	}
}

func openChrome() (*exec.Cmd, error) {
	// find default chrome executable
	chromePath, err := getExecutablePath()
	if err != nil {
		log.Errorf("Failed to get chrome executable path: %v", err)
		return nil, err
	}
	log.Debug("Chrome executable path: ", chromePath)

	// wait for copy to complete
	cpWg.Wait()
	if profileDataDirErr != nil {
		return nil, profileDataDirErr
	}

	args := []string{
		"--remote-debugging-port=9222",
		"--user-data-dir=" + profileDataDir,
	}

	cmd := exec.Command(chromePath, args...)
	if err := cmd.Start(); err != nil {
		log.Errorf("Failed to start chrome: %v", err)
		return nil, err
	}

	// wait for the browser debugging port to be available
	if err := waitForPort("localhost:9222"); err != nil {
		log.Errorf("Failed to wait for the debugging port: %v", err)
		return nil, err
	}

	return cmd, nil
}

func getDefaultChromeProfileDataDir() (string, error) {
	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf("%s/.config/google-chrome", os.Getenv("HOME")), nil
	case "windows":
		return fmt.Sprintf("%s\\AppData\\Local\\Google\\Chrome\\User Data", os.Getenv("USERPROFILE")), nil
	case "darwin":
		return fmt.Sprintf("%s/Library/Application Support/Google/Chrome", os.Getenv("HOME")), nil
	default:
		return "", ErrUnsupportedPlatform
	}
}

func deepCopyDefaultChromeProfile() (string, error) {
	defer cpWg.Done()

	log.Debug("Copying default chrome profile")
	tmpDir := fmt.Sprintf("%s/glide-chrome-profile/", os.TempDir())

	// check if tmp directory already exists
	// and skip the copy if it does
	if _, err := os.Stat(tmpDir); err == nil {
		log.Debugf("Temporary directory already exists: %s", tmpDir)
		return tmpDir, nil
	}

	profileDataDir, err := getDefaultChromeProfileDataDir()
	if err != nil {
		log.Errorf("Failed to get default chrome profile data dir: %v", err)
		return "", err
	}
	log.Debugf("Default chrome profile data dir: %s", profileDataDir)

	// copy the default chrome profile to a temporary directory
	if err := utils.CopyDirectory(profileDataDir, tmpDir, true); err != nil {
		log.Errorf("Failed to copy default chrome profile: %v", err)
		return "", err
	}

	log.Debugf("Copied default chrome profile to %s", tmpDir)
	return tmpDir, nil
}

func waitForPort(port string) error {
	timeout := time.After(60 * time.Second)
	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for port %s", port)
		default:
			conn, err := net.Dial("tcp", port)
			if err == nil {
				conn.Close()
				return nil
			}

			log.Debugf("Waiting for port %s", port)
			time.Sleep(1 * time.Second)
		}
	}
}
