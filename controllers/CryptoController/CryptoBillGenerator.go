package controllers

import (
	"fmt"
	"gokripto/Database"
	model "gokripto/Model"
	"log"
	"os"
	"time"

	generator "github.com/angelodlfrtr/go-invoice-generator"
)

// Update gerekiyor şu an için kullanım dışı

func BillGenerator(UserID string, price float64, cryptoname string, amount float64, transactionType string) {
	directoryFileTime := time.Now()
	var directoryUserName string
	if err := Database.DB.Model(&model.User{}).Where("id", UserID).Pluck("name", &directoryUserName).Error; err != nil {
		log.Fatal(err)
	}

	var walletAddress string
	if err := Database.DB.Model(&model.User{}).Where("id", UserID).Pluck("wallet_address", &walletAddress).Error; err != nil {
		log.Fatal(err)
	}

	var directoryFileName string
	directoryFileName = UserID + "_" + directoryUserName + "_" + directoryFileTime.Format("20060102_15041") + ".pdf"

	doc, _ := generator.New(generator.DeliveryNote, &generator.Options{
		TextTypeInvoice: "FACTUasdasdasRE",
		AutoPrint:       false,
	})

	doc.SetRef("testref")

	doc.SetDescription("Crypto Exchange")

	doc.SetDate(time.Now().Format("2006/01/02 15:04"))

	logoBytes, err := os.ReadFile("/home/umut/goKripto/Pdfs/Images/proxolab.jpg")
	if err != nil {
		log.Fatal(err)
	}

	doc.SetCompany(&generator.Contact{
		Name: "Proxolab",
		Logo: logoBytes,
		Address: &generator.Address{
			Address:    "Çamlaraltı, 6025.",
			Address2:   "Apartman No: 9",
			PostalCode: "20160",
			City:       "Denizli",
			Country:    "Türkiye",
		},
	})

	var customerAddress = &generator.Address{
		Address: walletAddress,
	}

	doc.SetCustomer(&generator.Contact{
		Name:    directoryUserName,
		Address: customerAddress,
	})

	unitCostStr := fmt.Sprintf("%.2f", price)
	quantityStr := fmt.Sprintf("%.2f", amount)

	doc.AppendItem(&generator.Item{
		Name:        cryptoname,
		Description: transactionType,
		UnitCost:    unitCostStr,
		Quantity:    quantityStr,
	})

	pdf, err := doc.Build()
	if err != nil {
		log.Fatal(err)
	}
	os.MkdirAll("/home/umut/goKripto/Pdfs/", os.ModePerm)
	err = pdf.OutputFileAndClose("/home/umut/goKripto/Pdfs/" + directoryFileName)
	if err != nil {
		log.Fatal(err)
	}
}
