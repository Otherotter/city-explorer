"use client";

import { useState, FormEvent } from "react";
import { useRouter } from "next/navigation";

export default function CitySearch() {
    
  const [city, setCity] = useState("");
  const router = useRouter();

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (!city.trim()) return;

    // Normalize to lowercase with hyphens for URL
    const slug = city.trim().toLowerCase().replace(/\s+/g, "-");
    router.push(`/cities/${slug}`);
  }

  return (
    <form
      onSubmit={handleSubmit}
      className="w-full flex flex-col gap-3"
    >
      <input
        type="text"
        value={city}
        onChange={(e) => setCity(e.target.value)}
        placeholder="Search a city... (e.g. New York City)"
        className="w-full px-5 py-4 rounded-xl bg-gray-800 border
                   border-gray-700 text-white placeholder-gray-500
                   text-lg focus:outline-none focus:border-blue-500
                   focus:ring-1 focus:ring-blue-500 transition"
      />
      <button
        type="submit"
        className="w-full px-5 py-4 rounded-xl bg-blue-600
                   hover:bg-blue-500 text-white font-semibold
                   text-lg transition"
      >
        Explore
      </button>
    </form>
  );
}