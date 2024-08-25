package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stargate",
	Short: "Stargate CLI: Your gateway to space updates",
	Long:  "Stargate CLI provides space updates including NASA's Astronomy Picture of the Day (APOD) and more.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
