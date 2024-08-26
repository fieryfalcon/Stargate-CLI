package rocket

import (
	"encoding/json"
	"fmt"
	"time"

	"stargate/pkg/models"
	"stargate/pkg/utils"
)

const rocketLaunchAPI = "https://fdo.rocketlaunch.live/json/launches/next/10"

// FetchUpcomingLaunches fetches the upcoming rocket launches
func FetchUpcomingLaunches() (*models.RocketLaunchResponse, error) {
	resp, err := utils.MakeGetRequest(rocketLaunchAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var launches models.RocketLaunchResponse
	if err := json.NewDecoder(resp.Body).Decode(&launches); err != nil {
		return nil, err
	}

	return &launches, nil
}

// CheckForTodayLaunches checks if there is any launch scheduled for today
func CheckForTodayLaunches(launches *models.RocketLaunchResponse) {
	today := time.Now().Format("2006-01-02")

	fmt.Println("Checking for launches today...")
	for _, launch := range launches.Result {
		if launch.WinOpen != "" {
			launchDate := launch.WinOpen[:10]
			if launchDate == today {
				fmt.Printf("Today's Launch: %s - %s\n", launch.Name, launch.LaunchDescription)
			}
		}
	}
}

// SetReminderForLaunch sets a reminder for a particular launch
func SetReminderForLaunch(launchID int, launches *models.RocketLaunchResponse) {
	for _, launch := range launches.Result {
		if launch.ID == launchID {
			fmt.Printf("Setting a reminder for launch: %s on %s\n", launch.Name, launch.WinOpen)
			// Logic to set reminder (e.g., using Windows Task Scheduler or a notification system)
			return
		}
	}
	fmt.Println("Launch ID not found.")
}
