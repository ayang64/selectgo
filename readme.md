# selectgo
This work is based off my [ColdFusion project](https://github.com/webRat/select). This assumes you know how to write proper SQL statement in the first place. This assembles a query statement to be passed into a 3rd party package like [sqlx](https://github.com/jmoiron/sqlx). Right now, this targets MySQL/MariaDB sql.

## Why would I use it?
- Conform to the same SQL standard across team
- Write faster query conditionals based on Go logic instead of creating a function for each query.

## What's missing?
- Tests. :/ I'm working on it.
- OR statements

## LICENSE
Apache
