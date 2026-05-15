package rank

type Stream struct {
	Length     int            `json:"length"`
	TermCounts map[string]int `json:"term_counts"`
}

func NewStream(tokens []string) (s Stream) {
	s.Length = len(tokens)
	s.TermCounts = make(map[string]int)
	for _, token := range tokens {
		s.TermCounts[token]++
	}
	return
}

type Document struct {
	Attachments map[string]string `json:"attachments"`
	Streams     map[Field]Stream  `json:"streams"`
}

func NewDocument() (d Document) {
	d.Attachments = make(map[string]string)
	d.Streams = make(map[Field]Stream)
	return
}

// SetAttachment sets attached data that is not parsed or used by BM25F.
func (d Document) SetAttachment(id string, text string) {
	d.Attachments[id] = text
}

func (d Document) SetStream(field Field, tokens []string) {
	d.Streams[field] = NewStream(tokens)
}

type Corpus struct {
	Documents map[string]Document `json:"documents"`

	// DocsWithTerm maps a term to the number of documents containing it.
	DocsWithTerm map[string]int `json:"docs_with_term"`

	// TotalLengths are the total lengths of each field across all documents.
	TotalLengths map[Field]int `json:"total_lengths"`
}

func NewCorpus() Corpus {
	return Corpus{
		Documents:    make(map[string]Document),
		DocsWithTerm: make(map[string]int),
		TotalLengths: make(map[Field]int),
	}
}

// Remove removes all data associated with a document.
func (c *Corpus) Remove(id string) {
	if doc, ok := c.Documents[id]; ok {
		delete(c.Documents, id)

		for field, stream := range doc.Streams {
			c.TotalLengths[field] -= stream.Length
		}

		for _, stream := range doc.Streams {
			for term := range stream.TermCounts {
				c.DocsWithTerm[term]--
				if c.DocsWithTerm[term] == 0 {
					delete(c.DocsWithTerm, term)
				}
			}
		}
	}
}

// Upsert processes and adds a document into the corpus.
func (c *Corpus) Upsert(id string, document Document) {
	c.Remove(id)

	c.Documents[id] = document

	for field, stream := range document.Streams {
		c.TotalLengths[field] += stream.Length
	}

	for _, stream := range document.Streams {
		for term := range stream.TermCounts {
			c.DocsWithTerm[term]++
		}
	}
}

func (c *Corpus) AvgStreamLength(field Field) float64 {
	if docCount := len(c.Documents); docCount > 0 {
		return float64(c.TotalLengths[field]) / float64(docCount)
	}
	return 0
}
