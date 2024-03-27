# Start The Project

```console
# from the project root
❯ go run cmd/glide/main.go

       **********************
    .************************.
  **************************
 ********.         .********         ClearX Open Gateway
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


  Explore ClearX Open Gateway Integration Layer, a one-stop infrastructure offering access to Network APIs, 5G, and edge resources from worldwide CSPs and
  perform time-sensitive SIM Swap checks using a unified interface and API

Usage:
  glide [flags]
  glide [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init          Connect your local development environment with ClearX Open Gateway on Google Cloud Platform.

Flags:
  -h, --help      help for glide
  -v, --version   version for glide

Use "glide [command] --help" for more information about a command.
```

# Glide Project Initialization

```console
❯ go run cmd/glide/main.go init

   ClearX OGI Project Initialization /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

 ┃ Choose the OGI endpoints you need                                                                    ╭──────────────────────────────────────────────────╮
 ┃ The selected endpoints will be added to your Google Cloud Platform project.                          │ Your ClearX OGI Setup                            │
 ┃ > • 1) Telco Finder - Find the telcom provider for the number                                        │ (None)                                           │
 ┃   • 2) SIM Swap Checker - Check if the SIM was swapped lately                                        │                                                  │
 ┃   • 3) Number Verify - Verify the number association to the network                                  │                                                  │
                                                                                                        │                                                  │
   Google Cloud Platform Project                                                                        │                                                  │
   Select the project you want to use for the OGI setup.                                                │                                                  │
   > opengatewayaggregation-public                                                                      │                                                  │
     kontax                                                                                             │                                                  │
     salfati-group-cloud                                                                                │                                                  │
     nopeus                                                                                             │                                                  │
     elon-private                                                                                       │                                                  │
                                                                                                        │                                                  │
   Write to .env?                                                                                       │                                                  │
   We'll write the .env to your project root with the relevant credentials.                             │                                                  │
                                                                                                        │                                                  │
     Yep     Nah, just print it.                                                                        │                                                  │
                                                                                                        │                                                  │
                                                                                                        │                                                  │
                                                                                                        │                                                  │
                                                                                                        │                                                  │
                                                                                                        │                                                  │
                                                                                                        ╰──────────────────────────────────────────────────╯

   x toggle • ↑ up • ↓ down • / filter • enter confirm ///////////////////////////////////////////////////////////////////////////////////////////////////
```
