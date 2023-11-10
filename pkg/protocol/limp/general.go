package limp

import (
	"encoding/hex"

	"github.com/thinkgos/gogate/pkg/ltl"
	"github.com/thinkgos/gogate/pkg/ltl/ltlspec"
	"github.com/thinkgos/gogate/pkg/numeric"
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

func BasicAttribute(rRec []ltl.ReadRspStatus) *GenerlBasicAttribute {
	o := &GenerlBasicAttribute{}
	for _, v := range rRec {
		switch v.AttrID {
		case ltlspec.ATTRID_BASIC_LTL_VERSION:
			if v.Status == ltl.LTL_SUCCESS {
				o.LTLVersion = v.MustUint8()
			}

		case ltlspec.ATTRID_BASIC_APPL_VERSION:
			if v.Status == ltl.LTL_SUCCESS {
				o.APPVersion = ltlspec.Version(v.MustUint16())
			}
		case ltlspec.ATTRID_BASIC_HW_VERSION:
			if v.Status == ltl.LTL_SUCCESS {
				o.HWVersion = ltlspec.Version(v.MustUint16())
			}
		case ltlspec.ATTRID_BASIC_MANUFACTURER_NAME:
			if v.Status == ltl.LTL_SUCCESS {
				o.ManufacturerName = v.MustString()
			}
		case ltlspec.ATTRID_BASIC_BUILDDATE_CODE:
			if v.Status == ltl.LTL_SUCCESS {
				o.BuildDateCode = ltlspec.BuildDate(v.MustUint32())
			}
		case ltlspec.ATTRID_BASIC_PRODUCT_ID:
			if v.Status == ltl.LTL_SUCCESS {
				o.ProductIdentifier = int(v.MustUint32())
			}
		case ltlspec.ATTRID_BASIC_SERIAL_NUMBER:
			if v.Status == ltl.LTL_SUCCESS {
				o.SerialNumber =
					hex.EncodeToString(numeric.ReverseBytes(v.MustArrayUint8()))
			}
		case ltlspec.ATTRID_BASIC_POWER_SOURCE:
			if v.Status == ltl.LTL_SUCCESS {
				o.PowerSource = ltlspec.PowerSource(v.MustUint8())
			}
		}
	}
	return o
}

type GenerlOnoffAttribute struct {
	OnOffStateValue bool
}

func OnoffAttribute(rRec []ltl.ReportRec) *GenerlOnoffAttribute {
	o := &GenerlOnoffAttribute{}
	if len(rRec) > 0 && rRec[0].AttrID == ltlspec.ATTRID_ONOFF_STATUS {
		o.OnOffStateValue = rRec[0].MustBool()
	}
	return o
}
