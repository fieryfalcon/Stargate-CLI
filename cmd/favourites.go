package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"stargate/internal/apod"
	"github.com/spf13/cobra"
)

var favsSetWallpaper bool

var favouritesCmd = &cobra.Command{
	Use:   "favourites",
	Short: "Manage APOD favorites",
	Long:  "Manage APOD favorites, including setting them as a wallpaper slideshow.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use one of the flags: -sw")
	},
}

func init() {
	favouritesCmd.Flags().BoolVarP(&favsSetWallpaper, "sw", "w", false, "Set the favorites folder as wallpaper slideshow")

	favouritesCmd.Run = func(cmd *cobra.Command, args []string) {
		if favsSetWallpaper {
			favoritesPath := filepath.Join("apod_images", "favorites")
			if _, err := os.Stat(favoritesPath); os.IsNotExist(err) {
				fmt.Println("No favorites found. Please add some APOD images to favorites first.")
				return
			}
			err := apod.SetSlideshowFromFavorites(favoritesPath)
			if err != nil {
				fmt.Println("Error setting slideshow from favorites:", err)
				return
			}
			fmt.Println("Slideshow set successfully.")
		}
	}

	rootCmd.AddCommand(favouritesCmd)
}
