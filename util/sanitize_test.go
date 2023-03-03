package util

import "testing"

func Benchmark_SanitizeCpfOrCnpj_50K(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 50_000; j++ {
			SanitizeCpfOrCnpj("79.379.491/0001-83")
		}
	}
}
