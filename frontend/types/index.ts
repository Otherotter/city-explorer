export type Category =
  | "food"
  | "nature"
  | "study"
  | "events"
  | "attraction"
  | "thrift"
  | "social"
  | "architecture";

export interface Place {
  id: number;
  name: string;
  subcategory: string;
  address: string;
  neighborhood: string;
  latitude: number;
  longitude: number;
  website: string;
  phone: string;
  source: string;
}

export interface CityResponse {
  city: string;
  count: number;
  places: Place[];
}

export const CATEGORIES: {
  slug: Category;
  label: string;
  emoji: string;
}[] = [
  { slug: "food",         label: "Food & Drink",   emoji: "🍜" },
  { slug: "nature",       label: "Nature",          emoji: "🌿" },
  { slug: "study",        label: "Study Spots",     emoji: "📚" },
  { slug: "events",       label: "Events",          emoji: "🎉" },
  { slug: "attraction",   label: "Attractions",     emoji: "🏛️" },
  { slug: "thrift",       label: "Thrift",          emoji: "🛍️" },
  { slug: "social",       label: "Social",          emoji: "🤝" },
  { slug: "architecture", label: "Architecture",    emoji: "🏗️" },
];