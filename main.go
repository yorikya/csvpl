package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/extrame/xls"
	_ "github.com/proullon/ramsql/driver"

	"github.com/yorikya/csvpl/parsers"
	"github.com/yorikya/csvpl/reports"
	"github.com/yorikya/csvpl/user"
)

func parseLunches(path string) []*user.User {
	xlFile, err := xls.Open(path, "utf-8")
	if err != nil {
		log.Println("failed to open file", path, "error:", err)
		return nil
	}

	if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
		users := make([]*user.User, sheet1.MaxRow+1)
		fmt.Println("Total Lines ", sheet1.MaxRow, sheet1.Name)
		for i := 1; i <= (int(sheet1.MaxRow)); i++ {
			u, err := parsers.ParseXSLRowToUser(sheet1.Row(i))
			if err != nil {
				log.Println("[ERROR] failed to parse line:", i, "error:", err)
				continue
			}
			users[i] = u
		}
		return users
	}
	return nil
}

var lunchesFilePath = flag.String("lunches", "", "the lunches file location")
var metadataFilePath = flag.String("metadata", "resources/metadata.xlsx", "metadata for users")

func main() {
	flag.Parse()
	users := parseLunches(*lunchesFilePath)
	reports.GenerateLunchReport("lunchreport.xls", users)

	reports.GenerateFinalReport("finalreport.xlsx", *metadataFilePath, users)
}
