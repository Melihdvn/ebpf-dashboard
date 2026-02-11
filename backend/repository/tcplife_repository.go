package repository

import (
	"database/sql"
	"ebpf-dashboard/models"
	"time"
)

type TCPLifeRepository struct {
	db *sql.DB
}

func NewTCPLifeRepository(db *sql.DB) *TCPLifeRepository {
	return &TCPLifeRepository{db: db}
}

// SaveTCPLifeEvents saves multiple TCP lifecycle events to the database
func (r *TCPLifeRepository) SaveTCPLifeEvents(events []models.TCPLifeEvent) error {
	if len(events) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO tcp_lifecycle (pid, comm, local_addr, local_port, remote_addr, remote_port, tx_kb, rx_kb, duration_ms)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, event := range events {
		_, err := stmt.Exec(
			event.PID,
			event.Comm,
			event.LocalAddr,
			event.LocalPort,
			event.RemoteAddr,
			event.RemotePort,
			event.TxKB,
			event.RxKB,
			event.DurationMS,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetRecentTCPLifeEvents retrieves the most recent TCP lifecycle events
func (r *TCPLifeRepository) GetRecentTCPLifeEvents(limit int) ([]models.TCPLifeEvent, error) {
	query := `
		SELECT id, timestamp, pid, comm, local_addr, local_port, remote_addr, remote_port, tx_kb, rx_kb, duration_ms
		FROM tcp_lifecycle
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.TCPLifeEvent
	for rows.Next() {
		var event models.TCPLifeEvent
		var timestamp string

		err := rows.Scan(
			&event.ID,
			&timestamp,
			&event.PID,
			&event.Comm,
			&event.LocalAddr,
			&event.LocalPort,
			&event.RemoteAddr,
			&event.RemotePort,
			&event.TxKB,
			&event.RxKB,
			&event.DurationMS,
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
				event.Timestamp = t
				break
			}
		}

		events = append(events, event)
	}

	return events, nil
}
