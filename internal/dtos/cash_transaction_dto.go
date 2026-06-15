package dtos

type CreateCashTransactionRequest struct {
	ReferenceNumber        string `json:"reference_number" validate:"required"`
	TransactionDescription string `json:"transaction_description" validate:"required"`
	TransactionType        string `json:"transaction_type" validate:"required"`
	TransactionDate        string `json:"transaction_date"`
	PaidBy                 string `json:"paid_by"`
	TransactionAmount      string `json:"transaction_amount"`
	CustomerType           string `json:"customer_type" validate:"omitempty,oneof=customer member transporter supplier guest vendor employee"`
	CustomerID             uint64 `json:"customer_id"`
	PaymentModeID          uint64 `json:"payment_mode_id"`
	PaymentType            string `json:"payment_type"`
	TransactionID          int64  `json:"transaction_id"`
}

type CashTransactionResponse struct {
	ID                     uint64 `json:"id"`
	ReferenceNumber        string `json:"reference_number"`
	TransactionDescription string `json:"transaction_description"`
	TransactionType        string `json:"transaction_type"`
	TransactionDate        string `json:"transaction_date"`
	PaidBy                 string `json:"paid_by"`
	TransactionAmount      string `json:"transaction_amount"`
	CustomerType           string `json:"customer_type"`
	CustomerID             uint64 `json:"customer_id"`
	PaymentModeID          uint64 `json:"payment_mode_id"`
	PaymentType            string `json:"payment_type"`
	TransactionID          int64  `json:"transaction_id"`
	CreatedAt              string `json:"created_at"`
	UpdatedAt              string `json:"updated_at"`
}
