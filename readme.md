# selectgo
This work is based off my [ColdFusion project](https://github.com/webRat/select). This assumes you know how to write proper SQL statement in the first place. This assembles a query statement to be passed into a 3rd party package like [sqlx](https://github.com/jmoiron/sqlx). Right now, this targets MySQL/MariaDB sql.

## Why would I use it?
- Conform to the same SQL standard across team
- Write faster query conditionals based on Go logic instead of creating a function for each query.

## LICENSE
Apache

## Examples
Basic Example:

```
package main

import (
    "fmt"

    . "github.com/webRat/selectgo"
)

func main() {
    sql := NewQueryStatement()
    sql.Select([]string{"*"}).
        From("TableName").
        Where("1 = 1")

    query, err := sql.Assemble()

    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(query)
}
```

Result:

```
onix:stuff webRat$ ./stuff
SELECT * FROM TableName WHERE 1 = 1 LIMIT 0
```
