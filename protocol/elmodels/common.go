package elmodels

type BaseSnPayload struct {
	ProductID int    `json:"productID"`
	Sn        string `json:"sn"`
}