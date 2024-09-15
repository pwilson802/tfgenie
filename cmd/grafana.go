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
	Short: "Get Configuration from Grafana",
	Long: `Can get configuration from Grafana including alerts and dashboards

	Ensure the API key variable is set:
		export GRAFANA_API_KEY=XXXXXXX
	
	Example of reading an alert:
		tfgenie grafana --hostname grafana-server.dev --resource alert --alertId aba370f4-ba77-4de6-93f1-4a32158cb2eb
`,
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
	grafanaCmd.PersistentFlags().String("resource", "", "The grafana resource to read - alert, dashboard")
	grafanaCmd.PersistentFlags().String("alertId", "", "The uid of the alert being read")
}
