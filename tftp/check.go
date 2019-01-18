package tftp

// 检查mode值是否合法
func CheckMode(mode string) bool {
	if mode == modeOctet || mode == modeMail || mode == modeNetAscii {
		return true
	}
	return false
}

// 检查opcode值是否合法
func CheckOpCode(code uint16) bool {
	if code == opRRQ || code == opWRQ || code == opDATA || code == opACK || code == opERR {
		return true
	}
	return false
}
