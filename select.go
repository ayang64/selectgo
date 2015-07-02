package selectgo

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

var (
	errSelectNoSelect       = errors.New("No Select statement prepared.")
	errSelectContainsBlanks = errors.New("Select statement contains blanks")
)

const (
	andState = iota
	orState
)

// QueryStatement is a simple container that builds out a simple select statement that's eventually going to grow.
type QueryStatement struct {
	selectcolumns       []string
	hasSelection        bool
	fromtable           string
	hasTable            bool
	innerjoins          []string
	hasJoins            bool
	where               string
	hasWhere            bool
	conditionalWhere    []wheretype
	hasConditionalWhere bool
	offset              int
	hasOffset           bool
	useOffset           bool
	hasRowcount         bool
	rowcount            int
	useRowcount         bool
}

// NewQueryStatement returns a new instance of the select statement
func NewQueryStatement() *QueryStatement {
	return &QueryStatement{
		useRowcount:         false,
		useOffset:           false,
		hasSelection:        false,
		hasTable:            false,
		hasJoins:            false,
		hasWhere:            false,
		hasConditionalWhere: false,
		hasOffset:           false,
		hasRowcount:         false,
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
		q.hasConditionalWhere = true
		q.addWhereConditional(andState, param)
	}

	return q
}

// Or continues to or conditional to the where statement
func (q *QueryStatement) Or(param string) *QueryStatement {
	if len(param) > 0 {
		q.hasConditionalWhere = true
		q.addWhereConditional(orState, param)
	}

	return q
}

func (q *QueryStatement) addWhereConditional(ctype int, value string) {
	q.conditionalWhere = append(q.conditionalWhere, wheretype{ctype: ctype, value: value})
}

type wheretype struct {
	ctype int
	value string
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
	if !q.hasSelection {
		return "", errSelectNoSelect
	}

	sql.WriteString("SELECT ")

	// Assemble Column selection
	numOfColumns := len(q.selectcolumns) - 1
	for i, col := range q.selectcolumns {
		if len(strings.TrimSpace(col)) > 0 {
			sql.WriteString(strings.TrimSpace(col))
			if i != numOfColumns {
				sql.WriteString(", ")
			}
		} else {
			return "", errSelectContainsBlanks
		}
	}

	// Assemble FROM statement
	if q.hasTable {
		sql.WriteString(" FROM ")
		sql.WriteString(q.fromtable)
	}

	// Assemble INNER JOINS
	if q.hasJoins {
		sql.WriteString(" INNER JOIN ")
		numnOfJoins := len(q.innerjoins) - 1
		for i, ij := range q.innerjoins {
			sql.WriteString(ij)
			if i != numnOfJoins {
				sql.WriteString(" INNER JOIN ")
			}
		}
	}

	// Assemble WHERE statement
	if q.hasWhere {
		sql.WriteString(" WHERE ")
		sql.WriteString(q.where)

		if q.hasConditionalWhere {
			for _, cw := range q.conditionalWhere {
				switch cw.ctype {
				case andState:
					sql.WriteString(" AND ")

				case orState:
					sql.WriteString(" OR ")
				}
				sql.WriteString(cw.value)
			}
		}
	}

	// Optional Stuff -- Offset, Rowcounts
	if !q.hasOffset || q.offset <= 0 {
		q.useOffset = false
	}

	// If we have a rowcount, does it conform to our ceiling?
	if q.hasRowcount && (q.rowcount < 1) {
		q.useRowcount = false
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
