package server

import (
	// "encoding/json"
	// "flag"
	// "io"
	"crypto/sha256"
	"encoding/hex"
	"path"

	// "fmt"
	"io"
	"log"
	"os"
	// "math/rand"
	// "net/http"
	// "github.com/gorilla/mux"
)

/*
Server gets filename, splits file into sizeable chunks and stores in /tmp directory
Server also stores a local db of chunk list and the order to rearrange them. It sends this record to client upon request.
Server loops and waits for clients to be available, then starts sending chunks to them.
Server keeps track of client it has sent chunk to, so as not to duplicate requests.
*/

// Get file record
// Send chunk record

func StartServer(filename string, chunkSize int)  {
	log.Println("Server")
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileSHA := sha256.New()
	fileSHA.Write([]byte(filename))

	fileHash := hex.EncodeToString(fileSHA.Sum(nil))

	// dirPath := path.Join("/tmp", fileHash)
	dirPath := path.Join(fileHash)

	err = os.Mkdir(dirPath, 0750)
	if err != nil {
		log.Fatalf("Error making new directory ==> %v", err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Error querying file info %v", err)
	}

	fileSize := fileInfo.Size()
	if chunkSize > int(fileSize) {
		log.Fatal("Expected chunksize to be less than filesize")
	}
	
	buffer := make([]byte, chunkSize)

	chunkList := []string{}

	for {
		_, err := file.Read(buffer)

		first := sha256.New()
		first.Write(buffer)

		checkSum := hex.EncodeToString(first.Sum(nil))

		if err != nil {
			if err != io.EOF {
				log.Fatalf("Error reading file into buffer %v", err)
			}
			break
		}

		chunkFile, err := os.Create(path.Join(dirPath, checkSum))
		if err != nil {
			log.Fatalf("Error creating chunkfile %v", err)
		}

		defer chunkFile.Close()

		chunkFile.Write(buffer)
		chunkFile.Close()

		chunkList = append(chunkList, checkSum)
	}
}
