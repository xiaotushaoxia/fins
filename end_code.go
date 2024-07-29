package fins

import "fmt"

// Data taken from Omron document Cat. No. W342-E1-15, pages 155-161
const (
	// EndCodeNormalCompletion End code: normal completion
	EndCodeNormalCompletion uint16 = 0x0000

	// EndCodeServiceInterrupted End code: normal completion; service was interrupted
	EndCodeServiceInterrupted uint16 = 0x0001

	// EndCodeLocalNodeNotInNetwork End code: local node error; local node not in network
	EndCodeLocalNodeNotInNetwork uint16 = 0x0101

	// EndCodeTokenTimeout End code: local node error; token timeout
	EndCodeTokenTimeout uint16 = 0x0102

	// EndCodeRetriesFailed End code: local node error; retries failed
	EndCodeRetriesFailed uint16 = 0x0103

	// EndCodeTooManySendFrames End code: local node error; too many send frames
	EndCodeTooManySendFrames uint16 = 0x0104

	// EndCodeNodeAddressRangeError End code: local node error; node address range error
	EndCodeNodeAddressRangeError uint16 = 0x0105

	// EndCodeNodeAddressRangeDuplication End code: local node error; node address range duplication
	EndCodeNodeAddressRangeDuplication uint16 = 0x0106

	// EndCodeDestinationNodeNotInNetwork End code: destination node error; destination node not in network
	EndCodeDestinationNodeNotInNetwork uint16 = 0x0201

	// EndCodeUnitMissing End code: destination node error; unit missing
	EndCodeUnitMissing uint16 = 0x0202

	// EndCodeThirdNodeMissing End code: destination node error; third node missing
	EndCodeThirdNodeMissing uint16 = 0x0203

	// EndCodeDestinationNodeBusy End code: destination node error; destination node busy
	EndCodeDestinationNodeBusy uint16 = 0x0204

	// EndCodeResponseTimeout End code: destination node error; response timeout
	EndCodeResponseTimeout uint16 = 0x0205

	// EndCodeCommunicationsControllerError End code: controller error; communication controller error
	EndCodeCommunicationsControllerError uint16 = 0x0301

	// EndCodeCPUUnitError End code: controller error; CPU unit error
	EndCodeCPUUnitError uint16 = 0x0302

	// EndCodeControllerError End code:  controller error; controller error
	EndCodeControllerError uint16 = 0x0303

	// EndCodeUnitNumberError End code: controller error; unit number error
	EndCodeUnitNumberError uint16 = 0x0304

	// EndCodeUndefinedCommand End code: service unsupported; undefined command
	EndCodeUndefinedCommand uint16 = 0x0401

	// EndCodeNotSupportedByModelVersion End code: service unsupported; not supported by model version
	EndCodeNotSupportedByModelVersion uint16 = 0x0402

	// EndCodeDestinationAddressSettingError End code: routing table error; destination address setting error
	EndCodeDestinationAddressSettingError uint16 = 0x0501

	// EndCodeNoRoutingTables End code: routing table error; no routing tables
	EndCodeNoRoutingTables uint16 = 0x0502

	// EndCodeRoutingTableError End code: routing table error; routing table error
	EndCodeRoutingTableError uint16 = 0x0503

	// EndCodeTooManyRelays End code: routing table error; too many relays
	EndCodeTooManyRelays uint16 = 0x0504

	// EndCodeCommandTooLong End code: command format error; command too long
	EndCodeCommandTooLong uint16 = 0x1001

	// EndCodeCommandTooShort End code: command format error; command too short
	EndCodeCommandTooShort uint16 = 0x1002

	// EndCodeElementsDataDontMatch End code: command format error; elements/data don't match
	EndCodeElementsDataDontMatch uint16 = 0x1003

	// EndCodeCommandFormatError End code: command format error; command format error
	EndCodeCommandFormatError uint16 = 0x1004

	// EndCodeHeaderError End code: command format error; header error
	EndCodeHeaderError uint16 = 0x1005

	// EndCodeAreaClassificationMissing End code: parameter error; classification missing
	EndCodeAreaClassificationMissing uint16 = 0x1101

	// EndCodeAccessSizeError End code: parameter error; access size error
	EndCodeAccessSizeError uint16 = 0x1102

	// EndCodeAddressRangeError End code: parameter error; address range error
	EndCodeAddressRangeError uint16 = 0x1103

	// EndCodeAddressRangeExceeded End code: parameter error; address range exceeded
	EndCodeAddressRangeExceeded uint16 = 0x1104

	// EndCodeProgramMissing End code: parameter error; program missing
	EndCodeProgramMissing uint16 = 0x1106

	// EndCodeRelationalError End code: parameter error; relational error
	EndCodeRelationalError uint16 = 0x1109

	// EndCodeDuplicateDataAccess End code: parameter error; duplicate data access
	EndCodeDuplicateDataAccess uint16 = 0x110a

	// EndCodeResponseTooBig End code: parameter error; response too big
	EndCodeResponseTooBig uint16 = 0x110b

	// EndCodeParameterError End code: parameter error
	EndCodeParameterError uint16 = 0x110c

	// EndCodeReadNotPossibleProtected End code: read not possible; protected
	EndCodeReadNotPossibleProtected uint16 = 0x2002

	// EndCodeReadNotPossibleTableMissing End code: read not possible; table missing
	EndCodeReadNotPossibleTableMissing uint16 = 0x2003

	// EndCodeReadNotPossibleDataMissing End code: read not possible; data missing
	EndCodeReadNotPossibleDataMissing uint16 = 0x2004

	// EndCodeReadNotPossibleProgramMissing End code: read not possible; program missing
	EndCodeReadNotPossibleProgramMissing uint16 = 0x2005

	// EndCodeReadNotPossibleFileMissing End code: read not possible; file missing
	EndCodeReadNotPossibleFileMissing uint16 = 0x2006

	// EndCodeReadNotPossibleDataMismatch End code: read not possible; data mismatch
	EndCodeReadNotPossibleDataMismatch uint16 = 0x2007

	// EndCodeWriteNotPossibleReadOnly End code: write not possible; read only
	EndCodeWriteNotPossibleReadOnly uint16 = 0x2101

	// EndCodeWriteNotPossibleProtected End code: write not possible; write protected
	EndCodeWriteNotPossibleProtected uint16 = 0x2102

	// EndCodeWriteNotPossibleCannotRegister End code: write not possible; cannot register
	EndCodeWriteNotPossibleCannotRegister uint16 = 0x2103

	// EndCodeWriteNotPossibleProgramMissing End code: write not possible; program missing
	EndCodeWriteNotPossibleProgramMissing uint16 = 0x2105

	// EndCodeWriteNotPossibleFileMissing End code: write not possible; file missing
	EndCodeWriteNotPossibleFileMissing uint16 = 0x2106

	// EndCodeWriteNotPossibleFileNameAlreadyExists End code: write not possible; file name already exists
	EndCodeWriteNotPossibleFileNameAlreadyExists uint16 = 0x2107

	// EndCodeWriteNotPossibleCannotChange End code: write not possible; cannot change
	EndCodeWriteNotPossibleCannotChange uint16 = 0x2108

	// EndCodeNotExecutableInCurrentModeNotPossibleDuringExecution End code: not executeable in current mode during execution
	EndCodeNotExecutableInCurrentModeNotPossibleDuringExecution uint16 = 0x2201

	// EndCodeNotExecutableInCurrentModeNotPossibleWhileRunning End code: not executeable in current mode while running
	EndCodeNotExecutableInCurrentModeNotPossibleWhileRunning uint16 = 0x2202

	// EndCodeNotExecutableInCurrentModeWrongPLCModeInProgram End code: not executeable in current mode; PLC is in PROGRAM mode
	EndCodeNotExecutableInCurrentModeWrongPLCModeInProgram uint16 = 0x2203

	// EndCodeNotExecutableInCurrentModeWrongPLCModeInDebug End code: not executeable in current mode; PLC is in DEBUG mode
	EndCodeNotExecutableInCurrentModeWrongPLCModeInDebug uint16 = 0x2204

	// EndCodeNotExecutableInCurrentModeWrongPLCModeInMonitor End code: not executeable in current mode; PLC is in MONITOR mode
	EndCodeNotExecutableInCurrentModeWrongPLCModeInMonitor uint16 = 0x2205

	// EndCodeNotExecutableInCurrentModeWrongPLCModeInRun End code: not executeable in current mode; PLC is in RUN mode
	EndCodeNotExecutableInCurrentModeWrongPLCModeInRun uint16 = 0x2206

	// EndCodeNotExecutableInCurrentModeSpecifiedNodeNotPollingNode End code: not executeable in current mode; specified node is not polling node
	EndCodeNotExecutableInCurrentModeSpecifiedNodeNotPollingNode uint16 = 0x2207

	// EndCodeNotExecutableInCurrentModeStepCannotBeExecuted End code: not executeable in current mode; step cannot be executed
	EndCodeNotExecutableInCurrentModeStepCannotBeExecuted uint16 = 0x2208

	// EndCodeNoSuchDeviceFileDeviceMissing End code: no such device; file device missing
	EndCodeNoSuchDeviceFileDeviceMissing uint16 = 0x2301

	// EndCodeNoSuchDeviceMemoryMissing End code: no such device; memory missing
	EndCodeNoSuchDeviceMemoryMissing uint16 = 0x2302

	// EndCodeNoSuchDeviceClockMissing End code: no such device; clock missing
	EndCodeNoSuchDeviceClockMissing uint16 = 0x2303

	// EndCodeCannotStartStopTableMissing End code: cannot start/stop; table missing
	EndCodeCannotStartStopTableMissing uint16 = 0x2401

	// EndCodeUnitErrorMemoryError End code: unit error; memory error
	EndCodeUnitErrorMemoryError uint16 = 0x2502

	// EndCodeUnitErrorIOError End code: unit error; IO error
	EndCodeUnitErrorIOError uint16 = 0x2503

	// EndCodeUnitErrorTooManyIOPoints End code: unit error; too many IO points
	EndCodeUnitErrorTooManyIOPoints uint16 = 0x2504

	// EndCodeUnitErrorCPUBusError End code: unit error; CPU bus error
	EndCodeUnitErrorCPUBusError uint16 = 0x2505

	// EndCodeUnitErrorIODuplication End code: unit error; IO duplication
	EndCodeUnitErrorIODuplication uint16 = 0x2506

	// EndCodeUnitErrorIOBusError End code: unit error; IO bus error
	EndCodeUnitErrorIOBusError uint16 = 0x2507

	// EndCodeUnitErrorSYSMACBUS2Error End code: unit error; SYSMAC BUS/2 error
	EndCodeUnitErrorSYSMACBUS2Error uint16 = 0x2509

	// EndCodeUnitErrorCPUBusUnitError End code: unit error; CPU bus unit error
	EndCodeUnitErrorCPUBusUnitError uint16 = 0x250a

	// EndCodeUnitErrorSYSMACBusNumberDuplication End code: unit error; SYSMAC bus number duplication
	EndCodeUnitErrorSYSMACBusNumberDuplication uint16 = 0x250d

	// EndCodeUnitErrorMemoryStatusError End code: unit error; memory status error
	EndCodeUnitErrorMemoryStatusError uint16 = 0x250f

	// EndCodeUnitErrorSYSMACBusTerminatorMissing End code: unit error; SYSMAC bus terminator missing
	EndCodeUnitErrorSYSMACBusTerminatorMissing uint16 = 0x2510

	// EndCodeCommandErrorNoProtection End code: command error; no protection
	EndCodeCommandErrorNoProtection uint16 = 0x2601

	// EndCodeCommandErrorIncorrectPassword End code: command error; incorrect password
	EndCodeCommandErrorIncorrectPassword uint16 = 0x2602

	// EndCodeCommandErrorProtected End code: command error; protected
	EndCodeCommandErrorProtected uint16 = 0x2604

	// EndCodeCommandErrorServiceAlreadyExecuting End code: command error; service already executing
	EndCodeCommandErrorServiceAlreadyExecuting uint16 = 0x2605

	// EndCodeCommandErrorServiceStopped End code: command error; service stopped
	EndCodeCommandErrorServiceStopped uint16 = 0x2606

	// EndCodeCommandErrorNoExecutionRight End code: command error; no execution right
	EndCodeCommandErrorNoExecutionRight uint16 = 0x2607

	// EndCodeCommandErrorSettingsNotComplete End code: command error; settings not complete
	EndCodeCommandErrorSettingsNotComplete uint16 = 0x2608

	// EndCodeCommandErrorNecessaryItemsNotSet End code: command error; necessary items not set
	EndCodeCommandErrorNecessaryItemsNotSet uint16 = 0x2609

	// EndCodeCommandErrorNumberAlreadyDefined End code: command error; number already defined
	EndCodeCommandErrorNumberAlreadyDefined uint16 = 0x260a

	// EndCodeCommandErrorErrorWillNotClear End code: command error; error will not clear
	EndCodeCommandErrorErrorWillNotClear uint16 = 0x260b

	// EndCodeAccessWriteErrorNoAccessRight End code: access write error; no access right
	EndCodeAccessWriteErrorNoAccessRight uint16 = 0x3001

	// EndCodeAbortServiceAborted End code: abort; service aborted
	EndCodeAbortServiceAborted uint16 = 0x4001
)

func EndCodeToMsg(u uint16) string {
	if s, ok := endCodeToMsg[u]; ok {
		return s
	}
	return fmt.Sprintf("End code: 0x%x: unknown", u)
}

var endCodeToMsg = map[uint16]string{
	0x0000: "end code 0x0000: normal completion",
	0x0001: "end code 0x0001: normal completion; service was interrupted",
	0x0101: "end code 0x0101: local node error; local node not in network",
	0x0102: "end code 0x0102: local node error; token timeout",
	0x0103: "end code 0x0103: local node error; retries failed",
	0x0104: "end code 0x0104: local node error; too many send frames",
	0x0105: "end code 0x0105: local node error; node address range error",
	0x0106: "end code 0x0106: local node error; node address range duplication",
	0x0201: "end code 0x0201: destination node error; destination node not in network",
	0x0202: "end code 0x0202: destination node error; unit missing",
	0x0203: "end code 0x0203: destination node error; third node missing",
	0x0204: "end code 0x0204: destination node error; destination node busy",
	0x0205: "end code 0x0205: destination node error; response timeout",
	0x0301: "end code 0x0301: controller error; communication controller error",
	0x0302: "end code 0x0302: controller error; CPU unit error",
	0x0303: "end code 0x0303: controller error; controller error",
	0x0304: "end code 0x0304: controller error; unit number error",
	0x0401: "end code 0x0401: service unsupported; undefined command",
	0x0402: "end code 0x0402: service unsupported; not supported by model version",
	0x0501: "end code 0x0501: routing table error; destination address setting error",
	0x0502: "end code 0x0502: routing table error; no routing tables",
	0x0503: "end code 0x0503: routing table error; routing table error",
	0x0504: "end code 0x0504: routing table error; too many relays",
	0x1001: "end code 0x1001: command format error; command too long",
	0x1002: "end code 0x1002: command format error; command too short",
	0x1003: "end code 0x1003: command format error; elements/data don't match",
	0x1004: "end code 0x1004: command format error; command format error",
	0x1005: "end code 0x1005: command format error; header error",
	0x1101: "end code 0x1101: parameter error; classification missing",
	0x1102: "end code 0x1102: parameter error; access size error",
	0x1103: "end code 0x1103: parameter error; address range error",
	0x1104: "end code 0x1104: parameter error; address range exceeded",
	0x1106: "end code 0x1106: parameter error; program missing",
	0x1109: "end code 0x1109: parameter error; relational error",
	0x110a: "end code 0x110a: parameter error; duplicate data access",
	0x110b: "end code 0x110b: parameter error; response too big",
	0x110c: "end code 0x110c: parameter error",
	0x2002: "end code 0x2002: read not possible; protected",
	0x2003: "end code 0x2003: read not possible; table missing",
	0x2004: "end code 0x2004: read not possible; data missing",
	0x2005: "end code 0x2005: read not possible; program missing",
	0x2006: "end code 0x2006: read not possible; file missing",
	0x2007: "end code 0x2007: read not possible; data mismatch",
	0x2101: "end code 0x2101: write not possible; read only",
	0x2102: "end code 0x2102: write not possible; write protected",
	0x2103: "end code 0x2103: write not possible; cannot register",
	0x2105: "end code 0x2105: write not possible; program missing",
	0x2106: "end code 0x2106: write not possible; file missing",
	0x2107: "end code 0x2107: write not possible; file name already exists",
	0x2108: "end code 0x2108: write not possible; cannot change",
	0x2201: "end code 0x2201: not executable in current mode during execution",
	0x2202: "end code 0x2202: not executable in current mode while running",
	0x2203: "end code 0x2203: not executable in current mode; PLC is in PROGRAM mode",
	0x2204: "end code 0x2204: not executable in current mode; PLC is in DEBUG mode",
	0x2205: "end code 0x2205: not executable in current mode; PLC is in MONITOR mode",
	0x2206: "end code 0x2206: not executable in current mode; PLC is in RUN mode",
	0x2207: "end code 0x2207: not executable in current mode; specified node is not polling node",
	0x2208: "end code 0x2208: not executable in current mode; step cannot be executed",
	0x2301: "end code 0x2301: no such device; file device missing",
	0x2302: "end code 0x2302: no such device; memory missing",
	0x2303: "end code 0x2303: no such device; clock missing",
	0x2401: "end code 0x2401: cannot start/stop; table missing",
	0x2502: "end code 0x2502: unit error; memory error",
	0x2503: "end code 0x2503: unit error; IO error",
	0x2504: "end code 0x2504: unit error; too many IO points",
	0x2505: "end code 0x2505: unit error; CPU bus error",
	0x2506: "end code 0x2506: unit error; IO duplication",
	0x2507: "end code 0x2507: unit error; IO bus error",
	0x2509: "end code 0x2509: unit error; SYSMAC BUS/2 error",
	0x250a: "end code 0x250a: unit error; CPU bus unit error",
	0x250d: "end code 0x250d: unit error; SYSMAC bus number duplication",
	0x250f: "end code 0x250f: unit error; memory status error",
	0x2510: "end code 0x2510: unit error; SYSMAC bus terminator missing",
	0x2601: "end code 0x2601: command error; no protection",
	0x2602: "end code 0x2602: command error; incorrect password",
	0x2604: "end code 0x2604: command error; protected",
	0x2605: "end code 0x2605: command error; service already executing",
	0x2606: "end code 0x2606: command error; service stopped",
	0x2607: "end code 0x2607: command error; no execution right",
	0x2608: "end code 0x2608: command error; settings not complete",
	0x2609: "end code 0x2609: command error; necessary items not set",
	0x260a: "end code 0x260a: command error; number already defined",
	0x260b: "end code 0x260b: command error; error will not clear",
	0x3001: "end code 0x3001: access write error; no access right",
	0x4001: "end code 0x4001: abort; service aborted",
}
