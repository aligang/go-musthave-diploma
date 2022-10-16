package customer_account

type AccountBalance struct {
	Current  float64 `json:"current"`
	Withdraw float64 `json:"withdrawn"`
}

type CustomerAccount struct {
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
	AccountBalance
}

type CustomerAccounts map[string]CustomerAccount

func New() *CustomerAccount {
	return &CustomerAccount{
		AccountBalance: AccountBalance{},
	}
}
