package scraper

import (
	"testing"
)

func TestGetReleases(t *testing.T) {
	nal := make([]string, 0)
	nnal := make([]string, 0)
	nnal = append(nnal, "Example Artist")

	var tests = []struct {
		al      []string
		co      Conf
		wantErr bool
		name    string
	}{
		{nal, Conf{"30d", "9"}, true, "empty artists list"},
		{nnal, Conf{"90d", "7"}, true, "incompatible time frame"},
		{nnal, Conf{"30d", "9"}, false, "comopatible paramaters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetReleases(tt.al, tt.co)
			if (err == nil) && tt.wantErr {
				t.Errorf("returned no error but one was expected")
			}
			if (err != nil) && !tt.wantErr {
				t.Errorf("returned an error but one was not expected")
			}
		})
	}
}
