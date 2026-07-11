package cli

import (
	"fmt"
	"os"

	"github.com/bhavya-dang/mstat/internal/listing"
	"github.com/bhavya-dang/mstat/internal/output"
	"github.com/bhavya-dang/mstat/internal/pathutil"
	"github.com/spf13/cobra"
)

var (
	noIcons     bool
	simpleIcons bool
	rootCmd     = &cobra.Command{
		Use:     "mstat [file...]",
		Version: "0.0.1",
		Short:   "Modern stat replacement",
		Long:    "mstat — a modern replacement for stat with bordered table output.",
		Args:    cobra.MinimumNArgs(1),
		RunE:    run,
	}
)

func init() {
	rootCmd.Flags().BoolVar(&noIcons, "no-icons", false, "disable Nerd Font icons")
	rootCmd.Flags().BoolVar(&simpleIcons, "simple-icons", false, "show only basic icons")

}

func run(cmd *cobra.Command, args []string) error {
	var entries []listing.Entry
	for _, arg := range args {
		path := pathutil.Expand(arg)
		e, err := listing.Stat(path, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mstat: %v\n", err)
			continue
		}
		entries = append(entries, e)
	}

	if len(entries) == 0 {
		return fmt.Errorf("no valid files")
	}

	opts := output.Options{
		Icons:       !noIcons,
		SimpleIcons: simpleIcons,
	}
	output.Render(os.Stdout, entries, opts)
	return nil
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
