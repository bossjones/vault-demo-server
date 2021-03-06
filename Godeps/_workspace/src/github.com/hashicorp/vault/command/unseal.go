package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/helper/password"
)

// UnsealCommand is a Command that unseals the vault.
type UnsealCommand struct {
	Meta

	// Key can be used to pre-seed the key. If it is set, it will not
	// be asked with the `password` helper.
	Key string
}

func (c *UnsealCommand) Run(args []string) int {
	var reset bool
	flags := c.Meta.FlagSet("unseal", FlagSetDefault)
	flags.BoolVar(&reset, "reset", false, "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	args = flags.Args()

	value := c.Key
	if len(args) > 0 {
		value = args[0]
	}
	if value == "" {
		fmt.Printf("Key (will be hidden): ")
		value, err = password.Read(os.Stdin)
		fmt.Printf("\n")
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error attempting to ask for password. The raw error message\n"+
					"is shown below, but the most common reason for this error is\n"+
					"that you attempted to pipe a value into unseal or you're\n"+
					"executing `vault unseal` from outside of a terminal.\n\n"+
					"You should use `vault unseal` from a terminal for maximum\n"+
					"security. If this isn't an option, the unseal key can be passed\n"+
					"in using the first parameter.\n\n"+
					"Raw error: %s", err))
			return 1
		}
	}

	status, err := client.Sys().Unseal(strings.TrimSpace(value))
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error attempting unseal: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf(
		"Sealed: %v\n"+
			"Key Shares: %d\n"+
			"Key Threshold: %d\n"+
			"Unseal Progress: %d",
		status.Sealed,
		status.N,
		status.T,
		status.Progress,
	))

	return 0
}

func (c *UnsealCommand) Synopsis() string {
	return "Unseals the vault server"
}

func (c *UnsealCommand) Help() string {
	helpText := `
Usage: vault unseal [options] [key]

  Unseal the vault by entering a portion of the master key. Once all
  portions are entered, the Vault will be unsealed.

  Every Vault server initially starts as sealed. It cannot perform any
  operation except unsealing until it is sealed. Secrets cannot be accessed
  in any way until the vault is unsealed. This command allows you to enter
  a portion of the master key to unseal the vault.

  The unseal key can be specified via the command line, but this is
  not recommended. The key may then live in your terminal history. This
  only exists to assist in scripting.

General Options:

  -address=TODO           The address of the Vault server.

  -ca-cert=path           Path to a PEM encoded CA cert file to use to
                          verify the Vault server SSL certificate.

  -ca-path=path           Path to a directory of PEM encoded CA cert files
                          to verify the Vault server SSL certificate. If both
                          -ca-cert and -ca-path are specified, -ca-path is used.

  -insecure               Do not verify TLS certificate. This is highly
                          not recommended. This is especially not recommended
                          for unsealing a vault.

Unseal Options:

  -reset                  Reset the unsealing process by throwing away
                          prior keys in process to unseal the vault.

`
	return strings.TrimSpace(helpText)
}
