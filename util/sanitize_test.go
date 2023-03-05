package util

import "testing"

func Test_SanitizeColumns(t *testing.T) {
	got := make([]any, 8)

	tests := []struct {
		cols    []string
		want    []any
		wantErr bool
	}{
		{
			cols:    []string{"987.654.321-00", "1", "0", "2020-10-05", "0,35", "NULL", "NULL", "57.545.869/0001-39"},
			want:    []any{"98765432100", true, false, "2020-10-05", "0.35", nil, nil, "57545869000139"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		err := SanitizeColumns(tt.cols, got)
		if (err != nil) != tt.wantErr {
			t.Errorf("wantErr: %v, err: %v", tt.wantErr, err)
		}

		for i := 0; i < 8; i++ {
			if got[i] != tt.want[i] {
				t.Errorf("at index %d, want: %v, got: %v", i, tt.want[i], got[i])
			}
		}
	}
}

func Test_SanitizeCpfOrCnpj(t *testing.T) {
	tests := []struct {
		val     string
		want    string
		wantErr bool
	}{
		{"987.654.321-00", "98765432100", false},
		{"987¡²³¢^³¢.654.321-00", "98765432100", false},
		{"987.654.321-70", "98765432170", true},
	}

	for _, tt := range tests {
		got, err := SanitizeCpfOrCnpj(tt.val)
		if (err != nil) != tt.wantErr {
			t.Errorf("wantErr: %v, err: %v", tt.wantErr, err)
		}

		if got != tt.want {
			t.Errorf("want: %s, got: %s", tt.want, got)
		}
	}
}

func Test_SanitizeTicket(t *testing.T) {
	tests := []struct {
		val  string
		want string
	}{
		{"12,90", "12.90"},
		{"0,00", "0.00"},
		{"1.500,47", "1500.47"},
	}

	for _, tt := range tests {
		got := SanitizeTicket(tt.val)
		if got != tt.want {
			t.Errorf("want: %s, got: %s", tt.want, got)
		}
	}
}

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

// antes de extrair o slice `res`:  36033203 ns/op	15200101 B/op   600000 allocs/op
// depois de extrair o slice `res`: 36005822 ns/op   8800089 B/op   550000 allocs/op
func Benchmark_SanitizeColumns_50K(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := make([]any, 8)

		for j := 0; j < 500_000; j++ {
			SanitizeColumns([]string{
				"058.189.421-98",
				"1",
				"0",
				"2023-03-03",
				"0,59",
				"1.000.000,12",
				"79.379.491/0001-83",
				"79.379.491/0001-83",
			}, res)
		}
	}
}
