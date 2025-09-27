package domain

type OutstandingData struct {
	LoanID      int64
	CustomerID  int64
	Outstanding int64
}

type PaymentData struct {
	PaymentID   int64
	Amount      int64
	PayAmount   int64
	Outstanding int64
	TermWeeks   int
	WeeksRemain int
}

type MakePaymentParam struct {
	LoanID int64
	Amount int64
}
