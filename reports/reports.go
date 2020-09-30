package reports

import (
	"fmt"
	"log"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/yorikya/csvpl/parsers"
	"github.com/yorikya/csvpl/user"
)

func loadUsersMetadata(metadataxlsx string) (map[string]*user.UserMetadata, error) {
	f, err := excelize.OpenFile(metadataxlsx)
	if err != nil {
		return nil, err
	}
	usersMeta := make(map[string]*user.UserMetadata)
	// Get all the rows in the Sheet1.
	rows := f.GetRows("Sheet1")
	for _, row := range rows {
		if row[4] == "" {
			log.Println("error: missing PlasgadID, user:", row[0:8])
			continue
		}
		usersMeta[row[4]] = parsers.ParseUserMetadata(row)
	}
	return usersMeta, nil
}

func GenerateFinalReport(path, metadataxlsx string, users []*user.User) {
	usersMeta, err := loadUsersMetadata(metadataxlsx)
	if err != nil {
		log.Println("[ERROR] failed to load users metadata, error:", err)
		return
	}

	sheetName := "Report"
	f := excelize.NewFile()
	f.NewSheet(sheetName)
	for k, v := range map[string]string{ //Line 1 / Header
		"A1": "קוד", "B1": "קידוד", "C1": "אתר",
		"D1": "מס עובד", "E1": "מס עובד בחדר אוכל", "F1": "",
		"G1": "מחיר ארוחה", "H1": "שם משפחה", "I1": "שם פרטי",
		"J1": "שווי", "K1": "חיוב", "L1": "סך הכל"} {
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
		metaUser, ok := usersMeta[u.LuncherID]
		if !ok {
			log.Printf("[ERROR] missing metadata for user: %+v\n", u)
			continue
		}
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", line), metaUser.Code)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", line), metaUser.SiteCode)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", line), metaUser.Site)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", line), metaUser.EmployeeIDPlasgad)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", line), metaUser.EmployeeIDFoodSite)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", line), metaUser.Column6)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", line), fmt.Sprintf("%d", metaUser.Price))
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", line), metaUser.LastName)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", line), metaUser.FirstName)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", line), fmt.Sprintf("%d", u.CoveredPrice))
		chargePrice := fmt.Sprintf("%.2f", u.ChargePrice)
		if u.ChargePrice <= 0 {
			chargePrice = ""
		}
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", line), chargePrice)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", line), fmt.Sprintf("%.2f", u.TotalPrice))
		total += u.TotalPrice
		line++
	}

	log.Println("create report with", line-2, "users")
	f.SetCellValue(sheetName, fmt.Sprintf("L%d", line), fmt.Sprintf("%5.2f", total))

	err = f.SaveAs(path)
	if err != nil {
		log.Println("get an error when save report file:", path, "error:", err)
	}
}

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
