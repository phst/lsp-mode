package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/sourcegraph/jsonrpc2"
)

func main() {
	flag.Parse()
	state := loadPersistentState()
	defer state.save()
	logger.Print("starting")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	objectStream := jsonrpc2.NewBufferedStream(stdio{}, jsonrpc2.VSCodeObjectCodec{})
	defer objectStream.Close()
	h := &handler{state, cancel}
	conn := jsonrpc2.NewConn(ctx, objectStream, jsonrpc2.HandlerWithError(h.handle), jsonrpc2.LogMessages(logger))
	defer conn.Close()
	logger.Print("waiting for quit")
	<-ctx.Done()
	logger.Print("quitting")
}

type persistentState struct {
	Crashes int
}

func loadPersistentState() *persistentState {
	b, err := ioutil.ReadFile(*stateFile)
	if err != nil {
		logger.Fatalf("couldn’t load state file: %v", err)
	}
	r := new(persistentState)
	if len(b) > 0 {
		if err := json.Unmarshal(b, &r); err != nil {
			logger.Fatalf("couldn’t unmarshal state %s: %v", b, err)
		}
	}
	logger.Printf("state loaded from %s: %+v", *stateFile, r)
	return r
}

func (s *persistentState) save() {
	b, err := json.Marshal(s)
	if err != nil {
		logger.Fatalf("couldn’t marshal state %+v: %v", s, err)
	}
	if err := ioutil.WriteFile(*stateFile, b, 0644); err != nil {
		logger.Fatalf("couldn’t write state file: %v", err)
	}
	logger.Printf("state saved to %s: %+v", *stateFile, s)
}

type handler struct {
	state  *persistentState
	cancel context.CancelFunc
}

func (h *handler) handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	switch req.Method {
	case "initialize":
		if *initialCrashes > h.state.Crashes {
			h.cancel()
			h.state.Crashes++
		}
		return struct {
			Capabilities struct{} `json:"capabilities"`
		}{}, nil
	}
	return nil, fmt.Errorf("unsupported method %s", req.Method)
}

type stdio struct{}

func (stdio) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (stdio) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (stdio) Close() error {
	return nil
}

var (
	logger         = log.New(os.Stderr, "[test server] ", log.LstdFlags)
	stateFile      = flag.String("state_file", "", "file to save persistent state in")
	initialCrashes = flag.Int("initial_crashes", 0, "number of times the test server will simulate a crash initially")
)
