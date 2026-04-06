'use client'

import { useQuery } from '@apollo/client'
import { GET_WORD } from '@/graphql/queries'

interface WordDetailProps {
  hanzi: string
}

interface WordRelation {
  word: {
    id: string
    hanzi: string
    pinyin: string
    meanings: string[]
  }
  relationshipType: string
  phrase?: string
  meaning: string
}

interface Word {
  id: string
  hanzi: string
  pinyin: string
  meanings: string[]
  hskLevel?: number
  frequency: number
  compounds: WordRelation[]
  phrases: WordRelation[]
  similarWords: WordRelation[]
  oppositeWords: WordRelation[]
}

export function WordDetail({ hanzi }: WordDetailProps) {
  const { data, loading, error } = useQuery(GET_WORD, {
    variables: { hanzi },
    skip: !hanzi
  })

  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-300 rounded w-3/4 mb-4"></div>
          <div className="space-y-2">
            <div className="h-4 bg-gray-300 rounded"></div>
            <div className="h-4 bg-gray-300 rounded w-2/3"></div>
          </div>
        </div>
      </div>
    )
  }

  if (error || !data?.word) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="text-center text-gray-500">
          <p>Word not found or error loading details</p>
          <p className="text-sm mt-2">Try searching for a different word</p>
        </div>
      </div>
    )
  }

  const word: Word = data.word

  const relationshipSections = [
    { title: 'Compounds', data: word.compounds, icon: '🔗' },
    { title: 'Phrases', data: word.phrases, icon: '💬' },
    { title: 'Similar Words', data: word.similarWords, icon: '≈' },
    { title: 'Opposite Words', data: word.oppositeWords, icon: '⚡' },
  ]

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-xl font-semibold text-gray-800 mb-4">Word Details</h2>
      
      {/* Main Word Info */}
      <div className="mb-6">
        <div className="flex items-center gap-4 mb-3">
          <span className="text-4xl font-bold text-gray-800">
            {word.hanzi}
          </span>
          <div>
            <div className="text-lg text-gray-600">{word.pinyin}</div>
            <div className="flex gap-2 mt-1">
              {word.hskLevel && (
                <span className="px-2 py-1 text-xs bg-orange-100 text-orange-800 rounded-full">
                  HSK {word.hskLevel}
                </span>
              )}
              <span className="px-2 py-1 text-xs bg-blue-100 text-blue-800 rounded-full">
                Frequency: {word.frequency}
              </span>
            </div>
          </div>
        </div>
        
        <div className="space-y-1">
          {word.meanings.map((meaning, index) => (
            <div key={index} className="text-gray-700">
              <span className="text-indigo-600 font-medium">{index + 1}.</span> {meaning}
            </div>
          ))}
        </div>
      </div>

      {/* Relationships */}
      {relationshipSections.map((section) => 
        section.data.length > 0 && (
          <div key={section.title} className="mb-4">
            <h3 className="text-lg font-medium text-gray-800 mb-2 flex items-center gap-2">
              <span>{section.icon}</span>
              {section.title}
            </h3>
            <div className="space-y-2">
              {section.data.map((relation, index) => (
                <div
                  key={index}
                  className="p-3 bg-gray-50 rounded-lg border border-gray-200"
                >
                  <div className="flex items-center gap-3">
                    <span className="text-xl font-semibold text-gray-800">
                      {relation.phrase || relation.word.hanzi}
                    </span>
                    <span className="text-sm text-gray-600">
                      {relation.word.pinyin}
                    </span>
                  </div>
                  <div className="text-sm text-gray-700 mt-1">
                    {relation.meaning}
                  </div>
                  {relation.phrase && (
                    <div className="text-xs text-gray-500 mt-1">
                      Contains: {relation.word.hanzi} ({relation.word.meanings.join(', ')})
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        )
      )}

      {/* Learning Actions */}
      <div className="mt-6 pt-4 border-t border-gray-200">
        <div className="flex flex-wrap gap-2">
          <button className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors text-sm">
            ✓ I Know This
          </button>
          <button className="px-4 py-2 bg-orange-600 text-white rounded-lg hover:bg-orange-700 transition-colors text-sm">
            📚 Study Later
          </button>
          <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm">
            🔄 Practice
          </button>
        </div>
      </div>
    </div>
  )
}