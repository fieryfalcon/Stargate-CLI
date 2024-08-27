package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"stargate/internal/apod"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var apodCmd = &cobra.Command{
	Use:   "apod",
	Short: "Astronomy Picture of the Day",
	Long:  "View NASA's Astronomy Picture of the Day (APOD), set it as wallpaper, add to favorites, and more.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use one of the flags: -view, -sw, -af, -details")
	},
}

var view bool
var setWallpaper bool
var addToFavs bool
var details bool

func init() {
	apodCmd.Flags().BoolVarP(&view, "view", "v", false, "View the APOD")
	apodCmd.Flags().BoolVarP(&setWallpaper, "sw", "w", false, "Set the APOD as wallpaper")
	apodCmd.Flags().BoolVarP(&addToFavs, "af", "f", false, "Add the APOD to favorites")
	apodCmd.Flags().BoolVarP(&details, "details", "d", false, "Get details about the APOD")

	apodCmd.Run = func(cmd *cobra.Command, args []string) {

		// Fetch APOD data and determine the file type
		apodData, err := apod.FetchTodayAPOD()
		if err != nil {
			fmt.Println("Error fetching APOD:", err)
			return
		}
		fileType := apod.DetermineFileType(apodData.URL)

		// Handle based on the file type
		switch fileType {
		case "image":
			today := time.Now().Format("2006-01-02")
			imagePath := filepath.Join("apod_images", today, "apod.jpg")

			// Handle view command
			if view {
				fmt.Println("View flag detected")
				if _, err := os.Stat(imagePath); os.IsNotExist(err) {
					fmt.Println("Downloading and opening today's APOD:", apodData.Title)
					err = apod.DownloadAPODAndSave(apodData.URL, today)
					if err != nil {
						fmt.Println("Error downloading APOD:", err)
						return
					}
				} else {
					fmt.Println("Image already exists, opening it...")
				}
				err := apod.OpenImage(imagePath)
				if err != nil {
					fmt.Println("Error opening APOD:", err)
					return
				}
			}

			// Handle setWallpaper command
			if setWallpaper {
				fmt.Println("Set wallpaper flag detected")
				if _, err := os.Stat(imagePath); os.IsNotExist(err) {
					fmt.Println("Downloading today's APOD:", apodData.Title)
					err = apod.DownloadAPODAndSave(apodData.URL, today)
					if err != nil {
						fmt.Println("Error downloading APOD:", err)
						return
					}
				}
				err = apod.SetWallpaper(imagePath)
				if err != nil {
					fmt.Println("Error setting wallpaper:", err)
					return
				}
				fmt.Println("Wallpaper set successfully.")
			}

			// Handle addToFavs command
			if addToFavs {
				fmt.Println("Add to favorites flag detected")
				favoritesPath := filepath.Join("apod_images", "favorites", today+"-apod.jpg")
				err := apod.CopyToFavorites(imagePath, favoritesPath)
				if err != nil {
					fmt.Println("Error adding to favorites:", err)
					return
				}
				fmt.Println("APOD added to favorites successfully.")
			}

		case "video":
			if view {
				fmt.Printf("Today's APOD is a video: %s\n", apodData.URL)
			}
			if setWallpaper {
				fmt.Printf("Today's APOD is a video: %s\nCannot set video as wallpaper.\n", apodData.URL)
			}
			if addToFavs {
				fmt.Printf("Today's APOD is a video: %s\nCannot add video to favorites.\n", apodData.URL)
			}

		default:
			fmt.Printf("Unknown file type found in the APOD URL: %s\n", apodData.URL)
		}

		if details {
			fmt.Println("Details flag detected")

			// Use colors and formatting for better readability
			header := color.New(color.FgCyan, color.Bold).PrintfFunc()
			subheader := color.New(color.FgGreen, color.Bold).PrintfFunc()
			text := color.New(color.FgWhite).PrintfFunc()
			url := color.New(color.FgBlue, color.Underline).PrintfFunc()

			// Display the APOD details
			header("\nAstronomy Picture of the Day (APOD) Details\n")
			header("===========================================\n")

			subheader("\nTitle: ")
			text("%s\n", apodData.Title)

			subheader("\nDate: ")
			text("%s\n", apodData.Date)

			subheader("\nExplanation: \n")
			text("%s\n", apodData.Explanation)

			subheader("\nURL: ")
			url("%s\n", apodData.URL)

			// Provide a visual separator for clarity
			fmt.Println("\n===========================================\n")
		}
	}

	rootCmd.AddCommand(apodCmd)
}
