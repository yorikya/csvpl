package reports

import (
	"fmt"
	"log"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/yorikya/csvpl/user"
)

func GenerateLunchReport(path string, users []*user.User) {
	sheetName := "LunchReport"
	f := excelize.NewFile()
	f.NewSheet(sheetName)
	for k, v := range map[string]string{ //Line 1 / Header
		"A1": "מחלקה", "B1": "שם", "C1": "מס תקציב",
		"D1": "מס אוכל", "E1": "סכום לחיוב", "F1": "שווי",
		"G1": "ניכוי"} {
		f.SetCellValue(sheetName, k, v)
	}

	line := 2
	total := 0.0
	for _, u := range users {
		if u == nil {
			log.Printf("skip user: %+v", u)
			continue
		}
		if u.Departament != "פלסגד" {
			log.Printf("non plasgad worker skip user: %+v", u)
			continue
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", line), u.Departament)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", line), u.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", line), u.BudgetID)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", line), u.LuncherID)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", line), u.TotalPrice)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", line), fmt.Sprintf("%d", u.CoveredPrice))
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", line), fmt.Sprintf("%.2f", u.ChargePrice))
		total += u.TotalPrice
		line++
	}
	log.Println("create report with", line-2, "users")
	f.SetCellValue(sheetName, fmt.Sprintf("E%d", line), fmt.Sprintf("%5.2f", total))

	err := f.SaveAs(path)
	if err != nil {
		log.Println("get an error when save report file:", path, "error:", err)
	}
}
