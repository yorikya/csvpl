package user

type User struct {
	Departament, Name, BudgetID, LuncherID              string
	BreakfestAmount, LunchAmount, CoveredPrice          int
	BreakfestPrice, LunchPrice, TotalPrice, ChargePrice float64
}

func NewUser(departament, name, budgetID, luncherID string,
	breakfestAmount, lunchAmount int,
	breakfestPrice, lunchPrice, totalPrice float64) *User {
	totalCompany := (breakfestAmount + lunchAmount) * 24 //The default price is 24 nis
	return &User{
		Departament:     departament,
		Name:            name,
		BudgetID:        budgetID,
		LuncherID:       luncherID,
		BreakfestAmount: breakfestAmount,
		LunchAmount:     lunchAmount,
		BreakfestPrice:  breakfestPrice,
		LunchPrice:      lunchPrice,
		TotalPrice:      totalPrice,
		CoveredPrice:    totalCompany,
		ChargePrice:     totalPrice - float64(totalCompany),
	}
}
