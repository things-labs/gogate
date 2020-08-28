package npi

const (
	/* 1st byte is the length of the data field, 2nd/3rd bytes are command field. */
	MT_RPC_FRAME_OVHD    = 2 // head and FCS
	MT_RPC_FRAME_HEAD_SZ = 1 // head size
	MT_RPC_FRAME_FCS_SZ  = 1 // FCS size

	MT_RPC_PDU_HDR_SZ  = 3 // length + cmd0 + cmd1
	MT_RPC_PDU_LEN_SZ  = 1 // length size
	MT_RPC_PDU_CMD0_SZ = 1 // cmd0 size
	MT_RPC_PDU_CMD1_SZ = 1 // cmd1 size

	/* position of fields in the general format frame */
	MT_RPC_POS_LEN  = 0
	MT_RPC_POS_CMD0 = 1
	MT_RPC_POS_CMD1 = 2
	MT_RPC_POS_DATA = 3

	// Start of frame character value
	MT_RPC_UART_SOF = 0xFE

	/* Maximum length of data in the general frame format. The upper limit is 255 because of the
	 * 1-byte length protocol. But the operation limit is lower for code size and ram savings so that
	 * the uart driver can use 256 byte rx/tx queues and so
	 * (MT_RPC_DATA_MAX + MT_RPC_FRAME_HDR_SZ + MT_UART_FRAME_OVHD) < 256
	 */
	MT_RPC_DATA_MAX = 250

	/* The 3 MSB's of the 1st command field byte are for command type. */
	MT_RPC_CMD_TYPE_MASK = 0xE000

	/* The 5 LSB's of the 1st command field byte are for the subsystem. */
	MT_RPC_SUBSYSTEM_MASK = 0x1FFF

	// type define
	MT_RPC_CMD_POLL = 0x0000
	MT_RPC_CMD_SREQ = 0x2000
	MT_RPC_CMD_AREQ = 0x4000
	MT_RPC_CMD_SRSP = 0x6000

	MT_RPC_SYS_SYS  = 0x0100
	MT_RPC_SYS_MAC  = 0x0200
	MT_RPC_SYS_NWK  = 0x0300
	MT_RPC_SYS_AF   = 0x0400
	MT_RPC_SYS_ZDO  = 0x0500
	MT_RPC_SYS_SAPI = 0x0600 /* Simple API. */
	MT_RPC_SYS_UTIL = 0x0700
	MT_RPC_SYS_DBG  = 0x0800
	MT_RPC_SYS_APP  = 0x0900 //0x09
	/***************/
	MT_RPC_SYS_OTA      = 0x0a00
	MT_RPC_SYS_ZNP      = 0x0b00
	MT_RPC_SYS_BOOT     = 0x0c00
	MT_RPC_SYS_UBL      = 0x0d00 //0x0d, 13 to be compatible with existing RemoTI.
	MT_RPC_SYS_RES14    = 0x0e00
	MT_RPC_SYS_APPCFG   = 0x0f00 // 0x0f APPconfig
	MT_RPC_SYS_RES16    = 0x1000
	MT_RPC_SYS_PROTOBUF = 0x1100
	MT_RPC_SYS_MAX      = 0x1200 /* Maximum value, must be last */
	MT_RPC_SYS_GP       = 0x1500 // GreenPower

	/* 18-32 available, not yet assigned. */
)
