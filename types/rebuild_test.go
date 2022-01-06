package types

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

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
		wantBFilter *BloomFilter
		wantErr     bool
	}{
		{
			name: "basic passing test - empty file",
			args: args{
				bits:     123,
				filepath: bytes.NewReader([]byte("")),
			},
			wantBFilter: NewBloom(123),
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
		filterLens map[int]int // length instead of []types.BloomFilter
		clearCache bool        // whether or not cache should be cleared before rebuild to create error cases
		wantErr    bool
	}{
		{
			name:     "happy path",
			nFilters: 4,
			filterLens: map[int]int{
				1_000_000: 12364167,
				100_000:   1236416,
				10_000:    123641,
				1000:      12364,
			},
			clearCache: false,
			wantErr:    false,
		},
		{
			name:       "sad path",
			nFilters:   0,
			filterLens: map[int]int{},
			clearCache: true,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.clearCache {
				err = filepath.Walk(CacheFolderPath, deletePWFiles)
				if err != nil {
					t.Fatalf(err.Error())
				}
			}

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

			if tt.clearCache {
				if _, err = DownloadLists(); err != nil {
					t.Fatalf("error rebuilding cache after clearing for test: %s", err.Error())
				}
			}
		})
	}
}

func deletePWFiles(path string, f os.FileInfo, err error) (e error) {
	if err != nil {
		return fmt.Errorf("error at call time (err parameter) to deletePWFiles: %s", err.Error())
	}
	if strings.HasSuffix(filepath.Base(path), "00.txt") {
		err := os.Remove(path)
		if err != nil {
			e = err
			return
		}
	}
	return nil
}
