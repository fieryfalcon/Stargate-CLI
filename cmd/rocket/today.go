package rocket

import (
	"fmt"
	"stargate/internal/rocket"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var TodayCmd = &cobra.Command{
	Use:   "today",
	Short: "Check if there are any launches today",
	Run: func(cmd *cobra.Command, args []string) {
		launches, err := rocket.FetchUpcomingLaunches()
		if err != nil {
			fmt.Println("Error fetching upcoming launches:", err)
			return
		}

		header := color.New(color.FgCyan, color.Bold).PrintfFunc()
		highlight := color.New(color.FgGreen).PrintfFunc()
		detail := color.New(color.FgWhite).PrintfFunc()

		header("Today's Rocket Launches:\n\n")

		today := rocket.CheckForTodayLaunches(launches)
		if len(today) == 0 {
			color.Red("No launches scheduled for today.\n")
			return
		}

		for _, launch := range today {
			highlight("ID: %d | %s\n", launch.ID, launch.Name)
			detail("  Launch Description: %s\n", launch.LaunchDescription)
			detail("  Provider: %s\n", launch.Provider.Name)
			detail("  Vehicle: %s\n", launch.Vehicle.Name)
			detail("  Launch Pad: %s, %s, %s\n", launch.Pad.Name, launch.Pad.Location.State, launch.Pad.Location.Country)
			detail("  Date: %s\n", launch.DateStr)
			detail("  More Info: https://rocketlaunch.live/launch/%s\n", launch.Slug)
			fmt.Println()  // Add a blank line between launches
		}
	},
}
