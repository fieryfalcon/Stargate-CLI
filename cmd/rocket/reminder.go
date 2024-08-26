package rocket

import (
	"fmt"
	"stargate/internal/rocket"
	"strconv"

	"github.com/spf13/cobra"
)

var ReminderCmd = &cobra.Command{
	Use:   "reminder",
	Short: "Set a reminder for a specific launch",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide a launch ID.")
			return
		}
		launchID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid launch ID:", err)
			return
		}
		launches, err := rocket.FetchUpcomingLaunches()
		if err != nil {
			fmt.Println("Error fetching upcoming launches:", err)
			return
		}
		rocket.SetReminderForLaunch(launchID, launches)
	},
}
