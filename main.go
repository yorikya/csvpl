package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/extrame/xls"
	_ "github.com/proullon/ramsql/driver"

	"github.com/yorikya/csvpl/parsers"
	"github.com/yorikya/csvpl/reports"
	"github.com/yorikya/csvpl/user"
)

func parse1Lunches(path string) error {
	log.Println("The Lunches file path:", path)
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	// Get all the rows in the Sheet1.
	rows := f.GetRows("1.2020")

	for i, row := range rows {
		if i < 2 || i == len(rows)-1 { //skip header and footer
			continue
		}
		// group_number TEXT, company TEXT, group_name TEXT, departament TEXT, full_name TEXT, employe_id TEXT, total_full TEXTS, total TEXT, company_part TEXT, company_part_amount TEXT, employe_part TEXT, employe_part_amount TEXT, deal_amount TEXT, deal_amount_copy TEXT, first_name TEXT, last_name TEXT, id TEXT
		var name string
		s := strings.Split(row[4], " ")
		name = strings.Join(s[1:], " ") + " " + s[0] //fix name order
		q := fmt.Sprintf(`INSERT INTO cybus (group_number, company, group_name, departament, full_name, employe_id, total_full, total, company_part, company_part_amount, employe_part, employe_part_amount, deal_amount, deal_amount_copy, first_name, last_name, id) VALUES ("%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s");`,
			row[0], row[1], row[2], row[3], name, row[5], row[6], row[7], row[8], row[9], row[10], row[11], row[12], row[13], row[14], row[15], row[16])

		log.Println("exec query:", q)

	}
	return nil
}

func parseLunches(path string) []*user.User {
	xlFile, err := xls.Open(path, "utf-8")
	if err != nil {
		log.Println("failed to open file", path, "error:", err)
		return nil
	}

	if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
		users := make([]*user.User, sheet1.MaxRow)
		fmt.Println("Total Lines ", sheet1.MaxRow, sheet1.Name)
		for i := 1; i <= (int(sheet1.MaxRow)); i++ {
			u, err := parsers.ParseXSLRowToUser(sheet1.Row(i))
			if err != nil {
				log.Println("failed to parse line:", i, "error:", err)
				continue
			}
			users[i] = u
		}
		return users
	}
	return nil
}

func main() {
	users := parseLunches("resources/lunches.xls")
	reports.GenerateLunchReport("lunchreport.xls", users)

}
