package cmd

import (
	"fmt"
	go_aspace "github.com/nyudlts/go-aspace"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var csvfile string
var checkCmd = &cobra.Command{

	Use:   "check",
	Short: "Check an archviesspace instance for malformed eadids",
	Run: func(cmd *cobra.Command, args []string) {

		err := GenerateCSV()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	checkCmd.PersistentFlags().StringVar(&csvfile, "csv", "", "default is output.csv")
	rootCmd.AddCommand(checkCmd)
}


func initConfig() {
	if csvfile == "" {
		csvfile = "output.csv"
	}
}

func GenerateCSV() error {
	ASpaceClient, err := go_aspace.NewClient(100)
	if err != nil {
		return err
	}

	csvFile, err := os.Create(csvfile)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	repositories := []int{2,3,6}

	for i := range repositories {
		repositoryId := repositories[i]
		resources, err := ASpaceClient.GetResourceIDsByRepository(repositoryId)
		if err != nil {
			return err
		}

		for j := range resources {
			resourceId := resources[j]
			resource, err := ASpaceClient.GetResourceByID(repositoryId, resourceId)
			if err != nil {
				return err
			}
			eadid := resource.EADID
			t := resource.Title
			t = strings.Replace(t, "\"", "", -1)
			title := strings.Replace(t, ",", "", -1)
			id0 := resource.ID0
			id1 := resource.ID1
			id2 := resource.ID2
			id3 := resource.ID3

			target := ""
			if id0 != "" {
				target = id0
			}
			if id1 != "" {
				target = target + "_" + id1
			}
			if id2 != "" {
				target = target + "_" + id2
			}
			if id3 != "" {
				target = target + "_" + id3
			}
			target = strings.ToLower(target)

			fmt.Print(title)
			if eadid != target {
				fmt.Println(": KO")
				writeString, err := csvFile.WriteString(fmt.Sprintf("%d,%d,%s,%s,%s\n", repositoryId, resourceId, title, eadid, target))
				if err != nil {
					log.Println(writeString)
					log.Fatal(err)
				}
			} else {
				fmt.Println(": OK")
			}
		}
	}

	return nil
}
