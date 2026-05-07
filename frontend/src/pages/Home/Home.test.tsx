import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import Home from './Home'

vi.mock('../../components/MapView/MapView', () => ({
  default: () => <div data-testid="map" />,
}))

vi.mock('../../lib/useAuth', () => ({
  useAuth: () => ({
    state: { status: 'authed', user: { id: 1, email: 'test@example.com' } },
    logout: vi.fn(),
  }),
}))

const apiMock = vi.fn()

vi.mock('../../lib/api', () => ({
  api: (...args: unknown[]) => apiMock(...args),
}))

beforeEach(() => {
  apiMock.mockReset()
})

describe('Home', () => {
  it('renders result after suggest', async () => {
    apiMock.mockResolvedValueOnce({
      origins: {
        a: { address: 'Alexanderplatz', point: { lat: 52.52, lng: 13.405 } },
        b: { address: 'Hermannplatz', point: { lat: 52.5, lng: 13.4 } },
      },
      best: {
        point: { lat: 52.51, lng: 13.402 },
        etaASeconds: 600,
        etaBSeconds: 650,
        maxEtaSeconds: 650,
        diffSeconds: 50,
        distanceAMeters: 800,
        distanceBMeters: 900,
      },
      alternatives: [],
    })

    render(<Home />)
    fireEvent.click(screen.getByRole('button', { name: /suggest meeting point/i }))

    expect(await screen.findByText(/Result/i)).toBeInTheDocument()
    expect(screen.getByTestId('map')).toBeInTheDocument()
    expect(screen.getByText(/Coordinates:/i)).toBeInTheDocument()
  })

  it('shows error when API fails', async () => {
    apiMock.mockRejectedValueOnce(new Error('boom'))

    render(<Home />)
    fireEvent.click(screen.getByRole('button', { name: /suggest meeting point/i }))

    expect(await screen.findByText(/boom/i)).toBeInTheDocument()
  })
})
