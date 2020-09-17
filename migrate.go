package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

/*
TODO:
- Create Postgre implementation
- Create migration file generator
*/

type MigrationConfig struct {
	//DBConfigPath is the .yml filepath for the database configuration
	//(default to db/dbconf.yml)
	DBConfigPath string

	//Env is the environment of database migration (default to development)
	Env string

	//db is the database instance
	db *sql.DB

	//Dialect is the type of Database (mysql)
	Dialect string

	//MigrationDir is the location of migration folder (default to db/migrations)
	MigrationDir string
}

func main() {

	confPath := flag.String("config-path", "db/dbconf.yml", "DB Configuration path")
	migrationEnv := flag.String("env", "development", "migration environment")
	migrationDir := flag.String("migration-dir", "db/migrations", "migration directory")

	upMigration := flag.Bool("up", false, "Up migration flag")
	downMigration := flag.Bool("down", false, "Down migration flag")
	versionMigration := flag.Bool("version", false, "Version of migration flag")

	newMigrationFile := flag.Bool("create", false, "Create new migration file")
	newMigrationFileName := flag.String("filename", "", "New migration file name")
	flag.Parse()

	if *newMigrationFile {
		if *newMigrationFileName == "" {
			fmt.Println("please specify migration file name with --filename")
			showHelp()
			return
		}

		//create new migration file
		err := createNewMigrationFile(*migrationDir, *newMigrationFileName)
		if err != nil {
			fmt.Println("failed to create migration file " + err.Error())
		}

		return
	}

	//check if at least up or down flag is specified
	if !(*upMigration || *downMigration || *versionMigration) {
		fmt.Println("please specify --up or --down for migration")
		showHelp()
		return
	}

	//check migration direction up/down
	if *upMigration && *downMigration {
		fmt.Println("use --up or --down at once only")
		showHelp()
		return
	}

	//setting db config
	migrationConf, err := NewMigrationConfig(*confPath, *migrationEnv, *migrationDir)
	if err != nil {
		log.Printf("error creating migration config : %s ", err.Error())
		return
	}
	defer migrationConf.db.Close()

	if *upMigration {
		err = migrateUp(migrationConf)
		if err != nil {
			log.Printf("error migrating Up database : %s ", err.Error())
			return
		}
	} else if *downMigration {
		err = migrateDown(migrationConf)
		if err != nil {
			log.Printf("error migrating Down database : %s ", err.Error())
			return
		}
	} else if *versionMigration {
		err = printMigrationVersion(migrationConf)
		if err != nil {
			log.Printf("error checking migration version : %s ", err.Error())
			return
		}
	}
}

func NewMigrationConfig(confPath, env, migrationDir string) (*MigrationConfig, error) {
	migrationConf := MigrationConfig{MigrationDir: migrationDir}

	yamlFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	//read YML file
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		return nil, err
	}

	//read according to env file
	dbConfig := m[env].(map[interface{}]interface{})
	dbDSN, ok := dbConfig["open"].(string)
	if !ok {
		return nil, errors.New("no database dsn specified (forgot to add 'open' ?)")
	}

	dbDriver, ok := dbConfig["driver"].(string)
	if !ok {
		return nil, errors.New("no database driver specified (forgot to add 'driver' ?)")
	}

	//open db connection based on driver
	switch dbDriver {
	case "mysql":
		migrationConf.Dialect = dbDriver
		migrationConf.db, err = sql.Open(dbDriver, dbDSN)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("error db driver is not found (currently mysql supported only)")
	}

	return &migrationConf, nil
}

//migrate up will migrate the database to the latest version
func migrateUp(config *MigrationConfig) error {
	fmt.Println("Migrating up database ...")
	driver, errDriver := mysql.WithInstance(config.db, &mysql.Config{})
	if errDriver != nil {
		return errDriver
	}

	migrateDatabase, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", config.MigrationDir),
		config.Dialect, driver)
	if err != nil {
		return err
	}

	err = migrateDatabase.Up()
	if err != nil {
		return err
	}

	fmt.Println("Migration done ...")

	//get latest version
	version, dirty, errVersion := migrateDatabase.Version()
	//ignore error in this line. Skip the version check
	if errVersion != nil {
		return errVersion
	}

	if dirty {
		fmt.Println("dirty migration. Please clean up database")
	}

	fmt.Printf("latest version is %d \n", version)
	return nil
}

//createNewMigrationFile will create new migration file to the specified page
func createNewMigrationFile(filePath, fileName string) error {

	//get current time
	currentTime := time.Now()

	currentYear := currentTime.Year()
	currentMonth := fmt.Sprintf("%02d", int(currentTime.Month()))
	currentDay := fmt.Sprintf("%02d", currentTime.Day())

	currentHour := fmt.Sprintf("%02d", currentTime.Hour())
	currentMinute := fmt.Sprintf("%02d", currentTime.Minute())
	currentSec := fmt.Sprintf("%02d", currentTime.Second())

	currentTimeFilePrefix := strconv.Itoa(currentYear) + currentMonth + currentDay +
		currentHour + currentMinute + currentSec

	fmt.Println(currentHour + currentMinute + currentSec)
	migrationUpFullPath := filePath + "/" +
		currentTimeFilePrefix + "_" + fileName + ".up.sql"
	err := ioutil.WriteFile(migrationUpFullPath, []byte(""), 0644)
	if err != nil {
		return err
	}

	migrationDownFullPath := filePath + "/" +
		currentTimeFilePrefix + "_" + fileName + ".down.sql"
	err = ioutil.WriteFile(migrationDownFullPath, []byte(""), 0644)
	if err != nil {
		return err
	}

	return nil
}

//migrate down will migrate the database to -1 of the version
func migrateDown(config *MigrationConfig) error {
	fmt.Println("Migrating down database ...")
	driver, errDriver := mysql.WithInstance(config.db, &mysql.Config{})
	if errDriver != nil {
		return errDriver
	}

	migrationDatabase, errMigrate := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", config.MigrationDir),
		config.Dialect, driver)
	if errMigrate != nil {
		return errMigrate
	}

	err := migrationDatabase.Down()
	if err != nil {
		return err
	}

	fmt.Println("Migration done ...")

	//get latest version
	version, dirty, errVersion := migrationDatabase.Version()
	if errVersion != nil {
		//ignore error in this line. Skip the version check
		return errVersion
	}

	if dirty {
		fmt.Println("dirty migration. Please clean up database")
	}

	fmt.Printf("latest version is %d \n", version)
	return nil
}

func printMigrationVersion(config *MigrationConfig) error {
	driver, errDriver := mysql.WithInstance(config.db, &mysql.Config{})
	if errDriver != nil {
		return errDriver
	}

	migrationDatabase, errMigrate := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", config.MigrationDir),
		config.Dialect, driver)
	if errMigrate != nil {
		return errMigrate
	}
	version, dirty, errVersion := migrationDatabase.Version()
	if errVersion != nil {
		return errVersion
	}
	fmt.Printf("Database migration version %d. Is dirty %v \n", version, dirty)

	return nil
}

func showHelp() {
	helpMessage := `
Database migration tool. Version 1.0

usage : 

migrate OPTIONS COMMAND

Options : 
 -config-path 		DB Configuration path in yml (default: db/dbconf.yml)
 -env 			Migration environment based on config path (default: development)
 -migration-dir		Migration directory files (default: db/migrations)

Command :
--up 		Up migration
--down		Down migration
--version   see migrations version


Creating migration file : 

--create --migration-dir [YOUR MIGRATION DIR] --filename [YOUR MIGRATION NAME]

--migration-dir 			Path of file (default: db/migrations)
--filename				Name of migration file  

`
	fmt.Print(helpMessage)
}
