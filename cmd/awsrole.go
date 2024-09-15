/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// awsroleCmd represents the awsrole command
var awsroleCmd = &cobra.Command{
	Use:   "awsrole",
	Short: "Can create terraform code froma refernce role or user via AWS Access Analyzer",
	Long: `Can create terraform code froma refernce role or user via AWS Access Analyzer`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("awsrole called")
	},
}

func init() {
	rootCmd.AddCommand(awsroleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// awsroleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// awsroleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
