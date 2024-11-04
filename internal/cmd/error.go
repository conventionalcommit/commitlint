package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

// handleError handles and logs errors in a consistent way
func handleError(err error, customMessage string) error {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s - %v\n", customMessage, err)
		return cli.Exit(customMessage, 1)
	}
	return nil
}
