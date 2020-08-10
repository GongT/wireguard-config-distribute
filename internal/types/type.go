package types

type SidType uint64

func (t SidType) Serialize() uint64 {
	return uint64(t)
}

func DeSerializeSidType(t uint64) SidType {
	return SidType(t)
}

type VpnIdType string

func (t VpnIdType) Serialize() string {
	return string(t)
}

func DeSerializeVpnIdType(t string) VpnIdType {
	return VpnIdType(t)
}
