package main

import (
	"os"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/sirupsen/logrus"
)

func main() {

	fileName := os.Args[1]
	readExcel(fileName)

}

func readExcel(fileName string) {

	logrus.Info("reading Excel ", fileName)

	f, err := excelize.OpenFile(fileName)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info(f)
	firstSheet := f.WorkBook.Sheets.Sheet[0].Name
	rows := f.GetRows(firstSheet)

	logrus.Info(rows)
	logrus.Info("closing excel")

}
