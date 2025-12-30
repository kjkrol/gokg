package spatial

type CellCodec interface {
	Encode(x, y uint32) int
	Decode(int) (x, y uint32)
}

type LinearCodeCodec struct {
	res  Resolution
	mask int
}

// Compile-time guard: ensures LinearCodeCodec implements CellCodec.
var _ CellCodec = (*LinearCodeCodec)(nil)

func NewLinearCodec(res Resolution) LinearCodeCodec {
	return LinearCodeCodec{
		res:  res,
		mask: int(res.MaxCoord()),
	}
}

func (lc LinearCodeCodec) Encode(x, y uint32) int {
	return (int(y) << lc.res) | int(x)
}

func (lc LinearCodeCodec) Decode(index int) (x, y uint32) {
	x = uint32(index & lc.mask)
	y = uint32(index >> lc.res)
	return
}

type MortonCodeCodec struct{}

// Compile-time guard: ensures MortonCodeCodec implements CellCodec.
var _ CellCodec = (*MortonCodeCodec)(nil)

func (mc MortonCodeCodec) Encode(x, y uint32) int {
	return int(NewMortonCode(x, y))
}

func (mc MortonCodeCodec) Decode(code int) (x, y uint32) {
	return MortonCode(code).Decode()
}
