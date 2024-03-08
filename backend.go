package knaudit_proxy

import (
	"database/sql/driver"
	"fmt"
	"io"

	go_ora "github.com/sijms/go-ora/v2"
)

type SendCloser interface {
	Open() error
	Send(data string) error
	Close() error
}

type OracleBackend struct {
	client *go_ora.Connection
}

func (b *OracleBackend) Open() error {
	err := b.client.Open()
	if err != nil {
		return fmt.Errorf("opening oracle database connection: %w", err)
	}

	return nil
}

func (b *OracleBackend) Send(data string) error {
	stmt := go_ora.NewStmt("begin dvh_dmo.knaudit_api.log(p_event_document => :1); end;", b.client)
	defer func() {
		_ = stmt.Close()
	}()

	rows, err := stmt.Query([]driver.Value{data})
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	_ = rows.Close()

	return nil
}

func (b *OracleBackend) Close() error {
	err := b.client.Close()
	if err != nil {
		return fmt.Errorf("closing oracle database connection: %w", err)
	}

	return nil
}

func NewOracleBackend(url string) (*OracleBackend, error) {
	client, err := go_ora.NewConnection(url)
	if err != nil {
		return nil, fmt.Errorf("creating oracle database connection: %w", err)
	}

	return &OracleBackend{
		client: client,
	}, nil
}

type WriterBackend struct {
	client io.Writer
}

func (s *WriterBackend) Open() error {
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
