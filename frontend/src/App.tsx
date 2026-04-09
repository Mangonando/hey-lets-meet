import { useEffect, useState } from 'react'
import './App.css'

type PingResponse = { message: string }

export default function App() {
  const [ping, setPing] = useState<string>('loading...')

  useEffect(() => {
    fetch('/api/ping')
      .then((r) => {
        if (!r.ok) throw new Error(`HTTP ${r.status}`)
        return r.json() as Promise<PingResponse>
      })
      .then((data) => setPing(data.message))
      .catch(() => setPing('error'))
  }, [])

  return (
    <div style={{ padding: 24 }}>
      <h1>hey lets meet</h1>
      <p>if i say ping backend says {ping}</p>
    </div>
  )
}
