package main

import (
	"fmt"
	"github.com/anjanashankar9/random-learning/go-sqlite/database"
	"github.com/anjanashankar9/random-learning/go-sqlite/database/client"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {

	db, err := client.New("sqlite3", ":memory:")

	if err != nil {
		log.Fatal(err)
	}

	dbc := database.New(db)
	err = dbinserts(err, dbc)

	err = dbc.RenameTable()
	if err != nil {
		log.Fatal(err)
	}

	resEvents, err := dbc.GetAll()
	fmt.Println(resEvents)

	fmt.Println("\n\n****")

	cnt, err := dbc.GetCount("All")
	fmt.Println(cnt)

	//resEvents, err = dbc.GetWithTag("true")
	resEvents, err = dbc.GetWithTag("False")
	fmt.Println(resEvents)

	fmt.Println("-------------------")
	dbc.ReadInBatches()
}

func dbinserts(err error, dbc database.DatabaseClient) error {
	err = dbc.CreateTable()
	if err != nil {
		log.Fatal(err)
	}

	eip1 := database.EIP{
		FI:              "dev1-uswest2",
		FD:              "foundation",
		Env:             "dev",
		Service:         "skywalker",
		ExternalIP:      "10.100.0.12",
		UsedForFdEgress: "true",
		Tags:            "k1:v1,k2:v2",
	}

	err = dbc.InsertIntoTable(eip1)
	if err != nil {
		log.Fatal(err)
	}

	eip2 := database.EIP{
		FI:              "aws-dev1-uswest2",
		FD:              "foundation2",
		Env:             "dev",
		Service:         "skywalker",
		ExternalIP:      "10.100.0.121",
		UsedForFdEgress: "true",
		Tags:            "k1:v1,k2:v2",
	}

	err = dbc.InsertIntoTable(eip2)
	if err != nil {
		log.Fatal(err)
	}

	eip3 := database.EIP{
		FI:              "aws-dev1-uswest2",
		FD:              "foundation2",
		Env:             "dev2",
		Service:         "skywalker",
		ExternalIP:      "10.100.0.121",
		UsedForFdEgress: "false",
		Tags:            "k1:v1,k2:v2",
	}

	err = dbc.InsertIntoTable(eip3)
	if err != nil {
		log.Fatal(err)
	}

	eip4 := database.EIP{
		FI:              "aws-dev1-uswest2",
		FD:              "foundation2",
		Env:             "dev2",
		Service:         "skywalker",
		ExternalIP:      "10.100.0.121",
		UsedForFdEgress: "true",
		Tags:            "k1:v1,k2:v2",
	}

	err = dbc.InsertIntoTable(eip4)
	if err != nil {
		log.Fatal(err)
	}

	return err
}
