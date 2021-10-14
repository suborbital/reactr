package rcap

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	ErrQueryNotFound     = errors.New("query not found")
	ErrQueryNotPrepared  = errors.New("query not prepared")
	ErrQueryTypeMismatch = errors.New("query type incorrect")
)

type DatabaseCapability interface {
	ExecInsertQuery(name string, vars map[string]string) (interface{}, error)
	Prepare(q *Query) error
}

type DatabaseConfig struct {
	Enabled          bool   `json:"enabled" yaml:"enabled"`
	ConnectionString string `json:"connectionString" yaml:"connectionString"`
}

type QueryType int

const (
	QueryTypeInsert QueryType = iota
)

type Query struct {
	Type     QueryType
	Name     string
	VarCount int
	Query    string

	stmt *sqlx.Stmt
}

// SqlDatabase is an SQL implementation of DatabaseCapability
type SqlDatabase struct {
	config *DatabaseConfig
	db     *sqlx.DB

	preparedQuery map[string]*Query
}

// NewSqlDatabase creates a new SQL database
func NewSqlDatabase(config *DatabaseConfig) DatabaseCapability {
	if !config.Enabled || config.ConnectionString == "" {
		return nil
	}

	db, err := sqlx.Connect("mysql", config.ConnectionString)
	if err != nil {
		fmt.Println("CONNECT FAILED", err)
		return nil
	}

	fmt.Println("connected!")

	s := &SqlDatabase{
		config:        config,
		db:            db,
		preparedQuery: map[string]*Query{},
	}

	q := &Query{
		Type:     QueryTypeInsert,
		Name:     "InsertTest",
		VarCount: 0,
		Query: `
		INSERT INTO users
			(uuid, email, created_at, state)
		VALUES
			('asdfdfghj', 'connor+1@suborbital.dev', NOW(), 'A')`,
	}

	if err := s.Prepare(q); err != nil {
		fmt.Println("FAILED TO PREPARE:", err)
		return nil
	}

	return s
}

// ExecInsertQuery executes a prepared Insert query
func (s *SqlDatabase) ExecInsertQuery(name string, vars map[string]string) (interface{}, error) {
	if !s.config.Enabled {
		return nil, ErrCapabilityNotEnabled
	}

	query, exists := s.preparedQuery[name]
	if !exists {
		return nil, ErrQueryNotFound
	}

	if query.Type != QueryTypeInsert {
		return nil, ErrQueryTypeMismatch
	}

	if query.stmt == nil {
		return nil, ErrQueryNotPrepared
	}

	result, err := query.stmt.Exec()
	if err != nil {
		return nil, errors.Wrap(err, "failed to Exec")
	}

	// no need to check error, if insertID is 0, that's fine
	insertID, _ := result.LastInsertId()

	return insertID, nil
}

func (s *SqlDatabase) Prepare(q *Query) error {
	stmt, err := s.db.Preparex(q.Query)
	if err != nil {
		return errors.Wrap(err, "failed to Prepare")
	}

	q.stmt = stmt

	s.preparedQuery[q.Name] = q

	return nil
}
