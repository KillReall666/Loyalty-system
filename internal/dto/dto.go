package dto

type Order struct {
	UserId      string
	OrderNumber string
	OrderDate   string
}

type FullOrder struct {
	OrderNumber string  `json:"number"`
	OrderStatus string  `json:"status"`
	Accrual     float32 `json:"accrual"`
	OrderDate   string  `json:"uploaded_at"`
}

type UserBalance struct {
	Current   float32
	Withdrawn float32
}
