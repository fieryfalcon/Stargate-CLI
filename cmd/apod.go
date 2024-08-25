package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"stargate/internal/apod"
	"time"

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
		fmt.Println("APOD command invoked")
		today := time.Now().Format("2006-01-02")
		imagePath := filepath.Join("apod_images", today, "apod.jpg")

		if view {
			fmt.Println("View flag detected")
			if _, err := os.Stat(imagePath); os.IsNotExist(err) {
				apodData, err := apod.FetchTodayAPOD()
				if err != nil {
					fmt.Println("Error fetching APOD:", err)
					return
				}
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

		if setWallpaper {
			fmt.Println("Set wallpaper flag detected")
			if _, err := os.Stat(imagePath); os.IsNotExist(err) {
				apodData, err := apod.FetchTodayAPOD()
				if err != nil {
					fmt.Println("Error fetching APOD:", err)
					return
				}
				fmt.Println("Downloading today's APOD:", apodData.Title)
				err = apod.DownloadAPODAndSave(apodData.URL, today)
				if err != nil {
					fmt.Println("Error downloading APOD:", err)
					return
				}
			}
			err := apod.SetWallpaper(imagePath)
			if err != nil {
				fmt.Println("Error setting wallpaper:", err)
				return
			}
			fmt.Println("Wallpaper set successfully.")
		}

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

		if details {
			fmt.Println("Details flag detected")
			apodData, err := apod.FetchTodayAPOD()
			if err != nil {
				fmt.Println("Error fetching APOD details:", err)
				return
			}
			fmt.Println("Title:", apodData.Title)
			fmt.Println("Date:", apodData.Date)
			fmt.Println("URL:", apodData.URL)
			fmt.Println("Explanation:", apodData.Explanation)
		}
	}

	rootCmd.AddCommand(apodCmd)
}
