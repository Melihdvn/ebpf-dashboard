package repository

import (
	"database/sql"
	"ebpf-dashboard/models"
)

type DiskRepository interface {
	SaveLatencySnapshot(latencies []models.DiskLatency) error
	GetLatestLatency(limit int) ([]models.DiskLatency, error)
}

type diskRepository struct {
	db *sql.DB
}

func NewDiskRepository(db *sql.DB) DiskRepository {
	return &diskRepository{db: db}
}

func (r *diskRepository) SaveLatencySnapshot(latencies []models.DiskLatency) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		"INSERT INTO disk_latency (range_min, range_max, count) VALUES (?, ?, ?)",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, lat := range latencies {
		if _, err := stmt.Exec(lat.RangeMin, lat.RangeMax, lat.Count); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *diskRepository) GetLatestLatency(limit int) ([]models.DiskLatency, error) {
	rows, err := r.db.Query(
		`SELECT id, timestamp, range_min, range_max, count 
		FROM disk_latency ORDER BY id DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.DiskLatency
	for rows.Next() {
		var lat models.DiskLatency
		if err := rows.Scan(&lat.ID, &lat.Timestamp, &lat.RangeMin, &lat.RangeMax, &lat.Count); err != nil {
			return nil, err
		}
		results = append(results, lat)
	}
	return results, nil
}
