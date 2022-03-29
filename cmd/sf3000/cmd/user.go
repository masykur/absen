package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/masykur/keico/pkg/entity"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// represents the time command
var userCommand = &cobra.Command{
	Use:   "user",
	Short: "Manage user",
	Long:  "View, enroll, and remove user from machine"}

var userGetCommand = &cobra.Command{
	Use:     "get [id]",
	Short:   "Obtain user information",
	Example: "sf3000 user get 12345678",
	Args:    cobra.ExactArgs(1),
	Run:     getUser}

var userSetCommand = &cobra.Command{
	Use:     "set",
	Short:   "Enroll user to machine",
	Example: `sf3000 user set --host 192.168.0.1 --nid 123 -d '{\"Id\":12345678,\"CardId\":45123,\"CardFacilityCode\":186,\"Fingerprint1\":\"\",\"Fingerprint2\":\"\"}'`,
	Args:    cobra.ExactArgs(0),
	Run:     setUser}

var userCountCommand = &cobra.Command{
	Use:   "count",
	Short: "Obtain number of users registered in the machine",
	Args:  cobra.ExactArgs(0),
	Run:   getUserCount}

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

// var outputJson bool
// var outputTable bool

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

// Obtain number of users registered in the machine
func getUserCount(cmd *cobra.Command, args []string) {
	if conn, device, ok := connect(); ok {
		defer conn.Close()
		if count, err := device.GetUserCount(); err == nil {
			fmt.Println(count)
		} else {
			log.Fatalln(err)
		}
	}
}

// Retrieve list of users registered in the machine
func getUsers(cmd *cobra.Command, args []string) {
	if conn, device, ok := connect(); ok {
		defer conn.Close()
		users, err := device.GetUsers()
		if err == nil {
			switch outputFormat {
			case "json":
				data, _ := json.Marshal(&users)
				fmt.Println(string(data))
			case "table":
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"No", "User ID", "Privilage", "Sensor", "Card ID"})
				for i, user := range users {
					table.Append([]string{strconv.Itoa(i + 1), strconv.Itoa(user.Id), strconv.Itoa(int(user.Level)), strconv.Itoa(int(user.Sensor)), strconv.Itoa(int(user.CardId))})
				}
				table.Render()
			default:
				log.Fatalln("Invalid output format")
			}
		} else {
			log.Fatalln(err)
		}
	}
}

// Obtain number of users registered in the machine
func getUser(cmd *cobra.Command, args []string) {
	if userId, err := strconv.Atoi(args[0]); err == nil {
		if conn, device, ok := connect(); ok {
			defer conn.Close()
			if user, err := device.GetEnrollData(int(userId)); err == nil {
				switch outputFormat {
				case "json":
					data, _ := json.Marshal(&user)
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
					table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
					table.SetAlignment(tablewriter.ALIGN_LEFT) // Set Alignment
					table.SetHeader([]string{"User ID", "Card ID", "CardFacilityCode", "Fingerprint1", "Fingerprint2"})
					table.Append([]string{strconv.Itoa(user.Id), strconv.Itoa(int(user.CardId)), strconv.Itoa(int(user.CardFacilityCode)), hex.EncodeToString(user.Fingerprint1), hex.EncodeToString(user.Fingerprint2)})
					table.Render()
				default:
					log.Fatalln("Invalid output format")
				}
			} else {
				log.Fatalln(err)
			}
		}
	}
}

// Obtain number of users registered in the machine
func setUser(cmd *cobra.Command, args []string) {
	if conn, device, ok := connect(); ok {
		defer conn.Close()
		if data != "" {
			var user entity.User
			if err := json.Unmarshal([]byte(data), &user); err == nil {
				ok, err := device.SetEnrollData(user)
				if err != nil {
					log.Fatal(err)
					os.Exit(1)
				}
				if !ok {
					cmd.PrintErrln("Failed to enroll user")
				}
			} else {
				cmd.PrintErrln("Invalid json format")
				cmd.Help()
				os.Exit(2)
			}
		} else if inputFile != "" {
			var f *os.File
			jsonText, err := os.ReadFile(inputFile)
			if err != nil {
				log.Fatal(err)
				os.Exit(2)
			}
			defer f.Close()
			var user entity.User
			json.Unmarshal(jsonText, &user)
			ok, err := device.SetEnrollData(user)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			if !ok {
				cmd.PrintErrln("Failed to enroll user")
			}
		} else {
			cmd.PrintErrln("No input data available")
			cmd.Help()
		}
	}
}
