package tracker_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Geeker1/p2p/tracker"
)

func ParseJSON(body io.ReadCloser) (resp tracker.UpdateResponse, err error) {
	data, err := io.ReadAll(body)
	if err != nil{
		return tracker.UpdateResponse{}, err
	}

	resp = tracker.UpdateResponse{}
	err = json.Unmarshal(data, &resp)

	if err != nil {
		return tracker.UpdateResponse{}, err
	}
	return resp, nil
}

func TestSetup(t *testing.T) {
	go func ()  {
		tracker.StartTracker()
	}()
}


func TestChunkHandler(t *testing.T)  {
	t.Cleanup(func() {
		tracker.CHUNK_STORE = map[string][]string{}
	})

	t.Run("chunk not found in store", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/get-peer/34")
		if err != nil {
			t.Errorf("Error fetching peer data: %v", err)
		}

		respI, err := ParseJSON(resp.Body)
		if err != nil {
			t.Errorf("Error parsing data %v", err)
		}

		if respI.Error != "Chunk hash not found" {
			t.Errorf("Expected Error `Chunk hash not found`")
		}

	})

	t.Run("chunk found in store", func(t *testing.T) {
		hashD := "829292JHSJHDD"
		tracker.CHUNK_STORE[hashD] = []string{"localhost:8080"}
		
		resp, err := http.Get("http://localhost:8080/get-peer/"+hashD)
		if err != nil {
			t.Errorf("Error fetching peer data: %v", err)
		}

		respI, err := ParseJSON(resp.Body)
		if err != nil {
			t.Errorf("Error parsing data %v", err)
		}
		
		peer, found := respI.Data.(map[string]interface{})["peer"]
		if found != true {
			t.Errorf("Peer data not returned")
		}

		if peer != "localhost:8080" {
			t.Errorf("Peer data returned is not accurate")
		}
	})
}

func MakeRequest(body string, clientHeader string) (*tracker.UpdateResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://localhost:8080/update-chunk", strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	if clientHeader != "" {
		req.Header.Add("X-Client", clientHeader)
	}

	resp, err := client.Do(req)
	if err != nil{
		return nil, err
	}

	respI, err := ParseJSON(resp.Body)
	return &respI, err
}

func TestUpdateHandler(t *testing.T) {
	t.Cleanup(func() {
		tracker.CHUNK_STORE = map[string][]string{}
	})

	peerId := "localhost:8080"
	hashID := "90293023JSSKSJS"

	t.Run("client header not set", func(t *testing.T) {
		respI, err := MakeRequest("", "")

		if err != nil {
			t.Errorf("Error making request %v", err)
		}

		if respI.Error != "Client URL is empty." {
			t.Errorf("Expected client header not to be set")
		}
	})

	t.Run("error parsing request body to json", func(t *testing.T) {
		respI, err := MakeRequest("sdsd", peerId)

		if err != nil {
			t.Errorf("Error making request %v", err)
		}

		if respI.Error != "Parsing data failed" {
			t.Errorf("Expected error parsing request body")
		}
	})

	t.Run("No hash key passed", func(t *testing.T) {
		respI, err := MakeRequest(`{"hs": 1}`, peerId)

		if err != nil {
			t.Errorf("Error making request %v", err)
		}

		if respI.Error != "Chunk hash invalid" {
			t.Errorf("Expected error validating chunk hash")
		}
	})

	t.Run("chunk is updated with new peer", func(t *testing.T) {
		body := fmt.Sprintf(`{"hash": "%s"}`, hashID)
		respI, err := MakeRequest(body, peerId)

		if err != nil {
			t.Errorf("Error making request %v", err)
		}

		if respI.Error != "" {
			t.Errorf("Error occured while updating chunk peer")
		}

		peerList := tracker.CHUNK_STORE[hashID]

		if len(peerList) == 0 {
			t.Errorf("Peer list empty.")
		}

		if peerList[0] != peerId {
			t.Errorf("Expected peer id, found "+peerList[0])
		}
	})
}
