package gqlscan_test

import (
	"testing"

	"github.com/romshark/gqlscan"
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
