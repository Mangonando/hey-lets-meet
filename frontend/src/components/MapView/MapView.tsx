import { useEffect, useMemo } from "react";
import { MapContainer, Marker, Popup, TileLayer, useMap } from "react-leaflet";
import L from 'leaflet'
import styles from './MapView.module.css'
import { colors } from '../../lib/colors'

export type LatLng = { lat: number; lng: number}

function ScrollWheelControl() {
    const map = useMap()

    useEffect(() => {
        map.scrollWheelZoom.disable()
        const container = map.getContainer()
        const enable = () => map.scrollWheelZoom.enable()
        const disable = () => map.scrollWheelZoom.disable()
        container.addEventListener('mouseenter', enable)
        container.addEventListener('mouseleave', disable)
        return () => {
            container.removeEventListener('mouseenter', enable)
            container.removeEventListener('mouseleave', disable)
        }
    }, [map])

    return null
}

function FitBounds ({ points}: {points: LatLng[]}) {
    const map = useMap()

    useEffect(() => {
        if (points.length === 0) return
        const bounds = L.latLngBounds(points.map((point) => [point.lat, point.lng] as [number, number]))
        map.fitBounds(bounds, { padding: [30, 30]})
    }, [map, points])
    return null
}

function pinIcon(fill: string, dot: string) {
    return L.divIcon({
        className: '',
        html: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 36" width="24" height="36">
          <path d="M12 0C5.373 0 0 5.373 0 12c0 9 12 24 12 24S24 21 24 12C24 5.373 18.627 0 12 0z" fill="${fill}"/>
          <circle cx="12" cy="12" r="5" fill="${dot}"/>
        </svg>`,
        iconSize: [24, 36],
        iconAnchor: [12, 36],
        popupAnchor: [0, -36],
    })
}

const oceanicIcon = pinIcon(colors.oceanic, colors.wheat)
const amberIcon   = pinIcon(colors.amber,   colors.wheat)

export default function MapView({
    a,
    b,
    best,
}: {
    a: LatLng
    b: LatLng
    best: LatLng
}) {
    const points = useMemo(() => [a, b, best], [a, b, best])
    const center: [number, number] = [best.lat, best.lng]
    return (
    <div className={styles.container}>
      <MapContainer center={center} zoom={13} className={styles.map}>
        <TileLayer
          attribution='&copy; OpenStreetMap contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />

        <ScrollWheelControl />
        <FitBounds points={points} />

        <Marker position={[a.lat, a.lng]} icon={oceanicIcon}>
          <Popup>
            <strong>Origin A</strong>
            <br />
            {a.lat.toFixed(6)}, {a.lng.toFixed(6)}
          </Popup>
        </Marker>

        <Marker position={[b.lat, b.lng]} icon={oceanicIcon}>
          <Popup>
            <strong>Origin B</strong>
            <br />
            {b.lat.toFixed(6)}, {b.lng.toFixed(6)}
          </Popup>
        </Marker>

        <Marker position={[best.lat, best.lng]} icon={amberIcon}>
          <Popup>
            <strong>Best meeting point</strong>
            <br />
            {best.lat.toFixed(6)}, {best.lng.toFixed(6)}
          </Popup>
        </Marker>
      </MapContainer>
    </div>
  )
}