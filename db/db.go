package db

import (
	"database/sql"
	"log"
)

func executeStatement(db *sql.DB, sql_statement string) {
	statement, err := db.Prepare(sql_statement) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
}

func CreatePackageTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE packages (
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
	executeStatement(db, createStudentTableSQL)
	log.Println("packages table created")
}

func CreateRepositoryTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE repositories (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"folder_id" integer NOT NULL ,
		"name" TEXT NOT NULL 
	  );` // SQL Statement for Create Table

	log.Println("Create repositories table...")
	executeStatement(db, createStudentTableSQL)
	log.Println("packages repositories created")
}

func CreateReleaseTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE releases (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"repository_id" integer NOT NULL ,
		"suite" TEXT NOT NULL 
	  );` // SQL Statement for Create Table

	log.Println("Create releaseas table...")
	executeStatement(db, createStudentTableSQL)
	log.Println("packages releaseas created")
}
func CreateStudentTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE student (
		"idStudent" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"code" TEXT,
		"name" TEXT,
		"program" TEXT		
	  );` // SQL Statement for Create Table

	log.Println("Create student table...")
	executeStatement(db, createStudentTableSQL)
	log.Println("student table created")
}

func InsertPackage(db *sql.DB, name, version, arch, file_path string, repo_id, release_id int) {
	log.Println("Inserting package record ...")
	insertStudentSQL := `INSERT INTO packages(repository_id, release_id, name, version, architecture, file_path) 
		VALUES (?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(repo_id, release_id, name, version, arch, file_path)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// We are passing db reference connection from main to our method with other parameters
func InsertStudent(db *sql.DB, code string, name string, program string) {
	log.Println("Inserting student record ...")
	insertStudentSQL := `INSERT INTO student(code, name, program) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(code, name, program)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func DisplayPackages(db *sql.DB) {
	log.Println("printing packages:")
	row, err := db.Query("SELECT * FROM packages ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var repo_id int
		var release_id int
		var name string
		var version string
		var arch string
		var file_path string
		var description string
		row.Scan(&id, &repo_id, &release_id, &name, &version, &arch, &file_path, &description)
		log.Printf("Pkg -  id: %d, repo: %d, release: %d, name: %s, ver: %s, arch: %s, path: %s, desc: %s \n",
			id, repo_id, release_id, name, version, arch, file_path, description)
	}
}

func DisplayStudents(db *sql.DB) {
	row, err := db.Query("SELECT * FROM student ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var code string
		var name string
		var program string
		row.Scan(&id, &code, &name, &program)
		log.Println("Student: ", code, " ", name, " ", program)
	}
}
