package defrag

import (
	"reflect"
	"strings"
	"testing"

	"github.com/jjtimmons/defrag/config"
)

func TestNewFeatureDB(t *testing.T) {
	db := NewFeatureDB()

	if len(db.features) < 1 {
		t.Fail()
	}
}

func Test_queryFeatures(t *testing.T) {
	type args struct {
		flags *Flags
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			"gather SV40 origin, p10 promoter, mEGFP",
			args{
				&Flags{
					in:  "SV40 origin,p10 promoter,mEGFP",
					dbs: []string{config.AddgeneDB, config.IGEMDB},
				},
			},
			[][]string{
				[]string{"SV40 origin", "ATCCCGCCCCTAACTCCGCCCAGTTCCGCCCATTCTCCGCCCCATGGCTGACTAATTTTTTTTATTTATGCAGAGGCCGAGGCCGCCTCGGCCTCTGAGCTATTCCAGAAGTAGTGAGGAGGCTTTTTTGGAGGCC"},
				[]string{"p10 promoter", "GACCTTTAATTCAACCCAACACAATATATTATAGTTAAATAAGAATTATTATCAAATCATTTGTATATTAATTAAAATACTATACTGTAAATTACATTTTATTTACAATC"},
				[]string{"mEGFP", "AGCAAGGGCGAGGAGCTGTTCACCGGGGTGGTGCCCATCCTGGTCGAGCTGGACGGCGACGTAAACGGCCACAAGTTCAGCGTGCGCGGCGAGGGCGAGGGCGATGCCACCAACGGCAAGCTGACCCTGAAGTTCATCTGCACCACCGGCAAGCTGCCCGTGCCCTGGCCCACCCTCGTGACCACCCTGACCTACGGCGTGCAGTGCTTCAGCCGCTACCCCGACCACATGAAGCAGCACGACTTCTTCAAGTCCGCCATGCCCGAAGGCTACGTCCAGGAGCGCACCATCTCCTTCAAGGACGACGGCACCTACAAGACCCGCGCCGAGGTGAAGTTCGAGGGCGACACCCTGGTGAACCGCATCGAGCTGAAGGGCATCGACTTCAAGGAGGACGGCAACATCCTGGGGCACAAGCTGGAGTACAACTTCAACAGCCACAACGTCTATATCACGGCCGACAAGCAGAAGAACGGCATCAAGGCGAACTTCAAGATCCGCCACAACGTCGAGGACGGCAGCGTGCAGCTCGCCGACCACTACCAGCAGAACACCCCCATCGGCGACGGCCCCGTGCTGCTGCCCGACAACCACTACCTGAGCACCCAGTCCAAGCTGAGCAAAGACCCCAACGAGAAGCGCGATCACATGGTCCTGCTGGAGTTCGTGACCGCCGCCGGGATCACTCTCGGCATGGACGAGCTGTACAAGTAG"},
			},
		},
		{
			"gather SV40 origin, p10 promoter, mEGFP:rev",
			args{
				&Flags{
					in:  "SV40 origin,p10 promoter,mEGFP:rev",
					dbs: []string{config.AddgeneDB, config.IGEMDB},
				},
			},
			[][]string{
				[]string{"SV40 origin", "ATCCCGCCCCTAACTCCGCCCAGTTCCGCCCATTCTCCGCCCCATGGCTGACTAATTTTTTTTATTTATGCAGAGGCCGAGGCCGCCTCGGCCTCTGAGCTATTCCAGAAGTAGTGAGGAGGCTTTTTTGGAGGCC"},
				[]string{"p10 promoter", "GACCTTTAATTCAACCCAACACAATATATTATAGTTAAATAAGAATTATTATCAAATCATTTGTATATTAATTAAAATACTATACTGTAAATTACATTTTATTTACAATC"},
				[]string{"mEGFP:REV", "CTACTTGTACAGCTCGTCCATGCCGAGAGTGATCCCGGCGGCGGTCACGAACTCCAGCAGGACCATGTGATCGCGCTTCTCGTTGGGGTCTTTGCTCAGCTTGGACTGGGTGCTCAGGTAGTGGTTGTCGGGCAGCAGCACGGGGCCGTCGCCGATGGGGGTGTTCTGCTGGTAGTGGTCGGCGAGCTGCACGCTGCCGTCCTCGACGTTGTGGCGGATCTTGAAGTTCGCCTTGATGCCGTTCTTCTGCTTGTCGGCCGTGATATAGACGTTGTGGCTGTTGAAGTTGTACTCCAGCTTGTGCCCCAGGATGTTGCCGTCCTCCTTGAAGTCGATGCCCTTCAGCTCGATGCGGTTCACCAGGGTGTCGCCCTCGAACTTCACCTCGGCGCGGGTCTTGTAGGTGCCGTCGTCCTTGAAGGAGATGGTGCGCTCCTGGACGTAGCCTTCGGGCATGGCGGACTTGAAGAAGTCGTGCTGCTTCATGTGGTCGGGGTAGCGGCTGAAGCACTGCACGCCGTAGGTCAGGGTGGTCACGAGGGTGGGCCAGGGCACGGGCAGCTTGCCGGTGGTGCAGATGAACTTCAGGGTCAGCTTGCCGTTGGTGGCATCGCCCTCGCCCTCGCCGCGCACGCTGAACTTGTGGCCGTTTACGTCGCCGTCCAGCTCGACCAGGATGGGCACCACCCCGGTGAACAGCTCCTCGCCCTTGCT"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := queryFeatures(tt.args.flags); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryFeatures() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_blastFeatures(t *testing.T) {
	type args struct {
		flags          *Flags
		targetFeatures [][]string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"blast a feature against the part databases",
			args{
				flags: &Flags{
					dbs:      []string{config.AddgeneDB, config.IGEMDB},
					filters:  []string{},
					identity: 100.0,
				},
				targetFeatures: [][]string{
					[]string{"SV40 origin", "ATCCCGCCCCTAACTCCGCCCAGTTCCGCCCATTCTCCGCCCCATGGCTGACTAATTTTTTTTATTTATGCAGAGGCCGAGGCCGCCTCGGCCTCTGAGCTATTCCAGAAGTAGTGAGGAGGCTTTTTTGGAGGCC"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := blastFeatures(tt.args.flags, tt.args.targetFeatures, config.New())

			matches := []match{}
			for _, ms := range got {
				for _, m := range ms {
					matches = append(matches, m.match)
				}
			}

			// confirm that the returned fragments sequences contain at least the full queried sequence
			for _, m := range matches {
				containsTargetSeq := false
				for _, wantedSeq := range tt.args.targetFeatures {
					if strings.Contains(m.seq, wantedSeq[1]) {
						containsTargetSeq = true
					}
				}

				if !containsTargetSeq {
					t.Fatalf("match with seq %s doesn't contain any of the target features", m.seq)
				}
			}
		})
	}
}
