'use client'

import { useQuery } from '@apollo/client'
import { GET_WORD_GRAPH } from '@/graphql/queries'
import { useState, useEffect } from 'react'

interface WordGraphProps {
  centerWord: string | null
  onWordClick: (hanzi: string) => void
}

interface WordNode {
  id: string
  word: {
    id: string
    hanzi: string
    pinyin: string
    meanings: string[]
    hskLevel?: number
  }
  distance: number
  importance: number
}

interface Edge {
  id: string
  source: string
  target: string
  type: string
  label: string
  weight: number
}

interface WordGraph {
  centerWord: {
    hanzi: string
    pinyin: string
    meanings: string[]
  }
  nodes: WordNode[]
  edges: Edge[]
  depth: number
}

export function WordGraph({ centerWord, onWordClick }: WordGraphProps) {
  const [selectedNode, setSelectedNode] = useState<string | null>(null)

  const { data, loading, error } = useQuery(GET_WORD_GRAPH, {
    variables: { hanzi: centerWord, depth: 2, maxNodes: 20 },
    skip: !centerWord
  })

  // Default empty state
  if (!centerWord) {
    return (
      <div className="bg-white rounded-lg shadow-md p-8 h-96 flex items-center justify-center">
        <div className="text-center text-gray-500">
          <div className="text-6xl mb-4">🌐</div>
          <h3 className="text-xl font-medium mb-2">Word Relationship Graph</h3>
          <p>Select a word from the search panel to see its connections and relationships</p>
          <div className="mt-4 text-sm text-gray-400">
            • See compounds and phrases
            <br />
            • Discover similar and opposite words  
            <br />
            • Explore learning paths
          </div>
        </div>
      </div>
    )
  }

  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow-md p-8 h-96 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Building word graph...</p>
        </div>
      </div>
    )
  }

  if (error || !data?.wordGraph) {
    return (
      <div className="bg-white rounded-lg shadow-md p-8 h-96 flex items-center justify-center">
        <div className="text-center text-red-500">
          <div className="text-4xl mb-4">⚠️</div>
          <p>Unable to load word graph</p>
          <p className="text-sm text-gray-500 mt-2">Try selecting a different word</p>
        </div>
      </div>
    )
  }

  const graph: WordGraph = data.wordGraph

  // Create a simple circular layout for the nodes
  const getNodePosition = (index: number, total: number, distance: number) => {
    const radius = distance === 0 ? 0 : 120 + distance * 80
    const angle = (index * 2 * Math.PI) / Math.max(total, 1)
    const x = 50 + (radius * Math.cos(angle)) / 4 // Scale down for container
    const y = 50 + (radius * Math.sin(angle)) / 4
    return { x, y }
  }

  // Group nodes by distance
  const nodesByDistance = graph.nodes.reduce((acc, node) => {
    const distance = node.distance
    if (!acc[distance]) acc[distance] = []
    acc[distance].push(node)
    return acc
  }, {} as Record<number, WordNode[]>)

  // Get relation type styling
  const getRelationStyle = (type: string) => {
    switch (type.toLowerCase()) {
      case 'compound':
        return { color: 'bg-green-500', label: 'Compound' }
      case 'phrase':
        return { color: 'bg-blue-500', label: 'Phrase' }
      case 'similar':
        return { color: 'bg-purple-500', label: 'Similar' }
      case 'opposite':
        return { color: 'bg-red-500', label: 'Opposite' }
      case 'modifier':
        return { color: 'bg-orange-500', label: 'Modifier' }
      default:
        return { color: 'bg-gray-500', label: type }
    }
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold text-gray-800">
          Word Graph: {centerWord}
        </h2>
        <div className="text-sm text-gray-600">
          {graph.nodes.length} connected words
        </div>
      </div>

      {/* Graph Visualization */}
      <div className="relative bg-gray-50 rounded-lg h-96 overflow-hidden">
        <svg className="absolute inset-0 w-full h-full" style={{ zIndex: 1 }}>
          {/* Draw edges */}
          {graph.edges.map((edge) => {
            const sourceNode = graph.nodes.find(n => n.id === edge.source)
            const targetNode = graph.nodes.find(n => n.id === edge.target)
            
            if (!sourceNode || !targetNode) return null

            const sourceIndex = graph.nodes.findIndex(n => n.id === edge.source)
            const targetIndex = graph.nodes.findIndex(n => n.id === edge.target)
            
            const sourcePos = sourceNode.distance === 0 
              ? { x: 50, y: 50 }
              : getNodePosition(sourceIndex, nodesByDistance[sourceNode.distance]?.length || 1, sourceNode.distance)
            
            const targetPos = targetNode.distance === 0
              ? { x: 50, y: 50 }
              : getNodePosition(targetIndex, nodesByDistance[targetNode.distance]?.length || 1, targetNode.distance)

            const relationStyle = getRelationStyle(edge.type)

            return (
              <g key={edge.id}>
                <line
                  x1={`${sourcePos.x}%`}
                  y1={`${sourcePos.y}%`}
                  x2={`${targetPos.x}%`}
                  y2={`${targetPos.y}%`}
                  stroke="#cbd5e0"
                  strokeWidth="2"
                  strokeOpacity="0.6"
                />
              </g>
            )
          })}
        </svg>

        {/* Draw nodes */}
        {Object.entries(nodesByDistance).map(([distance, nodes]) =>
          nodes.map((node, index) => {
            const position = parseInt(distance) === 0 
              ? { x: 50, y: 50 } 
              : getNodePosition(index, nodes.length, parseInt(distance))
            
            const isCenter = parseInt(distance) === 0
            const isSelected = selectedNode === node.id

            return (
              <div
                key={node.id}
                className={`absolute transform -translate-x-1/2 -translate-y-1/2 cursor-pointer transition-all duration-200 ${
                  isCenter 
                    ? 'z-20 scale-110' 
                    : isSelected 
                      ? 'z-10 scale-105' 
                      : 'z-5 hover:scale-105'
                }`}
                style={{
                  left: `${position.x}%`,
                  top: `${position.y}%`,
                }}
                onClick={() => {
                  setSelectedNode(node.id)
                  onWordClick(node.word.hanzi)
                }}
                onMouseEnter={() => setSelectedNode(node.id)}
                onMouseLeave={() => setSelectedNode(null)}
              >
                <div
                  className={`
                    px-3 py-2 rounded-lg shadow-md text-center min-w-16 transition-all
                    ${isCenter 
                      ? 'bg-indigo-600 text-white shadow-lg' 
                      : 'bg-white text-gray-800 border-2 border-gray-200 hover:border-indigo-300'
                    }
                    ${isSelected ? 'ring-2 ring-indigo-400' : ''}
                  `}
                >
                  <div className={`font-bold ${isCenter ? 'text-lg' : 'text-sm'}`}>
                    {node.word.hanzi}
                  </div>
                  <div className={`text-xs ${isCenter ? 'text-indigo-100' : 'text-gray-500'} truncate`}>
                    {node.word.pinyin}
                  </div>
                </div>
              </div>
            )
          })
        )}

        {/* Selected node tooltip */}
        {selectedNode && (
          <div className="absolute top-4 right-4 bg-white rounded-lg shadow-lg p-3 border max-w-xs z-30">
            {(() => {
              const node = graph.nodes.find(n => n.id === selectedNode)
              if (!node) return null
              
              return (
                <div>
                  <div className="font-semibold text-gray-800">
                    {node.word.hanzi} ({node.word.pinyin})
                  </div>
                  <div className="text-sm text-gray-600 mt-1">
                    {node.word.meanings.slice(0, 2).join(', ')}
                  </div>
                  {node.word.hskLevel && (
                    <div className="text-xs text-orange-600 mt-1">
                      HSK Level {node.word.hskLevel}
                    </div>
                  )}
                  <div className="text-xs text-gray-500 mt-2">
                    Distance: {node.distance} | Importance: {(node.importance * 100).toFixed(0)}%
                  </div>
                </div>
              )
            })()}
          </div>
        )}
      </div>

      {/* Legend */}
      <div className="mt-4 flex flex-wrap gap-2 text-xs">
        <div className="flex items-center gap-1">
          <div className="w-3 h-3 bg-indigo-600 rounded"></div>
          <span>Center Word</span>
        </div>
        {Array.from(new Set(graph.edges.map(e => e.type))).map(type => {
          const style = getRelationStyle(type)
          return (
            <div key={type} className="flex items-center gap-1">
              <div className={`w-3 h-3 ${style.color} rounded`}></div>
              <span>{style.label}</span>
            </div>
          )
        })}
      </div>
    </div>
  )
}