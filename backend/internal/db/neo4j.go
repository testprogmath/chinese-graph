package db

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jDB struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jDB(uri, username, password string) (*Neo4jDB, error) {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	// Verify connectivity
	ctx := context.Background()
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to verify neo4j connectivity: %w", err)
	}

	return &Neo4jDB{driver: driver}, nil
}

func (db *Neo4jDB) Close(ctx context.Context) error {
	return db.driver.Close(ctx)
}

func (db *Neo4jDB) Session(ctx context.Context, config neo4j.SessionConfig) neo4j.SessionWithContext {
	return db.driver.NewSession(ctx, config)
}

func (db *Neo4jDB) InitSchema(ctx context.Context) error {
	session := db.Session(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	// Create constraints and indexes
	queries := []string{
		// Unique constraint on hanzi
		"CREATE CONSTRAINT hanzi_unique IF NOT EXISTS FOR (w:Word) REQUIRE w.hanzi IS UNIQUE",
		
		// Index for faster lookups
		"CREATE INDEX word_pinyin IF NOT EXISTS FOR (w:Word) ON (w.pinyin)",
		"CREATE INDEX word_hsk IF NOT EXISTS FOR (w:Word) ON (w.hskLevel)",
		"CREATE INDEX word_frequency IF NOT EXISTS FOR (w:Word) ON (w.frequency)",
		
		// Index for learning stats
		"CREATE INDEX stats_strength IF NOT EXISTS FOR (s:LearningStats) ON (s.strength)",
	}

	for _, query := range queries {
		_, err := session.Run(ctx, query, nil)
		if err != nil {
			return fmt.Errorf("failed to execute schema query: %w", err)
		}
	}

	return nil
}