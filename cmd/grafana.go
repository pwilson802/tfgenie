/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/pwilson802/tfgenie/pkg/grafana"
	"github.com/spf13/cobra"
)

// grafanaCmd represents the grafana command
var grafanaCmd = &cobra.Command{
	Use:   "grafana",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("grafana called")
		hostname, _ := cmd.Flags().GetString("hostname")
		resource, _ := cmd.Flags().GetString("resource")
		alertId, _ := cmd.Flags().GetString("alertId")
		if resource != "alert" {
			fmt.Println("only alert resources are currently supported")
			return
		}
		if resource == "alert" {
			if hostname == "" {
				log.Fatal("Hostname is required")
			}
			if alertId == "" {
				log.Fatal("alertId is required")
			}
			gfClient := grafana.CreateNewGrafanaClient(hostname)
			_, err := gfClient.GetAlert(alertId)
			if err != nil {
				log.Fatal("Error getting Alert", err)
			}
			alert, err := gfClient.GetAlert(alertId)
			if err != nil {
				log.Fatal("Error getting Alert", err)
			}
			err = gfClient.ExportAlertToTerraform(alert, "tf-output.txt", "", "")
			if err != nil {
				log.Fatal("Error getting Alert", err)
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(grafanaCmd)
	grafanaCmd.PersistentFlags().String("hostname", "", "The grafana hostname")
	grafanaCmd.PersistentFlags().String("resource", "", "The grafana resource to read - Alert, Dashboard")
	grafanaCmd.PersistentFlags().String("alertId", "", "The uid of the alert being read")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// grafanaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// grafanaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
