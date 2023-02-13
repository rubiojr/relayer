package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fiatjaf/relayer"
	"github.com/fiatjaf/relayer/storage/sqlite"
	"github.com/kelseyhightower/envconfig"
	"github.com/nbd-wtf/go-nostr"
)

type Relay struct {
	SqliteDatabase string `envconfig:"SQLITE_DATABASE"`

	storage *sqlite.SqliteBackend
}

func (r *Relay) Name() string {
	return "BasicRelay"
}

func (r *Relay) Storage() relayer.Storage {
	return r.storage
}

func (r *Relay) OnInitialized(*relayer.Server) {}

func (r *Relay) Init() error {
	err := envconfig.Process("", r)
	if err != nil {
		return fmt.Errorf("couldn't process envconfig: %w", err)
	}

	// every hour, delete all very old events
	go func() {
		db := r.Storage().(*sqlite.SqliteBackend)

		for {
			time.Sleep(60 * time.Minute)
			db.DB.Exec(`DELETE FROM event WHERE created_at < $1`, time.Now().AddDate(0, -3, 0).Unix()) // 3 months
		}
	}()

	return nil
}

func (r *Relay) AcceptEvent(evt *nostr.Event) bool {
	// block events that are too large
	jsonb, _ := json.Marshal(evt)
	if len(jsonb) > 10000 {
		return false
	}

	return true
}

func main() {
	r := Relay{}
	if err := envconfig.Process("", &r); err != nil {
		log.Fatalf("failed to read from env: %v", err)
		return
	}

	if r.SqliteDatabase == "" {
		log.Print("Using :memory: SQLite")
		r.SqliteDatabase = ":memory:"
	}
	r.storage = &sqlite.SqliteBackend{DatabaseURL: r.SqliteDatabase}
	if err := relayer.Start(&r); err != nil {
		log.Fatalf("server terminated: %v", err)
	}
}
