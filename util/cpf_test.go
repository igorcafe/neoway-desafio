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
		{val: "66849734000", wantErr: true},
	}

	for _, tt := range tests {
		err := ValidateCpfOrCnpj(tt.val)
		if (err != nil) != tt.wantErr {
			t.Fatalf("wantErr: %v, err: %v", tt.wantErr, err)
		}
	}
}
