const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function fetchPlaces(
  city: string,
  category: string
): Promise<{ city: string; count: number; places: any[] }> {
  const encoded = encodeURIComponent(city);
  const res = await fetch(`${API_URL}/cities/${encoded}/${category}`, {
    // Revalidate every 5 minutes
    next: { revalidate: 300 },
  });

  if (!res.ok) {
    throw new Error(`Failed to fetch ${category} for ${city}`);
  }

  return res.json();
}