package gcp

import (
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
)

type Credentials struct {
	RedirectURI string
	ClientID    string
	ClientSecret string
}

func (c *Credentials) ToMap() map[string]string {
	return map[string]string{
		"GLIDE_REDIRECT_URI": c.RedirectURI,
		"GLIDE_CLIENT_ID":    c.ClientID,
		"GLIDE_CLIENT_SECRET": c.ClientSecret,
	}
}

func (c *Credentials) ToString() string {
	return "GLIDE_REDIRECT_URI=" + c.RedirectURI + "\n" +
		"GLIDE_CLIENT_ID=" + c.ClientID + "\n" +
		"GLIDE_CLIENT_SECRET=" + c.ClientSecret + "\n"
}

func fetchCredentials(
	pwi *PlaywrightInfo,
	actionLocator playwright.Locator,
) (*Credentials, error) {
	// click manage on provider
	if err := actionLocator.Click(); err != nil {
		log.Errorf("Failed to click on the button: %v", err)
		return nil, err
	}

	// complete checkout
	// click ok
	if err := pwi.clickButton("button", "OK"); err != nil {
		log.Errorf("Failed to click on the OK button: %v", err)
		return nil, err
	}

	// find page with domain https://dev.gateway-x.io/gcp-ogi
	internalPage, err := pwi.waitForPageLoad("https://dev.gateway-x.io/gcp-ogi")
	if err != nil {
		log.Errorf("Failed to wait for the page load: %v", err)
		return nil, err
	}

	groupLocator := internalPage.Locator(".input-group-wrapper")

	if err := groupLocator.WaitFor(); err != nil {
		log.Errorf("Failed to wait for the group locator: %v", err)
		return nil, err
	}

	return &Credentials{
		RedirectURI: pwi.getInputValue(groupLocator.Nth(0)),
		ClientID:    pwi.getInputValue(groupLocator.Nth(1)),
		ClientSecret: pwi.getInputValue(groupLocator.Nth(2)),
	}, nil
}
