package selectgo

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestUpdateAssemble(t *testing.T) {
	updateFields := map[string]string{
		"int":      "9001",
		"string":   `"THIS IS A STRING WOOT"`,
		"bit":      "0",
		"datetime": fmt.Sprintf("%q", time.Now()), // escape the quote
		"float":    "2.71",
	}

	should := struct {
		pass *UpdateStatement
		want string
	}{
		NewUpdateStatement().Update("user").Set(updateFields),
		fmt.Sprintf(`UPDATE user SET int = 1, string = "THIS IS A STRING WOOT", bool = true, datetime = %q, float = 2.71`, updateFields["datetime"]),
	}

	r, err := should.pass.Assemble()
	if err != nil {
		t.Errorf("Error should be nil\n got %v", err)
	}

	var test string
	for k, v := range updateFields {
		test = fmt.Sprintf("%v = %v", k, v)
		if !strings.Contains(r, test) {
			t.Errorf("SET should have %v", test)
		}
	}
}

func TestUpdateWhere(t *testing.T) {
	updateFields := map[string]string{
		"a": "1",
	}

	should := []struct {
		pass *UpdateStatement
		want string
	}{
		{
			NewUpdateStatement().Update("user").Set(updateFields).Where("a = 2"),
			"UPDATE user SET a = 1 WHERE a = 2",
		},
		{
			NewUpdateStatement().Update("user").Set(updateFields).Where("a = 2").And("b = 3"),
			"UPDATE user SET a = 1 WHERE a = 2 AND b = 3",
		},
		{
			NewUpdateStatement().Update("user").Set(updateFields).Where("a = 2").Or("b = 3"),
			"UPDATE user SET a = 1 WHERE a = 2 OR b = 3",
		},
	}

	for _, s := range should {
		r, err := s.pass.Assemble()
		if r != s.want {
			t.Errorf("Failed\n expected %q\n got %q", s.want, r)
		}
		if err != nil {
			t.Errorf("Error should be nil\n got %v", err)
		}
	}
}

func TestUpdateFailures(t *testing.T) {
	updateFields := map[string]string{
		"a": "1",
	}

	must := []struct {
		fail *UpdateStatement
		want string
	}{
		{
			NewUpdateStatement().Update(""),
			"",
		},
		{
			NewUpdateStatement().Update("").Set(updateFields),
			"",
		},
		{
			NewUpdateStatement().Update("user").Set(make(map[string]string)),
			"",
		},
	}

	for _, s := range must {
		r, err := s.fail.Assemble()
		if err == nil {
			t.Errorf("Error shouldn't be nil, failed on %q", r)
		}
	}
}
