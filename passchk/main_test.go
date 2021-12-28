package main

import (
	"reflect"
	"testing"

	"github.com/samiam2013/passwordcritic/critic"
)

// TODO : evaluate whether this mostly-generated code meets needs

// Test_checkEntropy _
func Test_checkEntropy(t *testing.T) {
	type args struct {
		pwCandPtr *string
	}
	tests := []struct {
		name          string
		args          args
		wantCandidate critic.PassCandidate
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			name: "testName?",
			args: args{pwCandPtr: strPtr("pointerToLiteral")},
			wantCandidate: critic.PassCandidate{
				StringVal:   "pointerToLiteral",
				Cardinality: 11,
				H:           3.375,
			},
			wantErr: true,
		},
		{
			name: "testName2?",
			args: args{pwCandPtr: strPtr("4D5f2A8E0fa3D9162dbAcfA543C730c80F980b92d60b833f2Ec97418c39e9")},
			wantCandidate: critic.PassCandidate{
				StringVal:   "4D5f2A8E0fa3D9162dbAcfA543C730c80F980b92d60b833f2Ec97418c39e9",
				Cardinality: 21,
				H:           4.184778,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCandidate, err := checkEntropy(tt.args.pwCandPtr)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkEntropy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCandidate, tt.wantCandidate) {
				t.Errorf("checkEntropy() = %v, want %v", gotCandidate, tt.wantCandidate)
			}
		})
	}
}

func strPtr(input string) (ptr *string) {
	ptr = &input
	return
}
