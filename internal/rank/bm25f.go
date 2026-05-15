package rank

import (
	"cmp"
	"math"
	"slices"
)

type BM25F struct {
	k1           float64
	fieldConfigs map[Field]FieldConfig
}

func NewBM25F(k1 float64, fieldConfigs map[Field]FieldConfig) BM25F {
	return BM25F{k1, fieldConfigs}
}

type Result struct {
	Name     string
	Document Document

	// score is not visible because it only has relative meaning for ranking.
	score float64
}

// Rank returns document results sorted by how well they match the query.
// The best match is first. Equal matches are sorted lexigraphically by id.
func (bm BM25F) Rank(corpus Corpus, query []string) []Result {
	// Deduplicate query
	slices.Sort(query)
	query = slices.Compact(query)

	// Init the results with document data and 0 scores.
	results := make([]Result, 0, len(corpus.Documents))
	for name, doc := range corpus.Documents {
		results = append(results, Result{
			Name:     name,
			Document: doc,
			score:    0.0,
		})
	}

	// A term's score for a document is its overall importance (idf) times its
	// saturation within the document. These scores are summed per document for
	// the final document scores.
	for _, term := range query {
		idf := bm.idf(corpus, term)
		for i := range results {
			result := &results[i]
			termFreq := bm.termFrequency(corpus, result.Document, term)
			saturation := termFreq / (bm.k1 * idf)
			result.score = saturation * idf
		}
	}

	// Sort the results descending by score.
	slices.SortFunc(results, func(a, b Result) int {
		if c := cmp.Compare(b.score, a.score); c != 0 {
			return c
		}
		return cmp.Compare(a.Name, b.Name)
	})

	return results
}

// idf returns the relative importance of a word based on its rarity.
func (bm BM25F) idf(corpus Corpus, term string) float64 {
	// For the IDF, we apply a modified Robertson/Sparck Jones formula across
	// all streams. There are rare scenarios where this does not yield good
	// results. We will ignore the problem until it shows itself in practice.
	docCount := float64(len(corpus.Documents))
	docFreq := float64(corpus.DocsWithTerm[term])
	return math.Log((docCount-docFreq+0.5)/(docFreq+0.5) + 1)
}

// termFrequency returns the normalized weighted frequency of a term within the
// document across all streams.
func (bm BM25F) termFrequency(corpus Corpus, doc Document, term string) (result float64) {
	for field, config := range bm.fieldConfigs {
		avgStreamLen := corpus.AvgStreamLength(field)
		if avgStreamLen == 0.0 {
			continue
		}

		stream := doc.Streams[field]

		// Normalize results when the stream length is far from average.
		streamLen := float64(stream.Length)
		lengthNorm := 1 - config.B + config.B*streamLen/avgStreamLen

		// Simple weighted summation with normalization.
		termFreq := float64(stream.TermCounts[term])
		result += config.Weight * termFreq / lengthNorm
	}
	return
}
