package spoon

import (
	"github.com/pkg/errors"
)

var (
	errIgnoreField = errors.New("error ignore this field")
)

const (
	defaultTagPrefix = "db"
	defaultIgnoreTag = "-"
)

type optionParam struct {
	tagPrefix string
	ignoreTag string
}

// Client is Google Cloud Spanner schema generator
type Client struct {
	param  *optionParam
	parser *parser
}

// New creates a Client and returns it.
func New(opts ...Option) (*Client, error) {
	op := &optionParam{
		tagPrefix: defaultTagPrefix,
		ignoreTag: defaultIgnoreTag,
	}

	for _, opt := range opts {
		if err := opt(op); err != nil {
			return nil, err
		}
	}

	c := &Client{
		parser: newParser(op.tagPrefix, op.ignoreTag),
	}

	return c, nil
}

// GenerateCreateTable outputs the `CREATE TABLE` schema of the specified Entity as a string.
func (c *Client) GenerateCreateTable(eb EntityBehavior) (string, error) {
	t, err := c.parser.Parse(eb)
	if err != nil {
		return "", err
	}

	return t.CreateTableSchema(), nil
}

// GenerateCreateTables outputs the `CREATE TABLE` schema of the specified Entity as a string slices.
func (c *Client) GenerateCreateTables(ebs []EntityBehavior) ([]string, error) {
	tables, err := c.parser.ParseMulti(ebs)
	if err != nil {
		return nil, err
	}

	ss := make([]string, 0, len(ebs))
	for i := range tables {
		t := tables[i]
		ss = append(ss, t.CreateTableSchema())
	}

	return ss, nil
}

// GenerateDropTable outputs the `DROP TABLE` schema of the specified Entity as a string.
func (c *Client) GenerateDropTable(eb EntityBehavior) (string, error) {
	t, err := c.parser.Parse(eb)
	if err != nil {
		return "", err
	}

	return t.DropTableSchema(), nil
}

// GenerateDropTables outputs the `DROP TABLE` schema of the specified Entity as a string slices.
func (c *Client) GenerateDropTables(ebs []EntityBehavior) ([]string, error) {
	tables, err := c.parser.ParseMulti(ebs)
	if err != nil {
		return nil, err
	}

	ss := make([]string, 0, len(ebs))
	for i := range tables {
		t := tables[i]
		ss = append(ss, t.DropTableSchema())
	}

	return ss, nil
}

// GenerateCreateIndexes outputs the `CREATE INDEX` schema of the specified Entity as a string slices.
func (c *Client) GenerateCreateIndexes(eb EntityBehavior) ([]string, error) {
	t, err := c.parser.Parse(eb)
	if err != nil {
		return nil, err
	}

	indexes := t.Indexes()
	ss := make([]string, 0, len(indexes))
	for i := range indexes {
		idx := indexes[i]
		ss = append(ss, idx.CreateIndexSchema())
	}

	return ss, nil
}

// GenerateDropIndexes outputs the `DROP INDEX` schema of the specified Entity as a string slices.
func (c *Client) GenerateDropIndexes(eb EntityBehavior) ([]string, error) {
	t, err := c.parser.Parse(eb)
	if err != nil {
		return nil, err
	}

	indexes := t.Indexes()
	ss := make([]string, 0, len(indexes))
	for i := range indexes {
		idx := indexes[i]
		ss = append(ss, idx.DropIndexSchema())
	}

	return ss, nil
}
