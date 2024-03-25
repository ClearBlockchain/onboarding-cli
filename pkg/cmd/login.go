package cmd

import (
	"fmt"

	"github.com/ClearBlockchain/onboarding-cli/pkg/gcp"
	"github.com/ClearBlockchain/onboarding-cli/pkg/ui"
	"github.com/ClearBlockchain/onboarding-cli/pkg/utils"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getGCPPorjectsAsHuhOptions(options *[]huh.Option[string]) {
	for _, project := range gcp.GetGCPProjects() {
		*options = append(*options, huh.NewOption(project, project))
	}
}

var LoginCmd = &cobra.Command{
	Use: "init",
	Short: ui.Paragraph(
		fmt.Sprintf("Connect your local development environment with %s on Google Cloud Platform.", ui.Keyword("ClearX Open Gateway")),
	),
	Args: cobra.NoArgs,
	Aliases: []string{"auth", "authenticate", "signin", "sign-in", "connect", "login"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(
			"\n%s\n",
			ui.Box(
				fmt.Sprintf(
					"%s\n%s",
					ui.Title("Welcome to ClearX Open Gateway!"),
					"\nConnect your local development environment with ClearX Open Gateway on Google Cloud Platform.",
				),
			),
		)

		var (
			selectedAPIs []string
			selectedGCPProject string
			// confirmGCPFlow bool = false
			generateEnvFile bool = false
		)

		// Get the desired products
		for ok := true; ok; ok = len(selectedAPIs) == 0 {
			huh.NewMultiSelect[string]().
				Options(
					huh.NewOption("1) Telco Finder - Identify the Telecom provider of a phone number", "telco-finder"),
					huh.NewOption("2) SIM Swap Checker - Perform time-sensitive SIM Swap checks", "sim-swap"),
					huh.NewOption("3) Number Verify - Authenticate users by verifying their number association to the mobile network", "number-verify"),
				).
				Title("Which ClearX OGI endpoint you need access to?").
				Value(&selectedAPIs).
				Run()
		}

		// fetch gcp projects
		var gcpProjects []huh.Option[string]
		spinner.New().
			Title("Fetching your Google Cloud Projects").
			Action(func() {
				getGCPPorjectsAsHuhOptions(&gcpProjects)
			}).
			Run()

		// get the desired GCP project
		huh.NewSelect[string]().
			Options(gcpProjects...).
			Title("Select the Google Cloud Platform project you want to use").
			Value(&selectedGCPProject).
			Run()

		// open the browser tabs
		var credentials *gcp.Credentials
		for _, api := range selectedAPIs {
			cred, err := gcp.PurchaseOGI(api, selectedGCPProject)
			if err != nil {
				log.Errorf("Failed to purchase %s: %v", api, err)
				continue
			}

			credentials = cred
		}

		// ask the user if they want to generate a .env file
		huh.NewConfirm().
			Title("Do you want to generate a .env file with the credentials?").
			Value(&generateEnvFile).
			Run()

		if generateEnvFile {
			// generate the .env file
			if err := utils.WriteCredsToEnv(credentials); err != nil {
				log.Errorf("Failed to write credentials to .env file: %v", err)
			}
		}
	},
}
