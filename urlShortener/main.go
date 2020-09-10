package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	_ "log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"urlShortener/urlshort"
)

//edit accordingly
var  mysqlUser = "username"
var mysqlPass = "password"

var(
	yamlPath = flag.String("yaml","","path to yaml file containing shorten URLs")
	jsonPath = flag.String("json","","path to json file containing shorten URLs")
	dbInfo = flag.String("db","","mysql database in the format database_name/table_name, table with field path,url")

)

func main() {

	flag.Parse()
	var handler http.HandlerFunc

	if *yamlPath !="" && filepath.Ext(*yamlPath)==".yaml" {
		yaml, err := getContentFromFile(*yamlPath)
		handleError(err)

		handler, err = urlshort.YAMLHandler(yaml, defaultMux())
		handleError(err)
	}else if *jsonPath !="" && filepath.Ext(*jsonPath)==".json" {

		json, err := getContentFromFile(*jsonPath)
		handleError(err)

		handler, err = urlshort.JSONHandler(json, defaultMux())
		handleError(err)
	} else if *dbInfo !=""  && strings.Contains(*dbInfo,"/"){
		 dbSlice := strings.Split(*dbInfo,"/")

		if len(dbSlice) !=2 {
			handleError(errors.New("Not in specified format"))
		}
		dbName := dbSlice[0]
		tableName:= dbSlice[1]
		db, err := sql.Open("mysql", mysqlUser+":"+mysqlPass+"@tcp(127.0.0.1:3306)/"+dbName)

		if err != nil {
			handleError(err)
		}

		handler, err = urlshort.DBHandler(db, tableName, defaultMux())
		handleError(err)

		defer db.Close()
	}else{
		handleError(errors.New("No flag is set"))
	}


    http.ListenAndServe(":80", handler)
}

func handleError(err error){
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func getContentFromFile(filePath string) ([]byte, error){

	data, err := ioutil.ReadFile(filePath)

	if err != nil {
		return []byte(""), err
	}

	return data , nil
}


func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}
func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Hello, world!")
}
