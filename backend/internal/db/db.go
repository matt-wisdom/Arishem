package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var placeholderRegex = regexp.MustCompile(`\$\d+`)

func convertPlaceholders(query string) string {
	// Convert Postgres placeholders ($1, $2, ...) to SQLite (?1, ?2, ...)
	query = placeholderRegex.ReplaceAllStringFunc(query, func(m string) string {
		return "?" + m[1:]
	})
	// Convert Postgres specific functions and defaults to SQLite equivalents
	query = strings.ReplaceAll(query, "DEFAULT gen_random_uuid()", "")
	query = strings.ReplaceAll(query, "gen_random_uuid()", "NULL")
	query = strings.ReplaceAll(query, "TEXT[]", "TEXT")
	query = strings.ReplaceAll(query, "NOW() - INTERVAL '1 day'", "datetime('now', '-1 day')")
	query = strings.ReplaceAll(query, "NOW()", "datetime('now')")
	query = strings.ReplaceAll(query, "= true", "= 1")
	query = strings.ReplaceAll(query, "= false", "= 0")
	return query
}

func convertArgs(args []any) []any {
	newArgs := make([]any, len(args))
	for i, arg := range args {
		if sa, ok := arg.([]string); ok {
			newArgs[i] = "{" + strings.Join(sa, ",") + "}"
		} else {
			newArgs[i] = arg
		}
	}
	return newArgs
}

type SQLiteResult struct {
	sql.Result
}

func (r SQLiteResult) RowsAffected() int64 {
	rows, _ := r.Result.RowsAffected()
	return rows
}

type SQLiteRow struct {
	*sql.Row
}

func (r *SQLiteRow) Scan(dest ...any) error {
	tempDest := make([]any, len(dest))
	stringArrayIndices := make(map[int]*[]string)
	jsonMapIndices := make(map[int]*map[string]any)

	for i, d := range dest {
		if sa, ok := d.(*[]string); ok {
			var s string
			tempDest[i] = &s
			stringArrayIndices[i] = sa
		} else if jm, ok := d.(*map[string]any); ok {
			var b []byte
			tempDest[i] = &b
			jsonMapIndices[i] = jm
		} else {
			tempDest[i] = d
		}
	}

	if err := r.Row.Scan(tempDest...); err != nil {
		return err
	}

	for i, sa := range stringArrayIndices {
		sptr := tempDest[i].(*string)
		if sptr == nil {
			*sa = []string{}
			continue
		}
		s := *sptr
		s = strings.Trim(s, "{}")
		if s == "" {
			*sa = []string{}
		} else {
			*sa = strings.Split(s, ",")
		}
	}

	for i, jm := range jsonMapIndices {
		bptr := tempDest[i].(*[]byte)
		if bptr == nil || len(*bptr) == 0 {
			*jm = make(map[string]any)
			continue
		}
		if err := json.Unmarshal(*bptr, jm); err != nil {
			return fmt.Errorf("failed to unmarshal json map: %w", err)
		}
	}
	return nil
}

type SQLiteRows struct {
	*sql.Rows
}

func (r *SQLiteRows) Scan(dest ...any) error {
	tempDest := make([]any, len(dest))
	stringArrayIndices := make(map[int]*[]string)
	jsonMapIndices := make(map[int]*map[string]any)

	for i, d := range dest {
		if sa, ok := d.(*[]string); ok {
			var s string
			tempDest[i] = &s
			stringArrayIndices[i] = sa
		} else if jm, ok := d.(*map[string]any); ok {
			var b []byte
			tempDest[i] = &b
			jsonMapIndices[i] = jm
		} else {
			tempDest[i] = d
		}
	}

	if err := r.Rows.Scan(tempDest...); err != nil {
		return err
	}

	for i, sa := range stringArrayIndices {
		sptr := tempDest[i].(*string)
		if sptr == nil {
			*sa = []string{}
			continue
		}
		s := *sptr
		s = strings.Trim(s, "{}")
		if s == "" {
			*sa = []string{}
		} else {
			*sa = strings.Split(s, ",")
		}
	}

	for i, jm := range jsonMapIndices {
		bptr := tempDest[i].(*[]byte)
		if bptr == nil || len(*bptr) == 0 {
			*jm = make(map[string]any)
			continue
		}
		if err := json.Unmarshal(*bptr, jm); err != nil {
			return fmt.Errorf("failed to unmarshal json map: %w", err)
		}
	}
	return nil
}

type SQLitePool struct {
	DB *sql.DB
}

func (p *SQLitePool) Exec(ctx context.Context, query string, args ...any) (SQLiteResult, error) {
	upper := strings.ToUpper(query)
	query = convertPlaceholders(query)
	convertedArgs := convertArgs(args)
	res, err := p.DB.ExecContext(ctx, query, convertedArgs...)
	if err != nil {
		if strings.Contains(upper, "CREATE TABLE") || strings.Contains(upper, "CREATE INDEX") || strings.Contains(upper, "ALTER TABLE") {
			return SQLiteResult{res}, nil
		}
		errStr := err.Error()
		if strings.Contains(errStr, "duplicate column name") || strings.Contains(errStr, "already exists") {
			return SQLiteResult{res}, nil
		}
	}
	return SQLiteResult{res}, err
}

func (p *SQLitePool) Query(ctx context.Context, query string, args ...any) (*SQLiteRows, error) {
	query = convertPlaceholders(query)
	convertedArgs := convertArgs(args)
	rows, err := p.DB.QueryContext(ctx, query, convertedArgs...)
	if err != nil {
		return nil, err
	}
	return &SQLiteRows{rows}, nil
}

func (p *SQLitePool) QueryRow(ctx context.Context, query string, args ...any) *SQLiteRow {
	query = convertPlaceholders(query)
	convertedArgs := convertArgs(args)
	row := p.DB.QueryRowContext(ctx, query, convertedArgs...)
	return &SQLiteRow{row}
}

func (p *SQLitePool) Close() {
	p.DB.Close()
}

func (p *SQLitePool) Ping(ctx context.Context) error {
	return p.DB.PingContext(ctx)
}

var Pool *SQLitePool

func Init() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" || strings.HasPrefix(databaseURL, "postgres") {
		databaseURL = "/mnt/C6EE65A1EE658B0F/WORKEST/Arishem/arishem.db"
	}

	dir := filepath.Dir(databaseURL)
	if dir != "" {
		_ = os.MkdirAll(dir, 0755)
	}

	dbConn, err := sql.Open("sqlite3", databaseURL)
	if err != nil {
		return fmt.Errorf("unable to open sqlite database: %w", err)
	}

	if err := dbConn.Ping(); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	Pool = &SQLitePool{DB: dbConn}

	// Enable foreign keys
	_, _ = Pool.DB.Exec("PRAGMA foreign_keys = ON;")

	schema := []string{
		`CREATE TABLE IF NOT EXISTS organizations (
			id TEXT PRIMARY KEY,
			clerk_org_id TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			clerk_user_id TEXT NOT NULL,
			org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
			role TEXT NOT NULL DEFAULT 'org:viewer',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(clerk_user_id, org_id)
		);`,
		`CREATE TABLE IF NOT EXISTS scans (
			id TEXT PRIMARY KEY,
			org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
			title TEXT DEFAULT '',
			type TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'queued',
			target TEXT NOT NULL,
			branch TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			completed_at TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS llm_pentest_runs (
			id TEXT PRIMARY KEY,
			org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
			title TEXT DEFAULT '',
			target_endpoint TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'queued',
			test_modules TEXT,
			logs TEXT DEFAULT '',
			docker INTEGER DEFAULT 0,
			config_mode TEXT DEFAULT 'default',
			api_key TEXT DEFAULT '',
			model TEXT DEFAULT '',
			llm_provider TEXT DEFAULT '',
			api_base TEXT DEFAULT '',
			mode TEXT DEFAULT '',
			budget INTEGER DEFAULT 0,
			concurrency INTEGER DEFAULT 0,
			version INTEGER DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			completed_at TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS findings (
			id TEXT PRIMARY KEY,
			org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
			scan_id TEXT REFERENCES scans(id) ON DELETE CASCADE,
			run_id TEXT REFERENCES llm_pentest_runs(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			severity TEXT NOT NULL,
			description TEXT,
			remediation TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS reports (
			id TEXT PRIMARY KEY,
			org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
			scan_id TEXT REFERENCES scans(id) ON DELETE SET NULL,
			run_id TEXT REFERENCES llm_pentest_runs(id) ON DELETE SET NULL,
			format TEXT NOT NULL,
			storage_key TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS integrations (
			id TEXT PRIMARY KEY,
			org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
			provider TEXT NOT NULL,
			credentials_encrypted TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(org_id, provider)
		);`,
		`CREATE TABLE IF NOT EXISTS alert_rules (
			id TEXT PRIMARY KEY,
			org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
			severity_threshold TEXT NOT NULL,
			channel TEXT NOT NULL,
			channel_config TEXT NOT NULL DEFAULT '{}',
			active INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, stmt := range schema {
		if _, err := Pool.DB.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute schema statement: %w", err)
		}
	}

	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_scans_org_id ON scans(org_id);",
		"CREATE INDEX IF NOT EXISTS idx_scans_status ON scans(status);",
		"CREATE INDEX IF NOT EXISTS idx_llm_pentest_runs_org_id ON llm_pentest_runs(org_id);",
		"CREATE INDEX IF NOT EXISTS idx_findings_scan_id ON findings(scan_id);",
		"CREATE INDEX IF NOT EXISTS idx_findings_run_id ON findings(run_id);",
		"CREATE INDEX IF NOT EXISTS idx_findings_org_id ON findings(org_id);",
		"CREATE INDEX IF NOT EXISTS idx_reports_org_id ON reports(org_id);",
		"CREATE INDEX IF NOT EXISTS idx_alert_rules_org_id ON alert_rules(org_id);",
	}
	for _, stmt := range indexes {
		_, _ = Pool.DB.Exec(stmt)
	}

	return nil
}

func Close() {
	if Pool != nil {
		Pool.Close()
		Pool = nil
	}
}

func GetPool() *SQLitePool {
	return Pool
}