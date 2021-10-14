package vm

type PointerType int

const (
	_ PointerType = iota
	_IllegalPointer

	_BasePointer
	_StackPointer
)

func (p PointerType) String() string {
	switch p {
	case _BasePointer:
		return "BasePointer"
	case _StackPointer:
		return "StackPointer"
	default:
		return "IllegalPointer"
	}
}
