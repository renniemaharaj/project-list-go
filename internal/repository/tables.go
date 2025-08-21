package repository

import (
	"context"
	"fmt"
)

// Intended to be called once only. Initialize tables (relational)
func InitializeDTables(ctx context.Context) error {
	dbx, err := GETDBX()
	if err != nil {
		return err
	}

	queries := []string{
		// 1. Consultants
		`CREATE TABLE IF NOT EXISTS consultants (
			id SERIAL PRIMARY KEY,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			profile_picture TEXT
		);`,

		// 2. Consultant Roles (one-to-many)
		`CREATE TABLE IF NOT EXISTS consultant_roles (
			consultant_id INTEGER REFERENCES consultants(id) ON DELETE CASCADE,
			role VARCHAR(50) NOT NULL,
			PRIMARY KEY (consultant_id, role)
		);`,

		// 3. Status
		`CREATE TABLE IF NOT EXISTS status (
			id SERIAL PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			date_created TIMESTAMP DEFAULT NOW(),
			description TEXT
		);`,

		// 4. Projects
		`CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			projected_start_date TIMESTAMP,
			start_date TIMESTAMP,
			projected_end_date TIMESTAMP,
			end_date TIMESTAMP,
			number VARCHAR(50) NOT NULL,
			name TEXT NOT NULL,
			hours_assigned NUMERIC(10,2),
			manager_id INTEGER REFERENCES consultants(id),
			description TEXT
		);`,

		// 5. Project tags (many-to-many string tags)
		`CREATE TABLE IF NOT EXISTS project_tags (
			project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
			tag VARCHAR(100) NOT NULL,
			PRIMARY KEY (project_id, tag)
		);`,

		// 6. Time entries
		`CREATE TABLE IF NOT EXISTS time_entries (
			id SERIAL PRIMARY KEY,
			hours NUMERIC(6,2) NOT NULL,
			title TEXT NOT NULL,
			description TEXT,
			consultant_id INTEGER REFERENCES consultants(id) ON DELETE SET NULL,
			project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
			entry_date TIMESTAMP DEFAULT NOW()
		);`,

		// 7. Project status history (project ↔ status timeline)
		`CREATE TABLE IF NOT EXISTS status_history (
			project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
			status_id INTEGER REFERENCES status(id),
			changed_at TIMESTAMP DEFAULT NOW(),
			PRIMARY KEY (project_id, status_id, changed_at)
		);`,
	}

	for _, query := range queries {
		if _, err := dbx.NewQuery(query).Execute(); err != nil {
			return fmt.Errorf("init tables error: %w", err)
		}
	}

	return nil
}
