"use client";

import { useState, useEffect } from "react";
import { CATEGORIES, Category, Place } from "@/types";
import { fetchPlaces } from "@/lib/api";
import CategoryTabs from "@/components/CategoryTabs";
import PlaceMap from "@/components/PlaceMap";
import PlaceCard from "@/components/PlaceCard";

interface Props {
  cityName: string;
}

export default function CityView({ cityName }: Props) {
  const [activeCategory, setActiveCategory] = useState<Category>("food");
  const [places, setPlaces] = useState<Place[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function load() {
      setLoading(true);
      setError(null);

      try {
        const data = await fetchPlaces(cityName, activeCategory);
        setPlaces(data.places || []);
      } catch (err) {
        setError("Could not load places. Try again.");
        setPlaces([]);
      } finally {
        setLoading(false);
      }
    }

    load();
  }, [cityName, activeCategory]);

  return (
    <div className="min-h-screen bg-gray-950 text-white">

      {/* Header */}
      <div className="bg-gray-900 border-b border-gray-800 px-4 py-6">
        <div className="max-w-6xl mx-auto">
          <h1 className="text-3xl font-bold capitalize mb-1">
            {cityName}
          </h1>
          <p className="text-gray-400">
            Explore what this city has to offer
          </p>
        </div>
      </div>

      {/* Category Tabs */}
      <div className="bg-gray-900 border-b border-gray-800 px-4">
        <div className="max-w-6xl mx-auto">
          <CategoryTabs
            active={activeCategory}
            onChange={setActiveCategory}
          />
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-6xl mx-auto px-4 py-6">

        {loading && (
          <div className="flex items-center justify-center py-24">
            <div className="text-gray-400 text-lg">
              Loading {activeCategory} spots...
            </div>
          </div>
        )}

        {error && (
          <div className="bg-red-900/30 border border-red-700
                          rounded-xl p-6 text-center text-red-300">
            {error}
          </div>
        )}

        {!loading && !error && (
          <div className="flex flex-col gap-6">

            {/* Map */}
            <div className="h-96 rounded-xl overflow-hidden
                            border border-gray-800">
              <PlaceMap places={places} />
            </div>

            {/* Count */}
            <p className="text-gray-400 text-sm">
              Showing {places.length} {activeCategory} spots
            </p>

            {/* Places Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2
                            lg:grid-cols-3 gap-4">
              {places.map((place) => (
                <PlaceCard key={place.id} place={place} />
              ))}
            </div>

            {places.length === 0 && (
              <div className="text-center py-16 text-gray-500">
                No {activeCategory} spots found for {cityName}.
              </div>
            )}

          </div>
        )}

      </div>
    </div>
  );
}