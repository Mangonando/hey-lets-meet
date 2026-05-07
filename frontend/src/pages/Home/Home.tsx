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
  debug?: { midpoint: LatLng}
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
  const [originA, setOriginA] = useState('Alexanderplatz, Berlin')
  const [originB, setOriginB] = useState('Hermannplatz, Berlin')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [result, setResult] = useState<MeetResponse | null>(null)

  if (state.status !== 'authed') return null

  async function onSubmit(e: { preventDefault(): void }) {
    e.preventDefault()
    setError(null)
    setLoading(true)
    setResult(null)

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
          <h1 className={styles.headerTitle}>Hey Let's Meet</h1>
          <p className={styles.headerEmail}>Hello, {state.user.email}</p>
        </div>

        <button onClick={() => void logout()}>Logout</button>
      </header>

      <hr className={styles.divider} />

      <section>
        <h2>Suggest a fair meeting point (walking)</h2>

        <form onSubmit={onSubmit} className={styles.form}>
          <label>
            Address A
            <input
              value={originA}
              onChange={(e) => setOriginA(e.target.value)}
              placeholder="e.g. Alexanderplatz"
              className={styles.fullWidth}
            />
          </label>

          <label>
            Address B
            <input
              value={originB}
              onChange={(e) => setOriginB(e.target.value)}
              placeholder="e.g. Hermannplatz"
              className={styles.fullWidth}
            />
          </label>

          <button type="submit" disabled={loading}>
            {loading ? 'Calculating…' : 'Suggest meeting point'}
          </button>

          {error && <div className={styles.errorMessage}>{error}</div>}
        </form>
      </section>

      {result && (
        <section className={styles.resultsSection}>
          <h2>Result</h2>

          <div className={styles.resultGrid}>
            <MapView
              a={result.origins.a.point}
              b={result.origins.b.point}
              best={result.best.point}
            />

            <div className={styles.bestPointCard}>
              <h3 className={styles.bestPointTitle}>Best point</h3>

              <p className={styles.cardRow}>
                <strong>Coordinates:</strong> {result.best.point.lat.toFixed(6)}, {result.best.point.lng.toFixed(6)}
              </p>

              <p className={styles.cardRow}>
                <strong>A ETA:</strong> {formatSeconds(result.best.etaASeconds)} ({formatMeters(result.best.distanceAMeters)})
                <br />
                <strong>B ETA:</strong> {formatSeconds(result.best.etaBSeconds)} ({formatMeters(result.best.distanceBMeters)})
              </p>

              <p className={styles.cardRow}>
                <strong>Fairness:</strong> max {formatSeconds(result.best.maxEtaSeconds)}, diff {formatSeconds(result.best.diffSeconds)}
              </p>

              <details className={styles.debugDetails}>
                <summary>Debug</summary>
                <pre className={styles.debugPre}>{JSON.stringify(result.debug ?? {}, null, 2)}</pre>
              </details>
            </div>
          </div>

          <h3 className={styles.alternativesTitle}>Alternatives</h3>
          {result.alternatives.length === 0 ? (
            <p>No alternatives returned.</p>
          ) : (
            <div className={styles.alternativesGrid}>
              {result.alternatives.map((alt, index) => (
                <div key={index} className={styles.alternativeCard}>
                  <div>
                    <strong>
                      {alt.point.lat.toFixed(6)}, {alt.point.lng.toFixed(6)}
                    </strong>
                  </div>
                  <div className={styles.alternativeCardRow}>
                    A: {formatSeconds(alt.etaASeconds)} ({formatMeters(alt.distanceAMeters)}) • B: {formatSeconds(alt.etaBSeconds)} (
                    {formatMeters(alt.distanceBMeters)})
                  </div>
                  <div className={styles.alternativeCardRow}>
                    max {formatSeconds(alt.maxEtaSeconds)} • diff {formatSeconds(alt.diffSeconds)}
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>
      )}
    </div>
  )
}
