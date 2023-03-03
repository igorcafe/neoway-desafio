package util

import "testing"

// antes do strings.Builder:   27740725 ns/op  9600070 B/op  1350000 allocs/op
// depois do strings.Builder:   4246351 ns/op  1200003 B/op    50000 allocs/op
func Benchmark_SanitizeCpfOrCnpj_50K(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 50_000; j++ {
			SanitizeCpfOrCnpj("79.379.491/0001-83")
		}
	}
}

// antes do strings.Builder:  11316073 ns/op  3200048 B/op   600000 allocs/op
// depois do strings.Builder:  2133852 ns/op   400000 B/op    50000 allocs/op
func Benchmark_SanitizeTicket_50K(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 50_000; j++ {
			SanitizeTicket("1634,00")
		}
	}
}
