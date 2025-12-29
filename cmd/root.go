package cmd

import (
	"fmt"
	"net/netip"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/suraiborys/subnetcalc/app/subnetcalc"
)

var rootCmd = &cobra.Command{
	Use:   "snc <cidr>",
	Short: "Calculate subnet information from CIDR notation",
	Long:  "Calculate subnet information from CIDR notation.",
	Example: `# calculate subnet information for 192.168.1.0/24
snc 192.168.1.0/24`,
	Version: "0.1.0",
	Args:    cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		prefix, err := netip.ParsePrefix(args[0])
		if err != nil {
			return fmt.Errorf("invalid prefix: %s", err)

		}

		result, err := subnetcalc.CalcSubnetInfo(prefix)
		if err != nil {
			return fmt.Errorf("error calculating subnet info: %s", err)
		}

		bold := color.New(color.Bold).SprintFunc()

		fmt.Printf("Network Address:    %s\n", bold(result.NetworkAddress))
		fmt.Printf("Broadcast Address:  %s\n", bold(result.BroadcastIP))
		fmt.Printf("Subnet Mask:        %s\n", bold(result.SubnetMask))
		fmt.Printf("Total IPs:          %s\n", bold(result.TotalIP))

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func Root() *cobra.Command { return rootCmd }

func init() {
	// Flags and configuration can go here
}
