mysqldiff -get diff between src and dst [badge]
-----

mysqldiff shows difference between database and sql file, sql files and databases.

[original mysqldiff](https://github.com/onishi/mysqldiff) is created by [@onishi](https://twitter.com/onishi).

mysqldiff get two <Schema> args , one is "source" and the other is "destination".
Schmea has two types:

* database
* sql.file

## Usage

```
% mysqldiff "-hlocalhost current" ./sql/newschema.sql
ALTER TABLE sushi ADD COLUMN `freshness` INT NOT NULL AFTER `price`;
CREATE TABLE yakizakana (
  id   INT UNSIGNED ...,
  name VARCHAR(255),
  ...
)
```

### Source Combination

* database - database
  * `mysqldiff "-hproduction sample" "-hlocalhost sample"`
* database - sql file
  * `mysqldiff "-hproduction sample" ./new.sql`
* sql file - sql file
  * `mysqldiff ./old.sql ./new.sql`

## Schema type detection

mysqldiff detects one from two shcema type by existance of given file.

## TODO

* write tests.
* use goroutin/channel to read schemas.
