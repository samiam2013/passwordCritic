package main

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/samiam2013/passwordcritic/types"
	"github.com/stretchr/testify/assert"
)

func Test_buildFilter(t *testing.T) {
	type args struct {
		bits     int
		filepath io.Reader
	}
	tests := []struct {
		name        string
		args        args
		wantBFilter *types.BloomFilter
		wantErr     bool
	}{
		{
			name: "basic passing test - empty file",
			args: args{
				bits:     123,
				filepath: bytes.NewReader([]byte("")),
			},
			wantBFilter: types.NewBloom(123),
			wantErr:     false,
		},
		// TODO failing cases
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBFilter, err := buildFilter(tt.args.bits, tt.args.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBFilter, tt.wantBFilter) {
				t.Errorf("buildFilter() = %v, want %v", gotBFilter, tt.wantBFilter)
			}
		})
	}
}

func TestRebuildFilters(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("pwd: " + wd)
	tests := []struct {
		name       string
		nFilters   int
		filterLens []int // length instead of []types.BloomFilter
		wantErr    bool
	}{
		{
			name:     "happy path",
			nFilters: 4,
			filterLens: []int{
				12364,
				123641,
				1236416,
				12364167,
			},
			wantErr: false,
		},
		// TODO failing cases
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotten, err := RebuildFilters()
			nFitlersGot := len(gotten)
			if (err != nil) != tt.wantErr {
				t.Errorf("RebuildFilters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(nFitlersGot, tt.nFilters) {
				t.Errorf("RebuildFilters() = %v, want %v", nFitlersGot, tt.nFilters)
			}
			for i, filter := range gotten {
				assert.Equal(t, len(filter.GetBitset()), tt.filterLens[i])
			}
		})
	}
}
