package config

import "testing"

func TestParseSupportedCustomResources(t *testing.T) {
	testCases := []struct {
		name       string
		flagValue  string
		expectErr  bool
		expectSize int
	}{
		{name: "empty", flagValue: "", expectErr: false, expectSize: 0},
		{name: "single entry", flagValue: "k8s.sentio.xyz/v1:DriverJob", expectErr: false, expectSize: 1},
		{name: "invalid format", flagValue: "invalid", expectErr: true, expectSize: 0},
	}

	for _, tc := range testCases {
		result, err := ParseSupportedCustomResources(tc.flagValue)
		if tc.expectErr {
			if err == nil {
				t.Fatalf("%s: expected error but got none", tc.name)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tc.name, err)
		}
		if len(result) != tc.expectSize {
			t.Fatalf("%s: expected size %d, got %d", tc.name, tc.expectSize, len(result))
		}
	}
}
