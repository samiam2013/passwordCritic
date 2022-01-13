package types

import (
	"log"
	"testing"
)

// global filters for testing
var list map[int]BloomFilter

func Init() {
	var err error
	list, err = RebuildFilters()
	if err != nil {
		log.Fatalf("failed building test case: %s", err.Error())
	}
}

func TestBitsetList_LoadFromRebuild(t *testing.T) {
	type fields struct {
		List map[int][]ZeroOneBool
	}
	type args struct {
		filters map[int]BloomFilter
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add non-'happy path' test cases.
		{
			name: "happy path",
			fields: fields{
				List: map[int][]ZeroOneBool{},
			},
			args: args{
				filters: list,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bl := &BitsetList{
				List: tt.fields.List,
			}
			if err := bl.LoadFromRebuild(tt.args.filters); (err != nil) != tt.wantErr {
				t.Errorf("BitsetList.LoadFromRebuild() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBitsetList_WriteToFile(t *testing.T) {
	type fields struct {
		List map[int][]ZeroOneBool
	}
	type args struct {
		pathToFile string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "happy path",
			fields: fields{
				List: map[int][]ZeroOneBool{4: {true, false, false, true}},
			},
			args: args{
				pathToFile: "../cache/test_write.json",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bl := &BitsetList{
				List: tt.fields.List,
			}
			if err := bl.WriteToFile(tt.args.pathToFile); (err != nil) != tt.wantErr {
				t.Errorf("BitsetList.WriteToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
