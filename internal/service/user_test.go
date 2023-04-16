package service

import (
	"errors"
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
		wantErr error
	}{
		{"Valid email", args{"example@email.ru", EmailRX}, nil},
		{"Invalid email missing .", args{"example@emailru", EmailRX}, ErrInvalidData},
		{"Invalid email missing @", args{"exampleemailru", EmailRX}, ErrInvalidData},
		{"Invalid email has only домен", args{"@mail.ru", EmailRX}, ErrInvalidData},
		{"Valid email", args{"example@gemail.com", EmailRX}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := matchesPattern(tt.args.value, tt.args.pattern); !errors.Is(err, tt.wantErr) {
				t.Errorf("matchesPattern() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
