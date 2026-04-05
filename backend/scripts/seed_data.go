package main

import (
	"context"
	"log"

	"github.com/chinese-graph/backend/internal/db"
	"github.com/chinese-graph/backend/internal/domain"
	"github.com/chinese-graph/backend/internal/repository"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type WordSeed struct {
	Hanzi     string
	Pinyin    string
	Meanings  []string
	HSKLevel  int
	Frequency int
}

type RelationSeed struct {
	From     string
	To       string
	Type     domain.RelationType
	Phrase   string
	Meaning  string
}

func main() {
	// Connect to Neo4j
	database, err := db.NewNeo4jDB("bolt://localhost:7687", "neo4j", "chinesegraph123")
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	defer database.Close(context.Background())

	ctx := context.Background()
	session := database.Session(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	repo := repository.NewNeo4jWordRepository(session)

	// Basic HSK1 words
	words := []WordSeed{
		// Pronouns and basic words
		{Hanzi: "我", Pinyin: "wǒ", Meanings: []string{"I", "me"}, HSKLevel: 1, Frequency: 100},
		{Hanzi: "你", Pinyin: "nǐ", Meanings: []string{"you"}, HSKLevel: 1, Frequency: 99},
		{Hanzi: "他", Pinyin: "tā", Meanings: []string{"he", "him"}, HSKLevel: 1, Frequency: 98},
		{Hanzi: "她", Pinyin: "tā", Meanings: []string{"she", "her"}, HSKLevel: 1, Frequency: 97},
		{Hanzi: "们", Pinyin: "men", Meanings: []string{"plural marker"}, HSKLevel: 1, Frequency: 96},
		
		// Basic adjectives
		{Hanzi: "好", Pinyin: "hǎo", Meanings: []string{"good", "well"}, HSKLevel: 1, Frequency: 95},
		{Hanzi: "大", Pinyin: "dà", Meanings: []string{"big", "large"}, HSKLevel: 1, Frequency: 94},
		{Hanzi: "小", Pinyin: "xiǎo", Meanings: []string{"small", "little"}, HSKLevel: 1, Frequency: 93},
		{Hanzi: "多", Pinyin: "duō", Meanings: []string{"many", "much"}, HSKLevel: 1, Frequency: 92},
		{Hanzi: "少", Pinyin: "shǎo", Meanings: []string{"few", "little"}, HSKLevel: 1, Frequency: 91},
		
		// Basic verbs
		{Hanzi: "是", Pinyin: "shì", Meanings: []string{"to be", "is"}, HSKLevel: 1, Frequency: 100},
		{Hanzi: "有", Pinyin: "yǒu", Meanings: []string{"to have", "there is"}, HSKLevel: 1, Frequency: 99},
		{Hanzi: "看", Pinyin: "kàn", Meanings: []string{"to look", "to see"}, HSKLevel: 1, Frequency: 90},
		{Hanzi: "听", Pinyin: "tīng", Meanings: []string{"to listen", "to hear"}, HSKLevel: 1, Frequency: 89},
		{Hanzi: "说", Pinyin: "shuō", Meanings: []string{"to speak", "to say"}, HSKLevel: 1, Frequency: 88},
		{Hanzi: "读", Pinyin: "dú", Meanings: []string{"to read"}, HSKLevel: 1, Frequency: 87},
		{Hanzi: "写", Pinyin: "xiě", Meanings: []string{"to write"}, HSKLevel: 1, Frequency: 86},
		
		// Common nouns
		{Hanzi: "人", Pinyin: "rén", Meanings: []string{"person", "people"}, HSKLevel: 1, Frequency: 98},
		{Hanzi: "天", Pinyin: "tiān", Meanings: []string{"day", "sky"}, HSKLevel: 1, Frequency: 85},
		{Hanzi: "年", Pinyin: "nián", Meanings: []string{"year"}, HSKLevel: 1, Frequency: 84},
		{Hanzi: "月", Pinyin: "yuè", Meanings: []string{"month", "moon"}, HSKLevel: 1, Frequency: 83},
		{Hanzi: "日", Pinyin: "rì", Meanings: []string{"day", "sun"}, HSKLevel: 1, Frequency: 82},
		{Hanzi: "时", Pinyin: "shí", Meanings: []string{"time", "hour"}, HSKLevel: 1, Frequency: 81},
		{Hanzi: "分", Pinyin: "fēn", Meanings: []string{"minute", "to divide"}, HSKLevel: 1, Frequency: 80},
		
		// Numbers
		{Hanzi: "一", Pinyin: "yī", Meanings: []string{"one"}, HSKLevel: 1, Frequency: 100},
		{Hanzi: "二", Pinyin: "èr", Meanings: []string{"two"}, HSKLevel: 1, Frequency: 99},
		{Hanzi: "三", Pinyin: "sān", Meanings: []string{"three"}, HSKLevel: 1, Frequency: 98},
		{Hanzi: "四", Pinyin: "sì", Meanings: []string{"four"}, HSKLevel: 1, Frequency: 97},
		{Hanzi: "五", Pinyin: "wǔ", Meanings: []string{"five"}, HSKLevel: 1, Frequency: 96},
		
		// Common particles
		{Hanzi: "的", Pinyin: "de", Meanings: []string{"possessive particle"}, HSKLevel: 1, Frequency: 100},
		{Hanzi: "了", Pinyin: "le", Meanings: []string{"completed action marker"}, HSKLevel: 1, Frequency: 99},
		{Hanzi: "吗", Pinyin: "ma", Meanings: []string{"question particle"}, HSKLevel: 1, Frequency: 98},
		{Hanzi: "不", Pinyin: "bù", Meanings: []string{"not", "no"}, HSKLevel: 1, Frequency: 100},
		{Hanzi: "很", Pinyin: "hěn", Meanings: []string{"very"}, HSKLevel: 1, Frequency: 97},
		
		// Additional useful words
		{Hanzi: "中", Pinyin: "zhōng", Meanings: []string{"middle", "center", "China"}, HSKLevel: 1, Frequency: 95},
		{Hanzi: "国", Pinyin: "guó", Meanings: []string{"country", "nation"}, HSKLevel: 1, Frequency: 94},
		{Hanzi: "家", Pinyin: "jiā", Meanings: []string{"home", "family"}, HSKLevel: 1, Frequency: 93},
		{Hanzi: "学", Pinyin: "xué", Meanings: []string{"to study", "to learn"}, HSKLevel: 1, Frequency: 92},
		{Hanzi: "生", Pinyin: "shēng", Meanings: []string{"life", "to give birth"}, HSKLevel: 1, Frequency: 91},
		{Hanzi: "上", Pinyin: "shàng", Meanings: []string{"up", "above", "to go up"}, HSKLevel: 1, Frequency: 90},
		{Hanzi: "下", Pinyin: "xià", Meanings: []string{"down", "below", "to go down"}, HSKLevel: 1, Frequency: 89},
	}

	// Create words
	wordMap := make(map[string]*domain.Word)
	for _, seed := range words {
		word := &domain.Word{
			Hanzi:     seed.Hanzi,
			Pinyin:    seed.Pinyin,
			Meanings:  seed.Meanings,
			HSKLevel:  &seed.HSKLevel,
			Frequency: seed.Frequency,
		}
		
		created, err := repo.Create(ctx, word)
		if err != nil {
			log.Printf("Failed to create word %s: %v", seed.Hanzi, err)
			continue
		}
		wordMap[seed.Hanzi] = created
		log.Printf("Created word: %s", seed.Hanzi)
	}

	// Define relationships
	relations := []RelationSeed{
		// Compounds with 好
		{From: "好", To: "看", Type: domain.RelationTypeCompound, Phrase: "好看", Meaning: "good-looking"},
		{From: "好", To: "听", Type: domain.RelationTypeCompound, Phrase: "好听", Meaning: "pleasant to hear"},
		{From: "好", To: "人", Type: domain.RelationTypeCompound, Phrase: "好人", Meaning: "good person"},
		
		// Phrases with 好
		{From: "你", To: "好", Type: domain.RelationTypePhrase, Phrase: "你好", Meaning: "hello"},
		{From: "很", To: "好", Type: domain.RelationTypeModifier, Phrase: "很好", Meaning: "very good"},
		{From: "不", To: "好", Type: domain.RelationTypeModifier, Phrase: "不好", Meaning: "not good"},
		
		// Plural forms
		{From: "我", To: "们", Type: domain.RelationTypeCompound, Phrase: "我们", Meaning: "we, us"},
		{From: "你", To: "们", Type: domain.RelationTypeCompound, Phrase: "你们", Meaning: "you (plural)"},
		{From: "他", To: "们", Type: domain.RelationTypeCompound, Phrase: "他们", Meaning: "they (male)"},
		{From: "她", To: "们", Type: domain.RelationTypeCompound, Phrase: "她们", Meaning: "they (female)"},
		
		// Country names
		{From: "中", To: "国", Type: domain.RelationTypeCompound, Phrase: "中国", Meaning: "China"},
		{From: "中", To: "国", Type: domain.RelationTypeCompound, Phrase: "中国人", Meaning: "Chinese person"},
		
		// Student
		{From: "学", To: "生", Type: domain.RelationTypeCompound, Phrase: "学生", Meaning: "student"},
		{From: "大", To: "学", Type: domain.RelationTypeCompound, Phrase: "大学", Meaning: "university"},
		{From: "大", To: "学", Type: domain.RelationTypeCompound, Phrase: "大学生", Meaning: "university student"},
		{From: "小", To: "学", Type: domain.RelationTypeCompound, Phrase: "小学", Meaning: "elementary school"},
		
		// Time compounds
		{From: "今", To: "天", Type: domain.RelationTypeCompound, Phrase: "今天", Meaning: "today"},
		{From: "昨", To: "天", Type: domain.RelationTypeCompound, Phrase: "昨天", Meaning: "yesterday"},
		{From: "明", To: "天", Type: domain.RelationTypeCompound, Phrase: "明天", Meaning: "tomorrow"},
		
		// Opposites
		{From: "大", To: "小", Type: domain.RelationTypeOpposite, Phrase: "", Meaning: "big vs small"},
		{From: "多", To: "少", Type: domain.RelationTypeOpposite, Phrase: "", Meaning: "many vs few"},
		{From: "上", To: "下", Type: domain.RelationTypeOpposite, Phrase: "", Meaning: "up vs down"},
		{From: "好", To: "坏", Type: domain.RelationTypeOpposite, Phrase: "", Meaning: "good vs bad"},
		
		// Similar meanings
		{From: "看", To: "见", Type: domain.RelationTypeSimilar, Phrase: "", Meaning: "to see"},
		{From: "听", To: "闻", Type: domain.RelationTypeSimilar, Phrase: "", Meaning: "to hear/perceive"},
	}

	// Create relationships
	for _, rel := range relations {
		fromWord, fromOk := wordMap[rel.From]
		toWord, toOk := wordMap[rel.To]
		
		if !fromOk || !toOk {
			log.Printf("Skipping relation %s -> %s: word not found", rel.From, rel.To)
			continue
		}

		relationInput := &domain.RelationInput{
			SourceID: fromWord.ID,
			TargetID: toWord.ID,
			Type:     rel.Type,
			Meaning:  rel.Meaning,
		}
		
		if rel.Phrase != "" {
			relationInput.Phrase = &rel.Phrase
		}

		err := repo.AddRelation(ctx, relationInput)
		if err != nil {
			log.Printf("Failed to create relation %s -> %s: %v", rel.From, rel.To, err)
			continue
		}
		log.Printf("Created relation: %s -> %s (%s)", rel.From, rel.To, rel.Type)
	}

	log.Println("Seed data loaded successfully!")
}