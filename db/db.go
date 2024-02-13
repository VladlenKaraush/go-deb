package db

import (
	"database/sql"
	"log"
)

type DbRepo struct {
	dbConn *sql.DB
}

type PkgDao struct {
	id int
	repository_id int
	release_id int
	name string
	version string
	architecture string
	file_path string
	description string
}


func CreateDbRepo() DbRepo {
	db_conn, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	repo := DbRepo { dbConn: db_conn }
	return repo
}

func (db *DbRepo)executeStatement(sql_statement string) {
	statement, err := db.dbConn.Prepare(sql_statement) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() 
}

func (db *DbRepo) CreatePackageTable() {
	createStudentTableSQL := `CREATE TABLE IF NOT EXISTS packages (
		"idPackage" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"repository_id" integer NOT NULL,
		"release_id" integer NOT NULL,		
		"name" TEXT NOT NULL ,
		"version" TEXT NOT NULL,
		"architecture" TEXT NOT NULL,
		"file_path" TEXT NOT NULL,
		"description" TEXT
	  );` // SQL Statement for Create Table

	log.Println("Create packages table...")
	db.executeStatement(createStudentTableSQL)
	log.Println("packages table created")
}

func (db *DbRepo) CreateRepositoryTable() {
	createStudentTableSQL := `CREATE TABLE IF NOT EXISTS repositories (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"folder_id" integer NOT NULL ,
		"name" TEXT NOT NULL 
	  );` // SQL Statement for Create Table

	log.Println("Create repositories table...")
	db.executeStatement(createStudentTableSQL)
	log.Println("packages repositories created")
}

func (db *DbRepo) CreateReleaseTable() {
	createStudentTableSQL := `CREATE TABLE IF NOT EXISTS releases (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"repository_id" integer NOT NULL ,
		"suite" TEXT NOT NULL 
	  );` // SQL Statement for Create Table

	log.Println("Create releaseas table...")
	db.executeStatement(createStudentTableSQL)
	log.Println("packages releaseas created")
}

func (db *DbRepo) InsertPackage(name, version, arch, file_path string, repo_id, release_id int) {
	log.Println("Inserting package record ...")
	insertStudentSQL := `INSERT INTO packages(repository_id, release_id, name, version, architecture, file_path) 
		VALUES (?, ?, ?, ?, ?, ?)`
	statement, err := db.dbConn.Prepare(insertStudentSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(repo_id, release_id, name, version, arch, file_path)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func (db *DbRepo) GetPackages(repoId, releaseId int) []PkgDao {
	row, err := db.dbConn.Query("SELECT * FROM packages WHERE repository_id = ? and release_id = ?", repoId, releaseId)
	if err != nil {
		log.Fatal(err)
	}
	var pkgs []PkgDao
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		pkg := PkgDao {}
		row.Scan(&pkg.id, &pkg.repository_id, &pkg.release_id, &pkg.name, 
			&pkg.version, &pkg.architecture, &pkg.file_path, &pkg.description)
		log.Println("pkg = ", pkg)
		pkgs = append(pkgs, pkg)
	}
	return pkgs
}
