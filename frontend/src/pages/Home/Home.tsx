import { useAuth } from '../../lib/useAuth'
import styles from './Home.module.css'

export default function Home() {
  const { state, logout } = useAuth()

  if (state.status !== 'authed') return null

  return (
    <div className={styles.page}>
      <h2>Hey Let's Meet</h2>
      <p>Hello, {state.user.email}</p>
      <button onClick={() => void logout()} className={styles.logoutButton}>
        Logout
      </button>
      <hr className={styles.divider} />
      <p>this is a protected area. endpoints and map will be here eventually</p>
    </div>
  )
}
