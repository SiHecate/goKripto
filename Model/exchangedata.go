package model

type ExchangeData struct {
	Adapter     string  `json:"adapter"`
	From        string  `json:"from"`
	FromNetwork string  `json:"fromNetwork"`
	To          string  `json:"to"`
	ToNetwork   string  `json:"toNetwork"`
	AmountFrom  float64 `json:"amountFrom"`
	AmountTo    float64 `json:"amountTo"`
	MinAmount   float64 `json:"minAmount"`
	MaxAmount   float64 `json:"maxAmount"`
	QuotaID     string  `json:"quotaId"`
	Time        int     `json:"time"`
}
