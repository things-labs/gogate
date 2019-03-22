package lint

import (
	"github.com/slzm40/gomo/ltl"
	"github.com/slzm40/gomo/ltl/attri"
)

type MsAttribute struct {
	MeasuredValue    float32 `json: ",string"`
	MinMeasuredValue float32 `json: ",string"`
	MaxMeasuredValue float32 `json: ",string"`
	Tolerance        float32 `json: ",string"`
}

type MsMeasure struct {
	MeasuredValue float32 `json: ",string"`
}

func MsMeasureAttribute(rRec []ltl.RcvReportRec) (*MsMeasure, error) {
	if len(rRec) > 0 &&
		(rRec[0].AttrID != attri.ATTRID_MS_TEMPERATURE_MEASURED_VALUE ||
			rRec[0].AttrID != attri.ATTRID_MS_RELATIVE_HUMIDITY_MEASURED_VALUE) {
		return &MsMeasure{float32(rRec[0].MustUint32()) / 100}, nil
	}

	return nil, ErrInvalidData

}
