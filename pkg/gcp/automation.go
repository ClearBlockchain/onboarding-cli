package gcp

import (
	"fmt"
	"time"

	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
)

type Credentials struct {
	RedirectURI string
	ClientID    string
	ClientSecret string
}

func PurchaseOGI(product, gcpProject string) (*Credentials, error) {
	pw, browser, page := getPlaywright()
	defer browser.Close()
	defer pw.Stop()

	marketplaceUrl := fmt.Sprintf(
		"https://console.cloud.google.com/marketplace/product/opengatewayaggregation-public/%s?hl=en&project=%s",
		product,
		gcpProject,
	)

	if _, err := page.Goto(marketplaceUrl); err != nil {
		log.Fatalf("Failed to navigate to the marketplace URL: %v", err)
		return nil, err
	}

	// locate button with data-prober="cloud-marketplace-request-product" and role link
	buttonLocator := page.GetByRole("link", playwright.PageGetByRoleOptions{
		Name: "manage on provider, with an",
	})
	// fmt.Printf("%+v\n", buttonLocator)
	if err := buttonLocator.WaitFor(); err != nil {
		log.Fatalf("Failed to wait for the button: %v", err)
		return nil, err
	}

	if err := buttonLocator.Click(); err != nil {
		log.Fatalf("Failed to click on the button: %v", err)
		return nil, err
	}

	// click ok
	okButtonLocator := page.GetByRole("button", playwright.PageGetByRoleOptions{
		Name: "OK",
	})
	// fmt.Printf("%+v\n", okButtonLocator)
	if err := okButtonLocator.WaitFor(); err != nil {
		log.Fatalf("Failed to wait for the OK button: %v", err)
		return nil, err
	}

	if err := okButtonLocator.Click(); err != nil {
		log.Fatalf("Failed to click on the OK button: %v", err)
		return nil, err
	}

	defaultContext := browser.Contexts()[0]
	// find page with domain https://dev.gateway-x.io/gcp-ogi
	var internalPage playwright.Page

	// wait until page starts
	for {
		if len(defaultContext.Pages()) > 1 {
			break
		}

		time.Sleep(3 * time.Second)

		for _, c := range browser.Contexts() {
			for _, p := range c.Pages() {
				if p.URL() == "https://dev.gateway-x.io/gcp-ogi" {
					internalPage = p
					break
				}
			}
		}
	}

	groupLocator := internalPage.Locator(".input-group-wrapper")
	return &Credentials{
		RedirectURI: getInputValue(groupLocator.Nth(0)),
		ClientID:    getInputValue(groupLocator.Nth(1)),
		ClientSecret: getInputValue(groupLocator.Nth(2)),
	}, nil
}

func getInputValue(selectorOrLocator playwright.Locator) (string) {
	inputLocator := selectorOrLocator.Locator("input")

	if err := inputLocator.WaitFor(); err != nil {
		log.Fatalf("Failed to wait for the input: %v", err)
		return ""
	}

	content, err := inputLocator.InputValue()
	if err != nil {
		log.Fatalf("Failed to get the input value: %v", err)
		return ""
	}

	return content

}
