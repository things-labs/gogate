package ltlspec

/*****************************************************************************/
/***    光照测量集id                                                            ***/
/*****************************************************************************/
// Illuminance Measurement Information attribute set
const (
	ATTRID_MS_ILLUMINANCE_MEASURED_VALUE = iota
	ATTRID_MS_ILLUMINANCE_MIN_MEASURED_VALUE
	ATTRID_MS_ILLUMINANCE_MAX_MEASURED_VALUE
	ATTRID_MS_ILLUMINANCE_TOLERANCE

	// 无效值定义
	MS_MS_ILLUMINANCE_INVALID_VALUE = 0xffff
)

/*****************************************************************************/
/***    光照水平感知配置集                                                          ***/
/*****************************************************************************/
const (
	// Illuminance Level Sensing Information attribute set
	ATTRID_MS_ILLUMINANCE_LEVEL_STATUS = 0x0000
	/***  Level Status attribute values  ***/
	MS_ILLUMINANCE_LEVEL_ON_TARGET    = 0x00
	MS_ILLUMINANCE_LEVEL_BELOW_TARGET = 0x01
	MS_ILLUMINANCE_LEVEL_ABOVE_TARGET = 0x02
	// Illuminance Level Sensing Settings attribute set
	ATTRID_MS_ILLUMINANCE_TARGET_LEVEL = 0x0010
)

/*****************************************************************************/
/***    温度测量集                                                              ***/
/*****************************************************************************/
const (
	// Temperature Measurement Information attributes set
	ATTRID_MS_TEMPERATURE_MEASURED_VALUE = iota
	ATTRID_MS_TEMPERATURE_MIN_MEASURED_VALUE
	ATTRID_MS_TEMPERATURE_MAX_MEASURED_VALUE
	ATTRID_MS_TEMPERATURE_TOLERANCE

	// 无效值定义
	MS_TEMPERATURE_INVALID_VALUE = 0x8000
)

/*****************************************************************************/
/***    压力测量集                                                              ***/
/*****************************************************************************/
const (
	// Pressure Measurement Information attribute set
	ATTRID_MS_PRESSURE_MEASUREMENT_MEASURED_VALUE = iota
	ATTRID_MS_PRESSURE_MEASUREMENT_MIN_MEASURED_VALUE
	ATTRID_MS_PRESSURE_MEASUREMENT_MAX_MEASURED_VALUE
	ATTRID_MS_PRESSURE_MEASUREMENT_TOLERANCE
	// 无效值定义
	MS_PRESSURE_MEASUREMENT_INVALID_VALUE = 0x8000
)

/*****************************************************************************/
/***        流量测量集                                                          ***/
/*****************************************************************************/
const (
	// Flow Measurement Information attribute set
	ATTRID_MS_FLOW_MEASUREMENT_MEASURED_VALUE = iota
	ATTRID_MS_FLOW_MEASUREMENT_MIN_MEASURED_VALUE
	ATTRID_MS_FLOW_MEASUREMENT_MAX_MEASURED_VALUE
	ATTRID_MS_FLOW_MEASUREMENT_TOLERANCE
	// 无效值定义
	MS_FLOW_MEASUREMENT_INVALID_VALUE = 0xffff
)

/*****************************************************************************/
/***        相对湿度测量集                                                        ***/
/*****************************************************************************/
const (
	// Relative Humidity Information attribute set
	ATTRID_MS_RELATIVE_HUMIDITY_MEASURED_VALUE = iota
	ATTRID_MS_RELATIVE_HUMIDITY_MIN_MEASURED_VALUE
	ATTRID_MS_RELATIVE_HUMIDITY_MAX_MEASURED_VALUE
	ATTRID_MS_RELATIVE_HUMIDITY_TOLERANCE

	// 无效值定义
	MS_RELATIVE_HUMIDITY_INVALID_VALUE = 0xffff
)

/*****************************************************************************/
/***         占有率                                                           ***/
/*****************************************************************************/
const (
	// Occupancy Sensor Configuration attribute set
	ATTRID_MS_OCCUPANCY_SENSING_CONFIG_OCCUPANCY = 0x0000
	// PIR Configuration attribute set
	ATTRID_MS_OCCUPANCY_SENSING_CONFIG_PIR_O_TO_U_DELAY = 0x0010
	ATTRID_MS_OCCUPANCY_SENSING_CONFIG_PIR_U_TO_O_DELAY = 0x0011
	// Ultrasonic Configuration attribute set
	ATTRID_MS_OCCUPANCY_SENSING_CONFIG_ULTRASONIC_O_TO_U_DELAY = 0x0020
	ATTRID_MS_OCCUPANCY_SENSING_CONFIG_ULTRASONIC_U_TO_O_DELAY = 0x0021
)
