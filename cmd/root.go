/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Yelsnik/caching-server/proxy"
	"github.com/spf13/cobra"
)

var (
	Port       string
	Origin     string
	ClearCache bool
	Proxy      *proxy.Proxy
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "caching-server",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		var port int
		var err error

		fmt.Println(cmd.Long)

		if Port != "" {
			port, err = strconv.Atoi(Port)
			if err != nil {
				fmt.Printf("Invalid port number: %v\n", err)
				return
			}
		}

		Proxy = proxy.NewProxyServer(Origin, ClearCache)

		server := &proxy.Server{
			Origin:     Origin,
			Port:       port,
			ClearCache: ClearCache,
			Proxy:      Proxy,
		}

		server.StartServer()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&Origin, "origin", "o", "", "origin url to be inputed by the user")
	rootCmd.PersistentFlags().BoolVarP(&ClearCache, "clear-cache", "c", false, "used to clear cache of the server")
	rootCmd.PersistentFlags().StringVarP(&Port, "port", "p", "", "port to be used to start the server")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
