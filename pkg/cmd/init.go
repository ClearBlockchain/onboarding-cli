package cmd

import (
	"fmt"
	"os"

	"github.com/ClearBlockchain/onboarding-cli/pkg/gcp"
	"github.com/ClearBlockchain/onboarding-cli/pkg/ui"
	"github.com/ClearBlockchain/onboarding-cli/pkg/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getGCPPorjectsAsHuhOptions(options *[]huh.Option[string]) {
	projects, err := gcp.GetGCPProjects()
	if err != nil {
		log.Errorf("Failed to fetch GCP projects: %v", err)
		return
	}

	for _, project := range projects {
		*options = append(*options, huh.NewOption(project, project))
	}
}

var LoginCmd = &cobra.Command{
	Use: "init",
	Short: ui.Paragraph.Render(
		fmt.Sprintf("Connect your local development environment with %s on Google Cloud Platform.", ui.Keyword.Render("Glide")),
	),
	Args: cobra.NoArgs,
	Aliases: []string{"auth", "authenticate", "signin", "sign-in", "connect", "login"},
	Run: func(cmd *cobra.Command, args []string) {
		// fetch gcp projects
		var gcpProjects []huh.Option[string]
		spinner.New().
			Title("Fetching your Google Cloud Projects").
			Action(func() {
				if !gcp.CheckGcloudExists() {
					var shouldInstallGcloud bool

					// ask the user if the want to install gcloud
					huh.NewConfirm().
						Title("gcloud not found").
						Description("gcloud command not found. Do you want to install it?").
						Affirmative("Sure").
						Negative("No").
						Value(&shouldInstallGcloud).
						Run()

					if shouldInstallGcloud {
						if err := gcp.InstallGcloud(); err != nil {
							log.Errorf("Failed to install gcloud: %v", err)
							return
						}
					} else {
						log.Error("gcloud command not found")
						os.Exit(1)
					}
				}

				getGCPPorjectsAsHuhOptions(&gcpProjects)
			}).
			Run()

		// run tui model
		model := ui.NewModel(
			huh.NewForm(
				huh.NewGroup(
					huh.NewMultiSelect[string]().
						Key("endpoints").
						Options(
							huh.NewOption("1) Telco Finder - Find the telcom provider for the number", "telco-finder"),
							huh.NewOption("2) SIM Swap Checker - Check if the SIM was swapped lately", "sim-swap"),
							huh.NewOption("3) Number Verify - Verify the number association to the network", "number-verify"),
						).
						Title("Choose the Glide endpoints you need").
						Description("The selected endpoints will be added to your Google Cloud Platform project.").
						Validate(func(value []string) error {
							if len(value) == 0 {
								return fmt.Errorf("Well, you must select at least one endpoint")
							}
							return nil
						}),

					huh.NewSelect[string]().
						Key("gcpProject").
						Options(gcpProjects...).
						Title("Google Cloud Platform Project").
						Description("Select the project you want to use for the Glide setup.").
						Validate(func(value string) error {
							if value == "" {
								return fmt.Errorf("We need to associate the Glide endpoints with a GCP project - please select one.")
							}
							return nil
						}),

					huh.NewConfirm().
						Key("done").
						Title("Write to .env?").
						Description("We'll write the .env to your project root with the relevant credentials.").
						Affirmative("Yep").
						Negative("Nah, just print it."),
				),
			),
		)

		if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
			log.Fatalf("Oh no: %v", err)
		}

		// find tea quit and return
		if model.Form.State != huh.StateCompleted {
			return
		}

		// open the browser tabs
		var credentials *gcp.Credentials
		for _, api := range model.Form.Get("endpoints").([]string) {
			cred, err := gcp.PurchaseOGI(api, model.Form.GetString("gcpProject"))
			if err != nil {
				log.Errorf("Failed to purchase %s: %v", api, err)
				return
			}

			credentials = cred
		}

		if model.Form.GetBool("done") {
			// generate the .env file
			if err := utils.WriteCredsToEnv(credentials.ToMap()); err != nil {
				log.Errorf("Failed to write credentials to .env file: %v", err)
			}
		} else {
			// print the credentials
			fmt.Printf(
				"\n%s\n",
				ui.Box.
					Width(80).
					BorderForeground(ui.Indigo).
					Align(lipgloss.Left).
					Render(
						fmt.Sprintf(
							"%s\n\n%s",
							ui.Highlight.Render("Your Glide Credentials:"),
							credentials.ToString(),
						),
					),
			)
		}
	},
}
