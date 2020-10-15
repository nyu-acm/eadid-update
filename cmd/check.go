package cmd

import (
	"fmt"
	go_aspace "github.com/nyudlts/go-aspace/lib"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var aspace = go_aspace.Client
var tsvfile string
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
	checkCmd.PersistentFlags().StringVar(&tsvfile, "tsv", "", "default is output.csv")
	rootCmd.AddCommand(checkCmd)
}

func initConfig() {
	if tsvfile == "" {
		tsvfile = "output.tsv"
	}
}

func GenerateCSV() error {

	tsvFile, err := os.Create(tsvfile)
	if err != nil {
		return err
	}
	defer tsvFile.Close()

	repositories := []int{2, 3, 6}

	for i := range repositories {
		repositoryId := repositories[i]
		resources, err := aspace.GetResourceIDsByRepository(repositoryId)
		if err != nil {
			return err
		}

		for j := range resources {
			resourceId := resources[j]
			resource, err := aspace.GetResourceByID(repositoryId, resourceId)
			if err != nil {
				return err
			}
			eadid := resource.EADID
			eadlocat := resource.EADLocation
			title := resource.Title
			title = strings.Replace(title, "\"", "", -1)
			title = strings.Replace(title, "\n", "", -1)
			title = strings.Replace(title, ",", "", -1)
			id0 := strings.Trim(resource.ID0, " ")
			id1 := strings.Trim(resource.ID1, " ")
			id2 := strings.Trim(resource.ID2, " ")
			id3 := strings.Trim(resource.ID3, " ")

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
			target = strings.ReplaceAll(target, " ", "")
			target = strings.ReplaceAll(target, ",", "")
			target = strings.ReplaceAll(target, ".", "")

			fmt.Print(title)
			if eadid != target {
				fmt.Println(": KO")
				writeString, err := tsvFile.WriteString(fmt.Sprintf("%d\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", repositoryId, resourceId, title, eadlocat, id0, id1, id2, id3, eadid, target))
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
