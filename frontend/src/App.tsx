import { useState } from 'react'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center">
      <div className="max-w-md mx-auto text-center space-y-8">
        <h1 className="text-4xl font-bold text-gray-900">
          Proompt
        </h1>
        <p className="text-lg text-gray-600">
          Prompt Management Tool
        </p>
        <div className="space-y-4">
          <button 
            onClick={() => setCount((count) => count + 1)}
            className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          >
            Count is {count}
          </button>
          <p className="text-sm text-gray-500">
            Tailwind CSS v4 + Vite + React + TypeScript
          </p>
        </div>
      </div>
    </div>
  )
}

export default App
