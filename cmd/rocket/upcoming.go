package rocket

import (
	"fmt"
	"stargate/internal/rocket"

	"github.com/spf13/cobra"
)

var UpcomingCmd = &cobra.Command{
	Use:   "upcoming",
	Short: "Get upcoming rocket launches",
	Run: func(cmd *cobra.Command, args []string) {
		launches, err := rocket.FetchUpcomingLaunches()
		if err != nil {
			fmt.Println("Error fetching upcoming launches:", err)
			return
		}
		for _, launch := range launches.Result {
			fmt.Printf("%s - %s (%s)\n", launch.Name, launch.LaunchDescription, launch.DateStr)
		}
	},
}
