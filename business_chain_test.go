package scriptstudio_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestBusinessChain38(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := listener.Addr().String()
	listener.Close()
	cmd := exec.Command("go", "run", "./script-backend/cmd/server", "-address", address, "-database", filepath.Join(t.TempDir(), "chain.db"))
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cmd.Process.Kill(); _, _ = cmd.Process.Wait() }()
	base := "http://" + address
	for attempt := 0; attempt < 100; attempt++ {
		response, requestErr := http.Get(base + "/api/scripts")
		if requestErr == nil {
			response.Body.Close()
			break
		}
		if attempt == 99 {
			t.Fatal(requestErr)
		}
		time.Sleep(20 * time.Millisecond)
	}
	payload := []byte(`{"requestKey":"author-submit-38","title":"The Last Station","logline":"A conductor chooses between duty and family.","genre":"drama"}`)
	ids := make([]string, 0, 2)
	for attempt := 0; attempt < 2; attempt++ {
		response, err := http.Post(base+"/api/scripts", "application/json", bytes.NewReader(payload))
		if err != nil {
			t.Fatal(err)
		}
		var item struct {
			ID string `json:"id"`
		}
		if err := json.NewDecoder(response.Body).Decode(&item); err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
		ids = append(ids, item.ID)
	}
	response, err := http.Get(base + "/api/scripts")
	if err != nil {
		t.Fatal(err)
	}
	var items []struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(response.Body).Decode(&items); err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if ids[0] != ids[1] {
		t.Fatalf("repeated request created %s and %s", ids[0], ids[1])
	}
	if len(items) != 1 {
		t.Fatal(fmt.Sprintf("repeated request produced %d scripts", len(items)))
	}
}
