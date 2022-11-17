package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	host string
	port uint16
	dsn  string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "sf3500",
	Short:   "Keico SF3500 http server",
	Long:    "Keico SF3500 is http server for receiving log data from machine with push function.",
	Example: "sf3500 --host 192.168.0.1 --port 9000 --dsn \"sqlserver://faceid:faceid-p4$$w0RD@192.168.0.2?database=FaceID\"",
	Version: "0.6.0"}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	v := viper.New()
	v.SetEnvPrefix("sf3500")
	v.SetDefault("port", 9009)
	v.AutomaticEnv()
	RootCmd.PersistentFlags().StringVarP(&host, "host", "o", v.GetString("host"), "Specify the host name or IP address that bind to")
	RootCmd.PersistentFlags().Uint16VarP(&port, "port", "p", uint16(v.GetUint("port")), "Specify the port number that listening on")
	RootCmd.PersistentFlags().StringVarP(&dsn, "dsn", "d", v.GetString("dsn"), "Data source name")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

}
