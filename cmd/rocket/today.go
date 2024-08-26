package rocket

import (
	"fmt"
	"stargate/internal/rocket"

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
		rocket.CheckForTodayLaunches(launches)
	},
}
