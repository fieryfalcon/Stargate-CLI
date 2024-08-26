package rocket

import (
	"fmt"
	"stargate/internal/rocket"

	"github.com/spf13/cobra"
)

var TomorrowCmd = &cobra.Command{
	Use:   "tomorrow",
	Short: "Check if there are any launches scheduled for tomorrow",
	Run: func(cmd *cobra.Command, args []string) {
		launches, err := rocket.FetchUpcomingLaunches()
		if err != nil {
			fmt.Println("Error fetching upcoming launches:", err)
			return
		}
		rocket.CheckForTomorrowLaunches(launches)
	},
}
