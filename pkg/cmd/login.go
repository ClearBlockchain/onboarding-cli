package cmd

import (
	"fmt"
	"time"

	"github.com/ClearBlockchain/onboarding-cli/pkg/gcp"
	"github.com/ClearBlockchain/onboarding-cli/pkg/ui"
	"github.com/ClearBlockchain/onboarding-cli/pkg/utils"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

func getGCPPorjectsAsHuhOptions(options *[]huh.Option[string]) {
	for _, project := range gcp.GetGCPProjects() {
		*options = append(*options, huh.NewOption(project, project))
	}
}

func waitForCredenetials() {
	// sleep for 5 seconds
	time.Sleep(10 * time.Second)
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
			confirmGCPFlow bool = false
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

		// confirm the user understand the flow
		for ok := true; ok; ok = !confirmGCPFlow {
			huh.NewConfirm().
				Title(fmt.Sprintf(
					"\n%s\n\n%s",
					ui.Paragraph(
						fmt.Sprintf(
							"We will now open a %d browser tabs to complete the checkout process with your Google Cloud account. To complete the checkout click on the 'Subscribe' button on the Google Cloud Marketplace page.",
							len(selectedAPIs),
						),
					),
					ui.Paragraph("Are you ready to continue?"),
				)).
				Value(&confirmGCPFlow).
				Run()
		}

		// open the browser tabs
		for _, api := range selectedAPIs {
			utils.OpenBrowser(
				fmt.Sprintf(
					"https://console.cloud.google.com/marketplace/product/opengatewayaggregation-public/%s?hl=en&project=%s",
					api,
					selectedGCPProject,
				),
			)
		}

		// Let the user know that once he completed the checkout process
		// he needs to click on the 'Manage on provider' button
		// and we will with for the event on localhost:9919
		spinner.New().
			Title(
				fmt.Sprintf(
					"%s\n\n%s",
					"Once you've completed the checkout process, please click on the 'Manage on provider' button.",
					ui.Title("Waiting for credentials..."),
				),
			).
			Action(waitForCredenetials).
			Run()

		// ask the user if they want to generate a .env file
		huh.NewConfirm().
			Title("Do you want to generate a .env file with the credentials?").
			Value(&generateEnvFile).
			Run()
	},
}
