package limp

import (
	"github.com/thinkgos/gogate/pkg/ltl"
	"github.com/thinkgos/gogate/pkg/ltl/ltlspec"
)

type MsAttribute struct {
	MeasuredValue    float32 `json:",string"`
	MinMeasuredValue float32 `json:",string"`
	MaxMeasuredValue float32 `json:",string"`
	Tolerance        float32 `json:",string"`
}

type MsMeasure struct {
	MeasuredValue float32 `json:",string"`
}

func MsMeasureAttrReport(trunkid uint16, rRec []ltl.ReportRec) (*MsMeasure, error) {
	var mstemp *MsMeasure

	if len(rRec) == 0 {
		return nil, ErrInvalidData
	}
	val, err := rRec[0].Uint16()
	if err != nil {
		return nil, err
	}
	switch trunkid {
	case ltl.TrunkID_MsTemperatureMeasurement:
		if rRec[0].AttrID != ltlspec.ATTRID_MS_TEMPERATURE_MEASURED_VALUE {
			return nil, ErrTrunkNotImplement
		}
		if val == ltlspec.MS_TEMPERATURE_INVALID_VALUE {
			return nil, ErrInvalidData
		}
		mstemp = &MsMeasure{float32(val) / 100}
	case ltl.TrunkID_MsRelativeHumidity:
		if rRec[0].AttrID != ltlspec.ATTRID_MS_RELATIVE_HUMIDITY_MEASURED_VALUE {
			return nil, ErrTrunkNotImplement
		}
		if val == ltlspec.MS_RELATIVE_HUMIDITY_INVALID_VALUE {
			return nil, ErrInvalidData
		}
		mstemp = &MsMeasure{float32(rRec[0].MustUint16()) / 100}
	default:
		return nil, ErrTrunkNotImplement
	}

	return mstemp, nil
}
