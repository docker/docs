go-testdb
=========

Framework for stubbing responses from go's driver.Driver interface.

This can be used to sit in place of your sql.Db so that you can stub responses for sql calls, and remove database dependencies for your test suite.

This project is in its infancy, but has worked well for all the use cases i've had so far, and continues to evolve as new scenarios are uncovered. Please feel free to send pull-requests, or submit feature requests if you have scenarios that are not accounted for yet.

## Setup
The only thing needed for setup is to include the go-testdb package, then you can create your db connection specifying "testdb" as your driver.
<pre>
import (
	"database/sql"
	_"github.com/erikstmartin/go-testdb"
)

db, _ := sql.Open("testdb", "")
</pre>

## Stubbing connection failure
You're able to set your own function to execute when the sql library calls sql.Open
<pre>
testdb.SetOpenFunc(func(dsn string) (driver.Conn, error) {
	return c, errors.New("failed to connect")
})
</pre>

## Stubbing queries
You're able to stub responses to known queries, unknown queries will trigger log errors so that you can see that queries were executed that were not stubbed.

Differences in whitespace, and case are ignored.

For convenience a method has been created for you to take a CSV string and turn it into a database result object (RowsFromCSVString).

<pre>
db, _ := sql.Open("testdb", "")

sql := "select id, name, age from users"
columns := []string{"id", "name", "age", "created"}
result := `
1,tim,20,2012-10-01 01:00:01
2,joe,25,2012-10-02 02:00:02
3,bob,30,2012-10-03 03:00:03
`
testdb.StubQuery(sql, testdb.RowsFromCSVString(columns, result))

res, err := db.Query(sql)
</pre>

If for some reason you need to specify another rune to split the columns, you can do it passing the rune that you want to use as `Comma` character as third argument to RowsFromCSVString

<pre>
db, _ := sql.Open("testdb", "")

sql := "select id, name, age, data from users"
columns := []string{"id", "name", "age", "data", "created"}
result := `
1|tim|20|part_1,part_2,part_3|2014-10-16 15:01:00
2|joe|25|part_4,part_5,part_6|2014-10-17 15:01:01
3|bob|30|part_7,part_8,part_9|2014-10-18 15:01:02
`
testdb.StunQuery(sql, RowsFromCSVString(columns, result, '|'))

res, err := db.Query(sql)
</pre>

## Stubbing Query function
Some times you need more control over Query being run, maybe you need to assert whether or not a particular query is run.

You can return either a driver.Rows for response (your own or utilize RowsFromCSV) or an error to be returned
<pre>
testdb.SetQueryFunc(func(query string) (result driver.Rows, err error) {
	columns := []string{"id", "name", "age", "created"}
	rows := `
1,tim,20,2012-10-01 01:00:01
2,joe,25,2012-10-02 02:00:02
3,bob,30,2012-10-03 03:00:03`

	// inspect query to ensure it matches a pattern, or anything else you want to do first
	return RowsFromCSVString(columns, rows), nil
})

db, _ := sql.Open("testdb", "")

res, err := db.Query("SELECT foo FROM bar")
</pre>

## Stubbing Parameterized Query
Sometimes you need control over the results of a parameterized query.

<pre>
testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (result driver.Rows, err error) {
	columns := []string{"id", "name", "age", "created"}

	rows := ""
	if args[0] == "joe" {
		rows = "2,joe,25,2012-10-02 02:00:02"
	}
	return testdb.RowsFromCSVString(columns, rows), nil
})

db, _ := sql.Open("testdb", "")

res, _ := db.Query("SELECT foo FROM bar WHERE name = $1", "joe")
</pre>

## Stubbing errors returned from queries
In case you need to stub errors returned from queries to ensure your code handles them properly

<pre>
db, _ := sql.Open("testdb", "")

sql := "select count(*) from error"
testdb.StubQueryError(sql, errors.New("test error"))

res, err := db.Query(sql)
</pre>

## Stubbing Parameterized Exec query
Sometimes you need control over the handling of a parameterized query that does not return any rows.

<pre>
type testResult struct{
	lastId int64
	affectedRows int64
}

func (r testResult) LastInsertId() (int64, error){
	return r.lastId, nil
}

func (r testResult) RowsAffected() (int64, error) {
	return r.affectedRows, nil
}
testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
	if args[0] == "joe" {
		return testResult{1, 1}, nil
	}
	return testResult{1, 0}, nil
})

db, _ := sql.Open("testdb", "")

res, _ := db.Exec("UPDATE bar SET name = 'foo' WHERE name = ?", "joe")
</pre>

## Stubbing Prepared Statements
You can use the same methods as `SetQueryFunc`, `SetQueryWithArgsFunc` for Prepared Statements

<pre>
testdb.SetQueryFunc(func(query string) (result driver.Rows, err error) {
	columns := []string{"id", "name", "age", "created"}
	rows := `
1,tim,20,2012-10-01 01:00:01
2,joe,25,2012-10-02 02:00:02
3,bob,30,2012-10-03 03:00:03`

	// inspect query to ensure it matches a pattern, or anything else you want to do first
	return RowsFromCSVString(columns, rows), nil
})

db, _ := sql.Open("testdb", "")

stmt, _ := db.Prepare("SELECT foo FROM bar")
res, err := stmt.Query("SELECT foo FROM bar")
</pre>

## Reset
At any point in your test, or as a defer you can remove all stubbed queries, errors, custom set Query or Open functions by calling the reset method.

<pre>
func TestMyDatabase(t *testing.T){
	defer testdb.Reset()
}
</pre>

#### TODO
Feel free to contribute and send pull requests
- Transactions

#### License
Copyright (c) 2013, Erik St. Martin
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
