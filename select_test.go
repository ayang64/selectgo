package selectgo

import "testing"

func TestAssemble(t *testing.T) {

	should := []struct {
		pass *QueryStatement
		want string
	}{
		{
			NewQueryStatement().Select([]string{"uuid()"}),
			"SELECT uuid()",
		},
		{
			NewQueryStatement().Select([]string{"*"}).From("table"),
			"SELECT * FROM table",
		},
		{
			NewQueryStatement().Select([]string{"id"}).From("table").Where("1 = 1"),
			"SELECT id FROM table WHERE 1 = 1",
		},
		{
			NewQueryStatement().Select([]string{"*"}).From("table"),
			"SELECT * FROM table",
		},

		{
			NewQueryStatement().Select([]string{"*"}).From("table").Limit(1, 50),
			"SELECT * FROM table LIMIT 50 OFFSET 1",
		},
		{
			NewQueryStatement().Select([]string{"*"}).From("table").Rowcount(50).Offset(1),
			"SELECT * FROM table LIMIT 50 OFFSET 1",
		},
		{
			NewQueryStatement().Select([]string{"*"}).From("table").Rowcount(50).Offset(0),
			"SELECT * FROM table LIMIT 50",
		},
		{
			NewQueryStatement().Select([]string{"*"}).From("table").Rowcount(0).Offset(0),
			"SELECT * FROM table",
		},
		{
			NewQueryStatement().Select([]string{"*"}).From("table").OrderBy("created ASC"),
			"SELECT * FROM table ORDER BY created ASC",
		},
		{
			NewQueryStatement().Select([]string{"*"}).From("table").OrderBy("created DESC, name ASC"),
			"SELECT * FROM table ORDER BY created DESC, name ASC",
		},
	}

	for _, s := range should {
		r, err := s.pass.Assemble()
		if r != s.want {
			t.Errorf("Failed, expected %q, got %q", s.want, r)
		}
		if err != nil {
			t.Errorf("Error should be nil, got %v", err)
		}
	}
}

func TestFailures(t *testing.T) {
	must := []struct {
		fail *QueryStatement
		want string
	}{
		{
			NewQueryStatement().Select([]string{}),
			"",
		},
		{
			NewQueryStatement().Select([]string{""}),
			"",
		},
		{
			NewQueryStatement().Select([]string{" "}).From("table"),
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
