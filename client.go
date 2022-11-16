package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

type Config struct {
	Project  string `json:"project"`
	Instance string `json:"instance"`
	Database string `json:"database"`
	Table    string `json:"table"`
}

type Name struct {
	Id        uuid.UUID `spanner:"id" json:"id"`
	FirstName string    `spanner:"first_name" json:"first_name"`
}

type Client struct {
	ctx        context.Context
	cfg        Config
	spannerURL string
	admin      *database.DatabaseAdminClient
	client     *spanner.Client
}

func (c *Client) Init(ctx context.Context, cfg *Config) error {
	c.ctx = ctx
	c.spannerURL = fmt.Sprintf(
		"projects/%s/instances/%s/databases/%s",
		c.cfg.Project,
		c.cfg.Instance,
		c.cfg.Database,
	)
	admin, err := database.NewDatabaseAdminClient(c.ctx)
	if err != nil {
		return err
	}
	c.admin = admin

	client, err := spanner.NewClient(c.ctx, c.spannerURL)
	if err != nil {
		return err
	}
	c.client = client

	return nil
}

func (c *Client) Teardown() {
	c.client.Close()
	c.admin.Close()
}

func (c *Client) AddNames() error {
	m := []*spanner.Mutation{
		spanner.Insert(
			c.cfg.Table,
			[]string{"id", "first_name"},
			[]interface{}{uuid.New(), "Alice"},
		),
		spanner.Insert(
			c.cfg.Table,
			[]string{"id", "first_name"},
			[]interface{}{uuid.New(), "Bob"},
		),
		spanner.Insert(
			c.cfg.Table,
			[]string{"id", "first_name"},
			[]interface{}{uuid.New(), "John"},
		),
	}
	_, err := c.client.Apply(c.ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetNames() ([]Name, error) {
	tx := c.client.ReadOnlyTransaction()
	defer tx.Close()

	iter := tx.Query(
		c.ctx,
		spanner.NewStatement(
			fmt.Sprintf(
				`SELECT * FROM %s`,
				c.cfg.Table,
			),
		),
	)
	defer iter.Stop()

	i := 0
	var names []Name
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, err
		}

		var ptr Name
		err = row.ToStruct(&ptr)
		if err != nil {
			return nil, err
		}

		names = append(names, ptr)

		// the i is in order not to pull too much data
		i += 1
		if i > 100 {
			break
		}

	}
	return names, nil
}

func (c *Client) CreateTable() error {
	createStatement := fmt.Sprintf(
		`CREATE TABLE %s (
                    id UUID,
                    first_name STRING(100)
                 ) PRIMARY KEY (id)`,
		c.cfg.Table,
	)
	op, err := c.admin.UpdateDatabaseDdl(c.ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database:   c.spannerURL,
		Statements: []string{createStatement},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(c.ctx); err != nil {
		log.Println("CreateTable: Table already exists")
	}
	return nil
}
