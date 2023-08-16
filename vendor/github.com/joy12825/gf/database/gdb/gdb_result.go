// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/joy12825/gf.

package gdb

import (
	"database/sql"

	"github.com/joy12825/gf/errors/gerror"
)

// SqlResult is execution result for sql operations.
// It also supports batch operation result for rowsAffected.
type SqlResult struct {
	Result   sql.Result
	Affected int64
}

// MustGetAffected returns the affected rows count, if any error occurs, it panics.
func (r *SqlResult) MustGetAffected() int64 {
	rows, err := r.RowsAffected()
	if err != nil {
		err = gerror.Wrap(err, `sql.Result.RowsAffected failed`)
		panic(err)
	}
	return rows
}

// MustGetInsertId returns the last insert id, if any error occurs, it panics.
func (r *SqlResult) MustGetInsertId() int64 {
	id, err := r.LastInsertId()
	if err != nil {
		err = gerror.Wrap(err, `sql.Result.LastInsertId failed`)
		panic(err)
	}
	return id
}

// RowsAffected returns the number of rows affected by an
// update, insert, or delete. Not every database or database
// driver may support this.
// Also, See sql.Result.
func (r *SqlResult) RowsAffected() (int64, error) {
	if r.Result == nil {
		return 0, nil
	}
	if r.Affected > 0 {
		return r.Affected, nil
	}
	if r.Result == nil {
		return 0, nil
	}
	return r.Result.RowsAffected()
}

// LastInsertId returns the integer generated by the database
// in response to a command. Typically, this will be from an
// "auto increment" column when inserting a new row. Not all
// databases support this feature, and the syntax of such
// statements varies.
// Also, See sql.Result.
func (r *SqlResult) LastInsertId() (int64, error) {
	if r.Result == nil {
		return 0, nil
	}
	return r.Result.LastInsertId()
}
