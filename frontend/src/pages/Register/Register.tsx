import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../../lib/useAuth'
import { useState } from 'react'
import { api, type User } from '../../lib/api'
import styles from './Register.module.css'

export default function Register() {
  const navigate = useNavigate()
  const { refresh } = useAuth()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function onSubmit(e: React.SyntheticEvent) {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      await api<User>('/auth/register', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
      })
      await refresh()
      navigate('/app')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Registration failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className={styles.page}>
      <h2>Register</h2>
      <form onSubmit={onSubmit} className={styles.form}>
        <label>
          Email
          <input
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            type="email"
            required
            autoComplete="email"
            className={styles.input}
          />
        </label>
        <label>
          Password
          <input
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            type="password"
            required
            minLength={8}
            autoComplete="new-password"
            className={styles.input}
          />
        </label>
        {error && <div className={styles.error}>{error}</div>}
        <button type="submit" disabled={loading}>
          {loading ? 'Creating account...' : 'Create account'}
        </button>
      </form>
      <p className={styles.footer}>
        Already have an account? <Link to="/login">Login</Link>
      </p>
    </div>
  )
}
