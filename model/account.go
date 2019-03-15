package model

// Account delcarition with built-in json serialization
type Account struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ServedBy string `json:"servedBy"`
}
