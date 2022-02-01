package types

import (
	"os"
	"strings"
	"testing"
)

func TestDownloadLists(t *testing.T) {
	tests := []struct {
		name        string
		wantPwLists map[int]string
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name: "happy path",
			wantPwLists: map[int]string{
				1000:    getCacheFolder()+"1000.txt",
				10_000:  getCacheFolder()+"10000.txt",
				100_000: getCacheFolder()+"100000.txt",
				// 1_000_000: "../cache/1000000.txt",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPwLists, err := DownloadLists()
			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadLists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for num, suffix := range tt.wantPwLists{
				if _, ok := gotPwLists[num]; !ok{
					t.Errorf("couldn't find list with %d elems!", num)
				}
				if !strings.HasSuffix(gotPwLists[num], suffix) {
					t.Errorf("got file path = %s did not contain prefix = %s", gotPwLists[num], tt.wantPwLists[num])
				}
				//t.Errorf("DownloadLists() = %v, want %v", gotPwLists, tt.wantPwLists)
			}
		})
	}
}

func Test_dlFileTo(t *testing.T) {
	type args struct {
		filepath string
		url      string
	}
	tests := []struct {
		name           string
		args           args
		wantMinWritten int64
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name: "happy path",
			args: args{
				filepath: "../cache/dummyFilePostmanEcho.json",
				url:      "https://postman-echo.com/get",
			},
			wantMinWritten: 100,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWritten, err := dlFileTo(tt.args.filepath, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("dlFileTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWritten < tt.wantMinWritten {
				t.Errorf("dlFileTo() = %v, want min %v", gotWritten, tt.wantMinWritten)
			}
		})
		if err := os.Remove(tt.args.filepath); err != nil {
			t.Errorf("failed removing file after testing dlFileTo(): %s", err.Error())
		}
	}
}
