package gcp

import (
	"time"

	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
)

func (p *PlaywrightInfo) waitForPageLoad(url string) (playwright.Page, error) {
	var internalPage playwright.Page
	for {
		if len(p.Context.Pages()) > 1 {
			break
		}
		time.Sleep(3 * time.Second)
		for _, c := range p.Browser.Contexts() {
			for _, page := range c.Pages() {
				if page.URL() == url {
					internalPage = page
					break
				}
			}
		}
	}
	return internalPage, nil
}

func (p *PlaywrightInfo) getInputValue(selectorOrLocator playwright.Locator) string {
	inputLocator := selectorOrLocator.Locator("input")

	if err := inputLocator.WaitFor(); err != nil {
		log.Errorf("Failed to wait for the input: %v", err)
		return ""
	}

	content, err := inputLocator.InputValue()
	if err != nil {
		log.Errorf("Failed to get the input value: %v", err)
		return ""
	}

	return content
}

func (p *PlaywrightInfo) clickButton(role playwright.AriaRole, name string) error {
	buttonLocator := p.Page.GetByRole(role, playwright.PageGetByRoleOptions{
		Name: name,
	})
	if err := buttonLocator.WaitFor(); err != nil {
		return err
	}

	if err := buttonLocator.Click(); err != nil {
		return err
	}
	return nil
}
