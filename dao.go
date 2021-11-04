package mysql

import (
	"database/sql"
)

type SqlQueryParams struct {
	CondItems []*SqlColQueryItem

	OrderBy string
	Offset  int64
	Cnt     int64
}

type SqlExecResult struct {
	Err          error
	LastInsertID int64
	RowsAffected int64
}

type Dao struct {
	*Client
}

func (d *Dao) Insert(tableName string, colNames []string, colsValues ...[]interface{}) *SqlExecResult {
	sqb := new(SqlQueryBuilder)
	sqb.Insert(tableName, colNames...).
		Values(colsValues...)

	return ConvertSqlResultToSqlExecResult(d.Exec(sqb.Query(), sqb.Args()...))
}

func (d *Dao) DeleteByIDs(tableName string, ids ...int64) *SqlExecResult {
	sqb := new(SqlQueryBuilder)

	sqb.Delete(tableName)
	if len(ids) == 1 {
		sqb.WhereConditionAnd(&SqlColQueryItem{"id", SqlCondEqual, ids[0]})
	} else {
		sqb.WhereConditionAnd(&SqlColQueryItem{"id", SqlCondIn, ids})
	}

	return ConvertSqlResultToSqlExecResult(d.Exec(sqb.Query(), sqb.Args()...))
}

func (d *Dao) UpdateByIDs(tableName string, updateFields map[string]interface{}, ids ...int64) *SqlExecResult {
	sqb := new(SqlQueryBuilder)

	sqb.Update(tableName).Set(updateFields)
	if len(ids) == 1 {
		sqb.WhereConditionAnd(&SqlColQueryItem{"id", SqlCondEqual, ids[0]})
	} else {
		sqb.WhereConditionAnd(&SqlColQueryItem{"id", SqlCondIn, ids})
	}

	return ConvertSqlResultToSqlExecResult(d.Exec(sqb.Query(), sqb.Args()...))
}

func (d *Dao) SelectByID(tableName string, what string, id int64) *sql.Row {
	sqb := new(SqlQueryBuilder)
	sqb.Select(what, tableName).
		WhereConditionAnd(&SqlColQueryItem{"id", SqlCondEqual, id})

	return d.QueryRow(sqb.Query(), sqb.Args()...)
}

func (d *Dao) SimpleQueryAnd(tableName string, what string, params *SqlQueryParams) (*sql.Rows, error) {
	sqb := new(SqlQueryBuilder)
	sqb.Select(what, tableName).
		WhereConditionAnd(params.CondItems...).
		OrderBy(params.OrderBy).
		Limit(params.Offset, params.Cnt)

	return d.Query(sqb.Query(), sqb.Args()...)
}

func (d *Dao) SimpleTotalAnd(tableName string, condItems ...*SqlColQueryItem) (int64, error) {
	sqb := new(SqlQueryBuilder)
	sqb.Select("count(1)", tableName).
		WhereConditionAnd(condItems...)

	var total int64
	err := d.QueryRow(sqb.Query(), sqb.Args()...).Scan(&total)

	return total, err
}

func ConvertSqlResultToSqlExecResult(sqlResult sql.Result, err error) *SqlExecResult {
	execResult := new(SqlExecResult)
	if err != nil {
		execResult.Err = err
	} else {
		lid, err := sqlResult.LastInsertId()
		if err != nil {
			execResult.Err = err
		} else {
			execResult.LastInsertID = lid
			ra, err := sqlResult.RowsAffected()
			if err != nil {
				execResult.Err = err
			} else {
				execResult.RowsAffected = ra
			}
		}
	}

	return execResult
}
