package service

import (
	"regexp"
	"testing"
)

func TestMatchesPatternTableDriven(t *testing.T) {
	type args struct {
		value   string
		pattern *regexp.Regexp
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Valid email", args{"example@email.ru", EmailRX}, false},
		{"Invalid email missing .", args{"example@emailru", EmailRX}, true},
		{"Invalid email missing @", args{"exampleemailru", EmailRX}, true},
		{"Invalid email has only домен", args{"@mail.ru", EmailRX}, true},
		{"Valid email", args{"example@gemail.com", EmailRX}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := matchesPattern(tt.args.value, tt.args.pattern); (err != nil) != tt.wantErr {
				t.Errorf("matchesPattern() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
