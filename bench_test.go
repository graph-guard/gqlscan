package gqlscan_test

import (
	"testing"

	"github.com/graph-guard/gqlscan"
)

func BenchmarkScan(b *testing.B) {
	for _, td := range testdata {
		b.Run("", func(b *testing.B) {
			in := []byte(td.input)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := gqlscan.Scan(in, func(*gqlscan.Iterator) (err bool) {
					return false
				}); err.IsErr() {
					panic(err)
				}
			}
		})
	}
}

func BenchmarkScanAll(b *testing.B) {
	for _, td := range testdata {
		b.Run("", func(b *testing.B) {
			in := []byte(td.input)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := gqlscan.ScanAll(in, func(*gqlscan.Iterator) {
				}); err.IsErr() {
					panic(err)
				}
			}
		})
	}
}

func BenchmarkScanErr(b *testing.B) {
	for _, td := range testdataErr {
		b.Run("", func(b *testing.B) {
			in := []byte(td.input)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := gqlscan.Scan(in, func(*gqlscan.Iterator) (err bool) {
					return false
				}); !err.IsErr() {
					panic("unexpected success")
				}
			}
		})
	}
}

func BenchmarkScanAllErr(b *testing.B) {
	for _, td := range testdataErr {
		b.Run("", func(b *testing.B) {
			in := []byte(td.input)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := gqlscan.ScanAll(in, func(*gqlscan.Iterator) {
				}); !err.IsErr() {
					panic("unexpected success")
				}
			}
		})
	}
}

var InterpretedBuffer []byte

func BenchmarkScanInterpreted(b *testing.B) {
	for _, td := range testdataBlockStrings {
		b.Run("", func(b *testing.B) {
			input := []byte(td.input)
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				c := 0
				if err := gqlscan.ScanAll(
					input,
					func(i *gqlscan.Iterator) {
						c++
						if c != td.tokenIndex {
							return
						}
						i.ScanInterpreted(
							td.buffer,
							func(buffer []byte) (stop bool) {
								InterpretedBuffer = buffer
								return false
							},
						)
					},
				); err.IsErr() {
					panic(err)
				}
			}
		})
	}
}
