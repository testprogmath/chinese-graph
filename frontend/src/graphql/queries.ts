import { gql } from '@apollo/client'

// Search for words by query string
export const SEARCH_WORDS = gql`
  query SearchWords($query: String!, $limit: Int = 20) {
    searchWords(query: $query, limit: $limit) {
      id
      hanzi
      pinyin
      meanings
      hskLevel
      frequency
    }
  }
`

// Get words by HSK level
export const GET_WORDS_BY_HSK = gql`
  query WordsByHSK($level: Int!, $limit: Int = 50) {
    wordsByHSK(level: $level, limit: $limit) {
      id
      hanzi
      pinyin
      meanings
      hskLevel
      frequency
    }
  }
`

// Get detailed word information
export const GET_WORD = gql`
  query GetWord($hanzi: String!) {
    word(hanzi: $hanzi) {
      id
      hanzi
      pinyin
      meanings
      hskLevel
      frequency
      compounds {
        word {
          id
          hanzi
          pinyin
          meanings
        }
        relationshipType
        phrase
        meaning
      }
      phrases {
        word {
          id
          hanzi
          pinyin
          meanings
        }
        relationshipType
        phrase
        meaning
      }
      similarWords {
        word {
          id
          hanzi
          pinyin
          meanings
        }
        relationshipType
        phrase
        meaning
      }
      oppositeWords {
        word {
          id
          hanzi
          pinyin
          meanings
        }
        relationshipType
        phrase
        meaning
      }
    }
  }
`

// Get word graph with relationships
export const GET_WORD_GRAPH = gql`
  query WordGraph($hanzi: String!, $depth: Int = 2, $maxNodes: Int = 30) {
    wordGraph(hanzi: $hanzi, depth: $depth, maxNodes: $maxNodes) {
      centerWord {
        id
        hanzi
        pinyin
        meanings
        hskLevel
      }
      nodes {
        id
        word {
          id
          hanzi
          pinyin
          meanings
          hskLevel
        }
        distance
        importance
      }
      edges {
        id
        source
        target
        type
        label
        weight
        bidirectional
      }
      depth
    }
  }
`

// Get daily practice words
export const GET_DAILY_WORDS = gql`
  query DailyWords($count: Int = 10, $hskLevel: Int) {
    dailyWords(count: $count, hskLevel: $hskLevel) {
      id
      hanzi
      pinyin
      meanings
      hskLevel
      frequency
    }
  }
`

// Get learning statistics for a word
export const GET_LEARNING_STATS = gql`
  query LearningStats($wordId: ID!) {
    learningStats(wordId: $wordId) {
      wordId
      timesReviewed
      correctCount
      lastReviewed
      strength
    }
  }
`

// Mark word as reviewed (mutation)
export const MARK_REVIEWED = gql`
  mutation MarkReviewed($wordId: ID!, $correct: Boolean!) {
    markReviewed(wordId: $wordId, correct: $correct) {
      wordId
      timesReviewed
      correctCount
      lastReviewed
      strength
    }
  }
`

// Add new word (mutation)
export const ADD_WORD = gql`
  mutation AddWord($input: WordInput!) {
    addWord(input: $input) {
      id
      hanzi
      pinyin
      meanings
      hskLevel
      frequency
    }
  }
`

// Add relationship between words (mutation)
export const ADD_RELATION = gql`
  mutation AddRelation($input: RelationInput!) {
    addRelation(input: $input)
  }
`