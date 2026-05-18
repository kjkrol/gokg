package spatial

import "errors"

var errCoordOutOfRange = errors.New("grid coordinates out of range")

type CellCodec interface {
	Encode(x, y uint32) (int, error)
	Decode(int) (x, y uint32)
}

type LinearCodeCodec struct {
	res      Resolution
	mask     int
	maxCoord uint32
}

// Compile-time guard: ensures LinearCodeCodec implements CellCodec.
var _ CellCodec = (*LinearCodeCodec)(nil)

func NewLinearCodec(res Resolution) LinearCodeCodec {
	maxCoord := res.MaxCoord()
	return LinearCodeCodec{
		res:      res,
		mask:     int(maxCoord),
		maxCoord: maxCoord,
	}
}

func (lc LinearCodeCodec) Encode(x, y uint32) (int, error) {
	if x > lc.maxCoord || y > lc.maxCoord {
		return -1, errCoordOutOfRange
	}
	return (int(y) << lc.res) | int(x), nil
}

func (lc LinearCodeCodec) Decode(index int) (x, y uint32) {
	x = uint32(index & lc.mask)
	y = uint32(index >> lc.res)
	return
}

type MortonCodeCodec struct {
	maxCoord uint32
}

// Compile-time guard: ensures MortonCodeCodec implements CellCodec.
var _ CellCodec = (*MortonCodeCodec)(nil)

func NewMortonCodec(res Resolution) MortonCodeCodec {
	return MortonCodeCodec{maxCoord: res.MaxCoord()}
}

func (mc MortonCodeCodec) Encode(x, y uint32) (int, error) {
	if x > mc.maxCoord || y > mc.maxCoord {
		return -1, errCoordOutOfRange
	}
	return int(NewMortonCode(x, y)), nil
}

func (mc MortonCodeCodec) Decode(code int) (x, y uint32) {
	return MortonCode(code).Decode()
}
