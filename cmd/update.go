package cmd

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLConf struct {
   Host string `json:"host"`
   Port int `json:"port"`
   Db string `json:"db"`
   Username string `json:"username"`
   Password string `json:"password"`
}

var csvFile string
var mySQLConf MySQLConf


var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update malformed eadids",
	Run: func(cmd *cobra.Command, args []string) {
		err := initMysqlConfig()
		if err != nil {
			log.Fatal(err)
		}
		update()
	},
}

func initMysqlConfig() error {
	viper.SetConfigFile("mysql.json")
	viper.AddConfigPath("github.com/nyudlts/eadid-update")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&mySQLConf); err != nil {
		return err
	}

	mySQLConf.Db = viper.GetString("db")
	mySQLConf.Host = viper.GetString("host")
	mySQLConf.Password = viper.GetString("password")
	mySQLConf.Username = viper.GetString("username")
	mySQLConf.Port = viper.GetInt("port")

	return nil
}

func init() {
	cobra.OnInitialize(updateConfig)
	rootCmd.PersistentFlags().StringVar(&csvFile, "csv", "", "config file Required")
	rootCmd.PersistentFlags().StringP("mysql", "m", "mysql.conf", "mysql config")
	viper.BindPFlag("mysql", rootCmd.PersistentFlags().Lookup("mysql"))
	rootCmd.AddCommand(updateCmd)
}

func updateConfig() {
	if len(csvFile) < 0 {
		log.Fatal("--csv flag is mandatory")
	}
}

func update() error {
	connection := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mySQLConf.Username, mySQLConf.Password, mySQLConf.Host, mySQLConf.Port, mySQLConf.Db)
	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Fatal(err)
	}
	csvFile, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		resourceId, err := strconv.Atoi(line[1])
		if err != nil {
			return err
		}

		eadid := line[4]

		log.Printf("updateing resource %d ead_id to: %s\n", resourceId, eadid)

		query := fmt.Sprintf("Update resource set ead_id = '%s' where id = '%d'", eadid, resourceId)

		result, err := db.Exec(query)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("success %v\n", result)


	}
	return nil
}