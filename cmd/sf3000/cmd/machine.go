package cmd

import (
	"fmt"
	"log"

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
		Use:     "product-code",
		Aliases: []string{"p", "prod", "product"},
		Short:   "Obtain machine product code",
		Args:    cobra.ExactArgs(0),
		Run:     getProductCode})
	machineGetCommand.AddCommand(&cobra.Command{
		Use:     "serial-number",
		Aliases: []string{"s", "sn"},
		Short:   "Obtain machine serial number",
		Args:    cobra.ExactArgs(0),
		Run:     getSerialNumber})

	RootCmd.AddCommand(machineCommand)
}

// Obtain product code
func getProductCode(cmd *cobra.Command, args []string) {
	if conn, device, ok := connect(); ok {
		if productCode, err := device.GetProductCode(); err == nil {
			fmt.Println(productCode)
		} else {
			log.Fatalln(err)
		}
		conn.Close()
	}
}

// Obtain serial number
func getSerialNumber(cmd *cobra.Command, args []string) {
	if conn, device, ok := connect(); ok {
		if serialNumber, err := device.GetSerialNumber(); err == nil {
			fmt.Println(serialNumber)
		} else {
			log.Fatalln(err)
		}
		conn.Close()
	}
}
