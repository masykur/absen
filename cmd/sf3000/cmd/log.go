package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/masykur/absen/pkg/sf3000"
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

func init() {
	logFetchCommand.Flags().StringVarP(&outputFile, "output-file", "o", "", "Write output to file")
	logFetchCommand.Flags().StringVarP(&outputFormat, "output-format", "f", "json", "Available format: json, table")
	logCommand.AddCommand(logFetchCommand)
	RootCmd.AddCommand(logCommand)
}
func fetchLog(cmd *cobra.Command, args []string) {
	servAddr := host + ":" + strconv.Itoa(port)
	device := new(sf3000.Sf3000)
	if ok, err := device.Connect(servAddr, nid, password, time.Duration(time.Second*20)); ok {
		if _, list, err := device.FetchAllLogs(); err == nil {
			device.Close()
			switch outputFormat {
			case "json":
				data, _ := json.Marshal(&list)
				if outputFile != "" {
					var f *os.File
					f, err := os.Create(outputFile)
					if err != nil {
						log.Fatal(err)
						os.Exit(2)
					}
					defer f.Close()
					_, err = f.Write(data)
					if err != nil {
						log.Fatal(err)
						os.Exit(2)
					}
				} else {
					fmt.Println(string(data))
				}
			case "table":
				var table *tablewriter.Table
				if outputFile != "" {
					var f *os.File
					f, err := os.Create(outputFile)
					if err != nil {
						log.Fatal(err)
						os.Exit(2)
					}
					defer f.Close()
					table = tablewriter.NewWriter(f)
				} else {
					table = tablewriter.NewWriter(os.Stdout)
				}
				table.SetHeader([]string{"No", "UserID", "Event", "Date Time", "UserType", "SensorType", "Mode", "FunctionKey", "FunctionNumber"})
				for i, logData := range list {
					table.Append([]string{strconv.Itoa(i + 1), strconv.Itoa(int(logData.UserID)), strconv.Itoa(int(logData.Event)), logData.DateTime.Format("2006-01-02 15:04:05"), strconv.Itoa(int(logData.UserType)), logData.SensorType.String(), logData.Mode.String(), logData.FunctionKey.String(), strconv.Itoa(int(logData.FunctionNumber))})
				}
				table.Render()
			default:
				log.Fatalln("Invalid output format")
			}

		} else {
			device.Close()
			log.Fatalln(err)
		}
	} else {
		device.Close()
		log.Fatalln(err)
	}
}
