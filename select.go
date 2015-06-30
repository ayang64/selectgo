package selectgo

import (
	"bytes"
	"strconv"
)

// QueryStatement is a simple container that builds out a simple select statement that's eventually going to grow.
type QueryStatement struct {
	selectcolumns   []string
	hasSelection    bool
	fromtable       string
	hasTable        bool
	innerjoins      []string
	hasJoins        bool
	where           string
	hasWhere        bool
	andwhere        []string
	hasAndWhere     bool
	offset          int
	hasOffset       bool
	useOffset       bool
	hasRowcount     bool
	rowcount        int
	rowcountCeiling int
	useRowcount     bool
}

// NewQueryStatement returns a new instance of the select statement
func NewQueryStatement() *QueryStatement {
	return &QueryStatement{
		rowcountCeiling: 100,
		useRowcount:     true,
		useOffset:       true,
		hasSelection:    false,
		hasTable:        false,
		hasJoins:        false,
		hasWhere:        false,
		hasAndWhere:     false,
		hasOffset:       false,
		hasRowcount:     false,
	}
}

// Select sets the columns to be selected
func (q *QueryStatement) Select(columns []string) *QueryStatement {
	if len(columns) > 0 {
		q.hasSelection = true
		q.selectcolumns = columns
	}

	return q
}

// From sets the columns to be selected
func (q *QueryStatement) From(table string) *QueryStatement {
	if len(table) > 0 {
		q.hasTable = true
		q.fromtable = table
	}

	return q
}

// InnerJoin continues to add an inner join to the where statement
func (q *QueryStatement) InnerJoin(join string) *QueryStatement {
	if len(join) > 0 {
		q.hasJoins = true
		q.innerjoins = append(q.innerjoins, join)
	}

	return q
}

// Where sets the columns to be selected
func (q *QueryStatement) Where(where string) *QueryStatement {
	if len(where) > 0 {
		q.hasWhere = true
		q.where = where
	}

	return q
}

// And continues to add to the where statement
func (q *QueryStatement) And(param string) *QueryStatement {
	if len(param) > 0 {
		q.andwhere = append(q.andwhere, param)
	}

	return q
}

// Offset adds offset to the query
func (q *QueryStatement) Offset(offset int) *QueryStatement {
	if offset > 0 {
		q.hasOffset = true
		q.offset = offset
	}
	return q
}

// Rowcount adds rowcount to the query
func (q *QueryStatement) Rowcount(rowcount int) *QueryStatement {
	if rowcount > 0 {
		q.hasRowcount = true
		q.rowcount = rowcount
	}
	return q
}

// Limit page and number of things we're limiting it too. This will overwrite Offset() / Rowcount() if you're not careful
func (q *QueryStatement) Limit(offset, rowcount int) *QueryStatement {
	q.Offset(offset)
	q.Rowcount(rowcount)
	return q
}

// Assemble it all together into something that makes sense
func (q *QueryStatement) Assemble() (string, error) {
	var sql bytes.Buffer

	sql.WriteString("SELECT ")
	if q.hasSelection {
		// Assemble Column selection
		numOfColumns := len(q.selectcolumns) - 1
		for i, col := range q.selectcolumns {
			sql.WriteString(col)
			if i != numOfColumns {
				sql.WriteString(", ")
			}
		}
	}

	// Assemble FROM statement
	if q.hasTable {
		sql.WriteString(" FROM ")
		sql.WriteString(q.fromtable)
	}

	// Assemble INNER JOINS
	if q.hasJoins {
		sql.WriteString(" ")
		numnOfJoins := len(q.innerjoins) - 1
		for i, ij := range q.innerjoins {
			sql.WriteString(ij)
			if i != numnOfJoins {
				sql.WriteString(" ")
			}
		}
		//sql.WriteString(innerJoins.String())
	}

	// Assemble WHERE statement
	if q.hasWhere {
		sql.WriteString(" WHERE ")
		sql.WriteString(q.where)

		if q.hasAndWhere {
			numOfANDs := len(q.andwhere) - 1
			sql.WriteString("AND ")
			for i, w := range q.andwhere {
				sql.WriteString(w)
				if i != numOfANDs {
					sql.WriteString(" AND ")
				}
			}
		}
	}

	// Optional Stuff -- Offset, Rowcounts
	if !q.hasOffset || q.offset <= 0 {
		q.useOffset = false
	}

	// If we have a rowcount, does it conform to our ceiling?
	if q.hasRowcount && (q.rowcount <= 0 || q.rowcount > q.rowcountCeiling) {
		q.useRowcount = false
		q.rowcount = q.rowcountCeiling
	}

	if q.useRowcount {
		sql.WriteString(" LIMIT ")
		sql.WriteString(strconv.Itoa(q.rowcount))
		if q.useOffset {
			sql.WriteString(" OFFSET ")
			sql.WriteString(strconv.Itoa(q.offset))
		}
	}

	return sql.String(), nil
}
