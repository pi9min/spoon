package spoon

import (
	"reflect"
	"strings"

	"golang.org/x/sync/errgroup"
)

// parser is
type parser struct {
	tagPrefix string
	ignoreTag string
}

func newParser(tagPrefix, ignoreTag string) *parser {
	return &parser{
		tagPrefix: tagPrefix,
		ignoreTag: ignoreTag,
	}
}

func (p *parser) Parse(eb EntityBehavior) (*Table, error) {
	columns, err := p.parseStruct(eb, p.tagPrefix)
	if err != nil {
		return nil, err
	}

	return newTable(eb.TableName(), columns, eb.PrimaryKey(), eb.Indexes()), nil
}

func (p *parser) ParseMulti(ebs []EntityBehavior) ([]*Table, error) {
	tables := make([]*Table, len(ebs))
	eg := errgroup.Group{}
	for i := range ebs {
		i := i
		eg.Go(func() error {
			columns, err := p.parseStruct(ebs[i], p.tagPrefix)
			if err != nil {
				return err
			}
			tables[i] = newTable(ebs[i].TableName(), columns, ebs[i].PrimaryKey(), ebs[i].Indexes())

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return tables, nil
}

func (p *parser) parseStruct(in interface{}, tp string) ([]*Column, error) {
	v := reflect.Indirect(reflect.ValueOf(in))
	t := v.Type()

	columns := make([]*Column, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.Type.Kind() == reflect.Ptr && sf.Type.Elem().Kind() == reflect.Struct {
			cols, err := p.parseStruct(v.Field(i).Interface(), tp)
			if err != nil {
				return nil, err
			}
			columns = append(columns, cols...)
		} else {
			col, err := p.parseField(sf, tp)
			if err != nil {
				if err == errIgnoreField {
					continue
				}

				return nil, err
			}

			columns = append(columns, col)
		}
	}

	return columns, nil
}

func (p *parser) parseField(field reflect.StructField, tp string) (*Column, error) {
	t := field.Tag.Get(tp)
	if t == "" {
		return newColumn(field.Name, map[string]string{}, field.Type), nil
	}

	tags := strings.Split(t, ",")
	ks := make(map[string]bool)
	ts := make([]string, 0, len(tags))
	for _, t := range tags {
		if t == p.ignoreTag {
			return nil, errIgnoreField
		}

		tag := strings.TrimSpace(t)
		// 重複してるものは入れない。先勝ち
		if _, ok := ks[tag]; !ok {
			ks[tag] = true
			ts = append(ts, tag)
		}
	}

	mts := p.mappingTag(ts)

	return newColumn(field.Name, mts, field.Type), nil
}

func (p *parser) mappingTag(tags []string) map[string]string {
	m := make(map[string]string)
	for _, elem := range tags {
		ss := strings.Split(elem, "=")
		switch len(ss) {
		case 1:
			m[ss[0]] = ""
		case 2:
			m[ss[0]] = ss[1]
		}
	}

	return m

}
