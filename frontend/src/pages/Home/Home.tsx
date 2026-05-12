import { useState } from 'react'
import { useAuth } from '../../lib/useAuth'
import styles from './Home.module.css'
import { api } from '../../lib/api'
import MapView from '../../components/MapView/MapView'


type LatLng = { lat: number, lng: number}

type MeetpointResult = {
  point: LatLng
  etaASeconds: number
  etaBSeconds: number
  maxEtaSeconds: number
  diffSeconds: number
  distanceAMeters: number
  distanceBMeters: number
}

type MeetResponse = {
  origins: {
    a: { address: string, point: LatLng}
    b: { address: string, point: LatLng}
  }
  best: MeetpointResult
  alternatives: MeetpointResult[]
}

function formatSeconds(seconds: number) {
  const m = Math.round(seconds / 60)
  if (m < 60) return `${m} min`
  const h = Math.floor(m / 60)
  const rest = m % 60
  return `${h}h ${rest}m`
}

function formatMeters(meters: number) {
  if (meters >= 1000) return `${(meters / 1000).toFixed(1)} km`
  return `${meters} m`
}

export default function Home() {
  const { state, logout } = useAuth()
  const [originA, setOriginA] = useState('')
  const [originB, setOriginB] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [result, setResult] = useState<MeetResponse | null>(null)
  const [selectedIndex, setSelectedIndex] = useState(0)

  if (state.status !== 'authed') return null

  async function onSubmit(e: { preventDefault(): void }) {
    e.preventDefault()
    setError(null)
    setLoading(true)
    setResult(null)
    setSelectedIndex(0)

    try {
      const meetResponse = await api<MeetResponse>('/api/meetpoints/suggest', {
        method: 'POST',
        body: JSON.stringify({originA, originB})
      })
      setResult(meetResponse)
    } catch (caughtError) {
      const message = caughtError instanceof Error ? caughtError.message : 'Request failed'
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className={styles.page}>
      <header className={styles.header}>
        <div>
          <h1 className={styles.headerTitle}>Hey {state.user.email.split('@')[0].charAt(0).toUpperCase() + state.user.email.split('@')[0].slice(1)}! Let's Meet</h1>
        </div>

        <button onClick={() => void logout()}>Logout</button>
      </header>
      <section>
        <p>Add your address and the one you want to meet and it will suggest you a fair meeting point</p>

        <form onSubmit={onSubmit} className={styles.form}>
          <label>
            <input
              value={originA}
              onChange={(e) => setOriginA(e.target.value)}
              placeholder="Address A"
              className={styles.fullWidth}
            />
          </label>

          <label>
            <input
              value={originB}
              onChange={(e) => setOriginB(e.target.value)}
              placeholder="Address B"
              className={styles.fullWidth}
            />
          </label>

          <button type="submit" disabled={loading}>
            {loading ? 'Calculating…' : 'Suggest meeting point'}
          </button>

          {error && <div className={styles.errorMessage}>{error}</div>}
        </form>
      </section>

      {result && (() => {
        const allPoints = [result.best, ...result.alternatives]
        const tabLabels = ['Best point', 'Alternative 1', 'Alternative 2', 'Alternative 3']
        const selected = allPoints[selectedIndex]
        return (
          <section className={styles.resultsSection}>
            <h2>Result</h2>

            <MapView
              a={result.origins.a.point}
              b={result.origins.b.point}
              best={selected.point}
            />

            <div className={styles.tabs}>
              {allPoints.map((_, index) => (
                <button
                  key={index}
                  onClick={() => setSelectedIndex(index)}
                  className={index === selectedIndex ? styles.tabActive : styles.tab}
                >
                  {tabLabels[index]}
                </button>
              ))}
            </div>

            <div className={styles.pointCard}>
              <p className={styles.cardRow}>
                <strong>Coordinates:</strong> {selected.point.lat.toFixed(6)}, {selected.point.lng.toFixed(6)}
              </p>

              <p className={styles.cardRow}>
                <strong>A ETA:</strong> {formatSeconds(selected.etaASeconds)} ({formatMeters(selected.distanceAMeters)})
                <br />
                <strong>B ETA:</strong> {formatSeconds(selected.etaBSeconds)} ({formatMeters(selected.distanceBMeters)})
              </p>

              <p className={styles.cardRow}>
                <strong>Fairness:</strong> max {formatSeconds(selected.maxEtaSeconds)}, diff {formatSeconds(selected.diffSeconds)}
              </p>

              <div className={styles.inlineButtons}>
                <button
                  className={styles.inlineButton}
                  onClick={() => void navigator.clipboard.writeText(`${selected.point.lat.toFixed(6)}, ${selected.point.lng.toFixed(6)}`)}
                >
                  Copy
                </button>
                <a
                  href={`https://www.google.com/maps?q=${selected.point.lat},${selected.point.lng}`}
                  target="_blank"
                  rel="noreferrer"
                  className={styles.inlineButton}
                >
                  Open in Google Maps
                </a>
                <a
                  href={`https://wa.me/?text=${encodeURIComponent(`Hey let's meet here! https://www.google.com/maps?q=${selected.point.lat},${selected.point.lng}`)}`}
                  target="_blank"
                  rel="noreferrer"
                  className={styles.inlineButton}
                >
                  Share on WhatsApp
                </a>
              </div>
            </div>
          </section>
        )
      })()}
    </div>
  )
}
