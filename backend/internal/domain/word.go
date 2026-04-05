package domain

import "time"

type Word struct {
	ID          string   `json:"id"`
	Hanzi       string   `json:"hanzi"`
	Traditional *string  `json:"traditional,omitempty"`
	Pinyin      string   `json:"pinyin"`
	Meanings    []string `json:"meanings"`
	HSKLevel    *int     `json:"hskLevel,omitempty"`
	Frequency   int      `json:"frequency"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type WordRelation struct {
	Word             *Word          `json:"word"`
	RelationshipType RelationType   `json:"relationshipType"`
	Phrase           *string        `json:"phrase,omitempty"`
	Meaning          string         `json:"meaning"`
	Pinyin           *string        `json:"pinyin,omitempty"`
	UsageExample     *string        `json:"usageExample,omitempty"`
	Context          *string        `json:"context,omitempty"`
}

type RelationType string

const (
	RelationTypeCompound  RelationType = "COMPOUND"
	RelationTypePhrase    RelationType = "PHRASE"
	RelationTypeSimilar   RelationType = "SIMILAR"
	RelationTypeOpposite  RelationType = "OPPOSITE"
	RelationTypeComponent RelationType = "COMPONENT"
	RelationTypeModifier  RelationType = "MODIFIER"
	RelationTypeMeasure   RelationType = "MEASURE"
)

type WordGraph struct {
	CenterWord *Word       `json:"centerWord"`
	Nodes      []*WordNode `json:"nodes"`
	Edges      []*Edge     `json:"edges"`
	Depth      int         `json:"depth"`
}

type WordNode struct {
	ID         string  `json:"id"`
	Word       *Word   `json:"word"`
	Distance   int     `json:"distance"`
	Importance float64 `json:"importance"`
}

type Edge struct {
	ID            string       `json:"id"`
	Source        string       `json:"source"`
	Target        string       `json:"target"`
	Type          RelationType `json:"type"`
	Label         string       `json:"label"`
	Weight        float64      `json:"weight"`
	Bidirectional bool         `json:"bidirectional"`
}

type LearningStats struct {
	WordID        string    `json:"wordId"`
	TimesReviewed int       `json:"timesReviewed"`
	CorrectCount  int       `json:"correctCount"`
	LastReviewed  *string   `json:"lastReviewed,omitempty"`
	Strength      float64   `json:"strength"`
}

type WordInput struct {
	Hanzi       string   `json:"hanzi"`
	Traditional *string  `json:"traditional,omitempty"`
	Pinyin      string   `json:"pinyin"`
	Meanings    []string `json:"meanings"`
	HSKLevel    *int     `json:"hskLevel,omitempty"`
	Frequency   *int     `json:"frequency,omitempty"`
}

type RelationInput struct {
	SourceID     string       `json:"sourceId"`
	TargetID     string       `json:"targetId"`
	Type         RelationType `json:"type"`
	Phrase       *string      `json:"phrase,omitempty"`
	Meaning      string       `json:"meaning"`
	Pinyin       *string      `json:"pinyin,omitempty"`
	UsageExample *string      `json:"usageExample,omitempty"`
}