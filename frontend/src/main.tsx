import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import 'leaflet/dist/leaflet.css'
import './lib/leaflet.ts'
import './index.css'
import { colors } from './lib/colors'

const root = document.documentElement
Object.entries(colors).forEach(([name, value]) => {
  root.style.setProperty(`--color-${name}`, value)
})

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
