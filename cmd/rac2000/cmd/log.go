package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/masykur/keico/pkg/rac2000"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// represents the time command
var logCommand = &cobra.Command{
	Use:   "log",
	Short: "Log data command"}
var logFetchCommand = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch log data from machine",
	Args:  cobra.ExactArgs(0),
	Run:   fetchLog}

var (
	outputFormat string
)

func init() {
	logFetchCommand.Flags().StringVarP(&outputFormat, "output-format", "f", "json", "Available format: json, table")
	logCommand.AddCommand(logFetchCommand)

	logCommand.AddCommand(&cobra.Command{
		Use:   "clear",
		Short: "Clear all log data from machine",
		Args:  cobra.ExactArgs(0),
		Run:   clearLog})
	logCommand.Flags().StringVarP(&outputFormat, "output-format", "f", "json", "Available format: json, table")

	RootCmd.AddCommand(logCommand)
}
func fetchLog(cmd *cobra.Command, args []string) {
	servAddr := host + ":" + strconv.Itoa(port)
	device := new(rac2000.Rac2000)
	if ok, err := device.Connect(servAddr, nid, time.Duration(time.Second*20)); ok {
		defer device.Close()
		if list, err := device.FetchLog(0); err == nil {
			switch outputFormat {
			case "json":
				data, _ := json.Marshal(&list)
				fmt.Println(string(data))
			case "table":
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"No", "Event", "Date Time", "Card Facility Code", "Card ID", "Device ID", "Reader ID"})
				for i, logData := range list {
					table.Append([]string{strconv.Itoa(i + 1), fmt.Sprintf("%X", logData.Event), logData.DateTime.Format("2006-01-02 15:04:05"), strconv.Itoa(int(logData.CardFacilityCode)), strconv.Itoa(int(logData.CardId)), strconv.Itoa(logData.DeviceId), strconv.Itoa(logData.ReaderId)})
				}
				table.Render()
			default:
				log.Fatalln("Invalid output format")
			}

		} else {
			log.Fatalln(err)
		}
	} else {
		log.Fatalln(err)
	}
}

func clearLog(cmd *cobra.Command, args []string) {
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
