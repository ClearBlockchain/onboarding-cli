package gcp

import (
	"fmt"
	"strings"

	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
)

const MANAGE_ON_PROVIDER = "MANAGE ON PROVIDER"
const SUBSCRIBE = "SUBSCRIBE"

var (
	ErrUnsupportedGCPMarketplaceState = fmt.Errorf("unsupported GCP Marketplace state")
)

func PurchaseOGI(product, gcpProject string) (*Credentials, error) {
	pwi, err := getPlaywright()
	if err != nil {
		log.Errorf("Failed to get playwright: %v", err)
		return nil, err
	}
	defer pwi.Close()


	if _, err := pwi.Page.Goto(getMarketPlaceUrl(product, gcpProject)); err != nil {
		log.Fatalf("Failed to navigate to the marketplace URL: %v", err)
		return nil, err
	}

	state, buttonLocator, err := findSubscribeOrManageButton(pwi.Page)
	if err != nil {
		log.Errorf("Failed to find the button: %v", err)
		return nil, err
	}

	if state == SUBSCRIBE {
		if err := pwi.InjectScript(); err != nil {
			log.Errorf("Failed to inject script: %v", err)
			return nil, err
		}

		if _, err := waitForManageOnProviderToAppear(pwi.Page); err != nil {
			log.Errorf("Failed to wait for the button: %v", err)
			return nil, err
		}

		// continue credential fetching flow
		return fetchCredentials(pwi, buttonLocator)
	} else if state == MANAGE_ON_PROVIDER {
		return fetchCredentials(pwi, buttonLocator)
	} else {
		log.Errorf("Unknown state: %s", state)
		return nil, fmt.Errorf("unknown state: %s", state)
	}
}

func getMarketPlaceUrl(product, gcpProject string) string {
	return fmt.Sprintf(
		"https://console.cloud.google.com/marketplace/product/opengatewayaggregation-public/%s?hl=en&project=%s",
		product,
		gcpProject,
	)
}

func waitForManageOnProviderToAppear(p playwright.Page) (playwright.Locator, error) {
	manageLocator := p.GetByRole("link", playwright.PageGetByRoleOptions{
		Name: "manage on provider, with an",
	})

	if err := manageLocator.WaitFor(
		playwright.LocatorWaitForOptions{
			State: playwright.WaitForSelectorStateAttached,
			Timeout: playwright.Float(0.0),
		},
	); err != nil {
		log.Errorf("Failed to wait for the button: %v", err)
		return nil, err
	}

	return manageLocator, nil
}

// check which button is available - manage on provider or subscribe
func findSubscribeOrManageButton(p playwright.Page) (string, playwright.Locator, error) {
	manageLocator := p.GetByRole("link", playwright.PageGetByRoleOptions{
		Name: "manage on provider, with an",
	})
	subscribeLocator := p.GetByRole("button", playwright.PageGetByRoleOptions{
		Name: "subscribe",
	})

	mergedLocator := manageLocator.Or(subscribeLocator)

	if err := mergedLocator.WaitFor(); err != nil {
		log.Errorf("Failed to wait for the button: %v", err)
		return "", nil, err
	}

	// get locator actual internal content
	mergedLocatorContent, err := mergedLocator.InnerText()
	if err != nil {
		log.Errorf("Failed to get the button content: %v", err)
		return "", nil, err
	}

	cleanedContent := strings.TrimSpace(mergedLocatorContent)

	if cleanedContent == MANAGE_ON_PROVIDER {
		return MANAGE_ON_PROVIDER, manageLocator, nil
	} else if cleanedContent == SUBSCRIBE {
		return SUBSCRIBE, subscribeLocator, nil
	} else {
		log.Errorf("Unknown button: %s", cleanedContent)
		return "", nil, ErrUnsupportedGCPMarketplaceState
	}
}
