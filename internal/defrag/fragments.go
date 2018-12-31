package defrag

import (
	"log"
	"os"
	"strings"

	"github.com/jjtimmons/defrag/config"
	"github.com/spf13/cobra"
)

// Fragments accepts a cobra.Command with flags for assembling a list of
// fragments together into a vector (in the order specified). Fragments
// without junctions for their neighbors are prepared via PCR
func Fragments(cmd *cobra.Command, args []string) {
	defer os.Exit(0)

	input, err := parseFlags(cmd)
	if err != nil {
		log.Fatalln(err)
	}

	fragments(input)
}

// fragments pieces together a list of fragments into a single vector
// with the fragments in the order and orientation specified
func fragments(input *flags) {
	conf := config.New()

	// read the target sequence (the first in the slice is used)
	inputFragments, err := read(input.in)
	if err != nil {
		log.Fatalf("failed to read in fasta files at %s: %v", input.in, err)
	}

	// try to find the target vector (sequence) and prepare the fragments to build it
	target, fragments := assembleFragments(inputFragments, &conf)

	// write the single list of fragments as a possible solution to the output file
	if err := write(input.out, target, [][]Fragment{fragments}); err != nil {
		log.Fatal(err)
	}
}

// assembleFragments takes a list of Fragments and returns the Vector we assume the user is
// trying to build as well as the Fragments (possibly prepared via PCR)
func assembleFragments(inputFragments []Fragment, conf *config.Config) (targetVector Fragment, fragments []Fragment) {
	if len(inputFragments) < 1 {
		log.Fatalln("failed: no fragments to assemble")
	}

	// convert the fragments to nodes (without a start and end)
	nodes := make([]*node, len(inputFragments))
	for i, f := range inputFragments {
		nodes[i] = &node{
			id:      f.ID,
			seq:     f.Seq,
			fullSeq: f.Seq,
			conf:    conf,
			start:   0,
			end:     0,
		}
	}

	// find out how much overlap the *last* node has with its next one
	// set the start, end, and vector sequence based on that
	//
	// add all of each nodes seq to the vector sequence, minus the region overlapping the next
	minHomology := conf.Fragments.MinHomology
	maxHomology := conf.Fragments.MaxHomology
	junction := nodes[len(nodes)-1].junction(nodes[0], minHomology, maxHomology)
	var vectorSeq strings.Builder
	for i, n := range nodes {
		// correct for this node's overlap with the last node
		n.start = vectorSeq.Len() - len(junction)
		n.end = n.start + len(n.seq) - 1

		// find the junction between this node and the next (if there is one)
		junction = n.junction(nodes[(i+1)%len(nodes)], minHomology, maxHomology)

		// add this node's sequence onto the accumulated vector sequence
		vectorSeq.WriteString(n.seq[0 : len(n.seq)-len(junction)])
	}

	// create the assumed vector object
	targetVector = Fragment{
		Seq:  vectorSeq.String(),
		Type: circular,
	}

	// create an assembly out of the nodes (to fill/convert to fragments with primers)
	a := assembly{nodes: nodes}
	fragments, err := a.fill(targetVector.Seq, conf)
	if err != nil {
		log.Fatalf("failed to fill in the nodes: %+v", err)
	}
	return targetVector, fragments
}
