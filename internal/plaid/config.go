package plaid

import (
	"github.com/plaid/plaid-go/v32/plaid"
)

type Config struct {
	ClientID      string
	Secret        string
	Env           plaid.Environment
	SandboxSecret string
	RedirectURL   string
}
