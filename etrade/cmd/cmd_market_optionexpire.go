package cmd

import (
	"fmt"
	"github.com/jerryryle/etrade-cli/pkg/etradelib/client/constants"
	"github.com/spf13/cobra"
)

type marketOptionExpireFlags struct {
	expiryType enumFlagValue[constants.OptionExpiryType]
}

type CommandMarketOptionexpire struct {
	Resources *CommandResources
	flags     marketOptionExpireFlags
}

func (c *CommandMarketOptionexpire) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "optionexpire [symbol]",
		Short: "Get option expire dates",
		Long:  "Get option expire dates for a specific underlying instrument",
		Args:  cobra.MatchAll(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.GetOptionExpireDates(args[0])
		},
	}

	// Initialize Enum Flag Values
	c.flags.expiryType = *newEnumFlagValue(expiryTypeMap, constants.OptionExpiryTypeNil)

	// Add Enum Flags
	cmd.Flags().VarP(
		&c.flags.expiryType, "expiry-type", "e",
		fmt.Sprintf("expiry type (%s)", c.flags.expiryType.JoinAllowedValues(", ")),
	)
	_ = cmd.RegisterFlagCompletionFunc(
		"expiry-type",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return c.flags.expiryType.AllowedValuesWithHelp(), cobra.ShellCompDirectiveDefault
		},
	)

	return cmd
}

func (c *CommandMarketOptionexpire) GetOptionExpireDates(symbol string) error {
	response, err := c.Resources.Client.GetOptionExpireDates(symbol, c.flags.expiryType.Value())
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(c.Resources.OFile, string(response))
	return nil
}

var expiryTypeMap = map[string]enumValueWithHelp[constants.OptionExpiryType]{
	"unspecified": {constants.OptionExpiryTypeUnspecified, "unspecified expiry type"},
	"daily":       {constants.OptionExpiryTypeDaily, "daily expiry type"},
	"weekly":      {constants.OptionExpiryTypeWeekly, "weekly expiry type"},
	"monthly":     {constants.OptionExpiryTypeMonthly, "monthly expiry type"},
	"quarterly":   {constants.OptionExpiryTypeQuarterly, "quarterly expiry type"},
	"vix":         {constants.OptionExpiryTypeVix, "VIX expiry type"},
	"all":         {constants.OptionExpiryTypeAll, "all expiry types"},
	"monthEnd":    {constants.OptionExpiryTypeMonthEnd, "month-end expiry type"},
}
