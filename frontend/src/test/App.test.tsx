import { render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import App from '../App'

describe('App', () => {
  beforeEach(() => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ message: 'pong' }),
    })
  })

  it('renders the login page when unauthenticated', async () => {
    render(<App />)
    expect(await screen.findByRole('heading', { name: 'Login' })).toBeInTheDocument()
  })
})
