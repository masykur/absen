package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/masykur/keico/pkg/sf3500"
	"github.com/masykur/keico/pkg/sf3500/models"
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
	device := new(sf3500.Sf3500)
	if ok, err := device.Connect(servAddr, time.Duration(time.Second*20)); ok {
		defer device.Close()
		logs := make([]models.LogData, 0)
		packageId := 0
		for {
			var logList []models.LogData
			if packageId, logList, err = device.GetLog(packageId, 0, time.Now().AddDate(-1, 0, 0), time.Now(), 0); err == nil {
				logs = append(logs, logList...)
				if packageId == 0 {
					break
				}
			} else {
				log.Fatalln(err)
				break
			}
		}
		switch outputFormat {
		case "json":
			data, _ := json.Marshal(&logs)
			fmt.Println(string(data))
		case "table":
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"No", "User ID", "Time", "VerifyMode"})
			for i, logData := range logs {
				table.Append([]string{strconv.Itoa(i + 1), logData.UserID, logData.Time.Time().Format("2006-01-02 15:04:05"), logData.VerifyMode})
			}
			table.Render()
		default:
			log.Fatalln("Invalid output format")
		}
	} else {
		log.Fatalln(err)
	}
}
