package selectgo

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errUpdateNoTable = errors.New("No table specified for update")
	errEmptySet      = errors.New("No values passed in to set")
)

type UpdateStatement struct {
	hasUpdate   bool
	updateTable string
	set         map[string]string
	hasWhere    bool
	where       string
}

func NewUpdateStatement() *UpdateStatement {
	return &UpdateStatement{
		hasUpdate:   false,
		updateTable: "",
		set:         make(map[string]string),
	}
}

func (q *UpdateStatement) Update(table string) *UpdateStatement {
	q.updateTable = table
	q.hasUpdate = true

	return q
}

func (q *UpdateStatement) Set(s map[string]string) *UpdateStatement {
	q.set = s

	return q
}

func (q *UpdateStatement) Where(f string) *UpdateStatement {
	q.where = f
	q.hasWhere = true

	return q
}

func (q *UpdateStatement) And(f string) *UpdateStatement {
	q.where = fmt.Sprintf("%v AND %v", q.where, f)

	return q
}

func (q *UpdateStatement) Or(f string) *UpdateStatement {
	q.where = fmt.Sprintf("%v OR %v", q.where, f)

	return q
}

// Assemble an UPDATE statement
func (q *UpdateStatement) Assemble() (string, error) {
	if len(q.updateTable) < 1 {
		return "", errUpdateNoTable
	}

	if len(q.set) < 1 {
		return "", errEmptySet
	}

	s := ""

	for k, v := range q.set {
		s = fmt.Sprintf("%v %v = %v, ", s, k, v)
	}

	s = strings.TrimSuffix(s, ", ")
	s = strings.TrimPrefix(s, " ")

	sql := fmt.Sprintf("UPDATE %v SET %v", q.updateTable, s)

	if q.hasWhere {
		sql += fmt.Sprintf(" WHERE %v", q.where)
	}

	return sql, nil
}
