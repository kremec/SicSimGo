package proc

type RelativeAddressingMode int
type AbsoluteAddressingMode int
type IndexAddressingMode bool

const (
	DirectRelativeAddressing RelativeAddressingMode = 0
	PCRelativeAddressing     RelativeAddressingMode = 1
	BaseRelativeAddressing   RelativeAddressingMode = 2
	UnkownRelativeAddressing RelativeAddressingMode = 3

	SICAbsoluteAddressing       AbsoluteAddressingMode = 0
	ImmediateAbsoluteAddressing AbsoluteAddressingMode = 1
	IndirectAbsoluteAddressing  AbsoluteAddressingMode = 2
	DirectAbsoluteAddressing    AbsoluteAddressingMode = 3
)

func GetRelativeAdressingModes(b, p bool) (RelativeAddressingMode, error) {
	if !b && !p {
		return DirectRelativeAddressing, nil
	} else if !b && p {
		return PCRelativeAddressing, nil
	} else if b && !p {
		return BaseRelativeAddressing, nil
	} else {
		return UnkownRelativeAddressing, InvalidAddressing()
	}
}
func (relativeAddressingMode RelativeAddressingMode) String() string {
	switch relativeAddressingMode {
	case DirectRelativeAddressing:
		return "Direct"
	case PCRelativeAddressing:
		return "PC-relative"
	case BaseRelativeAddressing:
		return "Base-relative"
	case UnkownRelativeAddressing:
		return "Unknown"
	}
	return "Not implemented"
}

func GetAbsoluteAdressingModes(n, i bool) AbsoluteAddressingMode {
	if !n && !i {
		return SICAbsoluteAddressing
	} else if !n && i {
		return ImmediateAbsoluteAddressing
	} else if n && !i {
		return IndirectAbsoluteAddressing
	} else {
		return DirectAbsoluteAddressing
	}
}
func (absoluteAddressingMode AbsoluteAddressingMode) String() string {
	switch absoluteAddressingMode {
	case SICAbsoluteAddressing:
		return "SIC"
	case ImmediateAbsoluteAddressing:
		return "Immediate"
	case IndirectAbsoluteAddressing:
		return "Indirect"
	case DirectAbsoluteAddressing:
		return "Direct"
	}
	return "Not implemented"
}

func GetIndexAdressingModes(x bool) IndexAddressingMode {
	if x {
		return true
	} else {
		return false
	}
}
