package repository

import (
	"database/sql"
	"ebpf-dashboard/models"
)

type NetworkRepository interface {
	SaveConnection(conn models.NetworkConnection) error
	GetRecentConnections(limit int) ([]models.NetworkConnection, error)
}

type networkRepository struct {
	db *sql.DB
}

func NewNetworkRepository(db *sql.DB) NetworkRepository {
	return &networkRepository{db: db}
}

func (r *networkRepository) SaveConnection(conn models.NetworkConnection) error {
	_, err := r.db.Exec(
		`INSERT INTO network_connections 
		(pid, comm, ip_version, source_addr, source_port, dest_addr, dest_port) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		conn.PID, conn.Comm, conn.IPVersion, conn.SourceAddr,
		conn.SourcePort, conn.DestAddr, conn.DestPort,
	)
	return err
}

func (r *networkRepository) GetRecentConnections(limit int) ([]models.NetworkConnection, error) {
	rows, err := r.db.Query(
		`SELECT id, timestamp, pid, comm, ip_version, source_addr, source_port, dest_addr, dest_port 
		FROM network_connections ORDER BY id DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.NetworkConnection
	for rows.Next() {
		var conn models.NetworkConnection
		if err := rows.Scan(
			&conn.ID, &conn.Timestamp, &conn.PID, &conn.Comm, &conn.IPVersion,
			&conn.SourceAddr, &conn.SourcePort, &conn.DestAddr, &conn.DestPort,
		); err != nil {
			return nil, err
		}
		results = append(results, conn)
	}
	return results, nil
}
