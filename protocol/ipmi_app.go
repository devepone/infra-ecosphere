package protocol

// port from OpenIPMI
// App Network Function
const (
	IPMI_CMD_GET_DEVICE_ID = 					0x01
	IPMI_CMD_BROADCAST_GET_DEVICE_ID = 			0x01
	IPMI_CMD_COLD_RESET = 					0x02
	IPMI_CMD_WARM_RESET = 					0x03
	IPMI_CMD_GET_SELF_TEST_RESULTS = 			0x04
	IPMI_CMD_MANUFACTURING_TEST_ON = 			0x05
	IPMI_CMD_SET_ACPI_POWER_STATE = 			0x06
	IPMI_CMD_GET_ACPI_POWER_STATE = 			0x07
	IPMI_CMD_GET_DEVICE_GUID = 				0x08
	IPMI_CMD_RESET_WATCHDOG_TIMER = 			0x22
	IPMI_CMD_SET_WATCHDOG_TIMER = 			0x24
	IPMI_CMD_GET_WATCHDOG_TIMER = 			0x25
	IPMI_CMD_SET_BMC_GLOBAL_ENABLES = 			0x2e
	IPMI_CMD_GET_BMC_GLOBAL_ENABLES = 			0x2f
	IPMI_CMD_CLEAR_MSG_FLAGS = 				0x30
	IPMI_CMD_GET_MSG_FLAGS = 				0x31
	IPMI_CMD_ENABLE_MESSAGE_CHANNEL_RCV = 		0x32
	IPMI_CMD_GET_MSG = 					0x33
	IPMI_CMD_SEND_MSG = 					0x34
	IPMI_CMD_READ_EVENT_MSG_BUFFER = 			0x35
	IPMI_CMD_GET_BT_INTERFACE_CAPABILITIES = 		0x36
	IPMI_CMD_GET_SYSTEM_GUID = 				0x37
	IPMI_CMD_GET_CHANNEL_AUTH_CAPABILITIES = 		0x38
	IPMI_CMD_GET_SESSION_CHALLENGE = 			0x39
	IPMI_CMD_ACTIVATE_SESSION = 				0x3a
	IPMI_CMD_SET_SESSION_PRIVILEGE = 			0x3b
	IPMI_CMD_CLOSE_SESSION = 					0x3c
	IPMI_CMD_GET_SESSION_INFO = 				0x3d

	IPMI_CMD_GET_AUTHCODE = 					0x3f
	IPMI_CMD_SET_CHANNEL_ACCESS = 				0x40
	IPMI_CMD_GET_CHANNEL_ACCESS = 				0x41
	IPMI_CMD_GET_CHANNEL_INFO = 				0x42
	IPMI_CMD_SET_USER_ACCESS = 				0x43
	IPMI_CMD_GET_USER_ACCESS = 				0x44
	IPMI_CMD_SET_USER_NAME = 				0x45
	IPMI_CMD_GET_USER_NAME = 				0x46
	IPMI_CMD_SET_USER_PASSWORD = 				0x47
	IPMI_CMD_ACTIVATE_PAYLOAD = 				0x48
	IPMI_CMD_DEACTIVATE_PAYLOAD = 				0x49
	IPMI_CMD_GET_PAYLOAD_ACTIVATION_STATUS = 		0x4a
	IPMI_CMD_GET_PAYLOAD_INSTANCE_INFO = 		0x4b
	IPMI_CMD_SET_USER_PAYLOAD_ACCESS = 			0x4c
	IPMI_CMD_GET_USER_PAYLOAD_ACCESS = 			0x4d
	IPMI_CMD_GET_CHANNEL_PAYLOAD_SUPPORT = 		0x4e
	IPMI_CMD_GET_CHANNEL_PAYLOAD_VERSION = 		0x4f
	IPMI_CMD_GET_CHANNEL_OEM_PAYLOAD_INFO = 		0x50

	IPMI_CMD_MASTER_READ_WRITE = 				0x52

	IPMI_CMD_GET_CHANNEL_CIPHER_SUITES = 			0x54
	IPMI_CMD_SUSPEND_RESUME_PAYLOAD_ENCRYPTION = 	0x55
	IPMI_CMD_SET_CHANNEL_SECURITY_KEY = 			0x56
	IPMI_CMD_GET_SYSTEM_INTERFACE_CAPABILITIES = 	0x57
)
