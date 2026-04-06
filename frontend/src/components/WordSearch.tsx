'use client'

import { useState, useEffect } from 'react'
import { useQuery } from '@apollo/client'
import { SEARCH_WORDS, GET_WORDS_BY_HSK } from '@/graphql/queries'

interface WordSearchProps {
  onWordSelect: (hanzi: string) => void
}

interface Word {
  id: string
  hanzi: string
  pinyin: string
  meanings: string[]
  hskLevel?: number
}

export function WordSearch({ onWordSelect }: WordSearchProps) {
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedHSK, setSelectedHSK] = useState<number | null>(null)

  // Search query
  const { data: searchData, loading: searchLoading } = useQuery(SEARCH_WORDS, {
    variables: { query: searchQuery, limit: 10 },
    skip: !searchQuery || searchQuery.length < 1
  })

  // HSK level query
  const { data: hskData, loading: hskLoading } = useQuery(GET_WORDS_BY_HSK, {
    variables: { level: selectedHSK, limit: 20 },
    skip: !selectedHSK
  })

  const words = searchQuery ? searchData?.searchWords || [] : hskData?.wordsByHSK || []

  // Popular starter words for initial display
  const starterWords = [
    { hanzi: '我', pinyin: 'wǒ', meanings: ['I', 'me'] },
    { hanzi: '你', pinyin: 'nǐ', meanings: ['you'] },
    { hanzi: '好', pinyin: 'hǎo', meanings: ['good'] },
    { hanzi: '中国', pinyin: 'zhōngguó', meanings: ['China'] },
    { hanzi: '学生', pinyin: 'xuéshēng', meanings: ['student'] },
    { hanzi: '大', pinyin: 'dà', meanings: ['big'] },
  ]

  const displayWords = words.length > 0 ? words : starterWords

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-xl font-semibold text-gray-800 mb-4">Word Search</h2>
      
      {/* Search Input */}
      <div className="mb-4">
        <input
          type="text"
          placeholder="Search words (汉字, pinyin, or English)..."
          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>

      {/* HSK Level Filter */}
      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Filter by HSK Level:
        </label>
        <div className="flex flex-wrap gap-2">
          <button
            onClick={() => setSelectedHSK(null)}
            className={`px-3 py-1 text-sm rounded-full ${
              selectedHSK === null
                ? 'bg-indigo-600 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            All
          </button>
          {[1, 2, 3, 4, 5, 6].map((level) => (
            <button
              key={level}
              onClick={() => setSelectedHSK(level)}
              className={`px-3 py-1 text-sm rounded-full ${
                selectedHSK === level
                  ? 'bg-indigo-600 text-white'
                  : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
              }`}
            >
              HSK {level}
            </button>
          ))}
        </div>
      </div>

      {/* Loading State */}
      {(searchLoading || hskLoading) && (
        <div className="text-center py-4">
          <div className="animate-pulse text-gray-500">Searching...</div>
        </div>
      )}

      {/* Word List */}
      <div className="space-y-2 max-h-96 overflow-y-auto">
        {displayWords.map((word: Word, index: number) => (
          <div
            key={word.id || index}
            onClick={() => onWordSelect(word.hanzi)}
            className="p-3 border border-gray-200 rounded-lg hover:bg-indigo-50 hover:border-indigo-300 cursor-pointer transition-colors"
          >
            <div className="flex items-center justify-between">
              <div>
                <div className="flex items-center gap-3">
                  <span className="text-2xl font-bold text-gray-800">
                    {word.hanzi}
                  </span>
                  <span className="text-sm text-gray-600">
                    {word.pinyin}
                  </span>
                  {word.hskLevel && (
                    <span className="px-2 py-1 text-xs bg-orange-100 text-orange-800 rounded-full">
                      HSK {word.hskLevel}
                    </span>
                  )}
                </div>
                <div className="text-sm text-gray-700 mt-1">
                  {word.meanings.join(', ')}
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Empty State */}
      {!searchLoading && !hskLoading && words.length === 0 && searchQuery && (
        <div className="text-center py-4 text-gray-500">
          No words found for &ldquo;{searchQuery}&rdquo;
        </div>
      )}
    </div>
  )
}