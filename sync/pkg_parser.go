package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

type packageIndex struct {
	contents []byte
	path     string
}

type pkgUrl struct {
	fullUrl string
	path    string
}

type release struct {
	index string
	url   string
}

type pkg struct {
	name     string
	arch     string
	version  string
	filepath string
	md5sum   string
	sha1     string
	sha256   string
}

func main() {
	start := time.Now()
	releaseUrl := "http://archive.ubuntu.com/ubuntu/dists/bionic/Release"
	body, err := readUrlBody(releaseUrl)
	if err != nil {
		panic(err)
	}

	pkgUrls := release{
		index: string(body),
		url:   releaseUrl,
	}.parseReleaseIndex()

	var pkgs [][]pkg
	for _, url := range pkgUrls {
		body, err := readUrlBody(url.fullUrl)
		if err != nil {
			log.Warn(err)
			continue
		}
		unzippedBody, err := unzip(body)
		if err != nil {
			log.Warnf("failed to unzip body from %s", url.fullUrl)
			continue
		}
		pkgList := packageIndex{
			contents: unzippedBody,
			path:     strings.Replace(url.path, ".gz", ".txt", 1),
		}.parse()
		pkgs = append(pkgs, pkgList)
	}
	for _, pkgList := range pkgs {
		for _, pkg := range pkgList {
			fmt.Println(pkg)
		}
	}
	elapsed := time.Since(start).Milliseconds()
	fmt.Printf("time passed in millis = %d\n", elapsed)

}

func (release release) parseReleaseIndex() []pkgUrl {
	lines := strings.Split(release.index, "\n")
	var pkgUrls []pkgUrl
	for _, line := range lines {
		if strings.HasSuffix(line, "Packages.gz") {
			pkgParts := strings.Split(line, " ")
			pkgPathSuffix := pkgParts[len(pkgParts)-1]
			pkgFullPath := strings.Replace(release.url, "Release", pkgPathSuffix, 1)
			pkgUrl := pkgUrl{
				fullUrl: pkgFullPath,
				path:    pkgPathSuffix,
			}
			pkgUrls = append(pkgUrls, pkgUrl)
		}
	}
	return pkgUrls
}

func unzip(content []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	unzippedContent, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return unzippedContent, nil
}

func readUrlBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("couldn't download body, url = %s, status = %d", url, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body, nil
}

func (pkg packageIndex) writeToFile(folder string) {
	filename := strings.Replace(pkg.path, "/", "-", len(pkg.path))
	log.Infof("writing pkg %s to file %s\n", pkg.path, filename)
	err := os.WriteFile(folder+filename, pkg.contents, 0644)
	if err != nil {
		panic(err)
	}
}

func (index packageIndex) parse() []pkg {
	var packages []pkg
	packageStrings := strings.Split(string(index.contents), "\n\n")
	for _, pkgString := range packageStrings {
		pkgLines := strings.Split(pkgString, "\n")
		pkgMap := make(map[string]string)
		for _, line := range pkgLines {
			args := strings.Split(line, ": ")
			if len(args) < 2 {
				log.Warn(args)
			} else {
				pkgMap[args[0]] = args[1]
			}
		}
		pkg := pkg{
			name:     pkgMap["Package"],
			arch:     pkgMap["Architecture"],
			version:  pkgMap["Version"],
			filepath: pkgMap["Filename"],
			md5sum:   pkgMap["MD5sum"],
			sha1:     pkgMap["SHA1"],
			sha256:   pkgMap["SHA256"],
		}
		packages = append(packages, pkg)
	}

	return packages
}
