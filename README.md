
#### Migration helper in Go based on golang-migrate/migrate

Set the configuration of database in `db/dbconf.yml`
and remember to add parameter `multiStatements=true` for MySQL DSN

Example of `db/dbconf.yml`
```sh
development:
  driver: mysql
  open: USERNAME:PASSWORD@/nobu_eform?parseTime=true&multiStatements=true
```
**important** : For MySQL, add `multiStatements=true` to support multi-statement SQL Execution

Navigate to folder `cmd/cli/migrate` and
run the following command to build migrate.go :

```sh
$ go build migrate.go
```

After compiling the CLI files, Run the following command in root folder

**Windows**
```sh
> ./cmd/cli/migrate/migrate.exe --up
```

**Linux/Mac**
```sh
$ ./cmd/cli/migrate/migrate --up
```


Parameter / Options:

| flag | description | default|
| ------ | ------ | ------ | 
| config-path |  DB Configuration path |  `db/dbconf.yml`|
| env |  migration env |  `development`|
| migration-dir | Migration files directory |  `db/migrations`|


Command:

| flag | description |
| ------ | ------ | 
| up| execute database migration to latest version|
| down| execute down migration|
| version| see current migration version|
| create| Create new migration file|


Options to create new migration file

| flag | description |
| ------ | ------ | 
| migration-dir| directory of the migration directory|
| filename| filename of the migration file|

example of creating new migration file


**Windows**

```sh
> ./cmd/cli/migrate/migrate.exe --create --migration-dir db/migrations --filename add_user_table
```

**Linux/Mac**

```sh
$ ./cmd/cli/migrate/migrate --create --migration-dir db/migrations --filename add_user_table
```

<br>
