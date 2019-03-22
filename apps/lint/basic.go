package lint

import (
	"encoding/hex"

	"github.com/astaxie/beego/logs"
	"github.com/slzm40/gomo/ltl"
	"github.com/slzm40/gomo/ltl/attri"
)

type GenerlBasicAttribute struct {
	LTLVersion        uint8 `json:"-"`
	APPVersion        string
	HWVersion         string
	ManufacturerName  string
	BuildDateCode     string
	ProductIdentifier int
	SerialNumber      string
	PowerSource       string
}

func BasicAttribute(rRec []ltl.RcvReportRec) *GenerlBasicAttribute {
	o := &GenerlBasicAttribute{}
	for _, v := range rRec {
		switch v.AttrID {
		case attri.ATTRID_BASIC_LTL_VERSION:
			o.LTLVersion = v.MustUint8()
		case attri.ATTRID_BASIC_APPL_VERSION:
			o.APPVersion = attri.Version(v.MustUint16())
		case attri.ATTRID_BASIC_HW_VERSION:
			o.HWVersion = attri.Version(v.MustUint16())
		case attri.ATTRID_BASIC_MANUFACTURER_NAME:
			o.ManufacturerName = v.MustString()
		case attri.ATTRID_BASIC_BUILDDATE_CODE:
			o.BuildDateCode = attri.BuildDate(v.MustUint32())
		case attri.ATTRID_BASIC_PRODUCT_ID:
			o.ProductIdentifier = int(v.MustUint32())
		case attri.ATTRID_BASIC_SERIAL_NUMBER:
			o.SerialNumber = hex.EncodeToString(v.MustArrayUint8())
		case attri.ATTRID_BASIC_POWER_SOURCE:
			//o.PowerSource = v.MustUint8()
		}
	}
	logs.Debug(o)
	return o
}
