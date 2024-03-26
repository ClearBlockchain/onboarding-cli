package gcp

import (
	_ "embed"

	log "github.com/sirupsen/logrus"
)

var (
	//go:embed scripts/inject.js
	injectScript string
)

func (p *PlaywrightInfo) InjectScript() error {
	_, err := p.Page.Evaluate(injectScript)
	if err != nil {
		log.Errorf("Failed to inject script: %v", err)
		return err
	}
	return nil
}
