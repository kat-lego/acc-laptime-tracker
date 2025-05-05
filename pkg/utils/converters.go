package utils

import "unicode/utf16"

// Convert UTF-16 string arrays to Go string
func Utf16ToString(u16 []uint16) string {
	for i, v := range u16 {
		if v == 0 {
			u16 = u16[:i]
			break
		}
	}
	return string(utf16.Decode(u16))
}

func Int32ToAccStatus(status int32) string {
	switch status {
	case 0:
		return "ACC_OFF"
	case 1:
		return "ACC_REPLAY"
	case 2:
		return "ACC_LIVE"
	case 3:
		return "ACC_PAUSE"
	default:
		return "UNKNOWN_STATUS"
	}
}

func Int32ToAccSession(session int32) string {
	switch session {
	case -1:
		return "ACC_UNKNOWN"
	case 0:
		return "ACC_PRACTICE"
	case 1:
		return "ACC_QUALIFY"
	case 2:
		return "ACC_RACE"
	case 3:
		return "ACC_HOTLAP"
	case 4:
		return "ACC_TIMEATTACK"
	case 5:
		return "ACC_DRIFT"
	case 6:
		return "ACC_DRAG"
	case 7:
		return "ACC_HOTSTINT"
	case 8:
		return "ACC_HOTSTINTSUPERPOLE"
	default:
		return "UNKNOWN_SESSION"
	}
}
