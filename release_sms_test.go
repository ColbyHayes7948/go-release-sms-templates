package main

import "testing"

func TestPublishReleaseValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   Release
		wantErr bool
	}{
		{"missing name", Release{Phone: "13800000000"}, true},
		{"missing phone", Release{Name: "1.2.3"}, true},
		{"complete", Release{Name: "1.2.3", Phone: "13800000000"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRelease(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error=%v wantErr=%v", err, tt.wantErr)
			}
		})
	}
}
