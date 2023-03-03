package util

import "testing"

func Test_ValidateCpfOrCnpj(t *testing.T) {

	tests := []struct {
		val     string
		wantErr bool
	}{
		{val: "", wantErr: true},
		{val: "0123456789", wantErr: true},
		{val: "012345678910", wantErr: true},
		{val: "0123456789101213", wantErr: true},

		{val: "66849734008", wantErr: false},
		{val: "66849734018", wantErr: true}, // alterar dígito 10 deve causar erro
		{val: "66849734005", wantErr: true}, // alterar dígito 11 deve causar erro

		{val: "45091647007", wantErr: false},
		{val: "45291647007", wantErr: true}, // alterar qualquer dígito deve causar erro
		{val: "45091607007", wantErr: true}, // alterar qualquer dígito deve causar erro

		// CPFs com dígitos iguais são matematicamente válidos, mas não devem ser aceitos
		{val: "11111111111", wantErr: true},
		{val: "22222222222", wantErr: true},
		{val: "99999999999", wantErr: true},
		{val: "00000000000", wantErr: true},
	}

	for _, tt := range tests {
		err := ValidateCpfOrCnpj(tt.val)
		if (err != nil) != tt.wantErr {
			t.Fatalf("val: %s, wantErr: %v, err: %v", tt.val, tt.wantErr, err)
		}
	}
}
