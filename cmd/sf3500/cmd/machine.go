package cmd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/masykur/absen/pkg/sf3500"
	"github.com/spf13/cobra"
)

// represents the time command
var machineCommand = &cobra.Command{
	Use:   "machine",
	Short: "Manage attandance machine",
	Long:  "Obtain product code and serial number of the machine"}

var machineGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Obtain attandance machine information",
	Long:  "Obtain product code and serial number of the machine",
}

func init() {
	machineCommand.AddCommand(machineGetCommand)
	machineGetCommand.AddCommand(&cobra.Command{
		Use:   "info",
		Short: "Obtain machine info",
		Args:  cobra.ExactArgs(0),
		Run:   getDeviceInfo})
	// machineGetCommand.AddCommand(&cobra.Command{
	// 	Use:     "serial-number",
	// 	Aliases: []string{"s", "sn"},
	// 	Short:   "Obtain machine serial number",
	// 	Args:    cobra.ExactArgs(0),
	// 	Run:     getSerialNumber})

	RootCmd.AddCommand(machineCommand)
}

// Obtain product code
func getDeviceInfo(cmd *cobra.Command, args []string) {
	servAddr := host + ":" + strconv.Itoa(port)
	device := new(sf3500.Sf3500)
	if ok, err := device.Connect(servAddr, time.Duration(time.Second*20)); ok {
		defer device.Close()
		if deviceInfo, err := device.GetDeviceInfo(); err == nil {
			fmt.Printf("%-25s: %s\n", "ID", deviceInfo.DeviceID)
			fmt.Printf("%-25s: %s\n", "Name", deviceInfo.Name)
			fmt.Printf("%-25s: %s\n", "Firmware", deviceInfo.Firmware)
			fmt.Printf("%-25s: %s\n", "Fingerprint Version", deviceInfo.FingerprintVersion)
			fmt.Printf("%-25s: %s\n", "Face Version", deviceInfo.FaceVersion)
			fmt.Printf("%-25s: %s\n", "Palmprint Version", deviceInfo.PalmprintVersion)
			fmt.Printf("%-25s: %d\n", "Maximum Buffer Length", deviceInfo.MaximumBufferLength)
			fmt.Printf("%-25s: %d\n", "User Limit", deviceInfo.UserLimit)
			fmt.Printf("%-25s: %d\n", "Fingerprint Limit", deviceInfo.FingerprintLimit)
			fmt.Printf("%-25s: %d\n", "Face Limit", deviceInfo.FaceLimit)
			fmt.Printf("%-25s: %d\n", "Password Limit", deviceInfo.PasswordLimit)
			fmt.Printf("%-25s: %d\n", "Card Limit", deviceInfo.CardLimit)
			fmt.Printf("%-25s: %d\n", "Log Limit", deviceInfo.LogLimit)
			fmt.Printf("%-25s: %d\n", "User Count", deviceInfo.UserCount)
			fmt.Printf("%-25s: %d\n", "Manager Count", deviceInfo.ManagerCount)
			fmt.Printf("%-25s: %d\n", "Fingerprint Count", deviceInfo.FingerprintCount)
			fmt.Printf("%-25s: %d\n", "Face Count", deviceInfo.FaceCount)
			fmt.Printf("%-25s: %d\n", "Password Count", deviceInfo.PasswordCount)
			fmt.Printf("%-25s: %d\n", "Card Count", deviceInfo.CardCount)
			fmt.Printf("%-25s: %d\n", "Log Count", deviceInfo.LogCount)
			fmt.Printf("%-25s: %d\n", "All Logs Count", deviceInfo.AllLogCount)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Fatalln(err)
	}
}
