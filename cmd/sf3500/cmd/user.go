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
var userCommand = &cobra.Command{
	Use:   "user",
	Short: "Manage user",
	Long:  "View, enroll, and remove user from machine"}

var userGetCommand = &cobra.Command{
	Use:     "get [id,..]",
	Short:   "Obtain user information",
	Example: "sf3500 user get 12345678,12345688",
	Run:     getUser}

var userSetCommand = &cobra.Command{
	Use:     "set",
	Short:   "Enroll user to machine",
	Example: `sf3000 user set --host 192.168.0.1 --nid 123 -d '{\"Id\":12345678,\"CardFacilityCode\":186,\"CardId\":45123,\"Fingerprint1\":\"\",\"Fingerprint2\":\"\"}'`,
	Args:    cobra.ExactArgs(0),
	Run:     getUsers}

var userCountCommand = &cobra.Command{
	Use:   "count",
	Short: "Obtain number of users registered in the machine",
	Args:  cobra.ExactArgs(0),
	Run:   getUsers}

var userListCommand = &cobra.Command{
	Use:   "list",
	Short: "Retrieve list of users registered in the machine",
	Args:  cobra.ExactArgs(0),
	Run:   getUsers}

var (
	outputFile   string
	outputFormat string
	data         string
	inputFile    string
)

func init() {
	userGetCommand.Flags().StringVarP(&outputFile, "output-file", "o", "", "Write output to file")
	userGetCommand.Flags().StringVarP(&outputFormat, "output-format", "f", "json", "Available format: json, table")
	userListCommand.Flags().StringVarP(&outputFormat, "output-format", "f", "json", "Available format: json, table")
	userSetCommand.Flags().StringVarP(&inputFile, "input-file", "i", "", "Read input from json file")
	userSetCommand.Flags().StringVarP(&data, "data", "d", "", "Input data in json format")
	userCommand.AddCommand(userGetCommand)
	userCommand.AddCommand(userCountCommand)
	userCommand.AddCommand(userListCommand)
	userCommand.AddCommand(userSetCommand)
	RootCmd.AddCommand(userCommand)
}

// Obtain user id list
func getUsers(cmd *cobra.Command, args []string) {
	servAddr := host + ":" + strconv.Itoa(port)
	device := new(sf3500.Sf3500)
	if ok, err := device.Connect(servAddr, time.Duration(time.Second*20)); ok {
		defer device.Close()
		users := make([]models.User, 0)
		packageId := 0
		for {
			var userList []models.User
			if packageId, userList, err = device.GetUserList(packageId); err == nil {
				users = append(users, userList...)
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
			data, _ := json.Marshal(&users)
			fmt.Println(string(data))
		case "table":
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"No", "User ID", "Name", "Card"})
			for i, user := range users {
				table.Append([]string{strconv.Itoa(i + 1), user.UserID, user.Name, user.Card})
			}
			table.Render()
		default:
			log.Fatalln("Invalid output format")
		}
	} else {
		log.Fatalln(err)
	}
}

// Obtain user id list
func getUser(cmd *cobra.Command, args []string) {
	servAddr := host + ":" + strconv.Itoa(port)
	device := new(sf3500.Sf3500)
	if ok, err := device.Connect(servAddr, time.Duration(time.Second*20)); ok {
		defer device.Close()
		users := make([]models.User, 0)
		packageId := 0
		for {
			var userList []models.User
			if packageId, userList, err = device.GetUserInfo(packageId, args...); err == nil {
				users = append(users, userList...)
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
			data, _ := json.Marshal(&users)
			fmt.Println(string(data))
		case "table":
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"No", "User ID", "Name", "Card"})
			for i, user := range users {
				table.Append([]string{strconv.Itoa(i + 1), user.UserID, user.Name, user.Card})
			}
			table.Render()
		default:
			log.Fatalln("Invalid output format")
		}
	} else {
		log.Fatalln(err)
	}
}
