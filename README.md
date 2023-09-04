# db_select
CLI tool for queriying Postgresql DB

## Prerequisites
Your DB should be open on localhost:5432

Provide as ENV variables your DB credentials and DB name:

export DB_NAME=db_name
export DB_USER=user
export DB_PASSWORD=db_password

## Usage

```
db_select clients                   - select * from clients limit 10
db_select clients "*" id 123        - select * from clients where id = 123 limit 10
db_select clients "id, name" id 123 - select id,name from clients where id = 123 limit 10
```