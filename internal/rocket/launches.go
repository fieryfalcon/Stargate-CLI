package rocket

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"stargate/pkg/models"
	"stargate/pkg/utils"
	"strings"
	"time"
)

const rocketLaunchAPI = "https://fdo.rocketlaunch.live/json/launches/next/5"

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
func CheckForTodayLaunches(launches *models.RocketLaunchResponse) []models.RocketLaunch {
	today := time.Now().Format("2006-01-02")

	fmt.Println("Checking for launches today...")
	var todayLaunches []models.RocketLaunch
	for _, launch := range launches.Result {
		if launch.WinOpen != "" {
			launchDate := launch.WinOpen[:10]
			if launchDate == today {
				todayLaunches = append(todayLaunches, launch)
			}
		}
	}
	return todayLaunches
}

// CheckForTomorrowLaunches checks if there is any launch scheduled for tomorrow
func CheckForTomorrowLaunches(launches *models.RocketLaunchResponse) []models.RocketLaunch {
	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")

	fmt.Println("Checking for launches tomorrow...")
	var tomorrowLaunches []models.RocketLaunch
	for _, launch := range launches.Result {
		if launch.WinOpen != "" {
			launchDate := launch.WinOpen[:10]
			if launchDate == tomorrow {
				tomorrowLaunches = append(tomorrowLaunches, launch)
			}
		}
	}
	return tomorrowLaunches
}

func processDateTime(input string) (string, string) {
	// Step 1: Remove the last character
	if len(input) > 0 {
		input = input[:len(input)-1]
	}

	// Step 2: Split the string at the 'T' character if it exists
	parts := strings.Split(input, "T")
	date := parts[0]
	timePart := ""
	if len(parts) > 1 {
		timePart = parts[1]
	}

	// Convert date from "yyyy-mm-dd" to "dd/mm/yyyy"
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return "", ""
	}
	formattedDate := parsedDate.Format("02/01/2006")

	return formattedDate, timePart
}

// SetReminderForLaunch sets a reminder for a particular launch
func SetReminderForLaunch(launchID int, launches *models.RocketLaunchResponse) {
	for _, launch := range launches.Result {
		if launch.ID == launchID {
			fmt.Printf("Setting a reminder for launch: %s on %s\n", launch.Name, launch.WinOpen)

			// Use the processDateTime function to get the formatted date and time
			startDate, startTime := processDateTime(launch.WinOpen)

			// Set the task name and message
			taskName := fmt.Sprintf("RocketLaunchReminder_%d", launchID)
			message := fmt.Sprintf("Reminder: Rocket launch '%s' is scheduled for %s", launch.Name, launch.WinOpen)

			// Schedule the task using schtasks
			err := scheduleTask(taskName, message, startTime, startDate)
			if err != nil {
				fmt.Println("Error setting reminder:", err)
				return
			}

			fmt.Println("Reminder set successfully.")
			return
		}
	}
	fmt.Println("Launch ID not found.")
}

// scheduleTask creates a scheduled task in Windows to show a reminder message at the specified time
func scheduleTask(taskName, message, startTime, startDate string) error {
	// Build the schtasks command to create the task
	cmd := exec.Command("schtasks", "/Create", "/TN", taskName, "/TR", fmt.Sprintf("msg * /time:60 \"%s\"", message), "/SC", "ONCE", "/ST", startTime, "/SD", startDate)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create scheduled task: %v, output: %s", err, string(output))
	}

	return nil
}
