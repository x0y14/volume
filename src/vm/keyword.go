package vm

type KeyWordType int

const (
	_ KeyWordType = iota
	_ILLEGALKeyWordType
	_REGISTER
	_POINTER
)

func (kw KeyWordType) String() string {
	switch kw {
	case _POINTER:
		return "POINTER"
	case _REGISTER:
		return "REGISTER"
	default:
		return "ILLEGALKeyWordType"
	}
}

func isReservedKeyWord(text string) bool {
	reserved := []string{"reg_a", "reg_b", "reg_c", "bp", "sp"}
	for _, reservedKeyword := range reserved {
		if text == reservedKeyword {
			return true
		}
	}
	return false
}

func CheckKeyWordType(text string) KeyWordType {
	registers := []string{"reg_a", "reg_b", "reg_c"}
	for _, keyword := range registers {
		if text == keyword {
			return _REGISTER
		}
	}

	pointers := []string{"bp", "sp"}
	for _, keyword := range pointers {
		if text == keyword {
			return _POINTER
		}
	}

	return _ILLEGALKeyWordType
}
