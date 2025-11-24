package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/emmanuelay/badger/internal/app"
	"github.com/emmanuelay/badger/internal/config"
	"github.com/spf13/cobra"
)

var (
	// Version information (set at build time)
	version = "dev"
	commit  = "none"
	date    = "unknown"

	// Configuration flags
	allCharacters bool
	alpha         bool
	alphaNumeric  bool
	numeric       bool
	customRange   string
	delay         int64
	tlds          string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "domainbadger [flags] <searchterms>",
	Short: "Find unregistered domains using wildcards & character combinations",
	Long: `Domainbadger is a fast and simple CLI tool used to find unregistered domains
using wildcards & characterset combinations.

Underscore (_) is treated as the wildcard character and is replaced with
different combinations of the selected character set.

Example:
  domainbadger -custom aoe -tld se,io,nu h_ll_ w_rld d_min_ti_n

This will perform WHOIS lookups on domain combinations using the characters
'aoe' in the search terms 'h_ll_', 'w_rld', and 'd_min_ti_n' with the TLDs
'se', 'io', and 'nu'.`,
	Example: `  # Search with custom characters
  domainbadger -c aoe -t se,io h_ll_ w_rld
  domainbadger --custom aoe --tld se,io h_ll_ w_rld

  # Search with alphabetic characters
  domainbadger -a -t com,net c_t d_g

  # Search with alphanumeric characters
  domainbadger -n -t io t_st

  # Search with numeric characters
  domainbadger -N -t com ap_ _pp`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Allow no args if only showing version
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			return nil
		}
		return cobra.MinimumNArgs(1)(cmd, args)
	},
	RunE: runDomainCheck,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute(ver, cmt, dt string) error {
	version = ver
	commit = cmt
	date = dt
	return rootCmd.Execute()
}

func init() {
	// Character set flags (mutually exclusive - if any specific flag is set, disable "all")
	rootCmd.Flags().BoolVarP(&alpha, "alpha", "a", false, "Use alphabetic range (a-z)")
	rootCmd.Flags().BoolVarP(&alphaNumeric, "alphanum", "n", false, "Use alphanumeric range (a-z, 0-9)")
	rootCmd.Flags().BoolVarP(&numeric, "numeric", "N", false, "Use numeric range (0-9)")
	rootCmd.Flags().BoolVarP(&allCharacters, "all", "A", false, "Use all possible characters (a-z, 0-9, -)")
	rootCmd.Flags().StringVarP(&customRange, "custom", "c", "", "Use a custom character range (ex. abc123)")

	// Lookup configuration flags
	rootCmd.Flags().Int64VarP(&delay, "delay", "d", 500, "Delay between lookup attempts, in milliseconds")
	rootCmd.Flags().StringVarP(&tlds, "tld", "t", "com", "TLDs to search. Use comma to add multiple (ex. com,org,net)")

	// Version flag
	rootCmd.Flags().BoolP("version", "v", false, "Show version information")
}

func runDomainCheck(cmd *cobra.Command, args []string) error {
	// Handle version flag
	versionFlag, _ := cmd.Flags().GetBool("version")
	if versionFlag {
		fmt.Printf("domainbadger version %s\ncommit: %s\nbuilt: %s\n", version, commit, date)
		return nil
	}

	// Build configuration from flags
	cfg := config.Configuration{
		AllCharacters:  allCharacters,
		Alpha:          alpha,
		AlphaNumeric:   alphaNumeric,
		Numeric:        numeric,
		CustomRange:    customRange,
		Delay:          delay,
		TLD:            strings.Split(tlds, ","),
		SearchPatterns: args,
	}

	// Validate configuration
	if err := config.ValidateConfiguration(cfg); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	// Run the application
	ctx := context.Background()
	app.Run(ctx, cfg)

	return nil
}
