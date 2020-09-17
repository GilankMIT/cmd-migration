
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

Navigate to folder `cmd/cli` and
run the following command to build migrate.go :

```sh
$ go build migrate/migrate.go
```

After compiling the CLI files, Run the following command in root folder

**Windows**
```sh
> ./cmd/cli/migrate.exe --up
```

**Linux/Mac**
```sh
$ ./cmd/cli/migrate --up
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

<br>
