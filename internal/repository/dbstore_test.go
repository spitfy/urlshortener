package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spitfy/urlshortener/internal/config"
	"github.com/stretchr/testify/assert"
)

// MockDB реализует интерфейс DB
type MockDB struct {
	ExecFunc      func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRowFunc  func(ctx context.Context, sql string, args ...any) pgx.Row
	QueryFunc     func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	PingFunc      func(ctx context.Context) error
	CloseFunc     func()
	BeginFunc     func(ctx context.Context) (pgx.Tx, error)
	BeginTxFunc   func(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	SendBatchFunc func(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

func (m *MockDB) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(ctx, sql, arguments...)
	}
	return pgconn.CommandTag{}, nil
}

func (m *MockDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(ctx, sql, args...)
	}
	return &MockRow{}
}

func (m *MockDB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, sql, args...)
	}
	return &MockRows{}, nil
}

func (m *MockDB) Ping(ctx context.Context) error {
	if m.PingFunc != nil {
		return m.PingFunc(ctx)
	}
	return nil
}

func (m *MockDB) Close() {
	if m.CloseFunc != nil {
		m.CloseFunc()
	}
}

func (m *MockDB) Begin(ctx context.Context) (pgx.Tx, error) {
	if m.BeginFunc != nil {
		return m.BeginFunc(ctx)
	}
	return &MockTx{}, nil
}

func (m *MockDB) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	if m.BeginTxFunc != nil {
		return m.BeginTxFunc(ctx, txOptions)
	}
	return &MockTx{}, nil
}

func (m *MockDB) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	if m.SendBatchFunc != nil {
		return m.SendBatchFunc(ctx, b)
	}
	return &MockBatchResults{}
}

// MockRow для pgx.Row
type MockRow struct {
	ScanFunc func(dest ...any) error
}

func (m *MockRow) Scan(dest ...any) error {
	if m.ScanFunc != nil {
		return m.ScanFunc(dest...)
	}
	return nil
}

// MockRows для pgx.Rows
type MockRows struct {
	CloseFunc func()
	ErrFunc   func() error
	NextFunc  func() bool
	ScanFunc  func(dest ...any) error
}

func (m *MockRows) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}

func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	return make([]pgconn.FieldDescription, 0)
}

func (m *MockRows) Values() ([]any, error) {
	return make([]any, 0), nil
}

func (m *MockRows) RawValues() [][]byte {
	return make([][]byte, 0)
}

func (m *MockRows) Conn() *pgx.Conn {
	return &pgx.Conn{}
}

func (m *MockRows) Close() {
	if m.CloseFunc != nil {
		m.CloseFunc()
	}
}

func (m *MockRows) Err() error {
	if m.ErrFunc != nil {
		return m.ErrFunc()
	}
	return nil
}

func (m *MockRows) Next() bool {
	if m.NextFunc != nil {
		return m.NextFunc()
	}
	return false
}

func (m *MockRows) Scan(dest ...any) error {
	if m.ScanFunc != nil {
		return m.ScanFunc(dest...)
	}
	return nil
}

// MockTx для pgx.Tx
type MockTx struct {
	ExecFunc     func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryFunc    func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRowFunc func(ctx context.Context, sql string, args ...any) pgx.Row
	CommitFunc   func(ctx context.Context) error
	RollbackFunc func(ctx context.Context) error
}

func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	//TODO implement me
	panic("implement me")
}

func (m *MockTx) LargeObjects() pgx.LargeObjects {
	//TODO implement me
	panic("implement me")
}

func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return &pgconn.StatementDescription{}, nil
}

func (m *MockTx) Conn() *pgx.Conn {
	//TODO implement me
	panic("implement me")
}

func (m *MockTx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(ctx, sql, arguments...)
	}
	return pgconn.CommandTag{}, nil
}

func (m *MockTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, sql, args...)
	}
	return &MockRows{}, nil
}

func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(ctx, sql, args...)
	}
	return &MockRow{}
}

func (m *MockTx) Commit(ctx context.Context) error {
	if m.CommitFunc != nil {
		return m.CommitFunc(ctx)
	}
	return nil
}

func (m *MockTx) Rollback(ctx context.Context) error {
	if m.RollbackFunc != nil {
		return m.RollbackFunc(ctx)
	}
	return nil
}

// MockBatchResults для pgx.BatchResults
type MockBatchResults struct {
	ExecFunc     func() (pgconn.CommandTag, error)
	QueryFunc    func() (pgx.Rows, error)
	QueryRowFunc func() pgx.Row
	CloseFunc    func() error
}

func (m *MockBatchResults) Exec() (pgconn.CommandTag, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc()
	}
	return pgconn.CommandTag{}, nil
}

func (m *MockBatchResults) Query() (pgx.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc()
	}
	return &MockRows{}, nil
}

func (m *MockBatchResults) QueryRow() pgx.Row {
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc()
	}
	return &MockRow{}
}

func (m *MockBatchResults) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func TestDBStore_Add(t *testing.T) {
	ctx := context.Background()

	t.Run("successful add", func(t *testing.T) {
		mockDB := &MockDB{
			ExecFunc: func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
				assert.Contains(t, sql, "INSERT INTO urls")
				assert.Equal(t, "abc123", arguments[0])
				assert.Equal(t, "https://example.com", arguments[1])
				assert.Equal(t, 1, arguments[2])
				return pgconn.NewCommandTag("INSERT 1"), nil
			},
		}

		store := &DBStore{
			conf: &config.Config{},
			pool: mockDB,
		}

		url := URL{
			Hash: "abc123",
			Link: "https://example.com",
		}

		hash, err := store.Add(ctx, url, 1)
		assert.NoError(t, err)
		assert.Equal(t, "abc123", hash)
	})

	t.Run("duplicate URL returns existing hash", func(t *testing.T) {
		execCallCount := 0
		mockDB := &MockDB{
			ExecFunc: func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
				execCallCount++
				if execCallCount == 1 {
					// Симулируем ошибку уникальности
					return pgconn.CommandTag{}, &pgconn.PgError{
						Code: "23505", // pgerrcode.UniqueViolation
					}
				}
				return pgconn.NewCommandTag("INSERT 1"), nil
			},
			QueryRowFunc: func(ctx context.Context, sql string, args ...any) pgx.Row {
				assert.Contains(t, sql, "SELECT hash FROM urls WHERE original_url=")
				assert.Equal(t, "https://duplicate.com", args[0])
				return &MockRow{
					ScanFunc: func(dest ...any) error {
						if hashPtr, ok := dest[0].(*string); ok {
							*hashPtr = "existing_hash_123"
						}
						return nil
					},
				}
			},
		}

		store := &DBStore{
			conf: &config.Config{},
			pool: mockDB,
		}

		url := URL{
			Hash: "new_hash_456",
			Link: "https://duplicate.com",
		}

		hash, err := store.Add(ctx, url, 1)
		assert.True(t, errors.Is(err, ErrExistsURL))
		assert.Equal(t, "existing_hash_123", hash)
	})

	t.Run("database error", func(t *testing.T) {
		mockDB := &MockDB{
			ExecFunc: func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
				return pgconn.CommandTag{}, errors.New("database connection failed")
			},
		}

		store := &DBStore{
			conf: &config.Config{},
			pool: mockDB,
		}

		url := URL{
			Hash: "test_hash",
			Link: "https://test.com",
		}

		_, err := store.Add(ctx, url, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection failed")
	})
}
