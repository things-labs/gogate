package elmodels

type BaseSnPayload struct {
	ProductID int    `json:"productID"`
	Sn        string `json:"sn"`
}

type DevicesInfo struct {
	ProductID int      `json:"productID"`
	Sn        []string `json:"sn"`
}

type ItemInfos struct {
	Pkid    int
	IsLocal bool // if local do not send message to up
	Val     interface{}
}
