package parsers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/yorikya/csvpl/user"

	"github.com/extrame/xls"
)

func ParseXSLRowToUser(row *xls.Row) (*user.User, error) {
	bamnt, err := strconv.Atoi(row.Col(6))
	if err != nil {
		return nil, fmt.Errorf("fail to parse brakefast amount, error: %s", err)
	}

	lamnt, err := strconv.Atoi(row.Col(8))
	if err != nil {
		return nil, fmt.Errorf("fail to parse lunch amount, error: %s", err)
	}

	bprice, err := strconv.ParseFloat(row.Col(7), 32)
	if err != nil {
		return nil, fmt.Errorf("fail to parse brakefast price, error: %s", err)
	}

	lprice, err := strconv.ParseFloat(row.Col(9), 32)
	if err != nil {
		return nil, fmt.Errorf("fail to parse lunch price, error: %s", err)
	}

	tprice, err := strconv.ParseFloat(row.Col(10), 32)
	if err != nil {
		return nil, fmt.Errorf("fail to parse total price, error: %s", err)
	}

	return user.NewUser(row.Col(0), row.Col(1), row.Col(2), row.Col(3),
		bamnt, lamnt, bprice, lprice, tprice), nil

}

func ParseUserMetadata(rowColums []string) *user.UserMetadata {
	price, err := strconv.Atoi(rowColums[6])
	if err != nil {
		log.Println("[ERROR] failed parse price from the metadata file, set to default 24, error:", err)
		price = 24
	}
	return user.NewUserMetadata(rowColums[0], rowColums[1], rowColums[2], rowColums[3],
		rowColums[4], rowColums[5], rowColums[7], rowColums[8], price)
}
