package cmd

import (
	"fmt"
	"stargate/cmd/rocket"

	"github.com/spf13/cobra"
)

var rocketCmd = &cobra.Command{
	Use:   "rocket",
	Short: "Commands related to rocket launches",
	Long:  "This command contains subcommands to fetch upcoming rocket launches, check if there are any launches today or tomorrow, and set reminders for specific launches.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use one of the subcommands: upcoming, today, tomorrow, reminder")
	},
}

func init() {
	rootCmd.AddCommand(rocketCmd)

	rocketCmd.AddCommand(rocket.UpcomingCmd)
	rocketCmd.AddCommand(rocket.TodayCmd)
	rocketCmd.AddCommand(rocket.TomorrowCmd)
	rocketCmd.AddCommand(rocket.ReminderCmd)
}
