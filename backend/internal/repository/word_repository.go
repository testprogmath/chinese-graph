package repository

import (
	"context"
	"fmt"

	"github.com/chinese-graph/backend/internal/domain"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type WordRepository interface {
	GetByHanzi(ctx context.Context, hanzi string) (*domain.Word, error)
	GetByID(ctx context.Context, id string) (*domain.Word, error)
	Create(ctx context.Context, word *domain.Word) (*domain.Word, error)
	Update(ctx context.Context, word *domain.Word) (*domain.Word, error)
	Delete(ctx context.Context, id string) error
	
	GetWordGraph(ctx context.Context, hanzi string, depth int, maxNodes int) (*domain.WordGraph, error)
	Search(ctx context.Context, query string, limit int) ([]*domain.Word, error)
	GetByHSKLevel(ctx context.Context, level int, limit int) ([]*domain.Word, error)
	
	AddRelation(ctx context.Context, relation *domain.RelationInput) error
	RemoveRelation(ctx context.Context, sourceID, targetID string, relationType domain.RelationType) error
}

type Neo4jWordRepository struct {
	session neo4j.SessionWithContext
}

func NewNeo4jWordRepository(session neo4j.SessionWithContext) WordRepository {
	return &Neo4jWordRepository{session: session}
}

func (r *Neo4jWordRepository) GetByHanzi(ctx context.Context, hanzi string) (*domain.Word, error) {
	query := `
		MATCH (w:Word {hanzi: $hanzi})
		RETURN w.id as id, w.hanzi as hanzi, w.traditional as traditional, 
		       w.pinyin as pinyin, w.meanings as meanings, w.hskLevel as hskLevel,
		       w.frequency as frequency
	`
	
	result, err := r.session.Run(ctx, query, map[string]interface{}{
		"hanzi": hanzi,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get word by hanzi: %w", err)
	}

	if result.Next(ctx) {
		record := result.Record()
		return r.recordToWord(record), nil
	}

	return nil, nil
}

func (r *Neo4jWordRepository) GetByID(ctx context.Context, id string) (*domain.Word, error) {
	query := `
		MATCH (w:Word {id: $id})
		RETURN w.id as id, w.hanzi as hanzi, w.traditional as traditional, 
		       w.pinyin as pinyin, w.meanings as meanings, w.hskLevel as hskLevel,
		       w.frequency as frequency
	`
	
	result, err := r.session.Run(ctx, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get word by id: %w", err)
	}

	if result.Next(ctx) {
		record := result.Record()
		return r.recordToWord(record), nil
	}

	return nil, nil
}

func (r *Neo4jWordRepository) Create(ctx context.Context, word *domain.Word) (*domain.Word, error) {
	query := `
		CREATE (w:Word {
			id: randomUUID(),
			hanzi: $hanzi,
			traditional: $traditional,
			pinyin: $pinyin,
			meanings: $meanings,
			hskLevel: $hskLevel,
			frequency: $frequency
		})
		RETURN w.id as id, w.hanzi as hanzi, w.traditional as traditional, 
		       w.pinyin as pinyin, w.meanings as meanings, w.hskLevel as hskLevel,
		       w.frequency as frequency
	`
	
	result, err := r.session.Run(ctx, query, map[string]interface{}{
		"hanzi":       word.Hanzi,
		"traditional": word.Traditional,
		"pinyin":      word.Pinyin,
		"meanings":    word.Meanings,
		"hskLevel":    word.HSKLevel,
		"frequency":   word.Frequency,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create word: %w", err)
	}

	if result.Next(ctx) {
		record := result.Record()
		return r.recordToWord(record), nil
	}

	return nil, fmt.Errorf("failed to create word")
}

func (r *Neo4jWordRepository) Update(ctx context.Context, word *domain.Word) (*domain.Word, error) {
	query := `
		MATCH (w:Word {id: $id})
		SET w.hanzi = $hanzi,
		    w.traditional = $traditional,
		    w.pinyin = $pinyin,
		    w.meanings = $meanings,
		    w.hskLevel = $hskLevel,
		    w.frequency = $frequency
		RETURN w.id as id, w.hanzi as hanzi, w.traditional as traditional, 
		       w.pinyin as pinyin, w.meanings as meanings, w.hskLevel as hskLevel,
		       w.frequency as frequency
	`
	
	result, err := r.session.Run(ctx, query, map[string]interface{}{
		"id":          word.ID,
		"hanzi":       word.Hanzi,
		"traditional": word.Traditional,
		"pinyin":      word.Pinyin,
		"meanings":    word.Meanings,
		"hskLevel":    word.HSKLevel,
		"frequency":   word.Frequency,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update word: %w", err)
	}

	if result.Next(ctx) {
		record := result.Record()
		return r.recordToWord(record), nil
	}

	return nil, fmt.Errorf("word not found")
}

func (r *Neo4jWordRepository) Delete(ctx context.Context, id string) error {
	query := `
		MATCH (w:Word {id: $id})
		DETACH DELETE w
	`
	
	_, err := r.session.Run(ctx, query, map[string]interface{}{
		"id": id,
	})
	return err
}

func (r *Neo4jWordRepository) GetWordGraph(ctx context.Context, hanzi string, depth int, maxNodes int) (*domain.WordGraph, error) {
	// Complex query to get the word graph with relationships
	query := `
		MATCH (center:Word {hanzi: $hanzi})
		CALL apoc.path.subgraphAll(center, {
			maxLevel: $depth,
			limit: $maxNodes
		})
		YIELD nodes, relationships
		WITH center, nodes, relationships
		UNWIND nodes as n
		WITH center, collect(DISTINCT n) as uniqueNodes, relationships
		UNWIND relationships as r
		WITH center, uniqueNodes, 
		     collect({
		         id: id(r),
		         source: startNode(r).id,
		         target: endNode(r).id,
		         type: type(r),
		         label: coalesce(r.phrase, r.meaning, type(r)),
		         weight: coalesce(r.weight, 1.0)
		     }) as edges
		RETURN center, uniqueNodes, edges
	`
	
	result, err := r.session.Run(ctx, query, map[string]interface{}{
		"hanzi":    hanzi,
		"depth":    depth,
		"maxNodes": maxNodes,
	})
	if err != nil {
		// Fallback query without APOC
		return r.getWordGraphFallback(ctx, hanzi, depth, maxNodes)
	}

	if result.Next(ctx) {
		// TODO: Parse result and build graph
		// record := result.Record()
		// Parse and return the graph
		// Implementation details...
	}

	return &domain.WordGraph{
		CenterWord: &domain.Word{Hanzi: hanzi},
		Nodes:      []*domain.WordNode{},
		Edges:      []*domain.Edge{},
		Depth:      depth,
	}, nil
}

func (r *Neo4jWordRepository) getWordGraphFallback(ctx context.Context, hanzi string, depth int, maxNodes int) (*domain.WordGraph, error) {
	// Simpler query without APOC procedures
	query := `
		MATCH path = (center:Word {hanzi: $hanzi})-[*0..$depth]-(connected:Word)
		WITH center, connected, path
		LIMIT $maxNodes
		WITH center, collect(DISTINCT connected) as nodes
		MATCH (n1:Word)-[r]-(n2:Word)
		WHERE n1 IN nodes AND n2 IN nodes
		RETURN center, nodes, collect(DISTINCT {
		    id: id(r),
		    source: n1.id,
		    target: n2.id,
		    type: type(r),
		    label: coalesce(r.phrase, r.meaning, type(r))
		}) as edges
	`
	
	_, err := r.session.Run(ctx, query, map[string]interface{}{
		"hanzi":    hanzi,
		"depth":    depth,
		"maxNodes": maxNodes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get word graph: %w", err)
	}

	// TODO: Parse results and build graph structure
	// For now, return a basic graph
	centerWord := &domain.Word{
		Hanzi:    hanzi,
		Pinyin:   "",
		Meanings: []string{},
	}
	
	return &domain.WordGraph{
		CenterWord: centerWord,
		Nodes:      []*domain.WordNode{},
		Edges:      []*domain.Edge{},
		Depth:      depth,
	}, nil
}

func (r *Neo4jWordRepository) Search(ctx context.Context, query string, limit int) ([]*domain.Word, error) {
	cypher := `
		MATCH (w:Word)
		WHERE w.hanzi CONTAINS $query 
		   OR w.pinyin CONTAINS $query
		   OR any(meaning IN w.meanings WHERE meaning CONTAINS $query)
		RETURN w.id as id, w.hanzi as hanzi, w.traditional as traditional, 
		       w.pinyin as pinyin, w.meanings as meanings, w.hskLevel as hskLevel,
		       w.frequency as frequency
		ORDER BY w.frequency DESC
		LIMIT $limit
	`
	
	result, err := r.session.Run(ctx, cypher, map[string]interface{}{
		"query": query,
		"limit": limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search words: %w", err)
	}

	var words []*domain.Word
	for result.Next(ctx) {
		record := result.Record()
		words = append(words, r.recordToWord(record))
	}

	return words, nil
}

func (r *Neo4jWordRepository) GetByHSKLevel(ctx context.Context, level int, limit int) ([]*domain.Word, error) {
	query := `
		MATCH (w:Word {hskLevel: $level})
		RETURN w.id as id, w.hanzi as hanzi, w.traditional as traditional, 
		       w.pinyin as pinyin, w.meanings as meanings, w.hskLevel as hskLevel,
		       w.frequency as frequency
		ORDER BY w.frequency DESC
		LIMIT $limit
	`
	
	result, err := r.session.Run(ctx, query, map[string]interface{}{
		"level": level,
		"limit": limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get words by HSK level: %w", err)
	}

	var words []*domain.Word
	for result.Next(ctx) {
		record := result.Record()
		words = append(words, r.recordToWord(record))
	}

	return words, nil
}

func (r *Neo4jWordRepository) AddRelation(ctx context.Context, relation *domain.RelationInput) error {
	query := `
		MATCH (source:Word {id: $sourceId}), (target:Word {id: $targetId})
		CREATE (source)-[r:` + string(relation.Type) + ` {
			phrase: $phrase,
			meaning: $meaning,
			pinyin: $pinyin,
			usageExample: $usageExample
		}]->(target)
	`
	
	_, err := r.session.Run(ctx, query, map[string]interface{}{
		"sourceId":     relation.SourceID,
		"targetId":     relation.TargetID,
		"phrase":       relation.Phrase,
		"meaning":      relation.Meaning,
		"pinyin":       relation.Pinyin,
		"usageExample": relation.UsageExample,
	})
	
	return err
}

func (r *Neo4jWordRepository) RemoveRelation(ctx context.Context, sourceID, targetID string, relationType domain.RelationType) error {
	query := `
		MATCH (source:Word {id: $sourceId})-[r:` + string(relationType) + `]->(target:Word {id: $targetId})
		DELETE r
	`
	
	_, err := r.session.Run(ctx, query, map[string]interface{}{
		"sourceId": sourceID,
		"targetId": targetID,
	})
	
	return err
}

func (r *Neo4jWordRepository) recordToWord(record *neo4j.Record) *domain.Word {
	word := &domain.Word{}
	
	if id, ok := record.Get("id"); ok && id != nil {
		word.ID = id.(string)
	}
	if hanzi, ok := record.Get("hanzi"); ok && hanzi != nil {
		word.Hanzi = hanzi.(string)
	}
	if traditional, ok := record.Get("traditional"); ok && traditional != nil {
		str := traditional.(string)
		word.Traditional = &str
	}
	if pinyin, ok := record.Get("pinyin"); ok && pinyin != nil {
		word.Pinyin = pinyin.(string)
	}
	if meanings, ok := record.Get("meanings"); ok && meanings != nil {
		word.Meanings = meanings.([]string)
	}
	if hskLevel, ok := record.Get("hskLevel"); ok && hskLevel != nil {
		level := int(hskLevel.(int64))
		word.HSKLevel = &level
	}
	if frequency, ok := record.Get("frequency"); ok && frequency != nil {
		word.Frequency = int(frequency.(int64))
	}
	
	return word
}