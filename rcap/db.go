package rcap

import (
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	ErrDatabaseTypeInvalid = errors.New("database type invalid")
	ErrQueryNotFound       = errors.New("query not found")
	ErrQueryNotPrepared    = errors.New("query not prepared")
	ErrQueryTypeMismatch   = errors.New("query type incorrect")
	ErrQueryTypeInvalid    = errors.New("query type invalid")
	ErrQueryVarsMismatch   = errors.New("number of variables incorrect")
)

type DatabaseCapability interface {
	ExecQuery(queryType int32, name string, vars []interface{}) ([]byte, error)
	Prepare(q *Query) error
}

type DatabaseConfig struct {
	Enabled          bool   `json:"enabled" yaml:"enabled"`
	DBType           string `json:"dbType" yaml:"dbType"`
	ConnectionString string `json:"connectionString" yaml:"connectionString"`
}

const (
	DBTypeMySQL    = "mysql"
	DBTypePostgres = "pgx"
)

type QueryType int32

const (
	QueryTypeInsert QueryType = QueryType(0)
	QueryTypeSelect QueryType = QueryType(1)
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

	queries map[string]*Query
}

// NewSqlDatabase creates a new SQL database
func NewSqlDatabase(config *DatabaseConfig) (DatabaseCapability, error) {
	if !config.Enabled || config.ConnectionString == "" {
		return nil, nil
	}

	if config.DBType != DBTypeMySQL && config.DBType != DBTypePostgres {
		return nil, ErrDatabaseTypeInvalid
	}

	db, err := sqlx.Connect(config.DBType, config.ConnectionString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Connect")
	}

	s := &SqlDatabase{
		config:  config,
		db:      db,
		queries: map[string]*Query{},
	}

	// q := &Query{
	// 	Type:     QueryTypeInsert,
	// 	Name:     "InsertTest",
	// 	VarCount: 2,
	// 	Query: `
	// 	INSERT INTO users
	// 		(uuid, email, created_at, state)
	// 	VALUES
	// 		(?, ?, NOW(), 'A')`,
	// }

	// if err := s.Prepare(q); err != nil {
	// 	fmt.Println("FAILED TO PREPARE:", err)
	// 	return nil, errors.Wrap(err, "failed to Prepare query")
	// }

	// q2 := &Query{
	// 	Type:     QueryTypeSelect,
	// 	Name:     "SelectUserWithEmail",
	// 	VarCount: 1,
	// 	Query: `
	// 	SELECT * FROM users
	// 	WHERE email = ?`,
	// }

	// if err := s.Prepare(q2); err != nil {
	// 	fmt.Println("FAILED TO PREPARE Q2:", err)
	// 	return nil, errors.Wrap(err, "failed to Prepare query")
	// }

	q3 := &Query{
		Type:     QueryTypeSelect,
		Name:     "PGSelectUserWithEmail",
		VarCount: 1,
		Query: `
		SELECT * FROM users
		WHERE email = $1`,
	}

	if err := s.Prepare(q3); err != nil {
		fmt.Println("FAILED TO PREPARE Q3:", err)
		return nil, errors.Wrap(err, "failed to Prepare query")
	}

	return s, nil
}

func (s *SqlDatabase) Prepare(q *Query) error {
	stmt, err := s.db.Preparex(q.Query)
	if err != nil {
		return errors.Wrap(err, "failed to Prepare")
	}

	q.stmt = stmt

	s.queries[q.Name] = q

	return nil
}

func (s *SqlDatabase) ExecQuery(queryType int32, name string, vars []interface{}) ([]byte, error) {
	// the returned data varies depending on the query type
	fmt.Println("QUERY TYPE:", queryType)

	switch QueryType(queryType) {
	case QueryTypeInsert:
		return s.execInsertQuery(name, vars)
	case QueryTypeSelect:
		return s.execSelectQuery(name, vars)
	}

	return nil, ErrQueryTypeInvalid
}

// execInsertQuery executes a prepared Insert query
func (s *SqlDatabase) execInsertQuery(name string, vars []interface{}) ([]byte, error) {
	if !s.config.Enabled {
		return nil, ErrCapabilityNotEnabled
	}

	query, exists := s.queries[name]
	if !exists {
		return nil, ErrQueryNotFound
	}

	if query.Type != QueryTypeInsert {
		return nil, ErrQueryTypeMismatch
	}

	if query.stmt == nil {
		return nil, ErrQueryNotPrepared
	}

	if query.VarCount != len(vars) {
		return nil, ErrQueryVarsMismatch
	}

	result, err := query.stmt.Exec(vars...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Exec")
	}

	// no need to check error, if insertID is 0, that's fine
	insertID, _ := result.LastInsertId()

	idBytes := make([]byte, binary.MaxVarintLen64)
	len := binary.PutVarint(idBytes, insertID)

	return idBytes[:len], nil
}

// execSelectQuery executes a prepared Select query
func (s *SqlDatabase) execSelectQuery(name string, vars []interface{}) ([]byte, error) {
	if !s.config.Enabled {
		return nil, ErrCapabilityNotEnabled
	}

	query, exists := s.queries[name]
	if !exists {
		return nil, ErrQueryNotFound
	}

	if query.Type != QueryTypeSelect {
		return nil, ErrQueryTypeMismatch
	}

	if query.stmt == nil {
		return nil, ErrQueryNotPrepared
	}

	if query.VarCount != len(vars) {
		return nil, ErrQueryVarsMismatch
	}

	rows, err := query.stmt.Query(vars...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to stmt.Query")
	}

	defer rows.Close()
	result, err := rowsToMap(rows)
	if err != nil {
		return nil, errors.Wrap(err, "failed to rowsToMap")
	}

	destJSON, err := json.Marshal(result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Marshal query result")
	}

	return destJSON, nil
}

func rowsToMap(rows *sql.Rows) ([]map[string]interface{}, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Columns from query result")
	}

	results := []map[string]interface{}{}

	for {
		if moreRows := rows.Next(); !moreRows {
			if rows.Err() != nil {
				return nil, errors.Wrap(err, "failed to rows.Next")
			}

			break
		}

		dest := make([]interface{}, len(cols))
		for i := range dest {
			var val []byte
			dest[i] = &val
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, errors.Wrap(err, "failed to Scan row")
		}

		result := map[string]interface{}{}

		for i, c := range cols {
			bytes, _ := dest[i].(*[]byte)
			result[c] = string(*bytes)
		}

		results = append(results, result)
	}

	return results, nil
}