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

	conf := config.New()

	input, err := parseFlags(cmd, conf)
	if err != nil {
		log.Fatalln(err)
	}

	fragments(input, conf)
}

// fragments pieces together a list of fragments into a single vector
// with the fragments in the order and orientation specified
func fragments(input *flags, conf *config.Config) {
	// read in the consituent fragments
	inputFragments, err := read(input.in)
	if err != nil {
		log.Fatalf("failed to read in fasta files at %s: %v", input.in, err)
	}

	// add in the backbone if it was provided
	if input.backbone.ID != "" {
		inputFragments = append([]Frag{input.backbone}, inputFragments...)
	}

	// piece together the adjacent fragments
	target, fragments := assembleFragments(inputFragments, conf)

	// write the single list of fragments as a possible solution to the output file
	if err := write(input.out, target, [][]Frag{fragments}); err != nil {
		log.Fatal(err)
	}
}

// assembleFragments takes a list of Fragments and returns the Vector we assume the user is
// trying to build as well as the Fragments (possibly prepared via PCR)
func assembleFragments(inputFragments []Frag, conf *config.Config) (targetVector Frag, fragments []Frag) {
	if len(inputFragments) < 1 {
		log.Fatalln("failed: no fragments to assemble")
	}

	// convert the fragments to frags (without a start and end and with the conf)
	frags := make([]*Frag, len(inputFragments))
	for i, f := range inputFragments {
		frags[i] = &Frag{
			ID:      f.ID,
			Seq:     f.Seq,
			fullSeq: f.Seq,
			conf:    conf,
			start:   0,
			end:     0,
			Type:    existing,
		}
	}

	// find out how much overlap the *last* Frag has with its next one
	// set the start, end, and vector sequence based on that
	//
	// add all of each frags seq to the vector sequence, minus the region overlapping the next
	minHomology := conf.Fragments.MinHomology
	maxHomology := conf.Fragments.MaxHomology
	junction := frags[len(frags)-1].junction(frags[0], minHomology, maxHomology)
	var vectorSeq strings.Builder
	for i, n := range frags {
		// correct for this Frag's overlap with the last Frag
		n.start = vectorSeq.Len() - len(junction)
		n.end = n.start + len(n.Seq) - 1

		// find the junction between this Frag and the next (if there is one)
		junction = n.junction(frags[(i+1)%len(frags)], minHomology, maxHomology)

		// add this Frag's sequence onto the accumulated vector sequence
		vectorSeq.WriteString(n.Seq[0 : len(n.Seq)-len(junction)])
	}

	// create the assumed vector object
	targetVector = Frag{
		Seq:  vectorSeq.String(),
		Type: circular,
	}

	// create an assembly out of the frags (to fill/convert to fragments with primers)
	a := assembly{frags: frags}
	fragments, err := a.fill(targetVector.Seq, conf)
	if err != nil {
		log.Fatalf("failed to fill in the frags: %+v", err)
	}
	return targetVector, fragments
}