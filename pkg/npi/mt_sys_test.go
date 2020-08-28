package npi

// func TestSys_ResetIndParse(t *testing.T) {
//     Convey("Sys_ResetIndParse", t, func() {
//         pdu := &Npdu{
//             CmdId:  MT_SYS_RESET_IND,
//             Data: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
//         }
//         expect := &SysResetInd_t{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}

//         actual, err := pdu.Sys_ResetIndParse()

//         So(err, ShouldBeNil)
//         So(reflect.DeepEqual(actual, expect), ShouldBeTrue)

//         pdu = &Npdu{
//             CmdId:  MT_SYS_RESET_IND,
//             Data: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
//         }

//         actual, err = pdu.Sys_ResetIndParse()

//         So(err, ShouldNotBeNil)
//         So(actual, ShouldBeNil)
//     })
// }
