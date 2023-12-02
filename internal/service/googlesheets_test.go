package service

import (
	"errors"
	"testing"
)

var (
	testValues = [][]any{
		{"name", "balance"},
		{"John", 100},
		{"Jane", 200},
	}
)

func TestTakeByValue(t *testing.T) {
	testcases := []struct {
		name              string
		values            [][]any
		searchHeaderIndex int
		takeHeaderIndex   int
		value             string
		caseSensitive     bool
		want              any
		wantErr           error
	}{
		{
			name:              "should return error when value not found",
			values:            testValues,
			searchHeaderIndex: 0,
			takeHeaderIndex:   1,
			value:             "Mary",
			wantErr:           ErrValueNotFound,
		},
		{
			name: "columns have different length",
			values: [][]any{
				{"name", "balance", "age"},
				{"John", 100},
				{"Jane", 200, 30},
			},
			searchHeaderIndex: 0,
			takeHeaderIndex:   2,
			value:             "Jane",
			want:              30,
		},
		{
			name: "case sensitive is true",
			values: [][]any{
				{"name", "balance"},
				{"John", 100},
				{"Jane", 200},
			},
			searchHeaderIndex: 0,
			takeHeaderIndex:   1,
			value:             "jane",
			caseSensitive:     true,
			wantErr:           ErrValueNotFound,
		},
		{
			name:              "case sensitive is false",
			values:            testValues,
			searchHeaderIndex: 0,
			takeHeaderIndex:   1,
			value:             "jane",
			caseSensitive:     false,
			want:              200,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := takeByValue(tc.values, tc.searchHeaderIndex, tc.takeHeaderIndex, tc.value, tc.caseSensitive)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("takeByValue() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if got != tc.want {
				t.Errorf("takeByValue() got = %v, want %v", got, tc.want)
			}
		})
	}
}
