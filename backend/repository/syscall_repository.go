package repository

import (
	"database/sql"
	"ebpf-dashboard/models"
	"time"
)

type SyscallRepository struct {
	db *sql.DB
}

func NewSyscallRepository(db *sql.DB) *SyscallRepository {
	return &SyscallRepository{db: db}
}

// SaveSyscallStats saves multiple syscall statistics to the database
func (r *SyscallRepository) SaveSyscallStats(stats []models.SyscallStat) error {
	if len(stats) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO syscall_stats (syscall_name, count)
		VALUES (?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, stat := range stats {
		_, err := stmt.Exec(stat.SyscallName, stat.Count)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetRecentSyscallStats retrieves recent syscall statistics
// It groups by syscall name and sums up the counts for the requested limit period
// Or returns raw entries depending on visualization needs.
// For a Pie Chart, we usually want aggregated data over the last X minutes.
// Here we return raw entries, aggregation can be done in frontend or via a different query.
func (r *SyscallRepository) GetRecentSyscallStats(limit int) ([]models.SyscallStat, error) {
	query := `
		SELECT id, timestamp, syscall_name, count
		FROM syscall_stats
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.SyscallStat
	for rows.Next() {
		var stat models.SyscallStat
		var timestamp string

		err := rows.Scan(
			&stat.ID,
			&timestamp,
			&stat.SyscallName,
			&stat.Count,
		)
		if err != nil {
			return nil, err
		}

		// Parse timestamp - support multiple formats
		formats := []string{
			"2006-01-02 15:04:05",
			time.RFC3339,
			"2006-01-02T15:04:05Z",
		}

		for _, format := range formats {
			if t, err := time.Parse(format, timestamp); err == nil {
				stat.Timestamp = t
				break
			}
		}

		stats = append(stats, stat)
	}

	return stats, nil
}
