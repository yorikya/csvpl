package main

import (
	"flag"
	"fmt"
	"log"

	"database/sql"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/extrame/xls"
	_ "github.com/proullon/ramsql/driver"
)

var (
	mainFilePath   = flag.String("main", "", "the main file csv path")
	cybusFilePsath = flag.String("cybus", "", "the cybus csv")
)

func init() {
	flag.Parse()
	if *mainFilePath == "" {
		panic("must provide the main file path")
	}
	if *cybusFilePsath == "" {
		panic("must provide the cibus file path")
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
			q := fmt.Sprintf(`INSERT INTO main (account, enter_noon, to_date, from_date, evning, noon, morning, amount, budget_num, recepie_id, name, budget_name, departament) VALUES ("%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s");`,
				row1.Col(12), row1.Col(11), row1.Col(10), row1.Col(9), row1.Col(8), row1.Col(7), row1.Col(6), row1.Col(5), row1.Col(4), row1.Col(3), row1.Col(2), row1.Col(1), row1.Col(0))
			if _, err := db.Exec(q); err != nil {
				log.Println("get erro:", err)
			}
		}
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
	rows, err := f.GetRows("Sheet1")
	for _, row := range rows {
		log.Printf("the row: %+v", row)
	}
	return nil
}

func initDBSQL(db *sql.DB) error {
	batch := []string{
		`CREATE TABLE main ( account TEXT, enter_noon TEXT, to_date TEXT, from_date TEXT, evning TEXT, noon TEXT, morning TEXT, amount TEXT, budget_num TEXT, recepie_id TEXT, name TEXT, budget_name TEXT, departament TEXT);`,
		`CREATE TABLE cybus (group_number INT, company TEXT, group_name TEXT, departament TEXT, full_name TEXT, employe_id INT, total_full FLOAT, total FLOAT, company_part FLOAT, company_part_amount FLOAT, employe_part FLOAT, employe_part_amount FLOAT, deal_amount INT, deal_amount_copy INT, first_name TEXT, last_name TEXT, id TEXT);`,
	}
	for _, b := range batch {
		if _, err := db.Exec(b); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	db, err := sql.Open("ramsql", "Main")
	if err != nil {
		log.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()

	initDBSQL(db)
	parseMain(*mainFilePath, db)
	// time.Sleep(3 * time.Second)
	rows, err := db.Query(`SELECT * FROM main`)
	if err != nil {
		log.Println("get wrror when run query, ", err)
	}

	for rows.Next() {
		var s1, s2, s3, s4, s5, s6, s7, s8, s9, s10, s11, s12, s13 string
		if err := rows.Scan(&s1, &s2, &s3, &s4, &s5, &s6, &s7, &s8, &s9, &s10, &s11, &s12, &s13); err != nil {
			log.Println("Get an error when scan row, error:", err)
		}
		log.Println("the row:", s1, s2, s3, s4, s5, s6, s7, s8, s9, s10, s11, s12, s13)
	}

	parseCybus(*cybusFilePsath, db)
}
