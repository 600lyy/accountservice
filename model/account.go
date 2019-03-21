package model

// Account delcarition with built-in json serialization
type Account struct {
	ID       uint64 `json:"id"`
	UserName string `json:"username"` //Primary Key
	Name	 string `json:"Name"`
	Passwd	 string `json:"password,omitempty"`
	ServedBy string `json:"servedBy"`
}
