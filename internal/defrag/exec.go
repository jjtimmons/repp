package defrag

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jjtimmons/defrag/config"
	"github.com/spf13/cobra"
)

// Execute is the root of the `defrag seq` functionality
//
// the goal is to find an "optimal" assembly vector with:
// 	1. the fewest fragments
// 	2. the lowest overall assembly cost ($)
// and, secondarily:
//	3. no duplicate end regions between Gibson fragments
// 	4. no inverted repeats in the junctions
// 	5. no off-target binding sites in the parent vectors
//	6. low primer3 penalty scores
func Execute(cmd *cobra.Command, args []string) {
	conf := config.New()

	in, err := cmd.Flags().GetString("in")
	if err != nil {
		log.Fatalf("Cannot parse an input path from args: %v", err)
	}

	output, err := cmd.Flags().GetString("out")
	if err != nil {
		log.Fatalf("Cannot parse an output path from args: %v", err)
	}

	dbs, err := cmd.Flags().GetString("dbs")
	if err != nil {
		log.Fatal(err)
	}

	addgene, err := cmd.Flags().GetBool("addgene")
	if err != nil {
		log.Fatal(err)
	}

	// no path to input file
	if in == "" {
		log.Fatal("Failed: no input file name set")
	}

	// read the target sequence (the first in the slice is used)
	fragments, err := Read(in)
	if err != nil {
		log.Fatalf("Failed to read in fasta files at %s: %v", in, err)
	}

	// set target fragment
	if len(fragments) > 1 {
		println(
			"Warning: %d building fragments were in %s. Only targeting the first: %s",
			len(fragments),
			in,
			fragments[0].ID,
		)
	}
	targetFrag := fragments[0]

	// read in the BLAST DB paths from config
	dbPaths, err := getDBs(dbs, config.Root, addgene)
	fmt.Printf("%+v\n", dbPaths)
	if err != nil {
		log.Fatalf("Failed to find a fragment database: %v", err)
	}

	// get all the matches against the fragment
	matches, err := BLAST(&targetFrag, dbPaths, conf.Fragments.MinHomology, conf.Vendors())
	if err != nil {
		log.Fatalf("Failed to blast %s against the BLAST DB: %v", targetFrag.ID, err)
	}

	// build up the assemblies
	builds := Assemble(matches, targetFrag.Seq, &conf)

	// try to write the JSON to the filepath
	if !filepath.IsAbs(output) {
		output, err = filepath.Abs(output)
		if err != nil {
			log.Fatalf("Failed to make output path absolute: %v", err)
		}
	}
	Write(output, targetFrag, builds)
}

// getDBs returns a list of absolute paths to BLAST databases used during a given run
func getDBs(dbsInput, root string, addgene bool) (paths []string, err error) {
	if addgene {
		addgenePath := path.Join(root, "assets", "addgene", "db", "addgene")
		return parseDBs(dbsInput + "," + addgenePath)
	}
	return parseDBs(dbsInput)
}

// parseDBs turns a single string of comma separated BLAST dbs into a
// slice of absolute paths to the BLAST dbs on the local fs
func parseDBs(dbList string) (paths []string, err error) {
	re := regexp.MustCompile(",")

	for _, db := range re.Split(dbList, -1) {
		dbTrimmed := strings.Trim(db, " ,")
		if dbTrimmed == "" {
			continue
		}

		absPath, err := filepath.Abs(dbTrimmed)

		if err != nil {
			return nil, fmt.Errorf("failed to create absolute path: %v", err)
		}
		paths = append(paths, absPath)
	}

	return
}
