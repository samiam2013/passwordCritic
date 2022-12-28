package types

import (
	"log"
	"reflect"
	"testing"
)

// global filters for testing
var list map[Rarity]BloomFilter

func Init() {
	var err error
	list, err = RebuildFilters()
	if err != nil {
		log.Fatalf("failed building test case: %s", err.Error())
	}
}

func TestBitSetMap_LoadFromRebuild(t *testing.T) {
	type fields struct {
		List map[Rarity]BitSet
	}
	type args struct {
		filters map[Rarity]BloomFilter
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
				List: map[Rarity]BitSet{},
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
		List map[Rarity]BitSet
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
				List: map[Rarity]BitSet{
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
			want:    []byte(`{"bitset":"Mdlprs"}`),
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

func TestBitSet_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Set []bool
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "happy path",
			fields: fields{
				Set: []bool{},
			},
			args: args{
				data: []byte(`Mdlprs`),
			},
			want: []bool{
				false, false, false, false, false, true,
				false, false, false, false, true, true,
				false, false, false, true, true, true,
				false, false, true, true, true, true,
				false, true, true, true, true, true,
				true, true, true, true, true, true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := &BitSet{
				Set: tt.fields.Set,
			}
			if err := bs.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("BitSet.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(bs.Set, tt.want) {
				t.Errorf("BitSet.UnmarshalJSON() got = %v, want = %v", bs.Set, tt.want)
			}
		})
	}
}

func TestBitSetMap_UnmarshalJSON(t *testing.T) {
	type fields struct {
		List map[Rarity]BitSet
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    BitSetMap
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			fields: fields{
				List: map[Rarity]BitSet{},
			},
			args: args{
				data: []byte(`{"list":{"bitset":{"10":"Mdlprs","100":"rpldMs"}}}`),
			},
			want: BitSetMap{
				List: map[Rarity]BitSet{
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bl := &BitSetMap{
				List: tt.fields.List,
			}
			if err := bl.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("BitSetMap.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(bl.List, tt.want.List) {
				t.Errorf("BitSetMap.UnmarshalJSON() got = %v, want = %v", bl.List, tt.want.List)
			}
		})
	}
}
