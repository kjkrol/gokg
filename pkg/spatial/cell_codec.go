package spatial

type CellCodec interface {
	Encode(x, y uint32) uint64
	Decode(uint64) (x, y uint32)
}

type LinearCodeCodec struct {
	res  Resolution
	mask uint64
}

// Compile-time guard: ensures LinearCodeCodec implements CellCodec.
var _ CellCodec = (*LinearCodeCodec)(nil)

func NewLinearCodec(res Resolution) LinearCodeCodec {
	return LinearCodeCodec{
		res:  res,
		mask: uint64(res.MaxCoord()),
	}
}

func (lc LinearCodeCodec) Encode(x, y uint32) uint64 {
	return (uint64(y) << lc.res) | uint64(x)
}

func (lc LinearCodeCodec) Decode(index uint64) (x, y uint32) {
	x = uint32(index & lc.mask)
	y = uint32(index >> lc.res)
	return
}

type MortonCodeCodec struct{}

// Compile-time guard: ensures MortonCodeCodec implements CellCodec.
var _ CellCodec = (*MortonCodeCodec)(nil)

func (mc MortonCodeCodec) Encode(x, y uint32) uint64 {
	return uint64(NewMortonCode(x, y))
}

func (mc MortonCodeCodec) Decode(code uint64) (x, y uint32) {
	return MortonCode(code).Decode()
}
