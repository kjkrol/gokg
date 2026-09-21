package spatial

type MortonCode uint64

func NewMortonCode(x, y uint32) MortonCode {
	return MortonCode((splitBy1(x)) | splitBy1(y)<<1)
}

func (c MortonCode) Decode() (x, y uint32) {
	code := uint64(c)
	x = compact1By1(code)
	y = compact1By1(code >> 1)
	return
}

func (c MortonCode) Offset(dx, dy int32) MortonCode {
	x, y := c.Decode()

	nx := uint32(int64(x) + int64(dx))
	ny := uint32(int64(y) + int64(dy))

	return NewMortonCode(nx, ny)
}

func MortonCodeArea(aabb AABB) []MortonCode {
	minX, minY, maxX, maxY, ok := cellBox(aabb)
	if !ok {
		return nil
	}

	width := maxX - minX + 1
	height := maxY - minY + 1

	count := uint64(width) * uint64(height)
	if count == 0 {
		return nil
	}

	res := make([]MortonCode, int(count))

	rowStart := NewMortonCode(minX, minY)

	idx := 0
	for range height {
		code := rowStart

		for range width {
			res[idx] = code
			idx++

			code = code.IncX()
		}

		rowStart = rowStart.IncY()
	}

	return res
}

func MortonCodeAreaConsume(aabb AABB, fn func(int, MortonCode)) {
	if aabb.BottomRight.X < aabb.TopLeft.X || aabb.BottomRight.Y < aabb.TopLeft.Y {
		return
	}

	minX, minY, maxX, maxY, ok := cellBox(aabb)
	if !ok {
		return
	}

	width := maxX - minX + 1
	height := maxY - minY + 1

	count := uint64(width) * uint64(height)
	if count == 0 {
		return
	}

	rowStart := NewMortonCode(minX, minY)

	idx := 0
	for range height {
		code := rowStart

		for range width {
			fn(idx, code)
			idx++
			code = code.IncX()
		}

		rowStart = rowStart.IncY()
	}
}

const (
	xMask uint64 = 0x5555555555555555 // bity X na pozycjach 0,2,4,...
	yMask uint64 = 0xAAAAAAAAAAAAAAAA // bity Y na pozycjach 1,3,5,...
)

func (c MortonCode) IncX() MortonCode {
	code := uint64(c)
	x := code & xMask
	y := code & yMask

	x = (x - xMask) & xMask

	return MortonCode(x | y)
}

func (c MortonCode) IncY() MortonCode {
	code := uint64(c)
	x := code & xMask
	y := (code & yMask) >> 1

	y = (y - xMask) & xMask
	y = y << 1

	return MortonCode(x | y)
}

// splitBy1 spreads the lower 32 bits of a so that a zero bit sits between each two.
func splitBy1(a uint32) uint64 {
	x := uint64(a)
	x = (x | (x << 16)) & 0x0000FFFF0000FFFF
	x = (x | (x << 8)) & 0x00FF00FF00FF00FF
	x = (x | (x << 4)) & 0x0F0F0F0F0F0F0F0F
	x = (x | (x << 2)) & 0x3333333333333333
	x = (x | (x << 1)) & 0x5555555555555555
	return x
}

func compact1By1(a uint64) uint32 {
	x := a & 0x5555555555555555
	x = (x | (x >> 1)) & 0x3333333333333333
	x = (x | (x >> 2)) & 0x0F0F0F0F0F0F0F0F
	x = (x | (x >> 4)) & 0x00FF00FF00FF00FF
	x = (x | (x >> 8)) & 0x0000FFFF0000FFFF
	x = (x | (x >> 16)) & 0x00000000FFFFFFFF
	return uint32(x)
}

// cellBox reads a box as the whole cells it addresses; ok is false when a corner has no cell.
func cellBox(aabb AABB) (minX, minY, maxX, maxY uint32, ok bool) {
	if !(aabb.TopLeft.X >= 0) || !(aabb.TopLeft.Y >= 0) {
		return 0, 0, 0, 0, false
	}
	if aabb.BottomRight.X < aabb.TopLeft.X || aabb.BottomRight.Y < aabb.TopLeft.Y {
		return 0, 0, 0, 0, false
	}
	return uint32(aabb.TopLeft.X), uint32(aabb.TopLeft.Y),
		uint32(aabb.BottomRight.X), uint32(aabb.BottomRight.Y), true
}
