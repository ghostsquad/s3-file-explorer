package cmds

import (
	_ "net/http"

	_ "github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/ghostsquad/s3-file-explorer/internal/clioptions"
	"github.com/ghostsquad/s3-file-explorer/internal/clioptions/iostreams"
)

func NewRootCmd(ioStreams iostreams.IOStreams) *cobra.Command {
	globalOpts := clioptions.GlobalOptions{
		IOStreams: ioStreams,
	}

	rootCmd := &cobra.Command{
		Use:           "s3-file-explorer",
		Short:         "a simple url based s3 file explorer",
		Args:          cobra.ArbitraryArgs,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.ParseFlags(args); err != nil {
				return err
			}

			return nil
		},
	}

	return rootCmd
}
