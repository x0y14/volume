package vvm

type KeyWordType int

const (
	_ KeyWordType = iota
	_ILLEGALKeyWordType
	_REGISTERKeyWord
	_POINTERKeyWord
)

func (kw KeyWordType) String() string {
	switch kw {
	case _POINTERKeyWord:
		return "POINTERKeyWord"
	case _REGISTERKeyWord:
		return "REGISTERKeyWord"
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
			return _REGISTERKeyWord
		}
	}

	pointers := []string{"bp", "sp"}
	for _, keyword := range pointers {
		if text == keyword {
			return _POINTERKeyWord
		}
	}

	return _ILLEGALKeyWordType
}
