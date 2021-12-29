package main

import (
	"bytes"
	"io"
	"reflect"
	"strings"
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
		wantErrPrefix string
	}{
		// TODO: Add test cases.
		{
			name: "meets minimums",
			args: args{pwCandPtr: strPtr("pointerToLiteral")},
			wantCandidate: critic.PassCandidate{
				StringVal:   "pointerToLiteral",
				Cardinality: 11,
				H:           3.375,
			},
			wantErr:       false,
			wantErrPrefix: "", //&types.HomogeneityError{},
		},
		{
			name: "entropy too low",
			args: args{pwCandPtr: strPtr("password")},
			wantCandidate: critic.PassCandidate{
				StringVal:   "password",
				Cardinality: 7,
				H:           2.75,
			},
			wantErr:       true,
			wantErrPrefix: "low entropy",
		},
		{
			name: "repitition too high",
			args: args{pwCandPtr: strPtr("ppsswwrrdd")},
			wantCandidate: critic.PassCandidate{
				StringVal:   "ppsswwrrdd",
				Cardinality: 5,
				H:           2.321928,
			},
			wantErr:       true,
			wantErrPrefix: "low variety", //"high repitition",
		},
		{
			name: "h too low & repitition too high",
			args: args{pwCandPtr: strPtr("12345566")},
			wantCandidate: critic.PassCandidate{
				StringVal:   "12345566",
				Cardinality: 6,
				H:           2.5,
			},
			wantErr:       true,
			wantErrPrefix: "low entropy",
		},
		{
			name: "good password",
			args: args{pwCandPtr: strPtr("4D5f2A8E0fa3D9162dbAcfA543C730c80F980b92d60b833f2Ec97418c39e9")},
			wantCandidate: critic.PassCandidate{
				StringVal:   "4D5f2A8E0fa3D9162dbAcfA543C730c80F980b92d60b833f2Ec97418c39e9",
				Cardinality: 21,
				H:           4.184778,
			},
			wantErr:       false,
			wantErrPrefix: "high repititon",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCandidate, err := checkEntropy(tt.args.pwCandPtr)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkEntropy() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !strings.HasPrefix(err.Error(), tt.wantErrPrefix) {
				t.Errorf("got unexpected error prefix: wanted '%s'; got type '%s';",
					tt.wantErrPrefix, err.Error())
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

func Test_getStdIn(t *testing.T) {
	const mockedInput = "the quick brown fox correct horse staple"
	var stdin bytes.Buffer
	stdin.Write([]byte(mockedInput))
	type args struct {
		file io.Reader
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: mockedInput,
			args: args{file: &stdin},
			want: mockedInput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStdIn(tt.args.file); got != tt.want {
				t.Errorf("getStdIn() = %v, want %v", got, tt.want)
			}
		})
	}
}
