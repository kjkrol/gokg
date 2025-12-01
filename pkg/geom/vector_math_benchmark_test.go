package geom

import "testing"

var (
	vecIntSink         Vec[int]
	vecFloatSink       Vec[float64]
	float64Sink        float64
	unsignedVecSink    Vec[uint32]
	signedLengthSink   int
	unsignedLengthSink uint32
)

func BenchmarkClampUnsignedInt(b *testing.B) {
	vm := UnsignedIntVectorMath[uint32]{}
	v := NewVec(uint32(0xFFFFFFF0), uint32(0xFFFFFFF8))
	size := NewVec(uint32(1024), uint32(1024))

	b.ReportAllocs()
	for b.Loop() {
		unsignedVecSink = vm.Clamp(v, size)
	}
}

func BenchmarkClampSignedInt(b *testing.B) {
	vm := SignedIntVectorMath[int]{}
	v := NewVec(15, -3)
	size := NewVec(10, 10)

	b.ReportAllocs()
	for b.Loop() {
		vecIntSink = vm.Clamp(v, size)
	}
}

func BenchmarkClampFloat(b *testing.B) {
	vm := FloatVectorMath[float64]{}
	v := NewVec(12.5, -3.25)
	size := NewVec(10.0, 10.0)

	b.ReportAllocs()
	for b.Loop() {
		vecFloatSink = vm.Clamp(v, size)
	}
}

func BenchmarkWrapSignedInt(b *testing.B) {
	vm := SignedIntVectorMath[int]{}
	v := NewVec(37, -9)
	size := NewVec(16, 16)

	b.ReportAllocs()
	for b.Loop() {
		vecIntSink = vm.Wrap(v, size)
	}
}

func BenchmarkWrapUnsignedInt(b *testing.B) {
	vm := UnsignedIntVectorMath[uint32]{}
	v := NewVec(uint32(0xFFFFFFF0), uint32(0xFFFFFFF8))
	size := NewVec(uint32(1024), uint32(1024))

	b.ReportAllocs()
	for b.Loop() {
		unsignedVecSink = vm.Wrap(v, size)
	}
}

func BenchmarkWrapFloat(b *testing.B) {
	vm := FloatVectorMath[float64]{}
	v := NewVec(37.5, -9.25)
	size := NewVec(16.0, 16.0)

	b.ReportAllocs()
	for b.Loop() {
		vecFloatSink = vm.Wrap(v, size)
	}
}

func BenchmarkLengthFloat(b *testing.B) {
	vm := FloatVectorMath[float64]{}
	v := NewVec(3.14, 2.71)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		float64Sink = vm.Length(v)
	}
}

func BenchmarkLengthSignedInt(b *testing.B) {
	vm := SignedIntVectorMath[int]{}
	v := NewVec(9, 12)

	b.ReportAllocs()
	for b.Loop() {
		signedLengthSink = vm.Length(v)
	}
}

func BenchmarkLengthUnsignedInt(b *testing.B) {
	vm := UnsignedIntVectorMath[uint32]{}
	v := NewVec(uint32(9), uint32(12))

	b.ReportAllocs()
	for b.Loop() {
		unsignedLengthSink = vm.Length(v)
	}
}
