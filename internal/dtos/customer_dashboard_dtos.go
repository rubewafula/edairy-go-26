package dtos

type CustomerDashboardResponse struct {
	CustomerBillingStat  CustomerDashboardBillingStat  `json:"customer_billing_stat"`
	InvoiceStat          CustomerDashboardInvoiceStat  `json:"invoice_stat"`
	CustomerPaymentStat  CustomerDashboardPaymentStat  `json:"customer_payment_stat"`
	BillingTrend         []CustomerDashboardTrendPoint `json:"billing_trend"`
}

type CustomerDashboardBillingStat struct {
	TotalInvoicesThisMonth int64   `json:"total_invoices_this_month"`
	TotalBillingsThisMonth float64 `json:"total_billings_this_month"`
	Currency               string  `json:"currency"`
}

type CustomerDashboardInvoiceStat struct {
	PendingInvoices int64 `json:"pending_invoices"`
	SettledInvoices int64 `json:"settled_invoices"`
}

type CustomerDashboardPaymentStat struct {
	TotalPaymentsThisMonth float64 `json:"total_payments_this_month"`
}

type CustomerDashboardTrendPoint struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}
