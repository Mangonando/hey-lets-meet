import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../../lib/useAuth'
import { useState } from 'react'
import { api, type User } from '../../lib/api'
import styles from './Login.module.css'

export default function Login() {
  const { refresh } = useAuth()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function onSubmit(e: React.SyntheticEvent) {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      await api<User>('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
      })
      await refresh()
      navigate('/app')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className={styles.page}>
      <h1 className={styles.title}>Welcome to Hey Let's Meet</h1>
      <p className={styles.subtitle}>Please login to continue</p>
      <form onSubmit={onSubmit} className={styles.form}>
        <input
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          type="email"
          required
          placeholder="Email"
          autoComplete="email"
          className={styles.input}
        />
        <input
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          type="password"
          required
          placeholder="Password"
          autoComplete="current-password"
          className={styles.input}
        />
        {error && <div className={styles.error}>{error}</div>}
        <button type="submit" disabled={loading}>
          {loading ? 'Logging in...' : 'Login'}
        </button>
      </form>
      <p className={styles.footer}>
        No account? <Link to="/register">Register</Link>
      </p>
    </div>
  )
}
