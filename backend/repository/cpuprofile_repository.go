package repository

import (
	"database/sql"
	"ebpf-dashboard/models"
	"time"
)

type CPUProfileRepository struct {
	db *sql.DB
}

func NewCPUProfileRepository(db *sql.DB) *CPUProfileRepository {
	return &CPUProfileRepository{db: db}
}

// SaveCPUProfiles saves multiple CPU profile samples to the database
func (r *CPUProfileRepository) SaveCPUProfiles(profiles []models.CPUProfile) error {
	if len(profiles) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO cpu_profiles (process_name, stack_trace, sample_count)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, profile := range profiles {
		_, err := stmt.Exec(profile.ProcessName, profile.StackTrace, profile.SampleCount)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetRecentCPUProfiles retrieves the most recent CPU profile samples
func (r *CPUProfileRepository) GetRecentCPUProfiles(limit int) ([]models.CPUProfile, error) {
	query := `
		SELECT id, timestamp, process_name, stack_trace, sample_count
		FROM cpu_profiles
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []models.CPUProfile
	for rows.Next() {
		var profile models.CPUProfile
		var timestamp string

		err := rows.Scan(
			&profile.ID,
			&timestamp,
			&profile.ProcessName,
			&profile.StackTrace,
			&profile.SampleCount,
		)
		if err != nil {
			return nil, err
		}

		// Parse timestamp - try multiple formats
		// SQLite might return different formats depending on configuration
		formats := []string{
			"2006-01-02 15:04:05",
			time.RFC3339,
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05",
		}

		for _, format := range formats {
			if t, err := time.Parse(format, timestamp); err == nil {
				profile.Timestamp = t
				break
			}
		}

		// If parsing failed for all formats, profile.Timestamp will be zero value

		profiles = append(profiles, profile)
	}

	return profiles, nil
}
