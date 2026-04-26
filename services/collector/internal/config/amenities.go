package config

// AmenitiesForCategory returns the OSM amenity types
// we care about for each category slug.
// These map directly to OpenStreetMap amenity tags.
func AmenitiesForCategory(slug string) []string {
	switch slug {
	case "food":
		return []string{
			"restaurant",
			"cafe",
			"food_court",
			"fast_food",
			"bar",
			"pub",
			"ice_cream",
			"bakery",
		}
	case "nature":
		return []string{
			"park",
		}
	case "study":
		return []string{
			"library",
			"university",
			"college",
		}
	case "social":
		return []string{
			"community_centre",
			"social_centre",
		}
	default:
		return []string{}
	}
}

// TTLForCategory returns how long collected data
// is considered fresh for each category.
// After this duration the API will re-collect.
func TTLForCategory(slug string) int {
	switch slug {
	case "events":
		return 1 // 1 day
	case "food", "thrift", "study":
		return 30 // 30 days
	case "nature", "attraction":
		return 90 // 90 days
	case "social":
		return 7 // 7 days
	default:
		return 30
	}
}
