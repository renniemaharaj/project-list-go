package repository

import (
	"context"
	"fmt"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// InitializeDatabaseTables initializes the database schema for the project-tracking domain.
//
// The creation is **idempotent** (uses CREATE TABLE IF NOT EXISTS) so it is safe
// to call multiple times, though it is intended to be called once on startup.
//
// ── Overview of the logical model ──────────────────────────────────────────────
//
// 1) Domain (dimension) – relation‑independent
//   - consultants               -> master record for a consultant (people catalog)
//
// 2) Domain (dimension) – relation‑dependent
//   - projects                  -> projects managed by/for consultants
//
// 3) Non‑domain (fact/association/operational) – relation‑dependent
//   - project_time_entries              -> hours logged by consultant(s) on project(s)
//   - project_statuses                  -> status notes per project by consultant(s)
//   - consultant_roles          -> 1‑to‑many roles per consultant (role catalog)
//   - project_tags              -> N‑to‑N string tags per project
//   - project_consultants       -> N‑to‑N assignment of consultants to projects
//
// Table creation order respects foreign‑key dependencies:
//
//	consultants -> projects -> (project_time_entries, project_statuses, consultant_roles, project_tags, project_consultants)
//
// Error handling: schema creation is wrapped in a transaction using UseTransaction API.
// Any failure aborts and rolls back changes automatically.
func InitializeDatabaseTables(ctx context.Context, r *repository) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		// 1) Domain (independent)
		if err := createDomainIndependentTables(ctx, tx); err != nil {
			return fmt.Errorf("init tables (domain independent) error: %w", err)
		}

		// 2) Domain (dependent)
		if err := createDomainDependentTables(ctx, tx); err != nil {
			return fmt.Errorf("init tables (domain dependent) error: %w", err)
		}

		// 3) Non‑domain (dependent)
		if err := createNonDomainDependentTables(ctx, tx); err != nil {
			return fmt.Errorf("init tables (non-domain dependent) error: %w", err)
		}

		return nil
	})
}

// createDomainIndependentTables creates tables that do not depend on any other
// relations. These are foundational catalogs/dimensions used by the rest of the
// schema.
func createDomainIndependentTables(ctx context.Context, tx *dbx.Tx) error {
	queries := []string{
		// consultants -- single source of truth for people (one row per person)
		`CREATE TABLE IF NOT EXISTS consultants (
			id              SERIAL PRIMARY KEY,
			first_name      VARCHAR(100) NOT NULL,
			last_name       VARCHAR(100) NOT NULL,
			email           VARCHAR(255) UNIQUE NOT NULL,
			profile_picture TEXT
		);`,
	}
	return runQueries(tx, queries)
}

// createDomainDependentTables creates domain relations that depend on domain
// catalogs already defined (e.g., projects referencing consultants as manager).
func createDomainDependentTables(ctx context.Context, tx *dbx.Tx) error {
	queries := []string{
		// projects -- core domain entity
		//
		// Notes:
		//	 manager_id is optional; ON DELETE default (RESTRICT) prevents orphaned
		//   managers from being deleted silently. Use explicit deletion policy as needed.
		`CREATE TABLE IF NOT EXISTS projects (
			id                   SERIAL PRIMARY KEY,
			manager_id           INTEGER REFERENCES consultants(id),
			number               VARCHAR(50) NOT NULL,
			name                 TEXT NOT NULL,
			start_date           TIMESTAMP,
			projected_start_date TIMESTAMP,
			end_date             TIMESTAMP,
			projected_end_date   TIMESTAMP,
			description          TEXT
		);`,
		// Useful uniqueness to prevent duplicate business numbers per project record.
		`CREATE UNIQUE INDEX IF NOT EXISTS ux_projects_number ON projects(number);`,
	}
	return runQueries(tx, queries)
}

// createNonDomainDependentTables creates operational/fact and association tables
// that depend on the domain entities.
func createNonDomainDependentTables(ctx context.Context, tx *dbx.Tx) error {
	queries := []string{
		// project_time_entries -- hours logged against a project (and usually by a consultant)
		`CREATE TABLE IF NOT EXISTS project_time_entries (
			id            SERIAL PRIMARY KEY,
			project_id    INTEGER REFERENCES projects(id) ON DELETE CASCADE,
			consultant_id INTEGER REFERENCES consultants(id) ON DELETE SET NULL,
			type          VARCHAR(10) NOT NULL,
			hours         NUMERIC(6,2) NOT NULL,
			title         TEXT NOT NULL,
			description   TEXT,
			entry_date    TIMESTAMP DEFAULT NOW()
		);`,

		// project_statuses -- status notes per project, optionally by a consultant
		`CREATE TABLE IF NOT EXISTS project_statuses (
			id            SERIAL PRIMARY KEY,
			project_id    INTEGER REFERENCES projects(id) ON DELETE CASCADE,
			consultant_id INTEGER REFERENCES consultants(id) ON DELETE SET NULL,
			title         VARCHAR(100) NOT NULL,
			date_created  TIMESTAMP DEFAULT NOW(),
			description   TEXT
		);`,

		// consultant_roles -- one‑to‑many roles per consultant (composite PK)
		// PRIMARY KEY (consultant_id, role) removed
		`CREATE TABLE IF NOT EXISTS consultant_roles (
			id            SERIAL PRIMARY KEY,
			consultant_id INTEGER REFERENCES consultants(id) ON DELETE CASCADE,
			role          VARCHAR(50) NOT NULL
		);`,

		// project_tags -- free‑text tags assigned to projects (composite PK)
		// PRIMARY KEY (project_id, tag) removed
		`CREATE TABLE IF NOT EXISTS project_tags (
			id            SERIAL PRIMARY KEY,
			project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
			tag        VARCHAR(100) NOT NULL
		);`,

		// project_consultants -- many‑to‑many relation between consultants and projects
		// PRIMARY KEY (consultant_id, project_id) removed
		`CREATE TABLE IF NOT EXISTS project_consultants (
			id            SERIAL PRIMARY KEY,
			project_id    INTEGER REFERENCES projects(id) ON DELETE CASCADE,
			consultant_id INTEGER REFERENCES consultants(id) ON DELETE CASCADE,
			role          VARCHAR(50) NOT NULL
		);`,
	}
	return runQueries(tx, queries)
}

// runQueries executes the DDL statements in order, returning the first error
// encountered. Each statement is executed standalone for clarity and easier
// error localization.
func runQueries(tx *dbx.Tx, queries []string) error {
	for i, stmt := range queries {
		if _, err := tx.NewQuery(stmt).Execute(); err != nil {
			return fmt.Errorf("ddl step %d failed: %w", i+1, err)
		}
	}
	return nil
}
