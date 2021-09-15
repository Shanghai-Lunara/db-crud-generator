package gen


var dbTemplate = `
// This file was autogenerated by db-generator. Do not edit it manually!

{{$isShard := gt .Shard 0}}package {{.OutputPackage}}

import (
	"context"
	"database/sql"
	"errors"
	{{if $isShard}}"fmt"{{end}}
	sq "github.com/Shanghai-Lunara/squirrel"
	model "github.com/Shanghai-Lunara/{{.Project}}/{{.PackagePath}}"
)

//

func Get{{.Name}}SchemaName({{if $isShard}}v {{.ShardCols.Type}}{{end}}) string { {{if $isShard}}return fmt.Sprintf("{{.SchemaName}}_%d", v%{{.Shard}}){{else}}return "{{.SchemaName}}"{{end}} }

{{if $isShard}}
func Get{{$.Name}}ShardIndex(v {{.ShardCols.Type}}) {{.ShardCols.Type}} { return v % 4 }
{{ end }}

type {{.Name}}Insert struct {
	{{if $isShard}}handlers map[{{.ShardCols.Type}}]sq.InsertBuilder{{else}}handler sq.InsertBuilder{{end}}
	cache    map[int32]map[string]interface{}
	cols    map[string]struct{}
}

func New{{.Name}}Insert() *{{.Name}}Insert {
	return &{{.Name}}Insert{
		{{if $isShard}}handlers: make(map[{{.ShardCols.Type}}]sq.InsertBuilder){{else}}handler: sq.Insert("{{.SchemaName}}"){{end}},
		cache:    make(map[int32]map[string]interface{}),
		cols:    make(map[string]struct{}),
	}
}

func (i *{{.Name}}Insert) setValue(index int32, k string, v interface{}) {
	m, ok := i.cache[index]
	if !ok || m == nil {
		m = make(map[string]interface{})
		i.cache[index] = m
	}
	i.cols[k] = struct{}{}
	m[k] = v
}

{{if gt $.Shard 0}}
func (i *{{$.Name}}Insert) build() error {
	if len(i.cache) <= 0 {
		return errors.New("not insert rows")
	}
	flag := false
	for _, argMap := range i.cache {
		row := make([]interface{}, 0, len(argMap))
		tmpValue, ok := argMap["{{$.ShardCols.SchemaName}}"]
		if !ok {
			return errors.New("shard key not found")
		}
		shardColsValue := tmpValue.({{$.ShardCols.Type}})
		shardIndex := Get{{$.Name}}ShardIndex(shardColsValue)
		handler, ok1 := i.handlers[shardIndex]
		if !ok1 {
			handler = sq.Insert(Get{{$.Name}}SchemaName(shardColsValue))
		}
		for k, v := range argMap {
			if !flag {
				handler = handler.Columns(k)
			}
			row = append(row, v)
		}
		handler = handler.Values(row...)
		i.handlers[shardIndex] = handler
		flag = true
	}
	return nil
}
{{else}}
func (i *{{$.Name}}Insert) build() error {
	if len(i.cache) <= 0 {
		return errors.New("not insert rows")
	}
	flag := false
	i.handler = sq.Insert("{{.SchemaName}}")

	var cols []string
	for k := range i.cols {
		cols = append(cols, k)
	}
	for _, argMap := range i.cache {
		row := make([]interface{}, 0, len(i.cols))
		for _, k := range cols {
			if !flag {
				i.handler = i.handler.Columns(k)
			}
			v, ok := argMap[k]
			if !ok {
				switch k {
{{range $k, $v := .Cols}}
				case "{{$v.SchemaName}}": v = {{if eq $v.Type "int32"}}0{{end}}{{if eq $v.Type "int64"}}0{{end}}{{if eq $v.Type "string"}}""{{end}}{{if eq $v.Type "bool"}}false{{end}}
{{end}}
				}
			}
			row = append(row, v)
		}
		i.handler = i.handler.Values(row...)
		flag = true
	}
	return nil
}
{{end}}

func (i *{{$.Name}}Insert) Exec(ctx context.Context, db *sql.DB) error {
	if err := i.build(); err != nil {
		return err
	}
    {{if gt $.Shard 0}}
	for _, handler := range i.handlers {
		sqlStr, args, err := handler.ToSql()
		if err != nil {
			return err
		}
		if _, err = db.ExecContext(ctx, sqlStr, args...); err != nil {
			return err
		}
	}
	{{else}}
	sqlStr, args, err := i.handler.ToSql()
    if err != nil {
    	return err
    }
    if _, err = db.ExecContext(ctx, sqlStr, args...); err != nil {
    	return err
    }
	{{end}}
	return nil
}

func (i *{{$.Name}}Insert) ExecTx(ctx context.Context, tx *sql.Tx) error {
	if err := i.build(); err != nil {
		tx.Rollback()
		return err
	}

    {{if gt $.Shard 0}}
	for _, handler := range i.handlers {
		sqlStr, args, err := handler.ToSql()
		if err != nil {
			tx.Rollback()
			return err
		}
		if _, err = tx.ExecContext(ctx, sqlStr, args...); err != nil {
			tx.Rollback()
			return err
		}
	}
	{{else}}
	sqlStr, args, err := i.handler.ToSql()
    if err != nil {
    	tx.Rollback()
    	return err
    }
    if _, err = tx.ExecContext(ctx, sqlStr, args...); err != nil {
    	tx.Rollback()
    	return err
    }
	{{end}}
	return nil
}

{{range $k, $v := .Cols}}
func (i *{{$.Name}}Insert) {{$v.Name}}(index int32, v {{$v.Type}}) *{{$.Name}}Insert {
	i.setValue(index, "{{$v.SchemaName}}", v)
	return i
}
{{end}}



type {{.Name}}Select struct {
	handler sq.SelectBuilder

	tmp      *model.{{.Name}}
	fieldMap map[string]interface{}
}

func New{{.Name}}Select() *{{.Name}}Select {
	return &{{.Name}}Select{
		handler:  sq.Select().From("{{.SchemaName}}"),
		tmp:      &model.{{.Name}}{},
		fieldMap: map[string]interface{}{},
	}
}

func (s *{{.Name}}Select) Count(db *sql.DB) (int, error) {
	sqlStr, args, err := s.handler.Column("COUNT(1)").ToSql()
	if err != nil {
		return 0, err
	}
	rows, err := db.QueryContext(context.Background(), sqlStr, args...)
	if err != nil {
		return 0, err
	}
	var dest int

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return 0, err
		}
		return 0, sql.ErrNoRows
	}
	err = rows.Scan(&dest)
	if err != nil {
		return 0, err
	}
	_ = rows.Close()
	return dest, nil
}

// page

func (s *{{.Name}}Select) Page(pageIndex, pageSize int) *{{.Name}}Select {
	s.handler = s.handler.Limit(uint64(pageSize)).Offset(uint64((pageIndex - 1) * pageSize))
	return s
}


func (s *{{.Name}}Select) Query(ctx context.Context, db *sql.DB) ([]*model.{{.Name}}, error) {
	sqlStr, args, err := s.handler.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	results := make([]*model.{{.Name}}, 0)
	columns, _ := rows.Columns()
	dest := make([]interface{}, 0, len(columns))

	for _, v := range columns {
		dest = append(dest, s.fieldMap[v])
	}
	for rows.Next() {
		result := &model.{{.Name}}{}
		err := rows.Scan(dest...)
		if err != nil {
			return nil, err
		}
		{{range $k, $v := .Cols}}
		result.{{$v.Name}} = s.tmp.{{$v.Name}}{{end}}

		results = append(results, result)
	}
	return results, nil
}

func (s *{{.Name}}Select) QueryRow(ctx context.Context, db *sql.DB) (*model.{{.Name}}, error) {
	sqlStr, args, err := s.handler.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	columns, _ := rows.Columns()
	dest := make([]interface{}, 0, len(columns))

	for _, v := range columns {
		dest = append(dest, s.fieldMap[v])
	}
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}
	err = rows.Scan(dest...)
	if err != nil {
		return nil, err
	}
	_ = rows.Close()
	result := &model.{{.Name}}{}
	{{range $k, $v := .Cols}}
	result.{{$v.Name}} = s.tmp.{{$v.Name}}{{end}}

	return result, nil
}

func (s *{{$.Name}}Select) Select() *{{$.Name}}Select {
	{{range $k, $v := .Cols}}s.handler = s.handler.Columns("`+"`{{$v.SchemaName}}`"+`")
	s.fieldMap["{{$v.SchemaName}}"] = &(s.tmp.{{$v.Name}})
	{{end}}
	return s
}

{{range $k, $v := .Cols}}
{{$isShardCols := false}}
{{if $.ShardCols}}
{{$isShardCols = eq $v.Name $.ShardCols.Name}}
{{end}}
{{$IsPrimary := eq $v.Name $.Primary.Name}}
func (s *{{$.Name}}Select) Select{{$v.Name}}() *{{$.Name}}Select {
	s.handler = s.handler.Columns("`+"`{{$v.SchemaName}}`"+`")
	s.fieldMap["{{$v.SchemaName}}"] = &(s.tmp.{{$v.Name}})
	return s
}

{{if or $isShardCols $v.IsIndex $IsPrimary}}
func (s *{{$.Name}}Select) Where{{$v.Name}}Eq(v {{$v.Type}}) *{{$.Name}}Select {
	s.handler = s.handler.Where(sq.Eq{"`+"`{{$v.SchemaName}}`"+`": v})
	return s
}

func (s *{{$.Name}}Select) Where{{$v.Name}}NotEq(v {{$v.Type}}) *{{$.Name}}Select {
	s.handler = s.handler.Where(sq.NotEq{"`+"`{{$v.SchemaName}}`"+`": v})
	return s
}

func (s *{{$.Name}}Select) Where{{$v.Name}}Gt(v {{$v.Type}}) *{{$.Name}}Select {
	s.handler = s.handler.Where(sq.Gt{"`+"`{{$v.SchemaName}}`"+`": v})
	return s
}

func (s *{{$.Name}}Select) Where{{$v.Name}}Lt(v {{$v.Type}}) *{{$.Name}}Select {
	s.handler = s.handler.Where(sq.Lt{"`+"`{{$v.SchemaName}}`"+`": v})
	return s
}

func (s *{{$.Name}}Select) Where{{$v.Name}}GtOrEq(v {{$v.Type}}) *{{$.Name}}Select {
	s.handler = s.handler.Where(sq.GtOrEq{"`+"`{{$v.SchemaName}}`"+`": v})
	return s
}

func (s *{{$.Name}}Select) Where{{$v.Name}}LtOrEq(v {{$v.Type}}) *{{$.Name}}Select {
	s.handler = s.handler.Where(sq.GtOrEq{"`+"`{{$v.SchemaName}}`"+`": v})
	return s
}

func (s *{{$.Name}}Select) Where{{$v.Name}}Like(v {{$v.Type}}) *{{$.Name}}Select {
	s.handler = s.handler.Where(sq.Like{"`+"`{{$v.SchemaName}}`"+`": v})
	return s
}
{{end}}
{{end}}





type {{.Name}}Update struct {
	handler sq.UpdateBuilder
}

func New{{.Name}}Update() *{{.Name}}Update {
	return &{{.Name}}Update{
		handler: sq.Update("{{.SchemaName}}"),
	}
}

func (u *{{.Name}}Update) Exec(ctx context.Context, db *sql.DB) error {
	sqlStr, args, err := u.handler.ToSql()
    if err != nil {
    	return err
    }
    if _, err := db.ExecContext(ctx, sqlStr, args...); err != nil {
    	return err
    }
	return nil
}

func (u *{{.Name}}Update) ExecTx(ctx context.Context, tx *sql.Tx) error {
	sqlStr, args, err := u.handler.ToSql()
    if err != nil {
		tx.Rollback()
    	return err
    }
    if _, err := tx.ExecContext(ctx, sqlStr, args...); err != nil {
		tx.Rollback()
    	return err
    }
	return nil
}

{{range $k, $v := .Cols}}
{{$isShardCols := false}}
{{if $.ShardCols}}
{{$isShardCols = eq $v.Name $.ShardCols.Name}}
{{end}}
{{$IsPrimary := eq $v.Name $.Primary.Name}}
func (u *{{$.Name}}Update) {{$v.Name}}(v {{$v.Type}}) *{{$.Name}}Update {
	u.handler = u.handler.Set("`+"`{{$v.SchemaName}}`"+`", v)
	return u
}

{{if or $isShardCols $v.IsIndex $IsPrimary}}
func (u *{{$.Name}}Update) Where{{$v.Name}}Eq(v {{$v.Type}}) *{{$.Name}}Update {
	u.handler = u.handler.Where(sq.Eq{"`+"`{{$v.SchemaName}}`"+`": v})
	return u
}

func (u *{{$.Name}}Update) Where{{$v.Name}}NotEq(v {{$v.Type}}) *{{$.Name}}Update {
	u.handler = u.handler.Where(sq.NotEq{"`+"`{{$v.SchemaName}}`"+`": v})
	return u
}

func (u *{{$.Name}}Update) Where{{$v.Name}}Gt(v {{$v.Type}}) *{{$.Name}}Update {
	u.handler = u.handler.Where(sq.Gt{"`+"`{{$v.SchemaName}}`"+`": v})
	return u
}

func (u *{{$.Name}}Update) Where{{$v.Name}}Lt(v {{$v.Type}}) *{{$.Name}}Update {
	u.handler = u.handler.Where(sq.Lt{"`+"`{{$v.SchemaName}}`"+`": v})
	return u
}

func (u *{{$.Name}}Update) Where{{$v.Name}}GtOrEq(v {{$v.Type}}) *{{$.Name}}Update {
	u.handler = u.handler.Where(sq.GtOrEq{"`+"`{{$v.SchemaName}}`"+`": v})
	return u
}

func (u *{{$.Name}}Update) Where{{$v.Name}}LtOrEq(v {{$v.Type}}) *{{$.Name}}Update {
	u.handler = u.handler.Where(sq.GtOrEq{"`+"`{{$v.SchemaName}}`"+`": v})
	return u
}

func (u *{{$.Name}}Update) Where{{$v.Name}}Like(v {{$v.Type}}) *{{$.Name}}Update {
	u.handler = u.handler.Where(sq.Like{"`+"`{{$v.SchemaName}}`"+`": v})
	return u
}
{{end}}
{{end}}
`