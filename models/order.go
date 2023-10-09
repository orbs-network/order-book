package models

type Order struct {
	Id        string
	Price     string
	Symbol    string
	Size      string
	Signature *string // EIP 712
	Pending   bool
}
