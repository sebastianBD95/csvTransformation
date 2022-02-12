package main

import (
	"fmt"
	"os"
	"strings"

	"encoding/csv"

	"github.com/sirupsen/logrus"
)

type Transaction struct {
	Fecha       string
	Oficina     string
	Descripcion string
	Referencia  string
	Valor       string
	Moneda      string
}

func (trans Transaction) ToSlice() []string {
	return []string{trans.Fecha, trans.Oficina, trans.Descripcion, trans.Referencia, trans.Valor, trans.Moneda}
}

func main() {

	fileName := "td.csv"
	readExcel(fileName)

}

func readExcel(fileName string) {

	logrus.Info("reading Excel ", fileName)

	f, err := os.Open(fileName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'
	if _, err := r.Read(); err != nil {
		logrus.Fatal(err)
	}

	records, err := r.ReadAll()
	transaction := make([]Transaction, 0)

	if err != nil {
		logrus.Fatal(err)
	}
	for _, record := range records {

		user := Transaction{
			Fecha:       record[0],
			Oficina:     record[2],
			Descripcion: record[3],
			Referencia:  record[4],
			Valor:       record[5],
			Moneda:      "COP",
		}
		transaction = append(transaction, user)
	}

	modifyCSV(transaction)
	writeCSV(transaction, fileName)
}

func modifyCSV(transactions []Transaction) {

	for i := range transactions {

		if strings.Contains(transactions[i].Descripcion, "ABONO") {
			transactions[i].Descripcion = "ABONO"
		} else if strings.Contains(transactions[i].Descripcion, "GENIUS SPORTS") {
			transactions[i].Descripcion = "SALARY"
		} else if strings.Contains(transactions[i].Descripcion, "BODYTECH") {
			transactions[i].Descripcion = "GYM"
		} else if strings.Contains(transactions[i].Descripcion, "CAFE") || strings.Contains(transactions[i].Descripcion, "CREPES") || strings.Contains(transactions[i].Descripcion, "PIZZ") {
			transactions[i].Descripcion = "COMIDA"
		} else if strings.Contains(transactions[i].Descripcion, "RETIRO") {
			transactions[i].Descripcion = "RETIRO"
		} else if strings.Contains(transactions[i].Descripcion, "CELULAR") {
			transactions[i].Descripcion = "PAGO CELULAR"
		}

		if strings.Contains(transactions[i].Referencia, "29907879915") {
			transactions[i].Descripcion = "TRANSFERENCIA LAURA"
		}
	}
	logrus.Info(transactions)
}

func writeCSV(transactions []Transaction, fileName string) {

	headers := []string{"FECHA", "OFICINA", "DESCRIPCION", "REFERENCIA", "VALOR", "MONEDA"}
	resultFileName := fmt.Sprintf("result_%s.csv", fileName)
	f, err := os.Create(resultFileName)

	if err != nil {
		logrus.Error(err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()
	w.Write(headers)

	for _, trans := range transactions {
		value := trans.ToSlice()
		if err := w.Write(value); err != nil {
			logrus.Error(err)
		}
	}

}
