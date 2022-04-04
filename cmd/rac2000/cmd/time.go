package cmd

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/masykur/keico/pkg/rac2000"
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
		Example: "To set machine date and time to specific value:\n\tsf3000 time set \"2006-01-02 15:04:05\"\nTo set machine date and time follow the client PC:\n\tsf3000 time set",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return errors.New("too many arguments")
			}
			if len(args) == 0 {
				return nil
			}
			if _, err := time.ParseInLocation("2006-01-02 15:04:05", args[0], time.Local); err == nil {
				return nil
			}
			return fmt.Errorf("invalid date time format: %s", args[0])
		},
		Run: setTime})

	RootCmd.AddCommand(timeCommand)
}
func getTime(cmd *cobra.Command, args []string) {
	servAddr := host + ":" + strconv.Itoa(port)
	device := new(rac2000.Rac2000)
	if ok, err := device.Connect(servAddr, nid, time.Duration(time.Second*20)); ok {
		defer device.Close()
		if dateTime, err := device.GetDateTime(); err == nil {
			fmt.Println(dateTime)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Fatalln(err)
	}
}

func setTime(cmd *cobra.Command, args []string) {
	var t time.Time
	if len(args) == 0 {
		t = time.Now()
	} else {
		t, _ = time.ParseInLocation("2006-01-02 15:04:05", args[0], time.Local)
	}
	servAddr := host + ":" + strconv.Itoa(port)
	device := new(rac2000.Rac2000)
	if ok, err := device.Connect(servAddr, nid, time.Duration(time.Second*20)); ok {
		defer device.Close()
		if ok, err := device.SetDateTime(t); ok {
			return
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Fatalln(err)
	}
}
