package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	host     string
	port     int
	nid      uint16
	password uint16
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "sf3500",
	Short:   "Keico SF3500 command line interface",
	Long:    "Application to manage Keico SF3500 attendance machine.",
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
	//RootCmd.PersistentFlags().StringVar(&host, "host", "", "config file (default is $HOME/.cobra-example.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.PersistentFlags().StringVar(&host, "host", "", "Specify the host name or IP address of the remote machine to connect to")
	RootCmd.MarkPersistentFlagRequired("host")
	RootCmd.PersistentFlags().IntVar(&port, "port", 5005, "Specify the port number of the remote machine to connect to")
	RootCmd.PersistentFlags().Uint16Var(&nid, "nid", 1, "Specify the machine number")
	RootCmd.PersistentFlags().Uint16Var(&password, "password", 0, "Specify the password to connect to remote machine")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

}

// // Open connection and send handshake command to machine
// func connect() (*net.TCPConn, *sf3000.Sf3000, bool) {
// 	servAddr := host + ":" + strconv.Itoa(port)
// 	dialer := net.Dialer{Timeout: time.Duration(time.Second * 20)}
// 	conn, err := dialer.Dial("tcp", servAddr)
// 	//conn, err := net.DialTCP("tcp", nil, tcpAddr)
// 	if err != nil {
// 		println("Dial failed:", err.Error())
// 		os.Exit(1)
// 	}
// 	if tcpConn, ok := conn.(*net.TCPConn); ok {
// 		device := new(sf3000.Sf3000)
// 		connected, err := device.Connect(tcpConn, uint16(nid), uint16(password))
// 		if err != nil {
// 			println(err.Error())
// 			os.Exit(1)
// 		}
// 		return tcpConn, device, connected
// 	}
// 	return nil, nil, false
// }
