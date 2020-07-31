package types

type SidType uint64

func (t SidType) Serialize() uint64 {
	return uint64(t)
}

func DeSerialize(t uint64) SidType {
	return SidType(t)
}
