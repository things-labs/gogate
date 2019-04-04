package elinkmd

type DevPropRspPy struct {
	ProductID int         `json:"productID"`
	Sn        string      `json:"sn"`
	Data      interface{} `json:"data"`
}

//type DevPropResponsePayload struct {
//	Payload DevPropRspPy `json:"payload,omitempty"`
//}
