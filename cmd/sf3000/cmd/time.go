package cmd

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

// represents the time command
var timeCommand = &cobra.Command{
	Use:   "time",
	Short: "Manage time",
	Long:  "Obtain and set current time of the machine"}

func init() {
	timeCommand.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Obtain machine date and time",
		Args:  cobra.ExactArgs(0),
		Run:   getTime})

	timeCommand.AddCommand(&cobra.Command{
		Use:     "set [value]",
		Short:   "Set machine date and time",
		Example: "sf3000 time set \"2006-01-02 15:04:05\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("requires a time argument")
			}
			if _, err := time.Parse("2006-01-02 15:04:05", args[0]); err == nil {
				return nil
			}
			return fmt.Errorf("invalid date time format: %s", args[0])
		},
		Run: setTime})

	RootCmd.AddCommand(timeCommand)
}
func getTime(cmd *cobra.Command, args []string) {
	if conn, device, ok := connect(); ok {
		if dateTime, err := device.GetDateTime(); err == nil {
			fmt.Println(dateTime)
		} else {
			log.Fatalln(err)
		}
		conn.Close()
	}
}

func setTime(cmd *cobra.Command, args []string) {
	fmt.Println("Set time", cmd.Args, args)
}
