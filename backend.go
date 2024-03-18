package knaudit_proxy

import (
	"database/sql"
	"fmt"
	"io"

	_ "github.com/sijms/go-ora/v2"
)

type SendCloser interface {
	Open() error
	Send(data string) error
	Ping() error
	Close() error
}

type OracleBackend struct {
	db            *sql.DB
	connectString string
}

func (g *OracleBackend) Open() error {
	db, err := sql.Open("oracle", g.connectString)
	if err != nil {
		return fmt.Errorf("opening oracle database connection: %w", err)
	}

	g.db = db

	return nil
}

func (g *OracleBackend) Ping() error {
	err := g.db.Ping()
	if err != nil {
		return fmt.Errorf("pinging oracle database: %w", err)
	}

	return nil
}

func (g *OracleBackend) Send(data string) error {
	_, err := g.db.Exec("begin dvh_dmo.knaudit_api.log(p_event_document => :1); end;", data)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	return nil
}

func (g *OracleBackend) Close() error {
	err := g.db.Close()
	if err != nil {
		return fmt.Errorf("closing oracle database connection: %w", err)
	}

	return nil
}

func NewOracleBackend(connectString string) *OracleBackend {
	return &OracleBackend{
		connectString: connectString,
	}
}

type WriterBackend struct {
	client io.Writer
}

func (s *WriterBackend) Open() error {
	return nil
}

func (s *WriterBackend) Ping() error {
	return nil
}

func (s *WriterBackend) Send(data string) error {
	_, err := fmt.Fprintf(s.client, "%s\n", data)
	if err != nil {
		return fmt.Errorf("writing to stream: %w", err)
	}

	return nil
}

func (s *WriterBackend) Close() error {
	return nil
}

func NewWriterBackend(client io.Writer) *WriterBackend {
	return &WriterBackend{
		client: client,
	}
}
