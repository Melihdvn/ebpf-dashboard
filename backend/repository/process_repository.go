package repository

import (
	"database/sql"
	"ebpf-dashboard/models"
)

type ProcessRepository interface {
	SaveProcess(p models.ProcessEvent) error
	SaveProcesses(processes []models.ProcessEvent) error
	GetRecentProcesses(limit int) ([]models.ProcessEvent, error)
}

type processRepository struct {
	db *sql.DB
}

func NewProcessRepository(db *sql.DB) ProcessRepository {
	return &processRepository{db: db}
}

func (r *processRepository) SaveProcess(p models.ProcessEvent) error {
	_, err := r.db.Exec(
		"INSERT INTO processes (time, pid, comm, args) VALUES (?, ?, ?, ?)",
		p.Time, p.PID, p.Comm, p.Args,
	)
	return err
}

func (r *processRepository) SaveProcesses(processes []models.ProcessEvent) error {
	if len(processes) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO processes (time, pid, comm, args) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, p := range processes {
		if _, err := stmt.Exec(p.Time, p.PID, p.Comm, p.Args); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *processRepository) GetRecentProcesses(limit int) ([]models.ProcessEvent, error) {
	rows, err := r.db.Query(
		"SELECT id, timestamp, time, pid, comm, args FROM processes ORDER BY id DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ProcessEvent
	for rows.Next() {
		var p models.ProcessEvent
		if err := rows.Scan(&p.ID, &p.Timestamp, &p.Time, &p.PID, &p.Comm, &p.Args); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return results, nil
}
