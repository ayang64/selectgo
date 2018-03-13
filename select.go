package selectgo

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

var (
	errSelectNoSelect       = errors.New("no select statement prepared")
	errSelectContainsBlanks = errors.New("select statement contains blanks")
)

const (
	andState = iota
	orState
	innerJoinState
	leftJoinState
)

// QueryStatement is a simple container that builds out a simple select statement that's eventually going to grow.
type QueryStatement struct {
	selectcolumns       []string
	hasSelection        bool
	fromtable           string
	hasTable            bool
	conditionalJoins    []jointype
	hasJoins            bool
	where               string
	hasWhere            bool
	hasOrderBy          bool
	orderby             string
	conditionalWhere    []wheretype
	hasConditionalWhere bool
	offset              int
	hasOffset           bool
	useOffset           bool
	hasRowcount         bool
	rowcount            int
	useRowcount         bool
	hasGroupBy          bool
	groupBy             string
}

// NewQueryStatement returns a new instance of the select statement
func NewQueryStatement() *QueryStatement {
	return &QueryStatement{
		hasSelection:        false,
		hasTable:            false,
		hasJoins:            false,
		hasWhere:            false,
		hasConditionalWhere: false,
		hasOffset:           false,
		hasRowcount:         false,
		hasOrderBy:          false,
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

// InnerJoin continues to add an inner join to the sql statement
func (q *QueryStatement) InnerJoin(join string) *QueryStatement {
	if len(join) > 0 {
		q.hasJoins = true
		q.addJoin(innerJoinState, join)
	}

	return q
}

// LeftJoin continues to add an left join to the sql statement
func (q *QueryStatement) LeftJoin(join string) *QueryStatement {
	if len(join) > 0 {
		q.hasJoins = true
		q.addJoin(leftJoinState, join)
	}

	return q
}

// InnerJoin continues to add an inner join to the where statement
func (q *QueryStatement) addJoin(jtype int, value string) {
	q.conditionalJoins = append(q.conditionalJoins, jointype{jtype: jtype, value: value})
}

// Where sets the columns to be selected
func (q *QueryStatement) Where(where string) *QueryStatement {
	if len(where) > 0 {
		q.hasWhere = true
		q.where = where
	}

	return q
}

// GroupBy will append a group by statement to the query
func (q *QueryStatement) GroupBy(groupBy string) *QueryStatement {
	if len(groupBy) > 0 {
		q.groupBy = groupBy
		q.hasGroupBy = true
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

type jointype struct {
	jtype int
	value string
}

// OrderBy adds "ORDER BY" to the query
func (q *QueryStatement) OrderBy(orderby string) *QueryStatement {
	if len(orderby) > 0 {
		q.hasOrderBy = true
		q.orderby = orderby
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
		for _, cj := range q.conditionalJoins {
			switch cj.jtype {
			case innerJoinState:
				sql.WriteString(" INNER ")

			case leftJoinState:
				sql.WriteString(" LEFT ")
			}
			sql.WriteString("JOIN ")
			sql.WriteString(cj.value)
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

	if q.hasGroupBy {
		sql.WriteString(" GROUP BY ")
		sql.WriteString(q.groupBy)
	}

	if q.hasOrderBy {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(q.orderby)
	}

	if q.hasRowcount {
		sql.WriteString(" LIMIT ")
		sql.WriteString(strconv.Itoa(q.rowcount))
		if q.hasOffset {
			sql.WriteString(" OFFSET ")
			sql.WriteString(strconv.Itoa(q.offset))
		}
	}

	return sql.String(), nil
}
