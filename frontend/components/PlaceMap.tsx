"use client";

import { useEffect, useRef } from "react";
import { Place } from "@/types";

interface Props {
  places: Place[];
}

export default function PlaceMap({ places }: Props) {
  const mapRef = useRef<HTMLDivElement>(null);
  const mapInstanceRef = useRef<any>(null);

  useEffect(() => {
    if (!mapRef.current) return;
    if (typeof window === "undefined") return;

    // Dynamically import Leaflet to avoid SSR issues
    import("leaflet").then((L) => {
      // Fix default marker icon path issue in Next.js
      delete (L.Icon.Default.prototype as any)._getIconUrl;
      L.Icon.Default.mergeOptions({
        iconRetinaUrl: "/leaflet/marker-icon-2x.png",
        iconUrl: "/leaflet/marker-icon.png",
        shadowUrl: "/leaflet/marker-shadow.png",
      });

      // Only initialize map once
      if (!mapInstanceRef.current) {
        mapInstanceRef.current = L.map(mapRef.current!).setView(
          [40.7128, -74.006],
          12
        );

        L.tileLayer(
          "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png",
          {
            attribution:
              '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
          }
        ).addTo(mapInstanceRef.current);
      }

      // Clear existing markers
      mapInstanceRef.current.eachLayer((layer: any) => {
        if (layer instanceof L.Marker) {
          mapInstanceRef.current.removeLayer(layer);
        }
      });

      // Add new markers
      const validPlaces = places.filter(
        (p) => p.latitude && p.longitude
      );

      validPlaces.forEach((place) => {
        L.marker([place.latitude, place.longitude])
          .addTo(mapInstanceRef.current)
          .bindPopup(`
            <strong>${place.name}</strong><br/>
            ${place.subcategory || ""}<br/>
            ${place.address || ""}
          `);
      });

      // Fit map to markers if we have any
      if (validPlaces.length > 0) {
        const bounds = L.latLngBounds(
          validPlaces.map((p) => [p.latitude, p.longitude])
        );
        mapInstanceRef.current.fitBounds(bounds, { padding: [40, 40] });
      }
    });

    // Add Leaflet CSS
    const link = document.createElement("link");
    link.rel = "stylesheet";
    link.href =
      "https://unpkg.com/leaflet@1.9.4/dist/leaflet.css";
    document.head.appendChild(link);

    return () => {
      if (mapInstanceRef.current) {
        mapInstanceRef.current.remove();
        mapInstanceRef.current = null;
      }
    };
  }, [places]);

  return <div ref={mapRef} className="w-full h-full" />;
}