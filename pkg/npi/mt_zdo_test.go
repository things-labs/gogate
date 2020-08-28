package npi

// func TestZdo_MgmtLeaveReq_Pack(t *testing.T) {
//     Convey("Zdo_MgmtLeaveReq_Pack", t, func() {
//         expect := &Npdu{
//             CmdId:  MT_ZDO_MGMT_LEAVE_REQ,
//             Data: []byte{0x34, 0x12, 0x78, 0x56, 0x34, 0x12, 0x78, 0x56, 0x34, 0x12, 0x01},
//         }
//         actual := &Npdu{}

//         actual.Zdo_MgmtLeaveReq_Pack(0x1234, 0x1234567812345678, 0x01)

//         So(reflect.DeepEqual(actual, expect), ShouldBeTrue)
//     })
// }

// func TestZdo_MgmtPermitJoinRspParse(t *testing.T) {
//     Convey("Zdo_MgmtPermitJoinRspParse", t, func() {
//         pdu := &Npdu{
//             CmdId:  MT_ZDO_MGMT_LEAVE_RSP,
//             Data: []byte{0x34, 0x12, 0x01},
//         }

//         expect := &ZdoMgmtPermitJionRsp_t{
//             SrcAddr: 0x1234,
//             Status:  0x01,
//         }

//         actual, err := pdu.Zdo_MgmtPermitJoinRspParse()

//         So(err, ShouldBeNil)
//         So(reflect.DeepEqual(actual, expect), ShouldBeTrue)
//     })
// }

// func TestZdo_StateChangeInd(t *testing.T) {
//     Convey("Zdo_StateChangeInd", t, func() {
//         pdu := &Npdu{
//             CmdId:  MT_ZDO_STATE_CHANGE_IND,
//             Data: []byte{0x34},
//         }

//         expect := byte(0x34)

//         actual, err := pdu.Zdo_StateChangeInd()

//         So(err, ShouldBeNil)
//         So(actual, ShouldEqual, expect)
//     })
// }

// func TestZdo_EnddeviceAnnceIndParse(t *testing.T) {
//     Convey("Zdo_EnddeviceAnnceIndParse", t, func() {
//         pdu := &Npdu{
//             CmdId:  MT_ZDO_END_DEV_ANNCE,
//             Data: []byte{0x34, 0x12, 0x78, 0x56, 0x78, 0x56, 0x34, 0x12, 0x78, 0x56, 0x34, 0x12, 0x01},
//         }

//         expect := &ZdoEnddeviceAnnceInd_t{
//             SrcAddr:      0x1234,
//             NwkAddr:      0x5678,
//             IeeeAddr:     0x1234567812345678,
//             Capabilities: 0x01,
//         }

//         actual, err := pdu.Zdo_EnddeviceAnnceIndParse()

//         So(err, ShouldBeNil)
//         So(reflect.DeepEqual(actual, expect), ShouldBeTrue)
//     })
// }

// func TestZdo_LeaveIndParse(t *testing.T) {
//     Convey("Zdo_LeaveIndParse", t, func() {
//         pdu := &Npdu{
//             CmdId:  MT_ZDO_LEAVE_IND,
//             Data: []byte{0x34, 0x12, 0x78, 0x56, 0x34, 0x12, 0x78, 0x56, 0x34, 0x12, 0x01, 0x00, 0x01},
//         }

//         expect := &ZdoLeaveInd_t{
//             SrcAddr: 0x1234,
//             ExtAddr: 0x1234567812345678,
//             Request: true,
//             Remove:  false,
//             Rejoin:  true,
//         }

//         actual, err := pdu.Zdo_LeaveIndParse()

//         So(err, ShouldBeNil)
//         So(reflect.DeepEqual(actual, expect), ShouldBeTrue)
//     })

// }
