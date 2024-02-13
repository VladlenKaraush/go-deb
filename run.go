package main

import (
	"log"
	"net/http"
	"pkg_registry/db"
	pkg_s3 "pkg_registry/s3"
	"pkg_registry/services"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db := dbInit()
	s3AuthCreds := pkg_s3.S3AuthCreds{Key: "minioadmin", Secret: "minioadmin"}
	s3Client := pkg_s3.GetS3Client(s3AuthCreds)
	uploadService := services.CreateUploadService(&db, &s3Client)

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", uploadService.UploadHandler)
	if err := http.ListenAndServe("localhost:8090", mux); err != nil {
		log.Fatal(err)
	}
}

func dbInit() db.DbRepo {
	repo := db.CreateDbRepo()
	repo.CreatePackageTable()
	repo.CreateReleaseTable()
	repo.CreateRepositoryTable()
	return repo
}
