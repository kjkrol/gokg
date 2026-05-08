package spatial

type EntryId uint64

func NewEntryID(id uint64, frag uint8) EntryId {
	return EntryId((id << 2) | uint64(frag&0x3))
}

func (e EntryId) OriginalID() uint64 {
	return uint64(e >> 2)
}

func (e EntryId) ExtractFrag() uint8 {
	return uint8(e & 0x3)
}
