package tracker

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

/*
Tracker is a server with a port of it's own
It keeps track of client servers and the chunks they have, so when client servers without those chunks request,
they are directed to where they can fetch chunk.

We create a key value store, mapping client IP-port to chunk ID

So periodically, clients send request to tracker to inform them of the chunk they have and tracker stores IDs
*/

var CHUNK_STORE = map[string] []string {}

type UpdateBody struct {
	ChunkHash string `json:"hash"`
}

type UpdateResponse struct {
	Status string `json:"status"`
	Error string `json:"error"`
	Data interface{} `json:"data"`
	code int
}

func sendJSON(rw http.ResponseWriter, data UpdateResponse)  {
	rw.Header().Set("Content-Type", "application/json")

	if data.code >= 400 {
		resp, err := json.Marshal(data)
		errMsg := string(resp)
		if err != nil { errMsg = "Error parsing Json" }
		http.Error(rw, errMsg, data.code)
	} else {
		json.NewEncoder(rw).Encode(data)
	}
}

func UpdateHandler(rw http.ResponseWriter, r *http.Request)  {
	uBody := UpdateBody{}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[ERROR] Failed Parsing request body", err)
		sendJSON(rw, UpdateResponse{Status: "fail", Error: "Failed Parsing request body", code: 400})
		return
	}

	clientUrl := r.Header.Get("X-Client")
	if clientUrl == "" {
		log.Println("[ERROR] Client URL is empty.")
		sendJSON(rw, UpdateResponse{Status: "fail", Error: "Client URL is empty.", code: 400})
		return
	}

	err = json.Unmarshal(data, &uBody)
	if err != nil {
		log.Println("[ERROR] Parsing data failed")
		sendJSON(rw, UpdateResponse{Status: "fail", Error: "Parsing data failed", code: 400})
		return
	}

	if uBody.ChunkHash == "" {
		log.Println("[ERROR] Chunk hash invalid")
		sendJSON(rw, UpdateResponse{Status: "fail", Error: "Chunk hash invalid", code: 400})
		return
	}

	clientData, present := CHUNK_STORE[uBody.ChunkHash]
	if present == false {
		CHUNK_STORE[uBody.ChunkHash] = []string{clientUrl}
	} else {
		var result bool = false
		for _, x := range clientData {
			if x == clientUrl {
				result = true
				break
			}
		}
		if result == false{
			CHUNK_STORE[uBody.ChunkHash] = append(clientData, clientUrl)
		}
	}

	sendJSON(rw, UpdateResponse{ Status: "ok", })
}

func ChunkHandler(rw http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)

	peerList, present := CHUNK_STORE[vars["chunk_id"]]
	if present == false {
		log.Println("[ERROR] Chunk hash not found")
		sendJSON(rw, UpdateResponse{Status: "fail", Error: "Chunk hash not found", code: 400})
		return
	}

	rand.Seed(80)

	idx := len(peerList) - 1
	if idx != 0 { 
		idx = rand.Intn(idx)
	}
	choice := peerList[idx]

	sendJSON(rw, UpdateResponse{ Status: "ok", Data: map[string] string {"peer": choice} })
}

func StartTracker()  {
	port := *flag.String("port", "8080", "Application port")
	flag.Parse()

	r := mux.NewRouter()
    r.HandleFunc("/update-chunk", UpdateHandler).Methods("POST")
    r.HandleFunc("/get-peer/{chunk_id}", ChunkHandler)

	log.Println("Listening on port ", port)

	http.ListenAndServe("127.0.0.1:"+port, r)
}

