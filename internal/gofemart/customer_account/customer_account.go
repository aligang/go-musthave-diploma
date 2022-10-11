package customer_account

type AccountBalance struct {
	Balance  float64 `json:"balance"`
	Withdraw float64 `json:"withdrawn"`
}

type CustomerAccount struct {
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
	AccountBalance
}

type CustomerAccounts map[string]CustomerAccount
