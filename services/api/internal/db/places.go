package db

import (
	"database/sql"
	"time"
)

// DataStatus describes the state of data
// for a given city and category.
type DataStatus struct {
	Count         int
	LastCollected time.Time
	IsEmpty       bool
	IsStale       bool
}

// CheckDataStatus tells the API whether it needs
// to trigger the Collector before returning data.
func CheckDataStatus(database *sql.DB, cityID int, categorySlug string, ttlDays int) (*DataStatus, error) {
	status := &DataStatus{}

	var lastCollected sql.NullTime

	err := database.QueryRow(`
		SELECT COUNT(*), MAX(p.created_at)
		FROM places p
		JOIN categories c ON p.category_id = c.id
		WHERE p.city_id = $1
		AND c.slug = $2
		AND p.source = 'overpass'
	`, cityID, categorySlug).Scan(&status.Count, &lastCollected)

	if err != nil {
		return nil, err
	}

	status.IsEmpty = status.Count == 0

	if lastCollected.Valid {
		status.LastCollected = lastCollected.Time
		ttl := time.Duration(ttlDays) * 24 * time.Hour
		status.IsStale = time.Since(lastCollected.Time) > ttl
	} else {
		status.IsStale = true
	}

	return status, nil
}
