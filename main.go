package main

import (
	"fmt"
	"os"
	"strconv"
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

	fileName_1 := "tc.csv"
	fileName_2 := "td.csv"
	transaction := readCredit(fileName_1)
	sumDebt(transaction)
	readDebit(fileName_2, transaction)
}

func sumDebt(transactionCredit []Transaction) {

	debt := 0
	for _, trans := range transactionCredit {
		stump := strings.Replace(trans.Valor, ".", "#", 2)
		fixCommas := strings.Replace(stump, ",", "", 2)
		fixdots := strings.Replace(fixCommas, "#", ",", 2)
		finalTrans := fixdots[0:len(fixdots)-3] + ""
		value, err := strconv.Atoi(finalTrans)
		if err != nil {
			logrus.Error(err)
		}
		debt = debt + value
	}
	logrus.Info("Total Debt: ", debt)
}

func readCredit(fileName string) []Transaction {
	logrus.Info("reading Credit")

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
			Oficina:     "CREDIT CARD VISA",
			Descripcion: record[1],
			Referencia:  " ",
			Valor:       fmt.Sprintf("-%s", record[4]),
			Moneda:      "COP",
		}
		transaction = append(transaction, user)
	}

	return transaction
}

func readDebit(fileName string, transactionsCredit []Transaction) {

	logrus.Info("reading Debit ", fileName)

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

	transaction = append(transaction, transactionsCredit...)
	transaction = modifyCSV(transaction)
	writeCSV(transaction)
}

func modifyCSV(transactions []Transaction) []Transaction {

	removerIndexest := []int{}
	for i := range transactions {

		// Descripcion
		if strings.Contains(transactions[i].Descripcion, "ABONO INTERESES") {
			transactions[i].Descripcion = "INCOME, INTERESTS"
		} else if strings.Contains(transactions[i].Descripcion, "GENIUS SPORTS") {
			transactions[i].Descripcion = "SALARY"
		} else if strings.Contains(transactions[i].Descripcion, "BODYTECH") {
			transactions[i].Referencia = transactions[i].Descripcion
			transactions[i].Descripcion = "GYM"
		} else if strings.Contains(transactions[i].Descripcion, "CAFE") ||
			strings.Contains(transactions[i].Descripcion, "CREPES") ||
			strings.Contains(transactions[i].Descripcion, "PIZZ") || strings.Contains(transactions[i].Descripcion, "RAPPI") {
			transactions[i].Referencia = transactions[i].Descripcion
			transactions[i].Descripcion = "COMIDA"
		} else if strings.Contains(transactions[i].Descripcion, "RETIRO") {
			transactions[i].Descripcion = "RETIRO"
		} else if strings.Contains(transactions[i].Descripcion, "CELULAR") {
			transactions[i].Descripcion = "PAGO CELULAR"
		} else if strings.Contains(transactions[i].Descripcion, "AMZN Mktp") {
			transactions[i].Referencia = transactions[i].Descripcion
			transactions[i].Descripcion = "SHOPPING"
		} else if strings.Contains(transactions[i].Descripcion, "Amazon Prime") || strings.Contains(transactions[i].Descripcion, "STAR PLUS") {
			transactions[i].Referencia = transactions[i].Descripcion
			transactions[i].Descripcion = "ENTRENTENIMIENTO"
		}else if strings.Contains(transactions[i].Oficina,"SERVICIOS ELCTR.") {
			transactions[i].Referencia = transactions[i].Descripcion
			transactions[i].Descripcion = "SHOPPING " + transactions[i].Referencia 
		}

		//Referencia
		if strings.Contains(transactions[i].Referencia, "29907879915") {
			transactions[i].Descripcion = "TRANSFERENCIA LAURA"
		}

		//Delete Rows

		if strings.Contains(transactions[i].Descripcion, "ABONO SUCURSAL VIRTUAL") ||
			strings.Contains(transactions[i].Descripcion, "PAGO SUC VIRT TC VISA") {
			removerIndexest = append(removerIndexest, i)
		}
	}

	deleted := 0
	for i := range removerIndexest {
		transactions = removeIndex(transactions, removerIndexest[i]-deleted)
		deleted++
	}

	logrus.Info(transactions)
	return transactions
}

func removeIndex(s []Transaction, index int) []Transaction {

	size := len(s)
	if index == size-1 {
		return s[:index]
	}
	return append(s[:index], s[index+1:]...)
}

func writeCSV(transactions []Transaction) {

	headers := []string{"FECHA", "OFICINA", "DESCRIPCION", "REFERENCIA", "VALOR", "MONEDA"}
	resultFileName := "EXTRACTO_ENERO.csv"
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
