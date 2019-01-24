package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var HOSTFILEPATH = "C:\\Windows\\System32\\drivers\\etc\\hosts"
var DOMAINBLOCKFILE = "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews/hosts"

var (
	fileInfo os.FileInfo
	err      error
)

func main() {
	log.Println("Main started")

	checkFile(HOSTFILEPATH)
	createBackup(HOSTFILEPATH)
	deleteFile(HOSTFILEPATH)

	var file = openFile(HOSTFILEPATH)
	log.Printf("File '%s' successfully loaded\n", HOSTFILEPATH)

	writeBlockedDomains(file, DOMAINBLOCKFILE)

	file.Close()

	exec.Command("cmd", "/C", "ipconfig", "/flushdns")
	log.Println("DNS Cache flushed🚽")

	log.Println("End Main")
}

func checkFile(file string) {
	fileInfo, err := os.Stat(file)
	if err != nil && os.IsNotExist(err) {
		log.Fatal("File does not exist.")
	}
	log.Println("File does exist. Filesize: ", fileInfo.Size()/1000, "kB")
}

func deleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}
}

func openFile(path string) *os.File {
	var file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("File opened for read and write.")
	return file
}

func writeBlockedDomains(filepath *os.File, url string) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to Server")
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("bad status: %s", resp.Status)
	}
	log.Println("Download OK")

	// Writer the body to file
	_, err = io.Copy(filepath, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Blocked Domains written to file")
}

func createBackup(filepath string) {
	checkFile(filepath)
	var bkpFilepath = filepath + ".bkp"
	nBytes, err := copy(filepath, bkpFilepath)
	if err != nil {
		fmt.Printf("The copy operation failed %q\n", err)
		log.Fatal(err)
	} else {
		fmt.Printf("Copied %d bytes!\n", nBytes)
	}
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
