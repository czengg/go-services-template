package banking

type consumer struct {
	ID            string
	PCID          string
	ExternalID    string
	IsActive      bool
	KYCStatus     string
	TaxIDType     string
	TaxIdentifier string
	CreatedAt     string
	UpdatedAt     string
	Deleted       bool
}
