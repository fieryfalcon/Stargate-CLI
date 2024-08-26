package cmd

import (
	"fmt"
	"stargate/cmd/rocket" // Import the subcommand packages

	"github.com/spf13/cobra"
)

var rocketCmd = &cobra.Command{
	Use:   "rocket",
	Short: "Commands related to rocket launches",
	Long:  "This command contains subcommands to fetch upcoming rocket launches, check if there are any launches today, and set reminders for specific launches.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use one of the subcommands: upcoming, today, reminder")
	},
}

func init() {
	// Register the rocket command as a root command
	rootCmd.AddCommand(rocketCmd)

	// Register the subcommands to the rocket command
	rocketCmd.AddCommand(rocket.UpcomingCmd)
	rocketCmd.AddCommand(rocket.TodayCmd)
	rocketCmd.AddCommand(rocket.ReminderCmd)
}
