package utils

import (
	"fmt"
	"os"

	"github.com/ClearBlockchain/onboarding-cli/pkg/ui"
	"golang.org/x/term"
)

const sigilThresholdWidth = 80

func GetLongDescription() string {
	var response string

	if shouldDisplayASCIIArt() {
		response = getLongDescriptionFull()
	} else {
		response = getLongDescriptionText()
	}

	return response
}

func shouldDisplayASCIIArt() bool {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return false
	}

	return width >= sigilThresholdWidth
}

func getLongDescriptionText() string {
	var response string

	response = "\n" + ui.Keyword("ClearX Open Gateway") + " - One API, Every Telecom Network\n\n"
	response += "Use the following commands to connect with Glide:\n"
	response += "  1) " + ui.Keyword("glide login") + " - Add OGI to your GCP account & click " + ui.Underline("Manage on provider") + " to complete the auth flow.\n"
	response += "  2) " + ui.Keyword("glide docs") + " - Explore our developer's documentation.\n"

	// TODO: glide ai; glide help via onscreen docs

	return response
}

func getLongDescriptionFull() string {
	var response string

	/*

      **********************
    .************************.
  **************************
 ********.         .********         Welcome to ClearX Open Gateway
********             ********        One API, Every Telecom Network
*******               *******
*******               *******
********             *******.        Use the following commands to connect with Glide:
 *********        **********            1) glide login - Add OGI to your GCP account & click Manage on provider to complete the auth flow.
  .***********************              2) glide docs - Explore our developer's documentation.
     ******************
         ...*****...


  ****************************
.****************************
**************************
	*/
	response = "\n       **********************\n"
	response += "    .************************.\n"
	response += "  **************************  \n"
	response += fmt.Sprintf(" ********.         .********         %s\n", ui.Keyword("ClearX Open Gateway"))
	response += fmt.Sprintf("********             ********        %s\n", ui.Underline("One API, Every Telecom Network"))
	response += "*******               ******* \n"
	response += "*******               ******* \n"
	response += "********             *******.        Use the following commands to connect with Glide:\n"
	response += " *********        **********            1) glide login - Add OGI to your GCP account & click Manage on provider to complete the auth flow.\n"
	response += "  .***********************              2) glide docs - Explore our developer's documentation.\n"
	response += "     ******************      \n"
	response += "         ...*****...          \n"
	response += "                              \n"
	response += "                              \n"
	response += "  ****************************\n"
	response += ".**************************** \n"
	response += "**************************\n\n"

	return response
}
