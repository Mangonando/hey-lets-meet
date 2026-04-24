import { Navigate } from 'react-router-dom'
import { useAuth } from '../../lib/useAuth'
import styles from './ProtectedRoute.module.css'

export default function ProtectedRoute({
  children,
}: {
  children: React.ReactNode
}) {
  const { state } = useAuth()

  if (state.status === 'loading') {
    return <div className={styles.loading}>Loading...</div>
  }
  if (state.status === 'anon') {
    return <Navigate to="/login" replace />
  }
  return children
}
