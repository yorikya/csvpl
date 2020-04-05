package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"database/sql"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/extrame/xls"
	_ "github.com/proullon/ramsql/driver"
)

var (
	baseFilePath  = flag.String("base", "", "the base file csv path")
	mainFilePath  = flag.String("main", "", "the main file csv path")
	cybusFilePath = flag.String("cybus", "", "the cybus csv")
	outFilePath   = flag.String("out", "out.xlsx", "the output file (default: out.xlsx)")
)

func init() {
	flag.Parse()
	if *mainFilePath == "" {
		panic("must provide the main file path")
	}
	if *baseFilePath == "" {
		panic("must provide the base file path")
	}

}

func parseMain(path string, db *sql.DB) error {
	log.Println("The main file path:", path)
	xlFile, err := xls.Open(path, "utf-8")
	if err != nil {
		return err
	}
	if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
		fmt.Print("Total Lines ", sheet1.MaxRow, sheet1.Name)
		for i := 1; i <= (int(sheet1.MaxRow)); i++ {
			row1 := sheet1.Row(i)
			var name string
			s := strings.Split(row1.Col(2), " ")
			name = strings.Join(strings.Split(row1.Col(2), " "), " ") //fix name encoding hebrew
			q := fmt.Sprintf(`INSERT INTO main (account, enter_noon, to_date, from_date, evning, noon, morning, amount, budget_num, recepie_id, name, budget_name, departament, first_name, last_name) VALUES ("%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s");`,
				row1.Col(12), row1.Col(11), row1.Col(10), row1.Col(9), row1.Col(8), row1.Col(7), row1.Col(6), row1.Col(5), row1.Col(4), row1.Col(3), name, row1.Col(1), row1.Col(0),
				s[0],                     //first name
				strings.Join(s[1:], " ")) //last name
			if _, err := db.Exec(q); err != nil {
				log.Println("insert main error:", err)
				continue
			}
			// log.Println("exec query:", q)
		}
	}
	return nil
}
func parseBase(path string, db *sql.DB) error {
	log.Println("the base file path:", path)
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	// Get all the rows in the Sheet1.
	rows := f.GetRows("Sheet1")
	for _, row := range rows {
		// code TEXT, site_code TEXT, site TEXT, employe_id TEXT, kibutz_id TEXT, launch_site_id TEXT, price TEXT, first_name TEXT, last_name TEXT
		q := fmt.Sprintf(`INSERT INTO base (code, site_code, site, employe_id, kibutz_id, launch_site_id, price, first_name, last_name, full_name) VALUES ("%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s");`,
			row[0], row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[8], fmt.Sprintf("%s %s", row[7], row[8]))
		if _, err := db.Exec(q); err != nil {
			log.Println("insert cybus error:", err)
			continue
		}
		log.Println("exec query:", q)

	}

	return nil
}

func parseCybus(path string, db *sql.DB) error {
	log.Println("The cybus file path:", path)
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	// Get all the rows in the Sheet1.
	rows := f.GetRows("Sheet1")
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
		if _, err := db.Exec(q); err != nil {
			log.Println("insert cybus error:", err)
			continue
		}
		log.Println("exec query:", q)

	}
	return nil
}

func initDBSQL(db *sql.DB) error {
	batch := []string{
		`CREATE TABLE main ( account TEXT, enter_noon TEXT, to_date TEXT, from_date TEXT, evning TEXT, noon TEXT, morning TEXT, amount TEXT, budget_num TEXT, recepie_id TEXT, name TEXT, budget_name TEXT, departament TEXT, first_name TEXT, last_name TEXT);`,
		`CREATE TABLE base ( code TEXT, site_code TEXT, site TEXT, employe_id TEXT, kibutz_id TEXT, launch_site_id TEXT, price TEXT, first_name TEXT, last_name TEXT, full_name TEXT);`,
		`CREATE TABLE cybus (group_number TEXT, company TEXT, group_name TEXT, departament TEXT, full_name TEXT, employe_id TEXT, total_full TEXTS, total TEXT, company_part TEXT, company_part_amount TEXT, employe_part TEXT, employe_part_amount TEXT, deal_amount TEXT, deal_amount_copy TEXT, first_name TEXT, last_name TEXT, id TEXT);`,
	}
	for _, b := range batch {
		if _, err := db.Exec(b); err != nil {
			return err
		}
	}
	return nil
}

func createSQLTables(basePath, mainPath, cybusPath string, db *sql.DB) error {
	if err := initDBSQL(db); err != nil {
		return err
	}

	if err := parseBase(basePath, db); err != nil {
		return err
	}

	if err := parseMain(mainPath, db); err != nil {
		return err
	}

	// if err := parseCybus(cybusPath, db); err != nil {
	// 	return err
	// }
	return nil
}

func logMissingUsers(db *sql.DB) {
	rows, err := db.Query(`SELECT main.name FROM main;`)
	if err != nil {
		log.Println("failed to find missing users, error:", err)
	}
	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			log.Println("Get an error when scan row, error:", err)
			continue
		}
		log.Println("missing user:", user)
	}

}

func cretaeOutpuFile(path string, db *sql.DB) error {
	hedaer := map[string]string{"A1": "employee_id", "B1": "full_name", "C1": "launch_price", "D1": "lunch_amount", "E1": "company_charge", "F1": "employe_charge", "G1": "lunch_total"}
	// values := map[string]int{"B2": 2, "C2": 3, "D2": 3, "B3": 5, "C3": 2, "D3": 4, "B4": 6, "C4": 7, "D4": 8}
	f := excelize.NewFile()
	for k, v := range hedaer {
		f.SetCellValue("Sheet1", k, v)
	}
	// code, site_code, site, employe_id, kibutz_id, launch_site_id, price, first_name, last_name
	rows, err := db.Query(`SELECT base.full_name, base.code, base.site_code, base.site, base.employe_id, base.kibutz_id, base.launch_site_id, base.price, base.first_name, base.last_name, main.amount, main.noon FROM base
	JOIN main ON base.full_name = main.name`)
	if err != nil {
		return err
	}

	i := 2
	users := []string{}
	for rows.Next() {
		var full_name, code, site_code, site, empl_id, kibutz_id, lunch_site_id, price, first_name, last_name, amount, noon string
		if err := rows.Scan(&full_name, &code, &site_code, &site, &empl_id, &kibutz_id, &lunch_site_id, &price, &first_name, &last_name, &amount, &noon); err != nil {
			log.Println("Get an error when scan row, error:", err)
			continue
		}
		users = append(users, full_name)
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i), code)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i), site_code)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i), site)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", i), empl_id)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", i), kibutz_id)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", i), lunch_site_id)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", i), price)
		f.SetCellValue("Sheet1", fmt.Sprintf("H%d", i), noon)
		f.SetCellValue("Sheet1", fmt.Sprintf("I%d", i), first_name)
		f.SetCellValue("Sheet1", fmt.Sprintf("J%d", i), last_name)
		noonFloat, err := strconv.ParseFloat(noon, 64)
		if err != nil {
			log.Println("failed to parse noon lunches nuber, error", err)
		}
		priceFloat, err := strconv.ParseFloat(price, 64)
		if err != nil {
			log.Printf("failed to parse price, for name: %s, %s, error: %s", first_name, last_name, err)
		}
		totalCompany := noonFloat * priceFloat
		f.SetCellValue("Sheet1", fmt.Sprintf("K%d", i), totalCompany)
		totalCharge, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			log.Println("failed to parse amount, error", err)
		}

		if c := totalCharge - totalCompany; c > 0 {
			f.SetCellValue("Sheet1", fmt.Sprintf("L%d", i), fmt.Sprintf("%.2f", c))
		} else {
			f.SetCellValue("Sheet1", fmt.Sprintf("L%d", i), "")
		}

		f.SetCellValue("Sheet1", fmt.Sprintf("M%d", i), amount)
		i++
		// log.Println("the row:", code, site_code, site, empl_id, kibutz_id, lunch_site_id, price, first_name, last_name, amount, noon)
	}

	rows, err = db.Query(`SELECT main.name FROM main;`)
	if err != nil {
		log.Println("failed to find missing users, error:", err)
	}

	missingUsers := []string{}
	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			log.Println("Get an error when scan row, error:", err)
			continue
		}
		if !contain(users, user) {
			missingUsers = append(missingUsers, user)
		}

	}
	log.Printf("create output file %s, with %d records", path, i-1)
	log.Printf("missing users in base file, users: %+v", missingUsers)
	return f.SaveAs(path)
}

func contain(col []string, key string) bool {
	for _, k := range col {
		if k == key {
			return true
		}
	}
	return false
}

func main() {
	db, err := sql.Open("ramsql", "Main")
	if err != nil {
		log.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()

	if err := createSQLTables(*baseFilePath, *mainFilePath, *cybusFilePath, db); err != nil {
		log.Panic(err)
	}

	if err := cretaeOutpuFile(*outFilePath, db); err != nil {
		log.Panic(err)
	}
}
