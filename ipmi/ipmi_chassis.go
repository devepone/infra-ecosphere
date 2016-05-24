package ipmi

import (
	"net"
	"log"
	"bytes"
	"github.com/rmxymh/infra-ecosphere/bmc"
	"github.com/rmxymh/infra-ecosphere/utils"
	"encoding/binary"
)

// port from OpenIPMI
// Chassis Network Function
const (
	IPMI_CMD_GET_CHASSIS_CAPABILITIES =	0x00
	IPMI_CMD_GET_CHASSIS_STATUS =		0x01
	IPMI_CMD_CHASSIS_CONTROL =		0x02
	IPMI_CMD_CHASSIS_RESET =		0x03
	IPMI_CMD_CHASSIS_IDENTIFY =		0x04
	IPMI_CMD_SET_CHASSIS_CAPABILITIES =	0x05
	IPMI_CMD_SET_POWER_RESTORE_POLICY =	0x06
	IPMI_CMD_GET_SYSTEM_RESTART_CAUSE =	0x07
	IPMI_CMD_SET_SYSTEM_BOOT_OPTIONS =	0x08
	IPMI_CMD_GET_SYSTEM_BOOT_OPTIONS =	0x09
	IPMI_CMD_GET_POH_COUNTER =		0x0f
)

func HandleIPMIUnsupportedChassisCommand(addr *net.UDPAddr, server *net.UDPConn, wrapper IPMISessionWrapper, message IPMIMessage) {
	log.Println("      IPMI App: This command is not supported currently, ignore.")
}

type IPMIGetChassisStatusResponse struct {
	CurrentPowerState uint8
	LastPowerEvent uint8
	MiscChassisState uint8
	FrontPanelButtonCapabilities uint8
}

const (
	CHASSIS_POWER_STATE_BITMASK_POWER_ON = 			0x01
	CHASSIS_POWER_STATE_BITMASK_POWER_OVERLOAD =		0x02
	CHASSIS_POWER_STATE_BITMASK_INTERLOCK = 		0x04
	CHASSIS_POWER_STATE_BITMASK_POWER_FAULT =		0x08
	CHASSIS_POWER_STATE_BITMASK_POWER_CONTROL_FAULT =	0x10

	// Bit 5 ~ 6
	CHASSIS_POWER_STATE_BITMASK_POWER_RESTORE_POWER_OFF =	0x00
	CHASSIS_POWER_STATE_BITMASK_POWER_RESTORE_RESTORE =	0x20
	CHASSIS_POWER_STATE_BITMASK_POWER_RESTORE_POWER_UP =	0x40
	CHASSIS_POWER_STATE_BITMASK_POWER_RESTORE_UNKNOWN =	0x60
)

const (
	CHASSIS_LAST_POWER_AC_FAILED =		0x01
	CHASSIS_LAST_POWER_DOWN_OVERLOAD =	0x02
	CHASSIS_LAST_POWER_DOWN_INTERLOCK =	0x04
	CHASSIS_LAST_PWOER_DOWN_FAULT =		0x08
	CHASSIS_LAST_POWER_ON_VIA_IPMI =	0x10
)

const (
	CHASSIS_MISC_INTRUCTION_ACTIVE =	0x01
	CHASSIS_MISC_FRONT_PANEL_LOCKOUT =	0x02
	CHASSIS_MISC_DRIVE_FAULT =		0x04
	CHASSIS_MISC_COOLING_FAULT =		0x08

	// Bit 4 ~ 5
	CHASSIS_MISC_IDENTIFY_OFF =		0x00
	CHASSIS_MISC_IDENTIFY_TEMPERARY =	0x10
	CHASSIS_MISC_IDENTIFY_INDEFINITE_ON =	0x20

	CHASSIS_MISC_IDENTIFY_SUPPORTED =	0x40
)

func HandleIPMIGetChassisStatus(addr *net.UDPAddr, server *net.UDPConn, wrapper IPMISessionWrapper, message IPMIMessage) {
	session, ok := GetSession(wrapper.SessionId)
	if ! ok {
		log.Printf("Unable to find session 0x%08x\n", wrapper.SessionId)
	} else {
		bmcUser := session.User
		code := GetAuthenticationCode(wrapper.AuthenticationType, bmcUser.Password, wrapper.SessionId, message, wrapper.SequenceNumber)
		if bytes.Compare(wrapper.AuthenticationCode[:], code[:]) == 0 {
			log.Println("      IPMI Authentication Pass.")
		} else {
			log.Println("      IPMI Authentication Failed.")
		}

		localIP := utils.GetLocalIP(server)
		bmc, ok := bmc.GetBMC(net.ParseIP(localIP))
		if ! ok {
			log.Printf("BMC %s is not found\n", localIP)
		} else {
			session.LocalSessionSequenceNumber += 1
			session.RemoteSessionSequenceNumber += 1

			response := IPMIGetChassisStatusResponse{}
			if bmc.VM.IsRunning() {
				response.CurrentPowerState |= CHASSIS_POWER_STATE_BITMASK_POWER_ON
			}
			response.LastPowerEvent = 0
			response.MiscChassisState = 0
			response.FrontPanelButtonCapabilities = 0

			dataBuf := bytes.Buffer{}
			binary.Write(&dataBuf, binary.BigEndian, response)

			responseWrapper, responseMessage := BuildResponseMessageTemplate(wrapper, message, (IPMI_NETFN_APP | IPMI_NETFN_RESPONSE), IPMI_CMD_GET_CHASSIS_STATUS)
			responseMessage.Data = dataBuf.Bytes()

			responseWrapper.SessionId = wrapper.SessionId
			responseWrapper.SequenceNumber = session.RemoteSessionSequenceNumber
			responseWrapper.AuthenticationCode = GetAuthenticationCode(wrapper.AuthenticationType, bmcUser.Password, responseWrapper.SessionId, responseMessage, responseWrapper.SequenceNumber)
			rmcp := BuildUpRMCPForIPMI()

			obuf := bytes.Buffer{}
			SerializeRMCP(&obuf, rmcp)
			SerializeIPMI(&obuf, responseWrapper, responseMessage)
			server.WriteToUDP(obuf.Bytes(), addr)
		}
	}
}

type IPMIChassisControlRequest struct {
	ChassisControl uint8
}

const (
	CHASSIS_CONTROL_POWER_DOWN =	0x00
	CHASSIS_CONTROL_POWER_UP =	0x01
	CHASSIS_CONTROL_POWER_CYCLE =	0x02
	CHASSIS_CONTROL_HARD_RESET =	0x03
	CHASSIS_CONTROL_PULSE = 	0x04
	CHASSIS_CONTROL_POWER_SOFT =	0x05
)

func HandleIPMIChassisControl(addr *net.UDPAddr, server *net.UDPConn, wrapper IPMISessionWrapper, message IPMIMessage) {
	buf := bytes.NewBuffer(message.Data)
	request := IPMIChassisControlRequest{}
	binary.Read(buf, binary.BigEndian, &request)

	session, ok := GetSession(wrapper.SessionId)
	if ! ok {
		log.Printf("Unable to find session 0x%08x\n", wrapper.SessionId)
	} else {
		bmcUser := session.User
		code := GetAuthenticationCode(wrapper.AuthenticationType, bmcUser.Password, wrapper.SessionId, message, wrapper.SequenceNumber)
		if bytes.Compare(wrapper.AuthenticationCode[:], code[:]) == 0 {
			log.Println("      IPMI Authentication Pass.")
		} else {
			log.Println("      IPMI Authentication Failed.")
		}

		localIP := utils.GetLocalIP(server)
		bmc, ok := bmc.GetBMC(net.ParseIP(localIP))
		if ! ok {
			log.Printf("BMC %s is not found\n", localIP)
		} else {
			switch request.ChassisControl {
			case CHASSIS_CONTROL_POWER_DOWN:
				bmc.PowerOff()
			case CHASSIS_CONTROL_POWER_UP:
				bmc.PowerOn()
			case CHASSIS_CONTROL_POWER_CYCLE:
				bmc.PowerOff()
				bmc.PowerOn()
			case CHASSIS_CONTROL_HARD_RESET:
				bmc.PowerOff()
				bmc.PowerOn()
			case CHASSIS_CONTROL_PULSE:
				// do nothing
			case CHASSIS_CONTROL_POWER_SOFT:
				bmc.PowerSoft()
			}

			session.LocalSessionSequenceNumber += 1
			session.RemoteSessionSequenceNumber += 1

			responseWrapper, responseMessage := BuildResponseMessageTemplate(wrapper, message, (IPMI_NETFN_APP | IPMI_NETFN_RESPONSE), IPMI_CMD_CHASSIS_CONTROL)

			responseWrapper.SessionId = wrapper.SessionId
			responseWrapper.SequenceNumber = session.RemoteSessionSequenceNumber
			responseWrapper.AuthenticationCode = GetAuthenticationCode(wrapper.AuthenticationType, bmcUser.Password, responseWrapper.SessionId, responseMessage, responseWrapper.SequenceNumber)
			rmcp := BuildUpRMCPForIPMI()

			obuf := bytes.Buffer{}
			SerializeRMCP(&obuf, rmcp)
			SerializeIPMI(&obuf, responseWrapper, responseMessage)
			server.WriteToUDP(obuf.Bytes(), addr)
		}
	}
}


func IPMI_CHASSIS_DeserializeAndExecute(addr *net.UDPAddr, server *net.UDPConn, wrapper IPMISessionWrapper, message IPMIMessage) {
	switch message.Command {
	case IPMI_CMD_GET_CHASSIS_CAPABILITIES:
		log.Println("      IPMI CHASSIS: Command = IPMI_CMD_GET_CHASSIS_CAPABILITIES")
		HandleIPMIUnsupportedChassisCommand(addr, server, wrapper, message)

	case IPMI_CMD_GET_CHASSIS_STATUS:
		log.Println("      IPMI CHASSIS: Command = IPMI_CMD_GET_CHASSIS_STATUS")
		HandleIPMIGetChassisStatus(addr, server, wrapper, message)

	case IPMI_CMD_CHASSIS_CONTROL:
		log.Println("      IPMI CHASSIS: Command = IPMI_CMD_CHASSIS_CONTROL")
		HandleIPMIChassisControl(addr, server, wrapper, message)

	case IPMI_CMD_CHASSIS_RESET:
		log.Println("      IPMI CHASSIS: Command = IPMI_CMD_CHASSIS_RESET")
		HandleIPMIUnsupportedChassisCommand(addr, server, wrapper, message)

	case IPMI_CMD_CHASSIS_IDENTIFY:
		log.Println("      IPMI CHASSIS: Command = IPMI_CMD_CHASSIS_IDENTIFY")
		HandleIPMIUnsupportedChassisCommand(addr, server, wrapper, message)

	case IPMI_CMD_SET_CHASSIS_CAPABILITIES:
		log.Println("      IPMI CHASSIS: Command = IPMI_CMD_SET_CHASSIS_CAPABILITIES")
		HandleIPMIUnsupportedChassisCommand(addr, server, wrapper, message)

	case IPMI_CMD_SET_POWER_RESTORE_POLICY:
		log.Println("      IPMI CHASSIS: Command = IPMI_CMD_SET_POWER_RESTORE_POLICY")
		HandleIPMIUnsupportedChassisCommand(addr, server, wrapper, message)

	case IPMI_CMD_GET_SYSTEM_RESTART_CAUSE:
		log.Println("      IPMI CHASSIS: Command = IPMI_CMD_GET_SYSTEM_RESTART_CAUSE")
		HandleIPMIUnsupportedChassisCommand(addr, server, wrapper, message)

	case IPMI_CMD_SET_SYSTEM_BOOT_OPTIONS:
		log.Println("      IPMI CHASSIS: Command = IPMI_CMD_SET_SYSTEM_BOOT_OPTIONS")
		IPMI_CHASSIS_SetBootOption_DeserializeAndExecute(addr, server, wrapper, message)

	case IPMI_CMD_GET_SYSTEM_BOOT_OPTIONS:
		log.Println("      IPMI CHASSIS: Command = IPMI_CMD_GET_SYSTEM_BOOT_OPTIONS")
		HandleIPMIUnsupportedChassisCommand(addr, server, wrapper, message)

	case IPMI_CMD_GET_POH_COUNTER:
		log.Println("      IPMI CHASSIS: Command = IPMI_CMD_GET_POH_COUNTER")
		HandleIPMIUnsupportedChassisCommand(addr, server, wrapper, message)

	}
}