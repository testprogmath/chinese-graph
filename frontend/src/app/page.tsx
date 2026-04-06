'use client'

import { useState } from 'react'
import { WordSearch } from '@/components/WordSearch'
import { WordGraph } from '@/components/WordGraph'
import { WordDetail } from '@/components/WordDetail'

export default function Home() {
  const [selectedWord, setSelectedWord] = useState<string | null>(null)

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="container mx-auto px-4 py-8">
        {/* Header */}
        <header className="text-center mb-8">
          <h1 className="text-5xl font-bold text-gray-800 mb-4">
            Chinese Learning Graph
            <span className="text-2xl block text-indigo-600 font-normal mt-2">
              中文学习图谱
            </span>
          </h1>
          <p className="text-lg text-gray-600 max-w-2xl mx-auto">
            Learn Chinese words through visual connections and relationships. 
            Discover how characters combine to form compounds, phrases, and express meaning.
          </p>
        </header>

        {/* Main Content */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left Panel - Search */}
          <div className="lg:col-span-1 space-y-6">
            <WordSearch onWordSelect={setSelectedWord} />
            {selectedWord && <WordDetail hanzi={selectedWord} />}
          </div>

          {/* Right Panel - Graph Visualization */}
          <div className="lg:col-span-2">
            <WordGraph 
              centerWord={selectedWord} 
              onWordClick={setSelectedWord}
            />
          </div>
        </div>

        {/* Footer */}
        <footer className="text-center mt-12 pt-8 border-t border-gray-200">
          <p className="text-sm text-gray-500">
            Explore word relationships • Practice with context • Learn through connections
          </p>
        </footer>
      </div>
    </div>
  )
}