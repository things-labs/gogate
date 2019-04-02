package elmodels

import (
	"github.com/thinkgos/gomo/protocol/elinkch/ctrl"
)

type DevPropReqPy struct {
	BaseSnPayload
	Params map[string]interface{} `json:"params,omitempty"`
}

type DevPropRequest struct {
	ctrl.BaseRequest
	Payload DevPropReqPy `json:"payload,omitempty"`
}

type DevPropRspPy struct {
	ProductID int         `json:"productID"`
	Sn        string      `json:"sn"`
	Data      interface{} `json:"data"`
}

//type DevPropResponsePayload struct {
//	Payload DevPropRspPy `json:"payload,omitempty"`
//}
