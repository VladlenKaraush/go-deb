package services

import (
	"archive/tar"
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"pkg_registry/db"
	"pkg_registry/s3"
	"strconv"
	"strings"

	"github.com/xi2/xz"
)

const CONTROL_BYTES_LEN = 132
const CONTROL_FILE_SIZE_OFFSET = 120
const CONTROL_FILE_SIZE_END = 130

type Pkg struct {
	Name    string
	Version string
	Arch    string
}

type UploadService struct {
	db      *db.DbRepo
	s3Client *s3.S3Client
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CreateUploadService(db *db.DbRepo, s3Client *s3.S3Client) UploadService {
	return UploadService{
		db: db,
		s3Client: s3Client,
	}
}

func (us UploadService) UploadHandler(w http.ResponseWriter, r *http.Request) {

	controlBytes, controlFile := readControlFile(r.Body)
	unzippedControlFile := unzipControlFile(controlFile)
	pkg := parseControlFile(string(unzippedControlFile))
	queryParams := r.URL.Query()
	repoId, err := strconv.Atoi(queryParams["repoId"][0])
	releaseId, err := strconv.Atoi(queryParams["releaseId"][0])
	check(err)
	us.savePkg(pkg, repoId, releaseId)

	bodyRest, err := io.ReadAll(r.Body)
	body := append(controlBytes, controlFile...)
	body = append(body, bodyRest...)
	shasum := calcHash(body)
	log.Println("sha512 sum = ", shasum)
	bucket := "deb-registry-" + strconv.Itoa(repoId)
	us.s3Client.UploadPackage(body, bucket, pkg.Name + "_" + pkg.Version)
	file, err := os.Create("debfile.deb")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("writing deb package")
	file.Write(body)
	defer file.Close()
}

func (us UploadService) savePkg(pkg Pkg, repoId, releaseId int) {
	us.db.InsertPackage(pkg.Name, pkg.Version, pkg.Arch, "filepath", repoId, releaseId)
}

func (us UploadService) displayPkgs(repoId, releaseId int) {
	pkgs := us.db.GetPackages(repoId, releaseId)
	for _, pkg := range pkgs {
		log.Println("pkg ", pkg)
	}
}

func calcHash(contents []byte) string{
	hasher := sha512.New()
	_, err := hasher.Write(contents)
	check(err)
	return hex.EncodeToString(hasher.Sum(nil))
}

func readControlFile(r io.Reader) ([]byte, []byte) {
	controlBytes := make([]byte, CONTROL_BYTES_LEN)
	_, err := r.Read(controlBytes)
	check(err)
	cfSize, err := strconv.Atoi(strings.Trim(string(controlBytes[CONTROL_FILE_SIZE_OFFSET:CONTROL_FILE_SIZE_END]), " "))
	check(err)
	log.Printf("control file size %d", cfSize)
	controlFile := make([]byte, cfSize)
	_, err = r.Read(controlFile)
	check(err)
	return controlBytes, controlFile
}

func unzipControlFile(controlFile []byte) []byte {

	reader, err := xz.NewReader(bytes.NewReader(controlFile), 0)
	tarReader := tar.NewReader(reader)
	check(err)
	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
		}
		switch header.Typeflag {
		case tar.TypeDir:
			log.Printf("found dir with name %s", header.Name)
		case tar.TypeReg:
			if header.Name == "./control" || header.Name == "control" {
				curFile, err := io.ReadAll(tarReader)
				if err != nil {
					log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
				}
				return curFile
			}
		default:
			log.Fatalf("ExtractTarGz: uknown type: %b in %s", header.Typeflag, header.Name)
		}
	}
	panic("control file not found:")
}

func parseControlFile(controlFile string) Pkg {
	pkgLines := strings.Split(controlFile, "\n")
	pkgMap := make(map[string]string)
	for _, line := range pkgLines {
		args := strings.Split(line, ": ")
		if len(args) == 2 {
			pkgMap[args[0]] = strings.Trim(args[1], " ")
		}
	}
	return Pkg{
		Name:    pkgMap["Package"],
		Version: pkgMap["Version"],
		Arch:    pkgMap["Architecture"],
	}
}
