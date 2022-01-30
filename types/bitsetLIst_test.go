package types

import (
	"log"
	"reflect"
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

func TestBitSetMap_LoadFromRebuild(t *testing.T) {
	type fields struct {
		List map[int]BitSet
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
				List: map[int]BitSet{},
			},
			args: args{
				filters: list,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bl := &BitSetMap{
				List: tt.fields.List,
			}
			if err := bl.LoadFromRebuild(tt.args.filters); (err != nil) != tt.wantErr {
				t.Errorf("BitSetMap.LoadFromRebuild() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBitSetMap_WriteToFile(t *testing.T) {
	type fields struct {
		List map[int]BitSet
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
				List: map[int]BitSet{
					10: {
						Set: []bool{
							false, false, false, false, false, true,
							false, false, false, false, true, true,
							false, false, false, true, true, true,
							false, false, true, true, true, true,
							false, true, true, true, true, true,
							true, true, true, true, true, true},
					},
					100: {
						Set: []bool{
							false, true, true, true, true, true,
							false, false, true, true, true, true,
							false, false, false, true, true, true,
							false, false, false, false, true, true,
							false, false, false, false, false, true,
							true, true, true, true, true, true},
					},
				},
			},
			args: args{
				pathToFile: "../cache/test_write.json",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bl := &BitSetMap{
				List: tt.fields.List,
			}
			if err := bl.WriteToFile(tt.args.pathToFile); (err != nil) != tt.wantErr {
				t.Errorf("BitSetMap.WriteToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBitSet_MarshalJSON(t *testing.T) {
	type fields struct {
		Set []bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Happy Path",
			fields: fields{
				Set: []bool{
					false, false, false, false, false, true,
					false, false, false, false, true, true,
					false, false, false, true, true, true,
					false, false, true, true, true, true,
					false, true, true, true, true, true,
					true, true, true, true, true, true,
				},
			},
			want:    []byte("{\"bitset\":\"P`hlno\"}"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := BitSet{
				Set: tt.fields.Set,
			}
			got, err := bs.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("BitSet.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				// binaryStr := ""
				// for _, byt := range got {
				// 	binaryStr += fmt.Sprintf("%08b\n", byt)
				// }
				t.Errorf("BitSet.MarshalJSON() = '%s', want '%v'", got, tt.want)
			}
		})
	}
}
