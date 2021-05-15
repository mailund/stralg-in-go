package gostr

import (
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

func testSuffixTree(
	algo string,
	construction func(string) SuffixTree,
	x string,
	t *testing.T) *SuffixTree {
	st := construction(x)

	leaves := LeafLabels(st.Root)
	noLeaves := 0
	prev, ok := <-leaves
	if ok {
		noLeaves++
		for i := range leaves {
			if st.String[prev:] >= st.String[i:] {
				t.Errorf(`We got the leaf "%s" before leaf "%s" in %s("%s").`,
					ReplaceSentinel(st.String[prev:]),
					ReplaceSentinel(st.String[i:]),
					algo, x)
			}
			noLeaves++
			prev = i
		}
	}
	if noLeaves != len(st.String) {
		t.Errorf(`%s("%s"): We got %d leaves but expected %d.\n`,
			algo, x, noLeaves, len(st.String))
	}

	return &st
}

func testSearchMatch(
	algo string,
	st *SuffixTree,
	p string,
	t *testing.T) {
	for i := range st.Search(p) {
		if st.String[i:i+len(p)] != p {
			t.Errorf(`%s("%s"): While searching for "%s" I found "%s".`,
				algo, ReplaceSentinel(st.String), p, ReplaceSentinel(st.String[i:]))
		}
	}
}

func testSearchMismatch(
	algo string,
	st *SuffixTree,
	p string,
	t *testing.T) {
	res := st.Search(p)
	if _, ok := <-res; ok {
		t.Errorf(`We shouldn't find "%s" in %s("%s").`,
			p, algo, ReplaceSentinel(st.String))
	}
}

func testSearchMississippi(
	algo string,
	st *SuffixTree,
	t *testing.T) {
	testSearchMatch(algo, st, "ssi", t)
	testSearchMismatch(algo, st, "x", t)
	testSearchMismatch(algo, st, "spi", t)
}

func TestNaiveConstruction(t *testing.T) {

	x := "mississippi"
	st := testSuffixTree("NaiveST", NaiveST, x, t)
	testSearchMississippi("NaiveST", st, t)
	f, err := os.Create("naive-dot.dot")
	if err != nil {
		log.Fatal(err)
	}
	st.ToDot(f)
	f.Close()
}

func TestMcCreightConstruction(t *testing.T) {

	x := "mississippi"
	st := testSuffixTree("McCreight", McCreight, x, t)
	testSearchMississippi("McCreight", st, t)
	f, err := os.Create("mccreight-dot.dot")
	if err != nil {
		log.Fatal(err)
	}
	st.ToDot(f)
	f.Close()
}

func randomString(n int, alpha string, rng *rand.Rand) string {
	runes := make([]byte, n)
	for i := 0; i < n; i++ {
		runes[i] = alpha[rng.Intn(len(alpha))]
	}
	return string(runes)
}

func TestSTRandomStrings(t *testing.T) {
	algos := []string{"NaiveST", "McCreight"}
	constructors := []func(string) SuffixTree{NaiveST, McCreight}

	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	n := 10       // testing 10 random strings
	maxlen := 100 // max length 100 (so we can still inspect them)
	for i := 0; i < n; i++ {
		slen := rng.Intn(maxlen)
		x := randomString(slen, "acgt", rng)
		t.Logf(`Testing string "%s".`, x)
		for i := range algos {
			testSuffixTree(algos[i], constructors[i], x, t)
		}
	}
}
