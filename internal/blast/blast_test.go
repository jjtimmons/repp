package blast

import (
	"fmt"
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jjtimmons/defrag/internal/defrag"
)

// test the ability to find test fragments in a mock database
// see test/blast/README.md for a description of where the subfragments
// in this test fragment's sequence came from (pieces from the 5 fragments)
// that make up the mock BLAST db
func Test_BLAST(t *testing.T) {
	// make path to test db
	testDB, _ := filepath.Abs(path.Join(conf.Root, "test", "blast", "db"))
	blastDir, _ := filepath.Abs(path.Join("..", "..", "bin", "blast"))

	// create mock test fragment
	f := defrag.Fragment{
		ID:  "test_target",
		Seq: "GGCCGCAATAAAATATCTTTATTTTCATTACATCTGTGTGTTGGTTTTTTGTGTGAATCGATAGTACTAACATGACCACCTTGATCTTCATGGTCTGGGTGCCCTCGTAGGGCTTGCCTTCGCCCTCGGATGTGCACTTGAAGTGGTGGTTGTTCACGGTGCCCTCCATGTACAGCTTCATGTGCATGTTCTCCTTGATCAGCTCGCTCATAGGTCCAGGGTTCTCCTCCACGTCTCCAGCCTGCTTCAGCAGGCTGAAGTTAGTAGCTCCGCTTCCGGATCCCCCGGGGAGCATGTCAAGGTCAAAATCGTCAAGAGCGTCAGCAGGCAGCATATCAAGGTCAAAGTCGTCAAGGGCATCGGCTGGGAgCATGTCTAAgTCAAAATCGTCAAGGGCGTCGGCCGGCCCGCCGCTTTcgcacGCCCTGGCAATCGAGATGCTGGACAGGCATCATACCCACTTCTGCCCCCTGGAAGGCGAGTCATGGCAAGACTTTCTGCGGAACAACGCCAAGTCATTCCGCTGTGCTCTCCTCTCACATCGCGACGGGGCTAAAGTGCATCTCGGCACCCGCCCAACAGAGAAACAGTACGAAACCCTGGAAAATCAGCTCGCGTTCCTGTGTCAGCAAGGCTTCTCCCTGGAGAACGCACTGTACGCTCTGTCCGCCGTGGGCCACTTTACACTGGGCTGCGTATTGGAGGATCAGGAGCATCAAGTAGCAAAAGAGGAAAGAGAGACACCTACCACCGATTCTATGCCTGACTGTGGCGGGTGAGCTTAGGGGGCCTCCGCTCCAGCTCGACACCGGGCAGCTGCTGAAGATCGCGAAGAGAGGGGGAGTAACAGCGGTAGAGGCAGTGCACGCCTGGCGCAATGCGCTCACCGGGGCCCCCTTGAACCTGACCCCAGACCAGGTAGTCGCAATCGCGAACAATAATGGGGGAAAGCAAGCCCTGGAAACCGTGCAAAGGTTGTTGCCGGTCCTTTGTCAAGACCACGGCCTTACACCGGAGCAAGTCGTGGCCATTGCAAGCAATGGGGGTGGCAAACAGGCTCTTGAGACGGTTCAGAGACTTCTCCCAGTTCTCTGTCAAGCCGTTGGAGTCCACGTTCTTTAATAGTGGACTCTTGTTCCAAACTGGAACAACACTCAACCCTATCTCGGTCTATTCTTTTGATTTATAAGGGATTTTGCCGATTTCGGCCTATTGGTTAAAAAATGAGCTGATTTAACAAAAATTTAACGCGAATTTTAACAAAATATTAACGCTTACAATTTAGGTGGCACTTTTCGGGGAAATGTGCGCGGAACCCCTATTTGTTTATTTTTCTAAATACATTCAAATATGTATCCGCTCATGAGACAATAACCCTGATAAATGCTTCAATAATATTGAAAAAGGAAGAGTATGAGTATTCAACATTTCCGTGTCGCCCTTATTCCCTTTTTTGCGGCATTTTGCCTTCCTGTTTTTGCTCACCCAGAAACGCTGGTGAAAGTAAAAGATGCTGAAGATCAGTTGGGTGCACGAGTGGGTTACATCGAACTGGATCTCAACAGCGGTAAGATCCTTGAGAGTTTTCGCCCCGAAGAACGTTTTCCAATGATGAGCACTTTTAAAGTTCTGCTATGTGGCGCGGTATTATCCCGTATTGACGCCGGGCAAGAGCAACTCGGTCGCCGCATACACTATTCTCAGAATGACTTGGTTGAGTACTCACCAGTCACAGAAAAGCATCTTACGGATGGCATGACAGTAAGAGAATTATGCAGTGCTGCCATAACCATGAGTGATAACACTGCGGCCAACTTACTTCTGACAACGATCGGAGGACCGAAGGAGCTAACCGCTTTTTTGCACAACATGGGGGATCATGTAACTCGCCTTGATCGTTGGGAACCGGAGCTGAATGAAGCCATACCAAACGACGAGCGTGACACCACGATGCCTGTAGCAATGGCAACAACGTTGCGCAAACTATTAACTGGCGAACTACTTACTCTAGCTTCCCGGCAACAATTAATAGACTGGATGGAGGCGGATAAAGTTGCAGGACCACTTCTGCGCTCGGCCCTTCCGGCTGGCTGGTTTATTGCTGATAAATCTGGAGCCGGTGAGCGTGGGTCTCGCGGTATCATTGCAGCACTGGGGCCAGATGGTAAGCCCTCCCGTATCGTAGTTATCTACACGACGGGGAGTCAGGCAACTATGGATGAACGAAATAGACAGATCGCTGAGATAGGTGCCTCACTGATTAAGCATTGGTAACTGTCAGACCAAGTTTACTCATATATACTTTAGATTGATTTAAAACTTCATTTTTAATTTAAAAGGATCTAGGTGAAGATCCTTTTTGATAATCTCATGACCAAAATCCCTTAACGTGAGTTTTCGTTCCACTGAGCGTCAGACCCCGTAGAA",
	}

	// run blast
	matches, err := BLAST(&f, testDB, blastDir, 10) // any match over 10 bp

	// check if it fails
	if err != nil {
		t.Errorf("failed to run BLAST: %v", err)
		return
	}

	// make sure matches are found
	if len(matches) < 1 {
		t.Error("failed to find any matches")
		return
	}

	matchesContain := func(targ defrag.Match) {
		for _, m := range matches {
			if targ.Entry == m.Entry && targ.Start == m.Start && targ.End == m.End {
				return
			}
		}

		t.Errorf("failed to find match %v in fragment matches", targ)
	}

	matchesContain(defrag.Match{
		Entry: "gnl|addgene|107006(circular)",
		Start: 0,
		End:   72,
	})
}

func Test_parseDBs(t *testing.T) {
	db1 := "../exampleDir/exampleDB.fa"
	dbAbs1, _ := filepath.Abs(db1)

	db2 := "otherBLASTDir/otherDB.fa"
	dbAbs2, _ := filepath.Abs(db2)

	type args struct {
		dbList string
	}
	tests := []struct {
		name      string
		args      args
		wantPaths []string
		wantError error
	}{
		{
			"single blast path",
			args{
				dbList: db1,
			},
			[]string{dbAbs1},
			nil,
		},
		{
			"multi fasta db paths",
			args{
				dbList: fmt.Sprintf("%s, %s", db1, db2),
			},
			[]string{dbAbs1, dbAbs2},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPaths, _ := parseDBs(tt.args.dbList); !reflect.DeepEqual(gotPaths, tt.wantPaths) {
				t.Errorf("parseDBs() = %v, want %v", gotPaths, tt.wantPaths)
			}
		})
	}
}
