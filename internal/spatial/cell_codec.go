package spatial

// cellCodec numbers the cells of a square grid row by row.
type cellCodec struct {
	res  Resolution
	mask int
}

func newCellCodec(res Resolution) cellCodec {
	return cellCodec{res: res, mask: int(res.MaxCoord())}
}

// Encode is the index of cell (x, y); both must be within the grid.
func (c cellCodec) Encode(x, y uint32) int { return (int(y) << c.res) | int(x) }

// Decode is the cell at index.
func (c cellCodec) Decode(index int) (x, y uint32) {
	return uint32(index & c.mask), uint32(index >> c.res)
}
