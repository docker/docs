package testdb

type Result struct {
	lastInsertId      int64
	lastInsertIdError error
	rowsAffected      int64
	rowsAffectedError error
}

func NewResult(lastId int64, lastIdError error, rowsAffected int64, rowsAffectedError error) (res *Result) {
	return &Result{
		lastInsertId:      lastId,
		lastInsertIdError: lastIdError,
		rowsAffected:      rowsAffected,
		rowsAffectedError: rowsAffectedError,
	}
}

func (res *Result) LastInsertId() (int64, error) {
	return res.lastInsertId, res.lastInsertIdError
}

func (res *Result) RowsAffected() (int64, error) {
	return res.rowsAffected, res.rowsAffectedError
}
