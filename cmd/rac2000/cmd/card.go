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
var cardCommand = &cobra.Command{
	Use:   "card",
	Short: "Card management command"}
var cardListCommand = &cobra.Command{
	Use:   "list",
	Short: "Get registered cards list from machine",
	Args:  cobra.ExactArgs(0),
	Run:   getCards}
var cardAddCommand = &cobra.Command{
	Use:   "add",
	Short: "Register new card to machine",
	Args:  cobra.ExactArgs(0),
	Run:   addCard}
var cardDelCommand = &cobra.Command{
	Use:   "del",
	Short: "Remove card from machine",
	Args:  cobra.ExactArgs(0),
	Run:   delCard}

// var (
// 	outputFormat string
// )
var (
	cardFacilityCode uint8
	cardId           int
	cardPassword     int
	cardTimezone     int8
	cardStatus       int8
)

func init() {
	cardListCommand.Flags().StringVarP(&outputFormat, "output-format", "f", "json", "Available format: json, table")
	cardCommand.AddCommand(cardListCommand)
	cardAddCommand.Flags().Uint8VarP(&cardFacilityCode, "card-facility-code", "f", 0, "Card facility code")
	cardAddCommand.Flags().IntVarP(&cardId, "card-id", "i", 0, "Card ID")
	cardAddCommand.Flags().IntVarP(&cardPassword, "card-password", "w", 0, "Card password")
	cardAddCommand.Flags().Int8VarP(&cardTimezone, "card-timezone", "t", 0, "Card timezone")
	cardAddCommand.Flags().Int8VarP(&cardStatus, "card-status", "s", 0, "Card status")
	cardAddCommand.MarkFlagRequired("card-facility-code")
	cardAddCommand.MarkFlagRequired("card-id")
	cardCommand.AddCommand(cardAddCommand)
	cardDelCommand.Flags().Uint8VarP(&cardFacilityCode, "card-facility-code", "f", 0, "Card facility code")
	cardDelCommand.Flags().IntVarP(&cardId, "card-id", "i", 0, "Card ID")
	cardCommand.AddCommand(cardDelCommand)

	// cardCommand.AddCommand(&cobra.Command{
	// 	Use:   "clear",
	// 	Short: "Clear all log data from machine",
	// 	Args:  cobra.ExactArgs(0),
	// 	Run:   clearLog})
	cardCommand.Flags().StringVarP(&outputFormat, "output-format", "f", "json", "Available format: json, table")

	RootCmd.AddCommand(cardCommand)
}
func getCards(cmd *cobra.Command, args []string) {
	servAddr := host + ":" + strconv.Itoa(port)
	device := new(rac2000.Rac2000)
	if ok, err := device.Connect(servAddr, nid, time.Duration(time.Second*20)); ok {
		defer device.Close()
		if list, err := device.GetCards(); err == nil {
			switch outputFormat {
			case "json":
				data, _ := json.Marshal(&list)
				fmt.Println(string(data))
			case "table":
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"No", "Facility Code", "Card ID", "Password", "Timezone", "Status"})
				for i, card := range list {
					table.Append([]string{strconv.Itoa(i + 1), strconv.Itoa(int(card.FacilityCode)), strconv.Itoa(int(card.Id)), strconv.Itoa(int(card.Password)), strconv.Itoa(int(card.Timezone)), strconv.Itoa(int(card.Status))})
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

func addCard(cmd *cobra.Command, args []string) {
	servAddr := host + ":" + strconv.Itoa(port)
	device := new(rac2000.Rac2000)
	if ok, err := device.Connect(servAddr, nid, time.Duration(time.Second*20)); ok {
		defer device.Close()
		if _, err := device.AddCard(rac2000.Card{
			FacilityCode: uint8(cardFacilityCode),
			Id:           uint16(cardId),
			Password:     uint32(cardPassword),
			Timezone:     byte(cardTimezone),
			Status:       byte(cardStatus)}); err == nil {
			os.Exit(0)
		}
	} else {
		log.Fatalln(err)
	}
}

func delCard(cmd *cobra.Command, args []string) {
	servAddr := host + ":" + strconv.Itoa(port)
	device := new(rac2000.Rac2000)
	if ok, err := device.Connect(servAddr, nid, time.Duration(time.Second*20)); ok {
		defer device.Close()
		if _, err := device.DelCard(uint8(cardFacilityCode), uint16(cardId)); err == nil {
			os.Exit(0)
		}
	} else {
		log.Fatalln(err)
	}
}
