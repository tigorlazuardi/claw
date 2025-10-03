package dblogger

import (
	"context"
	"database/sql/driver"
	"log/slog"
	"strconv"
	"strings"

	"github.com/networkteam/go-sqllogger"
)

type skipLogKey struct{}

func ContextWithSkipLog(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipLogKey{}, struct{}{})
}

var _ sqllogger.SQLLogger = (*DBLogger)(nil)

type DBLogger struct {
	Logger *slog.Logger
	Level  slog.Level
}

// Connect is called on DB connect with a generated connection id.
func (sl DBLogger) Connect(ctx context.Context, connID int64) {
	sl.Logger.Log(ctx, sl.Level, "DB Connect",
		slog.Int64("conn_id", connID),
	)
}

func (sl DBLogger) ConnBegin(ctx context.Context, connID int64, txID int64, opts driver.TxOptions) {
	sl.Logger.Log(ctx, sl.Level, "DB Tx Begin",
		slog.Int64("conn_id", connID),
		slog.Int64("tx_id", txID),
		slog.Int("isolation", int(opts.Isolation)),
		slog.Bool("read_only", opts.ReadOnly),
	)
}

type namedValueLogger []driver.NamedValue

func (va namedValueLogger) LogValue() slog.Value {
	attrs := make([]slog.Attr, len(va))
	for i, v := range va {
		if v.Name != "" {
			attrs[i] = slog.Any(v.Name, v.Value)
		} else {
			attrs[i] = slog.Any(strconv.Itoa(v.Ordinal), v.Value)
		}
	}
	return slog.GroupValue(attrs...)
}

type valueLogger []driver.Value

func (va valueLogger) LogValue() slog.Value {
	attrs := make([]slog.Attr, len(va))
	for i, v := range va {
		attrs[i] = slog.Any(strconv.Itoa(i+1), v)
	}
	return slog.GroupValue(attrs...)
}

func (sl DBLogger) ConnPrepare(ctx context.Context, connID int64, stmtID int64, query string) {
}

func (sl DBLogger) ConnPrepareContext(ctx context.Context, connID int64, stmtID int64, query string) {
}

func (sl DBLogger) ConnQuery(ctx context.Context, connID int64, rowsID int64, query string, args []driver.Value) {
	sl.Logger.Log(ctx, sl.Level, strings.TrimSpace(query), "args", valueLogger(args), "conn_id", connID, "rows_id", rowsID)
}

func (sl DBLogger) ConnQueryContext(ctx context.Context, connID int64, rowsID int64, query string, args []driver.NamedValue) {
	sl.Logger.Log(ctx, sl.Level, strings.TrimSpace(query), "args", namedValueLogger(args), "conn_id", connID, "rows_id", rowsID)
}

func (sl DBLogger) ConnExec(ctx context.Context, connID int64, query string, args []driver.Value) {
	sl.Logger.Log(ctx, sl.Level, strings.TrimSpace(query), "args", valueLogger(args), "conn_id", connID)
}

func (sl DBLogger) ConnExecContext(ctx context.Context, connID int64, query string, args []driver.NamedValue) {
	sl.Logger.Log(ctx, sl.Level, strings.TrimSpace(query), "args", namedValueLogger(args), "conn_id", connID)
}

func (sl DBLogger) ConnClose(ctx context.Context, connID int64) {
}

func (sl DBLogger) StmtExec(ctx context.Context, stmtID int64, query string, args []driver.Value) {
	sl.Logger.Log(ctx, sl.Level, strings.TrimSpace(query), "args", valueLogger(args), "stmt_id", stmtID)
}

// StmtExecContext is called on an exec with context on a statement with the statement id.
func (sl DBLogger) StmtExecContext(ctx context.Context, stmtID int64, query string, args []driver.NamedValue) {
	sl.Logger.Log(ctx, sl.Level, strings.TrimSpace(query), "args", namedValueLogger(args), "stmt_id", stmtID)
}

// StmtQuery is called on a query on a statement with the statement id and generated rows id.
// Note: ctx is only for sqllogger metadata since StmtQuery does not receive a context.
func (sl DBLogger) StmtQuery(ctx context.Context, stmtID int64, rowsID int64, query string, args []driver.Value) {
	sl.Logger.Log(ctx, sl.Level, strings.TrimSpace(query), "args", valueLogger(args), "stmt_id", stmtID, "rows_id", rowsID)
}

// StmtQueryContext is called on a query with context on a statement with the statement id and generated rows id.
func (sl DBLogger) StmtQueryContext(ctx context.Context, stmtID int64, rowsID int64, query string, args []driver.NamedValue) {
	if ctx.Value(skipLogKey{}) != nil {
		return
	}
	sl.Logger.Log(ctx, sl.Level, strings.TrimSpace(query), "args", namedValueLogger(args), "stmt_id", stmtID, "rows_id", rowsID)
}

// StmtClose is called on a close on a statement with the statement id.
// Note: ctx is only for sqllogger metadata since StmtClose does not receive a context.
func (sl DBLogger) StmtClose(ctx context.Context, stmtID int64) {
}

// RowsClose is called on a close on rows with the rows id.
// Note: ctx is only for sqllogger metadata since RowsClose does not receive a context.
func (sl DBLogger) RowsClose(ctx context.Context, rowsID int64) {
}

// TxCommit is called on a commit on a transaction with the transaction id.
// Note: ctx is only for sqllogger metadata since TxCommit does not receive a context.
func (sl DBLogger) TxCommit(ctx context.Context, txID int64) {
}

// TxRollback is called on a rollback on a transaction with the transaction id.
// Note: ctx is only for sqllogger metadata since TxRollback does not receive a context.
func (sl DBLogger) TxRollback(ctx context.Context, txID int64) {
	sl.Logger.Log(ctx, sl.Level, "DB Tx Rollback", "tx_id", txID)
}
