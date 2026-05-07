import { useEffect, useMemo } from "react";
import { MapContainer, Marker, Popup, TileLayer, useMap } from "react-leaflet";
import L from 'leaflet'
import styles from './MapView.module.css'

export type LatLng = { lat: number; lng: number}

function FitBounds ({ points}: {points: LatLng[]}) {
    const map = useMap()

    useEffect(() => {
        if (points.length === 0) return
        const bounds = L.latLngBounds(points.map((p) => [p.lat, p.lng] as [number, number]))
        map.fitBounds(bounds, { padding: [30, 30]})
    }, [map, points])
    return null
}

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
      <MapContainer center={center} zoom={13} className={styles.map} scrollWheelZoom>
        <TileLayer
          attribution='&copy; OpenStreetMap contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />

        <FitBounds points={points} />

        <Marker position={[a.lat, a.lng]}>
          <Popup>
            <strong>Origin A</strong>
            <br />
            {a.lat.toFixed(6)}, {a.lng.toFixed(6)}
          </Popup>
        </Marker>

        <Marker position={[b.lat, b.lng]}>
          <Popup>
            <strong>Origin B</strong>
            <br />
            {b.lat.toFixed(6)}, {b.lng.toFixed(6)}
          </Popup>
        </Marker>

        <Marker position={[best.lat, best.lng]}>
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