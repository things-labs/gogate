package npi

const (
	/***************************************************************************************************
	 * MACROS
	 ***************************************************************************************************/

	/***************************************************************************************************
	 * SYS COMMANDS
	 ***************************************************************************************************/

	/* AREQ from host */
	MT_SYS_RESET_REQ = (MT_RPC_SYS_SYS | 0x00)

	/* SREQ/SRSP */
	MT_SYS_PING               = (MT_RPC_SYS_SYS | 0x01)
	MT_SYS_VERSION            = (MT_RPC_SYS_SYS | 0x02)
	MT_SYS_SET_EXTADDR        = (MT_RPC_SYS_SYS | 0x03)
	MT_SYS_GET_EXTADDR        = (MT_RPC_SYS_SYS | 0x04)
	MT_SYS_RAM_READ           = (MT_RPC_SYS_SYS | 0x05)
	MT_SYS_RAM_WRITE          = (MT_RPC_SYS_SYS | 0x06)
	MT_SYS_OSAL_NV_ITEM_INIT  = (MT_RPC_SYS_SYS | 0x07)
	MT_SYS_OSAL_NV_READ       = (MT_RPC_SYS_SYS | 0x08)
	MT_SYS_OSAL_NV_WRITE      = (MT_RPC_SYS_SYS | 0x09)
	MT_SYS_OSAL_START_TIMER   = (MT_RPC_SYS_SYS | 0x0A)
	MT_SYS_OSAL_STOP_TIMER    = (MT_RPC_SYS_SYS | 0x0B)
	MT_SYS_RANDOM             = (MT_RPC_SYS_SYS | 0x0C)
	MT_SYS_ADC_READ           = (MT_RPC_SYS_SYS | 0x0D)
	MT_SYS_GPIO               = (MT_RPC_SYS_SYS | 0x0E)
	MT_SYS_STACK_TUNE         = (MT_RPC_SYS_SYS | 0x0F)
	MT_SYS_SET_TIME           = (MT_RPC_SYS_SYS | 0x10)
	MT_SYS_GET_TIME           = (MT_RPC_SYS_SYS | 0x11)
	MT_SYS_OSAL_NV_DELETE     = (MT_RPC_SYS_SYS | 0x12)
	MT_SYS_OSAL_NV_LENGTH     = (MT_RPC_SYS_SYS | 0x13)
	MT_SYS_SET_TX_POWER       = (MT_RPC_SYS_SYS | 0x14)
	MT_SYS_JAMMER_PARAMETERS  = (MT_RPC_SYS_SYS | 0x15)
	MT_SYS_SNIFFER_PARAMETERS = (MT_RPC_SYS_SYS | 0x16)

	/* AREQ to host */
	MT_SYS_RESET_IND          = (MT_RPC_SYS_SYS | 0x80)
	MT_SYS_OSAL_TIMER_EXPIRED = (MT_RPC_SYS_SYS | 0x81)
	MT_SYS_JAMMER_IND         = (MT_RPC_SYS_SYS | 0x82)

	MT_SYS_RESET_HARD     = 0
	MT_SYS_RESET_SOFT     = 1
	MT_SYS_RESET_SHUTDOWN = 2

	MT_SYS_SNIFFER_DISABLE     = 0
	MT_SYS_SNIFFER_ENABLE      = 1
	MT_SYS_SNIFFER_GET_SETTING = 2

	/***************************************************************************************************
	 * MAC COMMANDS
	 ***************************************************************************************************/
	/* SREQ/SRSP */
	MT_MAC_RESET_REQ        = (MT_RPC_SYS_MAC | 0x01)
	MT_MAC_INIT             = (MT_RPC_SYS_MAC | 0x02)
	MT_MAC_START_REQ        = (MT_RPC_SYS_MAC | 0x03)
	MT_MAC_SYNC_REQ         = (MT_RPC_SYS_MAC | 0x04)
	MT_MAC_DATA_REQ         = (MT_RPC_SYS_MAC | 0x05)
	MT_MAC_ASSOCIATE_REQ    = (MT_RPC_SYS_MAC | 0x06)
	MT_MAC_DISASSOCIATE_REQ = (MT_RPC_SYS_MAC | 0x07)
	MT_MAC_GET_REQ          = (MT_RPC_SYS_MAC | 0x08)
	MT_MAC_SET_REQ          = (MT_RPC_SYS_MAC | 0x09)
	MT_MAC_GTS_REQ          = (MT_RPC_SYS_MAC | 0x0a)
	MT_MAC_RX_ENABLE_REQ    = (MT_RPC_SYS_MAC | 0x0b)
	MT_MAC_SCAN_REQ         = (MT_RPC_SYS_MAC | 0x0c)
	MT_MAC_POLL_REQ         = (MT_RPC_SYS_MAC | 0x0d)
	MT_MAC_PURGE_REQ        = (MT_RPC_SYS_MAC | 0x0e)
	MT_MAC_SET_RX_GAIN_REQ  = (MT_RPC_SYS_MAC | 0x0f)

	/* Security PIB SREQ/SRSP */
	MT_MAC_SECURITY_GET_REQ = (MT_RPC_SYS_MAC | 0x10)
	MT_MAC_SECURITY_SET_REQ = (MT_RPC_SYS_MAC | 0x11)

	/* Enhanced Active Scan request */
	MT_MAC_ENHANCED_ACTIVE_SCAN_REQ = (MT_RPC_SYS_MAC | 0x12)
	MT_MAC_ENHANCED_ACTIVE_SCAN_RSP = (MT_RPC_SYS_MAC | 0x13)

	/* Enhanced MAC interface SREQ/SRSP */
	MT_MAC_SRC_MATCH_ENABLE       = (MT_RPC_SYS_MAC | 0x14)
	MT_MAC_SRC_MATCH_ADD_ENTRY    = (MT_RPC_SYS_MAC | 0x15)
	MT_MAC_SRC_MATCH_DELETE_ENTRY = (MT_RPC_SYS_MAC | 0x16)
	MT_MAC_SRC_MATCH_ACK_ALL      = (MT_RPC_SYS_MAC | 0x17)
	MT_MAC_SRC_CHECK_ALL          = (MT_RPC_SYS_MAC | 0x18)

	/* AREQ from Host */
	MT_MAC_ASSOCIATE_RSP = (MT_RPC_SYS_MAC | 0x50)
	MT_MAC_ORPHAN_RSP    = (MT_RPC_SYS_MAC | 0x51)

	/* AREQ to host */
	MT_MAC_SYNC_LOSS_IND     = (MT_RPC_SYS_MAC | 0x80)
	MT_MAC_ASSOCIATE_IND     = (MT_RPC_SYS_MAC | 0x81)
	MT_MAC_ASSOCIATE_CNF     = (MT_RPC_SYS_MAC | 0x82)
	MT_MAC_BEACON_NOTIFY_IND = (MT_RPC_SYS_MAC | 0x83)
	MT_MAC_DATA_CNF          = (MT_RPC_SYS_MAC | 0x84)
	MT_MAC_DATA_IND          = (MT_RPC_SYS_MAC | 0x85)
	MT_MAC_DISASSOCIATE_IND  = (MT_RPC_SYS_MAC | 0x86)
	MT_MAC_DISASSOCIATE_CNF  = (MT_RPC_SYS_MAC | 0x87)
	MT_MAC_GTS_CNF           = (MT_RPC_SYS_MAC | 0x88)
	MT_MAC_GTS_IND           = (MT_RPC_SYS_MAC | 0x89)
	MT_MAC_ORPHAN_IND        = (MT_RPC_SYS_MAC | 0x8a)
	MT_MAC_POLL_CNF          = (MT_RPC_SYS_MAC | 0x8b)
	MT_MAC_SCAN_CNF          = (MT_RPC_SYS_MAC | 0x8c)
	MT_MAC_COMM_STATUS_IND   = (MT_RPC_SYS_MAC | 0x8d)
	MT_MAC_START_CNF         = (MT_RPC_SYS_MAC | 0x8e)
	MT_MAC_RX_ENABLE_CNF     = (MT_RPC_SYS_MAC | 0x8f)
	MT_MAC_PURGE_CNF         = (MT_RPC_SYS_MAC | 0x90)
	MT_MAC_POLL_IND          = (MT_RPC_SYS_MAC | 0x91)

	/***************************************************************************************************
	 * NWK COMMANDS
	 ***************************************************************************************************/

	/* AREQ from host */
	MT_NWK_INIT = (MT_RPC_SYS_NWK | 0x00)

	/* SREQ/SRSP */
	MT_NLDE_DATA_REQ              = (MT_RPC_SYS_NWK | 0x01)
	MT_NLME_NETWORK_FORMATION_REQ = (MT_RPC_SYS_NWK | 0x02)
	MT_NLME_PERMIT_JOINING_REQ    = (MT_RPC_SYS_NWK | 0x03)
	MT_NLME_JOIN_REQ              = (MT_RPC_SYS_NWK | 0x04)
	MT_NLME_LEAVE_REQ             = (MT_RPC_SYS_NWK | 0x05)
	MT_NLME_RESET_REQ             = (MT_RPC_SYS_NWK | 0x06)
	MT_NLME_GET_REQ               = (MT_RPC_SYS_NWK | 0x07)
	MT_NLME_SET_REQ               = (MT_RPC_SYS_NWK | 0x08)
	MT_NLME_NETWORK_DISCOVERY_REQ = (MT_RPC_SYS_NWK | 0x09)
	MT_NLME_ROUTE_DISCOVERY_REQ   = (MT_RPC_SYS_NWK | 0x0A)
	MT_NLME_DIRECT_JOIN_REQ       = (MT_RPC_SYS_NWK | 0x0B)
	MT_NLME_ORPHAN_JOIN_REQ       = (MT_RPC_SYS_NWK | 0x0C)
	MT_NLME_START_ROUTER_REQ      = (MT_RPC_SYS_NWK | 0x0D)

	/* AREQ to host */
	MT_NLDE_DATA_CONF              = (MT_RPC_SYS_NWK | 0x80)
	MT_NLDE_DATA_IND               = (MT_RPC_SYS_NWK | 0x81)
	MT_NLME_NETWORK_FORMATION_CONF = (MT_RPC_SYS_NWK | 0x82)
	MT_NLME_JOIN_CONF              = (MT_RPC_SYS_NWK | 0x83)
	MT_NLME_JOIN_IND               = (MT_RPC_SYS_NWK | 0x84)
	MT_NLME_LEAVE_CONF             = (MT_RPC_SYS_NWK | 0x85)
	MT_NLME_LEAVE_IND              = (MT_RPC_SYS_NWK | 0x86)
	MT_NLME_POLL_CONF              = (MT_RPC_SYS_NWK | 0x87)
	MT_NLME_SYNC_IND               = (MT_RPC_SYS_NWK | 0x88)
	MT_NLME_NETWORK_DISCOVERY_CONF = (MT_RPC_SYS_NWK | 0x89)
	MT_NLME_START_ROUTER_CONF      = (MT_RPC_SYS_NWK | 0x8A)

	/***************************************************************************************************
	 * AF COMMANDS
	 ***************************************************************************************************/

	/* SREQ/SRSP */
	MT_AF_REGISTER            = (MT_RPC_SYS_AF | 0x00)
	MT_AF_DATA_REQUEST        = (MT_RPC_SYS_AF | 0x01) /* AREQ optional, but no AREQ response. */
	MT_AF_DATA_REQUEST_EXT    = (MT_RPC_SYS_AF | 0x02) /* AREQ optional, but no AREQ response. */
	MT_AF_DATA_REQUEST_SRCRTG = (MT_RPC_SYS_AF | 0x03)
	MT_AF_DELETE              = (MT_RPC_SYS_AF | 0x04)

	MT_AF_INTER_PAN_CTL   = (MT_RPC_SYS_AF | 0x10)
	MT_AF_DATA_STORE      = (MT_RPC_SYS_AF | 0x11)
	MT_AF_DATA_RETRIEVE   = (MT_RPC_SYS_AF | 0x12)
	MT_AF_APSF_CONFIG_SET = (MT_RPC_SYS_AF | 0x13)
	MT_AF_APSF_CONFIG_GET = (MT_RPC_SYS_AF | 0x14)

	/* AREQ to host */
	MT_AF_DATA_CONFIRM     = (MT_RPC_SYS_AF | 0x80)
	MT_AF_INCOMING_MSG     = (MT_RPC_SYS_AF | 0x81)
	MT_AF_INCOMING_MSG_EXT = (MT_RPC_SYS_AF | 0x82)
	MT_AF_REFLECT_ERROR    = (MT_RPC_SYS_AF | 0x83)

	/***************************************************************************************************
	 * ZDO COMMANDS
	 ***************************************************************************************************/

	/* SREQ/SRSP */
	MT_ZDO_NWK_ADDR_REQ           = (MT_RPC_SYS_ZDO | 0x00)
	MT_ZDO_IEEE_ADDR_REQ          = (MT_RPC_SYS_ZDO | 0x01)
	MT_ZDO_NODE_DESC_REQ          = (MT_RPC_SYS_ZDO | 0x02)
	MT_ZDO_POWER_DESC_REQ         = (MT_RPC_SYS_ZDO | 0x03)
	MT_ZDO_SIMPLE_DESC_REQ        = (MT_RPC_SYS_ZDO | 0x04)
	MT_ZDO_ACTIVE_EP_REQ          = (MT_RPC_SYS_ZDO | 0x05)
	MT_ZDO_MATCH_DESC_REQ         = (MT_RPC_SYS_ZDO | 0x06)
	MT_ZDO_COMPLEX_DESC_REQ       = (MT_RPC_SYS_ZDO | 0x07)
	MT_ZDO_USER_DESC_REQ          = (MT_RPC_SYS_ZDO | 0x08)
	MT_ZDO_END_DEV_ANNCE          = (MT_RPC_SYS_ZDO | 0x0A)
	MT_ZDO_USER_DESC_SET          = (MT_RPC_SYS_ZDO | 0x0B)
	MT_ZDO_SERVICE_DISC_REQ       = (MT_RPC_SYS_ZDO | 0x0C)
	MT_ZDO_END_DEVICE_TIMEOUT_REQ = (MT_RPC_SYS_ZDO | 0x0D)
	MT_ZDO_END_DEV_BIND_REQ       = (MT_RPC_SYS_ZDO | 0x20)
	MT_ZDO_BIND_REQ               = (MT_RPC_SYS_ZDO | 0x21)
	MT_ZDO_UNBIND_REQ             = (MT_RPC_SYS_ZDO | 0x22)

	MT_ZDO_SET_LINK_KEY             = (MT_RPC_SYS_ZDO | 0x23)
	MT_ZDO_REMOVE_LINK_KEY          = (MT_RPC_SYS_ZDO | 0x24)
	MT_ZDO_GET_LINK_KEY             = (MT_RPC_SYS_ZDO | 0x25)
	MT_ZDO_NWK_DISCOVERY_REQ        = (MT_RPC_SYS_ZDO | 0x26)
	MT_ZDO_JOIN_REQ                 = (MT_RPC_SYS_ZDO | 0x27)
	MT_ZDO_SEND_DATA                = (MT_RPC_SYS_ZDO | 0x28)
	MT_ZDO_NWK_ADDR_OF_INTEREST_REQ = (MT_RPC_SYS_ZDO | 0x29)

	MT_ZDO_MGMT_NWKDISC_REQ     = (MT_RPC_SYS_ZDO | 0x30)
	MT_ZDO_MGMT_LQI_REQ         = (MT_RPC_SYS_ZDO | 0x31)
	MT_ZDO_MGMT_RTG_REQ         = (MT_RPC_SYS_ZDO | 0x32)
	MT_ZDO_MGMT_BIND_REQ        = (MT_RPC_SYS_ZDO | 0x33)
	MT_ZDO_MGMT_LEAVE_REQ       = (MT_RPC_SYS_ZDO | 0x34)
	MT_ZDO_MGMT_DIRECT_JOIN_REQ = (MT_RPC_SYS_ZDO | 0x35)
	MT_ZDO_MGMT_PERMIT_JOIN_REQ = (MT_RPC_SYS_ZDO | 0x36)
	MT_ZDO_MGMT_NWK_UPDATE_REQ  = (MT_RPC_SYS_ZDO | 0x37)

	/* AREQ optional, but no AREQ response. */
	MT_ZDO_MSG_CB_REGISTER  = (MT_RPC_SYS_ZDO | 0x3E)
	MT_ZDO_MSG_CB_REMOVE    = (MT_RPC_SYS_ZDO | 0x3F)
	MT_ZDO_STARTUP_FROM_APP = (MT_RPC_SYS_ZDO | 0x40)

	/* AREQ from host */
	MT_ZDO_AUTO_FIND_DESTINATION_REQ    = (MT_RPC_SYS_ZDO | 0x41)
	MT_ZDO_SEC_ADD_LINK_KEY             = (MT_RPC_SYS_ZDO | 0x42)
	MT_ZDO_SEC_ENTRY_LOOKUP_EXT         = (MT_RPC_SYS_ZDO | 0x43)
	MT_ZDO_SEC_DEVICE_REMOVE            = (MT_RPC_SYS_ZDO | 0x44)
	MT_ZDO_EXT_ROUTE_DISC               = (MT_RPC_SYS_ZDO | 0x45)
	MT_ZDO_EXT_ROUTE_CHECK              = (MT_RPC_SYS_ZDO | 0x46)
	MT_ZDO_EXT_REMOVE_GROUP             = (MT_RPC_SYS_ZDO | 0x47)
	MT_ZDO_EXT_REMOVE_ALL_GROUP         = (MT_RPC_SYS_ZDO | 0x48)
	MT_ZDO_EXT_FIND_ALL_GROUPS_ENDPOINT = (MT_RPC_SYS_ZDO | 0x49)
	MT_ZDO_EXT_FIND_GROUP               = (MT_RPC_SYS_ZDO | 0x4A)
	MT_ZDO_EXT_ADD_GROUP                = (MT_RPC_SYS_ZDO | 0x4B)
	MT_ZDO_EXT_COUNT_ALL_GROUPS         = (MT_RPC_SYS_ZDO | 0x4C)
	MT_ZDO_EXT_RX_IDLE                  = (MT_RPC_SYS_ZDO | 0x4D)
	MT_ZDO_EXT_UPDATE_NWK_KEY           = (MT_RPC_SYS_ZDO | 0x4E)
	MT_ZDO_EXT_SWITCH_NWK_KEY           = (MT_RPC_SYS_ZDO | 0x4F)
	MT_ZDO_EXT_NWK_INFO                 = (MT_RPC_SYS_ZDO | 0x50)
	MT_ZDO_EXT_SEC_APS_REMOVE_REQ       = (MT_RPC_SYS_ZDO | 0x51)
	MT_ZDO_FORCE_CONCENTRATOR_CHANGE    = (MT_RPC_SYS_ZDO | 0x52)
	MT_ZDO_EXT_SET_PARAMS               = (MT_RPC_SYS_ZDO | 0x53)

	/* AREQ to host */
	MT_ZDO_AREQ_TO_HOST    = (MT_RPC_SYS_ZDO | 0x80) /* Mark the start of the ZDO CId AREQs to host. */
	MT_ZDO_NWK_ADDR_RSP    = (MT_RPC_SYS_ZDO | 0x80) // =((uint8)NWK_addr_req | 0x80)
	MT_ZDO_IEEE_ADDR_RSP   = (MT_RPC_SYS_ZDO | 0x81) // =((uint8)IEEE_addr_req | 0x80)
	MT_ZDO_NODE_DESC_RSP   = (MT_RPC_SYS_ZDO | 0x82) // =((uint8)Node_Desc_req | 0x80)
	MT_ZDO_POWER_DESC_RSP  = (MT_RPC_SYS_ZDO | 0x83) // =((uint8)Power_Desc_req | 0x80)
	MT_ZDO_SIMPLE_DESC_RSP = (MT_RPC_SYS_ZDO | 0x84) // =((uint8)Simple_Desc_req | 0x80)
	MT_ZDO_ACTIVE_EP_RSP   = (MT_RPC_SYS_ZDO | 0x85) // =((uint8)Active_EP_req | 0x80)
	MT_ZDO_MATCH_DESC_RSP  = (MT_RPC_SYS_ZDO | 0x86) // =((uint8)Match_Desc_req | 0x80)

	MT_ZDO_COMPLEX_DESC_RSP = (MT_RPC_SYS_ZDO | 0x87)
	MT_ZDO_USER_DESC_RSP    = (MT_RPC_SYS_ZDO | 0x88)
	//                                        /* =( MT_RPC_SYS_ZDO | 0x92 )*/ =((uint8)Discovery_Cache_req | 0x80)
	MT_ZDO_USER_DESC_CONF  = (MT_RPC_SYS_ZDO | 0x89)
	MT_ZDO_SERVER_DISC_RSP = (MT_RPC_SYS_ZDO | 0x8A)

	MT_ZDO_END_DEVICE_TIMEOUT_RSP = (MT_RPC_SYS_ZDO | 0x9F) //  =((uint8)End_Device_Timeout_req | 0x80)

	MT_ZDO_END_DEVICE_BIND_RSP = (MT_RPC_SYS_ZDO | 0xA0) // =((uint8)End_Device_Bind_req | 0x80)
	MT_ZDO_BIND_RSP            = (MT_RPC_SYS_ZDO | 0xA1) // =((uint8)Bind_req | 0x80)
	MT_ZDO_UNBIND_RSP          = (MT_RPC_SYS_ZDO | 0xA2) // =((uint8)Unbind_req | 0x80)

	MT_ZDO_MGMT_NWK_DISC_RSP    = (MT_RPC_SYS_ZDO | 0xB0) // =((uint8)Mgmt_NWK_Disc_req | 0x80)
	MT_ZDO_MGMT_LQI_RSP         = (MT_RPC_SYS_ZDO | 0xB1) // =((uint8)Mgmt_Lqi_req | 0x80)
	MT_ZDO_MGMT_RTG_RSP         = (MT_RPC_SYS_ZDO | 0xB2) // =((uint8)Mgmt_Rtg_req | 0x80)
	MT_ZDO_MGMT_BIND_RSP        = (MT_RPC_SYS_ZDO | 0xB3) // =((uint8)Mgmt_Bind_req | 0x80)
	MT_ZDO_MGMT_LEAVE_RSP       = (MT_RPC_SYS_ZDO | 0xB4) // =((uint8)Mgmt_Leave_req | 0x80)
	MT_ZDO_MGMT_DIRECT_JOIN_RSP = (MT_RPC_SYS_ZDO | 0xB5) // =((uint8)Mgmt_Direct_Join_req | 0x80)
	MT_ZDO_MGMT_PERMIT_JOIN_RSP = (MT_RPC_SYS_ZDO | 0xB6) // =((uint8)Mgmt_Permit_Join_req | 0x80)

	//                                        /* =( MT_RPC_SYS_ZDO | 0xB8 ) */ =((uint8)Mgmt_NWK_Update_req | 0x80)

	MT_ZDO_STATE_CHANGE_IND     = (MT_RPC_SYS_ZDO | 0xC0)
	MT_ZDO_END_DEVICE_ANNCE_IND = (MT_RPC_SYS_ZDO | 0xC1)
	MT_ZDO_MATCH_DESC_RSP_SENT  = (MT_RPC_SYS_ZDO | 0xC2)
	MT_ZDO_STATUS_ERROR_RSP     = (MT_RPC_SYS_ZDO | 0xC3)
	MT_ZDO_SRC_RTG_IND          = (MT_RPC_SYS_ZDO | 0xC4)
	MT_ZDO_BEACON_NOTIFY_IND    = (MT_RPC_SYS_ZDO | 0xC5)
	MT_ZDO_JOIN_CNF             = (MT_RPC_SYS_ZDO | 0xC6)
	MT_ZDO_NWK_DISCOVERY_CNF    = (MT_RPC_SYS_ZDO | 0xC7)
	MT_ZDO_CONCENTRATOR_IND_CB  = (MT_RPC_SYS_ZDO | 0xC8)
	MT_ZDO_LEAVE_IND            = (MT_RPC_SYS_ZDO | 0xC9)
	MT_ZDO_TC_DEVICE_IND        = (MT_RPC_SYS_ZDO | 0xCA)
	MT_ZDO_PERMIT_JOIN_IND      = (MT_RPC_SYS_ZDO | 0xCB)

	MT_ZDO_MSG_CB_INCOMING = (MT_RPC_SYS_ZDO | 0xFF)

	// Some arbitrarily chosen value for a default error status msg.
	MtZdoDef_rsp = (MT_RPC_SYS_ZDO | 0x0040)

	/***************************************************************************************************
	 * SAPI COMMANDS
	 ***************************************************************************************************/

	// SAPI MT Command Identifiers
	/* AREQ from Host */
	MT_SAPI_SYS_RESET = (MT_RPC_SYS_SAPI | 0x09)
	MT_SAPI_START_NWK = (MT_RPC_SYS_SAPI | 0x5a) // add by mo

	/* SREQ/SRSP */
	MT_SAPI_START_REQ        = (MT_RPC_SYS_SAPI | 0x00)
	MT_SAPI_BIND_DEVICE_REQ  = (MT_RPC_SYS_SAPI | 0x01)
	MT_SAPI_ALLOW_BIND_REQ   = (MT_RPC_SYS_SAPI | 0x02)
	MT_SAPI_SEND_DATA_REQ    = (MT_RPC_SYS_SAPI | 0x03)
	MT_SAPI_READ_CFG_REQ     = (MT_RPC_SYS_SAPI | 0x04)
	MT_SAPI_WRITE_CFG_REQ    = (MT_RPC_SYS_SAPI | 0x05)
	MT_SAPI_GET_DEV_INFO_REQ = (MT_RPC_SYS_SAPI | 0x06)
	MT_SAPI_FIND_DEV_REQ     = (MT_RPC_SYS_SAPI | 0x07)
	MT_SAPI_PMT_JOIN_REQ     = (MT_RPC_SYS_SAPI | 0x08)
	MT_SAPI_APP_REGISTER_REQ = (MT_RPC_SYS_SAPI | 0x0a)

	MT_SAPI_RESET_FACTORY        = (MT_RPC_SYS_SAPI | 0x70) // add by mo
	MT_SAPI_GET_DEV_ALL_INFO_REQ = (MT_RPC_SYS_SAPI | 0x71) // add by mo
	MT_SAPI_WRITE_LOGICAL_TYPE   = (MT_RPC_SYS_SAPI | 0x72) // add by mo
	/* AREQ to host */
	MT_SAPI_START_CNF      = (MT_RPC_SYS_SAPI | 0x80)
	MT_SAPI_BIND_CNF       = (MT_RPC_SYS_SAPI | 0x81)
	MT_SAPI_ALLOW_BIND_CNF = (MT_RPC_SYS_SAPI | 0x82)
	MT_SAPI_SEND_DATA_CNF  = (MT_RPC_SYS_SAPI | 0x83)
	MT_SAPI_READ_CFG_RSP   = (MT_RPC_SYS_SAPI | 0x84)
	MT_SAPI_FIND_DEV_CNF   = (MT_RPC_SYS_SAPI | 0x85)
	MT_SAPI_DEV_INFO_RSP   = (MT_RPC_SYS_SAPI | 0x86)
	MT_SAPI_RCV_DATA_IND   = (MT_RPC_SYS_SAPI | 0x87)

	/***************************************************************************************************
	 * UTIL COMMANDS
	 ***************************************************************************************************/
	/* SREQ/SRSP: */
	MT_UTIL_GET_DEVICE_INFO  = (MT_RPC_SYS_UTIL | 0x00)
	MT_UTIL_GET_NV_INFO      = (MT_RPC_SYS_UTIL | 0x01)
	MT_UTIL_SET_PANID        = (MT_RPC_SYS_UTIL | 0x02)
	MT_UTIL_SET_CHANNELS     = (MT_RPC_SYS_UTIL | 0x03)
	MT_UTIL_SET_SECLEVEL     = (MT_RPC_SYS_UTIL | 0x04)
	MT_UTIL_SET_PRECFGKEY    = (MT_RPC_SYS_UTIL | 0x05)
	MT_UTIL_CALLBACK_SUB_CMD = (MT_RPC_SYS_UTIL | 0x06)
	MT_UTIL_KEY_EVENT        = (MT_RPC_SYS_UTIL | 0x07)
	MT_UTIL_TIME_ALIVE       = (MT_RPC_SYS_UTIL | 0x09)
	MT_UTIL_LED_CONTROL      = (MT_RPC_SYS_UTIL | 0x0A)

	MT_UTIL_TEST_LOOPBACK = (MT_RPC_SYS_UTIL | 0x10)
	MT_UTIL_DATA_REQ      = (MT_RPC_SYS_UTIL | 0x11)

	MT_UTIL_GPIO_SET_DIRECTION = (MT_RPC_SYS_UTIL | 0x14)
	MT_UTIL_GPIO_READ          = (MT_RPC_SYS_UTIL | 0x15)
	MT_UTIL_GPIO_WRITE         = (MT_RPC_SYS_UTIL | 0x16)

	MT_UTIL_SRC_MATCH_ENABLE            = (MT_RPC_SYS_UTIL | 0x20)
	MT_UTIL_SRC_MATCH_ADD_ENTRY         = (MT_RPC_SYS_UTIL | 0x21)
	MT_UTIL_SRC_MATCH_DEL_ENTRY         = (MT_RPC_SYS_UTIL | 0x22)
	MT_UTIL_SRC_MATCH_CHECK_SRC_ADDR    = (MT_RPC_SYS_UTIL | 0x23)
	MT_UTIL_SRC_MATCH_ACK_ALL_PENDING   = (MT_RPC_SYS_UTIL | 0x24)
	MT_UTIL_SRC_MATCH_CHECK_ALL_PENDING = (MT_RPC_SYS_UTIL | 0x25)

	MT_UTIL_ADDRMGR_EXT_ADDR_LOOKUP  = (MT_RPC_SYS_UTIL | 0x40)
	MT_UTIL_ADDRMGR_NWK_ADDR_LOOKUP  = (MT_RPC_SYS_UTIL | 0x41)
	MT_UTIL_APSME_LINK_KEY_DATA_GET  = (MT_RPC_SYS_UTIL | 0x44)
	MT_UTIL_APSME_LINK_KEY_NV_ID_GET = (MT_RPC_SYS_UTIL | 0x45)
	MT_UTIL_ASSOC_COUNT              = (MT_RPC_SYS_UTIL | 0x48)
	MT_UTIL_ASSOC_FIND_DEVICE        = (MT_RPC_SYS_UTIL | 0x49)
	MT_UTIL_ASSOC_GET_WITH_ADDRESS   = (MT_RPC_SYS_UTIL | 0x4A)
	MT_UTIL_APSME_REQUEST_KEY_CMD    = (MT_RPC_SYS_UTIL | 0x4B)

	// MT_UTIL_SRNG_GENERATE                =( MT_RPC_SYS_UTIL | 0x4C )

	MT_UTIL_BIND_ADD_ENTRY = (MT_RPC_SYS_UTIL | 0x4D)

	MT_UTIL_ZCL_KEY_EST_INIT_EST = (MT_RPC_SYS_UTIL | 0x80)
	MT_UTIL_ZCL_KEY_EST_SIGN     = (MT_RPC_SYS_UTIL | 0x81)

	/* AREQ from/to host */
	MT_UTIL_SYNC_REQ              = (MT_RPC_SYS_UTIL | 0xE0)
	MT_UTIL_ZCL_KEY_ESTABLISH_IND = (MT_RPC_SYS_UTIL | 0xE1)

	/***************************************************************************************************
	 * DEBUG COMMANDS
	 ***************************************************************************************************/

	/* SREQ/SRSP: */
	MT_DEBUG_SET_THRESHOLD = (MT_RPC_SYS_DBG | 0x00)

	MT_DEBUG_MAC_DATA_DUMP = (MT_RPC_SYS_DBG | 0x10)

	/* AREQ */
	MT_DEBUG_MSG = (MT_RPC_SYS_DBG | 0x80)

	/***************************************************************************************************
	 * APP CONFIG COMMANDS
	 ***************************************************************************************************/

	MT_APP_CNF_SET_DEFAULT_REMOTE_ENDDEVICE_TIMEOUT   = (MT_RPC_SYS_APPCFG | 0x01)
	MT_APP_CNF_SET_ENDDEVICETIMEOUT                   = (MT_RPC_SYS_APPCFG | 0x02)
	MT_APP_CNF_SET_ALLOWREJOIN_TC_POLICY              = (MT_RPC_SYS_APPCFG | 0x03)
	MT_APP_CNF_BDB_ADD_INSTALLCODE                    = (MT_RPC_SYS_APPCFG | 0x04)
	MT_APP_CNF_BDB_START_COMMISSIONING                = (MT_RPC_SYS_APPCFG | 0x05)
	MT_APP_CNF_BDB_SET_JOINUSESINSTALLCODEKEY         = (MT_RPC_SYS_APPCFG | 0x06)
	MT_APP_CNF_BDB_SET_ACTIVE_DEFAULT_CENTRALIZED_KEY = (MT_RPC_SYS_APPCFG | 0x07)
	MT_APP_CNF_BDB_SET_CHANNEL                        = (MT_RPC_SYS_APPCFG | 0x08)
	MT_APP_CNF_BDB_SET_TC_REQUIRE_KEY_EXCHANGE        = (MT_RPC_SYS_APPCFG | 0x09)
	MT_APP_CNF_BDB_ZED_ATTEMPT_RECOVER_NWK            = (MT_RPC_SYS_APPCFG | 0x0A)

	MT_APP_CNF_BDB_COMMISSIONING_NOTIFICATION = (MT_RPC_SYS_APPCFG | 0x80)
	//Application debug commands
	MT_APP_CNF_SET_NWK_FRAME_COUNTER = (MT_RPC_SYS_APPCFG | 0xFF)

	/***************************************************************************************************
	 * GP COMMANDS
	 ***************************************************************************************************/
	//From Host to ZNP
	MT_GP_DATA_REQ = (MT_RPC_SYS_GP | 0x01)
	MT_GP_SEC_RSP  = (MT_RPC_SYS_GP | 0x02)

	//From ZNP to Host
	MT_GP_SEC_REQ  = (MT_RPC_SYS_GP | 0x03)
	MT_GP_DATA_IND = (MT_RPC_SYS_GP | 0x04)
	MT_GP_DATA_CNF = (MT_RPC_SYS_GP | 0x05)

	/***************************************************************************************************
	 * APP COMMANDS
	 ***************************************************************************************************/

	/* SREQ/SRSP: */
	MT_APP_MSG       = (MT_RPC_SYS_APP | 0x00)
	MT_APP_USER_TEST = (MT_RPC_SYS_APP | 0x01)

	/* SRSP */
	MT_APP_RSP        = (MT_RPC_SYS_APP | 0x80)
	MT_APP_ZLL_TL_IND = (MT_RPC_SYS_APP | 0x81)

	/***************************************************************************************************
	 * FILE SYSTEM COMMANDS
	 ***************************************************************************************************/
	MT_OTA_FILE_READ_REQ = (MT_RPC_SYS_OTA | 0x00)
	MT_OTA_NEXT_IMG_REQ  = (MT_RPC_SYS_OTA | 0x01)

	MT_OTA_FILE_READ_RSP = (MT_RPC_SYS_OTA | 0x80)
	MT_OTA_NEXT_IMG_RSP  = (MT_RPC_SYS_OTA | 0x81)
	MT_OTA_STATUS_IND    = (MT_RPC_SYS_OTA | 0x82)

	/***************************************************************************************************
	 * UBL COMMANDS
	 ***************************************************************************************************/

	/* SREQ/SRSP: */

	//efine SB_TGT_BOOTLOAD                       =( MT_RPC_SYS_UBL | 0x10 )// Optional command in Target Application.
	MT_UBL_ERASE = (MT_RPC_SYS_UBL | 0x10)

	MT_UBL_WRITE     = (MT_RPC_SYS_UBL | 0x11)
	MT_UBL_READ      = (MT_RPC_SYS_UBL | 0x12)
	MT_UBL_ENABLE    = (MT_RPC_SYS_UBL | 0x13)
	MT_UBL_HANDSHAKE = (MT_RPC_SYS_UBL | 0x14)
	MT_UBL_FORCE     = (MT_RPC_SYS_UBL | 0x15)

	/***************************************************************************************************
	 * ZNP COMMANDS
	 ***************************************************************************************************/

	MT_ZNP_BASIC_CFG = (MT_RPC_SYS_ZNP | 0x00)
	MT_ZNP_ZCL_CFG   = (MT_RPC_SYS_ZNP | 0x10)
	MT_ZNP_SE_CFG    = (MT_RPC_SYS_ZNP | 0x20)

	/***************************************************************************************************
	 * ZNP RESPONSES
	 ***************************************************************************************************/

	MT_ZNP_BASIC_RSP = (MT_RPC_SYS_ZNP | 0x80)
	MT_ZNP_ZCL_RSP   = (MT_RPC_SYS_ZNP | 0x90)
	MT_ZNP_SE_RSP    = (MT_RPC_SYS_ZNP | 0xA0)

	/* Capabilities - PING Response */
	/* Capabilities - PING Response */
	MT_CAP_SYS   = 0x0001
	MT_CAP_MAC   = 0x0002
	MT_CAP_NWK   = 0x0004
	MT_CAP_AF    = 0x0008
	MT_CAP_ZDO   = 0x0010
	MT_CAP_SAPI  = 0x0020
	MT_CAP_UTIL  = 0x0040
	MT_CAP_DEBUG = 0x0080
	MT_CAP_APP   = 0x0100
	MT_CAP_ZOAD  = 0x1000
)
