package cli

import (
	"fmt"
	"os"

	"github.com/bhavya-dang/mstat/internal/git"
	"github.com/bhavya-dang/mstat/internal/listing"
	"github.com/bhavya-dang/mstat/internal/output"
	"github.com/bhavya-dang/mstat/internal/pathutil"
	"github.com/spf13/cobra"
)

var (
	noIcons      bool
	simpleIcons  bool
	briefView    bool
	extendedView bool
	noColor      bool
	noGit        bool
	porcelain    bool
	fullPath     bool
	rootCmd      = &cobra.Command{
		Use:   "mstat [file...]",
		Short: "Modern stat replacement",
		Long:  "mstat — a modern replacement for stat with bordered table output.",
		Args:  cobra.MinimumNArgs(1),
		RunE:  run,
	}
)

// SetVersion sets the version shown by mstat --version.
func SetVersion(v string) {
	rootCmd.Version = v
}

func init() {
	// icons
	rootCmd.Flags().BoolVarP(&noIcons, "no-icons", "n", false, "disable Nerd Font icons")
	rootCmd.Flags().BoolVarP(&simpleIcons, "simple-icons", "s", false, "show only basic icons")

	// views
	rootCmd.Flags().BoolVarP(&briefView, "brief", "b", false, "show minimal output (name, size, type)")
	rootCmd.Flags().BoolVarP(&extendedView, "extended", "x", false, "show extended output with all details")

	// color
	rootCmd.Flags().BoolVar(&noColor, "no-color", false, "disable colored output")

	// git
	rootCmd.Flags().BoolVar(&noGit, "no-git", false, "disable git status column")
	rootCmd.Flags().BoolVar(&porcelain, "porcelain", false, "machine-readable output (no borders, no colors)")

	// paths
	rootCmd.Flags().BoolVar(&fullPath, "full-path", false, "show full absolute paths")
}

func run(cmd *cobra.Command, args []string) error {
	var entries []listing.Entry
	var paths []string
	for _, arg := range args {
		path := pathutil.Expand(arg)
		e, err := listing.Stat(path, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mstat: %v\n", err)
			continue
		}
		entries = append(entries, e)
		paths = append(paths, path)
	}

	if len(entries) == 0 {
		return fmt.Errorf("no valid files")
	}

	// git status
	var gitMap map[string]git.Status
	if !noGit {
		if root := git.RepoRoot("."); root != "" {
			gitMap = git.StatusMap(root)
		}
	}

	opts := output.Options{
		Icons:        !noIcons,
		SimpleIcons:  simpleIcons,
		BriefView:    briefView,
		ExtendedView: extendedView,
		NoColor:      noColor,
		NoGit:        noGit,
		Porcelain:    porcelain,
		FullPath:     fullPath,
		GitMap:       gitMap,
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
