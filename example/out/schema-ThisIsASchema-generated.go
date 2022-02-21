// This file was autogenerated by db-generator. Do not edit it manually!

package out

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	model "github.com/Shanghai-Lunara/db-crud-generator/example/input"
	sq "github.com/Shanghai-Lunara/squirrel"
	"strings"
	// custom model imports
	"time"
)

//

func GetThisIsASchemaSchemaName() string { return "this_is_a_schema" }

func GetThisIsASchemaColsNameId() string {
	return "id"
}

func GetThisIsASchemaColsNameThisIsAnIndexCols() string {
	return "thisIsAnIndexCols"
}

func GetThisIsASchemaColsNameIgnoreCols() string {
	return "ignoreCols"
}

func GetThisIsASchemaColsNameCreateTime() string {
	return "createTime"
}

type ThisIsASchemaInsert struct {
	handler sq.InsertBuilder
	cache   map[int32]map[string]interface{}
	cols    map[string]struct{}
}

func NewThisIsASchemaInsert() *ThisIsASchemaInsert {
	return &ThisIsASchemaInsert{
		handler: sq.Insert("this_is_a_schema"),
		cache:   make(map[int32]map[string]interface{}),
		cols:    make(map[string]struct{}),
	}
}

func (i *ThisIsASchemaInsert) setValue(index int32, k string, v interface{}) {
	m, ok := i.cache[index]
	if !ok || m == nil {
		m = make(map[string]interface{})
		i.cache[index] = m
	}
	i.cols[k] = struct{}{}
	m[k] = v
}

func (i *ThisIsASchemaInsert) build() error {
	if len(i.cache) <= 0 {
		return errors.New("not insert rows")
	}
	flag := false
	i.handler = sq.Insert("this_is_a_schema")

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

				case "id":
					v = 0
					break

				case "thisIsAnIndexCols":
					v = ""
					break

				case "ignoreCols":
					v = false
					break

				case "createTime":
					v = time.Unix(0, 0)
					break

				}
			}
			row = append(row, v)
		}
		i.handler = i.handler.Values(row...)
		flag = true
	}
	return nil
}

func (i *ThisIsASchemaInsert) Exec(ctx context.Context, db *sql.DB) (sql.Result, error) {
	if err := i.build(); err != nil {
		return nil, err
	}
	var result sql.Result
	var err error

	sqlStr, args, err1 := i.handler.ToSql()
	if err1 != nil {
		return nil, err1
	}
	result, err = db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (i *ThisIsASchemaInsert) ExecTx(ctx context.Context, tx *sql.Tx) (sql.Result, error) {
	if err := i.build(); err != nil {
		tx.Rollback()
		return nil, err
	}
	var result sql.Result
	var err error

	sqlStr, args, err1 := i.handler.ToSql()
	if err1 != nil {
		tx.Rollback()
		return nil, err1
	}
	result, err = tx.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return result, nil
}

func (i *ThisIsASchemaInsert) Id(index int32, v int32) *ThisIsASchemaInsert {
	i.setValue(index, "id", v)
	return i
}

func (i *ThisIsASchemaInsert) ThisIsAnIndexCols(index int32, v string) *ThisIsASchemaInsert {
	i.setValue(index, "thisIsAnIndexCols", v)
	return i
}

func (i *ThisIsASchemaInsert) IgnoreCols(index int32, v bool) *ThisIsASchemaInsert {
	i.setValue(index, "ignoreCols", v)
	return i
}

func (i *ThisIsASchemaInsert) CreateTime(index int32, v time.Time) *ThisIsASchemaInsert {
	i.setValue(index, "createTime", v)
	return i
}

type ThisIsASchemaSelect struct {
	handler sq.SelectBuilder

	tmp      *model.ThisIsASchema
	fieldMap map[string]interface{}
}

func NewThisIsASchemaSelect() *ThisIsASchemaSelect {
	return &ThisIsASchemaSelect{
		handler:  sq.Select().From("this_is_a_schema"),
		tmp:      &model.ThisIsASchema{},
		fieldMap: map[string]interface{}{},
	}
}

func (s *ThisIsASchemaSelect) Count(db *sql.DB) (int, error) {
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

func (s *ThisIsASchemaSelect) Page(pageIndex, pageSize int) *ThisIsASchemaSelect {
	s.handler = s.handler.Limit(uint64(pageSize)).Offset(uint64((pageIndex - 1) * pageSize))
	return s
}

func (s *ThisIsASchemaSelect) OrderByRandom() *ThisIsASchemaSelect {
	s.handler = s.handler.OrderBy("rand()")
	return s
}

func (s *ThisIsASchemaSelect) Query(ctx context.Context, db *sql.DB) ([]*model.ThisIsASchema, error) {
	sqlStr, args, err := s.handler.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	results := make([]*model.ThisIsASchema, 0)
	columns, _ := rows.Columns()
	dest := make([]interface{}, 0, len(columns))

	for _, v := range columns {
		dest = append(dest, s.fieldMap[v])
	}
	for rows.Next() {
		result := &model.ThisIsASchema{}
		err := rows.Scan(dest...)
		if err != nil {
			return nil, err
		}

		result.Id = s.tmp.Id
		result.ThisIsAnIndexCols = s.tmp.ThisIsAnIndexCols
		result.IgnoreCols = s.tmp.IgnoreCols
		result.CreateTime = s.tmp.CreateTime

		results = append(results, result)
	}
	return results, nil
}

func (s *ThisIsASchemaSelect) QueryTx(ctx context.Context, tx *sql.Tx) ([]*model.ThisIsASchema, error) {
	sqlStr, args, err := s.handler.ToSql()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	rows, err := tx.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	results := make([]*model.ThisIsASchema, 0)
	columns, _ := rows.Columns()
	dest := make([]interface{}, 0, len(columns))

	for _, v := range columns {
		dest = append(dest, s.fieldMap[v])
	}
	for rows.Next() {
		result := &model.ThisIsASchema{}
		err := rows.Scan(dest...)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}

		result.Id = s.tmp.Id
		result.ThisIsAnIndexCols = s.tmp.ThisIsAnIndexCols
		result.IgnoreCols = s.tmp.IgnoreCols
		result.CreateTime = s.tmp.CreateTime

		results = append(results, result)
	}
	return results, nil
}

func (s *ThisIsASchemaSelect) QueryRow(ctx context.Context, db *sql.DB) (*model.ThisIsASchema, error) {
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
	result := &model.ThisIsASchema{}

	result.Id = s.tmp.Id
	result.ThisIsAnIndexCols = s.tmp.ThisIsAnIndexCols
	result.IgnoreCols = s.tmp.IgnoreCols
	result.CreateTime = s.tmp.CreateTime

	return result, nil
}

func (s *ThisIsASchemaSelect) QueryRowTx(ctx context.Context, tx *sql.Tx) (*model.ThisIsASchema, error) {
	sqlStr, args, err := s.handler.ToSql()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	rows, err := tx.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	columns, _ := rows.Columns()
	dest := make([]interface{}, 0, len(columns))

	for _, v := range columns {
		dest = append(dest, s.fieldMap[v])
	}
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		return nil, sql.ErrNoRows
	}
	err = rows.Scan(dest...)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	_ = rows.Close()
	result := &model.ThisIsASchema{}

	result.Id = s.tmp.Id
	result.ThisIsAnIndexCols = s.tmp.ThisIsAnIndexCols
	result.IgnoreCols = s.tmp.IgnoreCols
	result.CreateTime = s.tmp.CreateTime

	return result, nil
}

func (s *ThisIsASchemaSelect) Select() *ThisIsASchemaSelect {
	s.handler = s.handler.Columns("`id`")
	s.fieldMap["id"] = &(s.tmp.Id)
	s.handler = s.handler.Columns("`thisIsAnIndexCols`")
	s.fieldMap["thisIsAnIndexCols"] = &(s.tmp.ThisIsAnIndexCols)
	s.handler = s.handler.Columns("`ignoreCols`")
	s.fieldMap["ignoreCols"] = &(s.tmp.IgnoreCols)
	s.handler = s.handler.Columns("`createTime`")
	s.fieldMap["createTime"] = &(s.tmp.CreateTime)

	return s
}

func (s *ThisIsASchemaSelect) SelectId() *ThisIsASchemaSelect {
	s.handler = s.handler.Columns("`id`")
	s.fieldMap["id"] = &(s.tmp.Id)
	return s
}

func (s *ThisIsASchemaSelect) WhereIdEq(v int32) *ThisIsASchemaSelect {
	s.handler = s.handler.Where(sq.Eq{"`id`": v})
	return s
}

func (s *ThisIsASchemaSelect) WhereIdNotEq(v int32) *ThisIsASchemaSelect {
	s.handler = s.handler.Where(sq.NotEq{"`id`": v})
	return s
}

func (s *ThisIsASchemaSelect) WhereIdGt(v int32) *ThisIsASchemaSelect {
	s.handler = s.handler.Where(sq.Gt{"`id`": v})
	return s
}

func (s *ThisIsASchemaSelect) WhereIdLt(v int32) *ThisIsASchemaSelect {
	s.handler = s.handler.Where(sq.Lt{"`id`": v})
	return s
}

func (s *ThisIsASchemaSelect) WhereIdGtOrEq(v int32) *ThisIsASchemaSelect {
	s.handler = s.handler.Where(sq.GtOrEq{"`id`": v})
	return s
}

func (s *ThisIsASchemaSelect) WhereIdLtOrEq(v int32) *ThisIsASchemaSelect {
	s.handler = s.handler.Where(sq.GtOrEq{"`id`": v})
	return s
}

func (s *ThisIsASchemaSelect) WhereIdLike(v int32) *ThisIsASchemaSelect {
	s.handler = s.handler.Where(sq.Like{"`id`": v})
	return s
}

func (s *ThisIsASchemaSelect) OrderById(desc bool) *ThisIsASchemaSelect {
	if desc {
		s.handler = s.handler.OrderBy("id DESC")
	} else {
		s.handler = s.handler.OrderBy("id ASC")
	}
	return s
}

func (s *ThisIsASchemaSelect) WhereIdIn(args ...int32) *ThisIsASchemaSelect {
	sb := strings.Builder{}
	sb.WriteString("id IN (")
	for index, v := range args {
		sb.WriteString(fmt.Sprintf("%v", v))
		if index < len(args)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
	s.handler = s.handler.Where(sq.Expr(sb.String()))
	return s
}

func (s *ThisIsASchemaSelect) SelectThisIsAnIndexCols() *ThisIsASchemaSelect {
	s.handler = s.handler.Columns("`thisIsAnIndexCols`")
	s.fieldMap["thisIsAnIndexCols"] = &(s.tmp.ThisIsAnIndexCols)
	return s
}

func (s *ThisIsASchemaSelect) WhereThisIsAnIndexColsEq(v string) *ThisIsASchemaSelect {
	s.handler = s.handler.Where(sq.Eq{"`thisIsAnIndexCols`": v})
	return s
}

func (s *ThisIsASchemaSelect) WhereThisIsAnIndexColsNotEq(v string) *ThisIsASchemaSelect {
	s.handler = s.handler.Where(sq.NotEq{"`thisIsAnIndexCols`": v})
	return s
}

func (s *ThisIsASchemaSelect) WhereThisIsAnIndexColsGt(v string) *ThisIsASchemaSelect {
	s.handler = s.handler.Where(sq.Gt{"`thisIsAnIndexCols`": v})
	return s
}

func (s *ThisIsASchemaSelect) WhereThisIsAnIndexColsLt(v string) *ThisIsASchemaSelect {
	s.handler = s.handler.Where(sq.Lt{"`thisIsAnIndexCols`": v})
	return s
}

func (s *ThisIsASchemaSelect) WhereThisIsAnIndexColsGtOrEq(v string) *ThisIsASchemaSelect {
	s.handler = s.handler.Where(sq.GtOrEq{"`thisIsAnIndexCols`": v})
	return s
}

func (s *ThisIsASchemaSelect) WhereThisIsAnIndexColsLtOrEq(v string) *ThisIsASchemaSelect {
	s.handler = s.handler.Where(sq.GtOrEq{"`thisIsAnIndexCols`": v})
	return s
}

func (s *ThisIsASchemaSelect) WhereThisIsAnIndexColsLike(v string) *ThisIsASchemaSelect {
	s.handler = s.handler.Where(sq.Like{"`thisIsAnIndexCols`": v})
	return s
}

func (s *ThisIsASchemaSelect) OrderByThisIsAnIndexCols(desc bool) *ThisIsASchemaSelect {
	if desc {
		s.handler = s.handler.OrderBy("thisIsAnIndexCols DESC")
	} else {
		s.handler = s.handler.OrderBy("thisIsAnIndexCols ASC")
	}
	return s
}

func (s *ThisIsASchemaSelect) WhereThisIsAnIndexColsIn(args ...string) *ThisIsASchemaSelect {
	sb := strings.Builder{}
	sb.WriteString("thisIsAnIndexCols IN (")
	for index, v := range args {
		sb.WriteString(fmt.Sprintf("%v", v))
		if index < len(args)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
	s.handler = s.handler.Where(sq.Expr(sb.String()))
	return s
}

func (s *ThisIsASchemaSelect) SelectIgnoreCols() *ThisIsASchemaSelect {
	s.handler = s.handler.Columns("`ignoreCols`")
	s.fieldMap["ignoreCols"] = &(s.tmp.IgnoreCols)
	return s
}

func (s *ThisIsASchemaSelect) SelectCreateTime() *ThisIsASchemaSelect {
	s.handler = s.handler.Columns("`createTime`")
	s.fieldMap["createTime"] = &(s.tmp.CreateTime)
	return s
}

type ThisIsASchemaUpdate struct {
	handler   sq.UpdateBuilder
	whereFlag bool
}

func NewThisIsASchemaUpdate() *ThisIsASchemaUpdate {
	return &ThisIsASchemaUpdate{
		handler:   sq.Update("this_is_a_schema"),
		whereFlag: false,
	}
}

func (u *ThisIsASchemaUpdate) Exec(ctx context.Context, db *sql.DB) error {
	if !u.whereFlag {
		return errors.New("update no where clause")
	}
	sqlStr, args, err := u.handler.ToSql()
	if err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, sqlStr, args...); err != nil {
		return err
	}
	return nil
}

func (u *ThisIsASchemaUpdate) ExecTx(ctx context.Context, tx *sql.Tx) error {
	if !u.whereFlag {
		return errors.New("update no where clause")
	}
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

func (u *ThisIsASchemaUpdate) Id(v int32) *ThisIsASchemaUpdate {
	u.handler = u.handler.Set("`id`", v)
	return u
}

func (u *ThisIsASchemaUpdate) WhereIdEq(v int32) *ThisIsASchemaUpdate {
	u.handler = u.handler.Where(sq.Eq{"`id`": v})
	u.whereFlag = true
	return u
}

func (u *ThisIsASchemaUpdate) WhereIdNotEq(v int32) *ThisIsASchemaUpdate {
	u.handler = u.handler.Where(sq.NotEq{"`id`": v})
	u.whereFlag = true
	return u
}

func (u *ThisIsASchemaUpdate) WhereIdGt(v int32) *ThisIsASchemaUpdate {
	u.handler = u.handler.Where(sq.Gt{"`id`": v})
	u.whereFlag = true
	return u
}

func (u *ThisIsASchemaUpdate) WhereIdLt(v int32) *ThisIsASchemaUpdate {
	u.handler = u.handler.Where(sq.Lt{"`id`": v})
	u.whereFlag = true
	return u
}

func (u *ThisIsASchemaUpdate) WhereIdGtOrEq(v int32) *ThisIsASchemaUpdate {
	u.handler = u.handler.Where(sq.GtOrEq{"`id`": v})
	u.whereFlag = true
	return u
}

func (u *ThisIsASchemaUpdate) WhereIdLtOrEq(v int32) *ThisIsASchemaUpdate {
	u.handler = u.handler.Where(sq.GtOrEq{"`id`": v})
	u.whereFlag = true
	return u
}

func (u *ThisIsASchemaUpdate) WhereIdLike(v int32) *ThisIsASchemaUpdate {
	u.handler = u.handler.Where(sq.Like{"`id`": v})
	u.whereFlag = true
	return u
}

func (u *ThisIsASchemaUpdate) WhereIdIn(args ...int32) *ThisIsASchemaUpdate {
	sb := strings.Builder{}
	sb.WriteString("id IN (")
	for index, v := range args {
		sb.WriteString(fmt.Sprintf("%v", v))
		if index < len(args)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
	u.handler = u.handler.Where(sq.Expr(sb.String()))
	u.whereFlag = true
	return u
}

func (u *ThisIsASchemaUpdate) ThisIsAnIndexCols(v string) *ThisIsASchemaUpdate {
	u.handler = u.handler.Set("`thisIsAnIndexCols`", v)
	return u
}

func (u *ThisIsASchemaUpdate) WhereThisIsAnIndexColsEq(v string) *ThisIsASchemaUpdate {
	u.handler = u.handler.Where(sq.Eq{"`thisIsAnIndexCols`": v})
	u.whereFlag = true
	return u
}

func (u *ThisIsASchemaUpdate) WhereThisIsAnIndexColsNotEq(v string) *ThisIsASchemaUpdate {
	u.handler = u.handler.Where(sq.NotEq{"`thisIsAnIndexCols`": v})
	u.whereFlag = true
	return u
}

func (u *ThisIsASchemaUpdate) WhereThisIsAnIndexColsGt(v string) *ThisIsASchemaUpdate {
	u.handler = u.handler.Where(sq.Gt{"`thisIsAnIndexCols`": v})
	u.whereFlag = true
	return u
}

func (u *ThisIsASchemaUpdate) WhereThisIsAnIndexColsLt(v string) *ThisIsASchemaUpdate {
	u.handler = u.handler.Where(sq.Lt{"`thisIsAnIndexCols`": v})
	u.whereFlag = true
	return u
}

func (u *ThisIsASchemaUpdate) WhereThisIsAnIndexColsGtOrEq(v string) *ThisIsASchemaUpdate {
	u.handler = u.handler.Where(sq.GtOrEq{"`thisIsAnIndexCols`": v})
	u.whereFlag = true
	return u
}

func (u *ThisIsASchemaUpdate) WhereThisIsAnIndexColsLtOrEq(v string) *ThisIsASchemaUpdate {
	u.handler = u.handler.Where(sq.GtOrEq{"`thisIsAnIndexCols`": v})
	u.whereFlag = true
	return u
}

func (u *ThisIsASchemaUpdate) WhereThisIsAnIndexColsLike(v string) *ThisIsASchemaUpdate {
	u.handler = u.handler.Where(sq.Like{"`thisIsAnIndexCols`": v})
	u.whereFlag = true
	return u
}

func (u *ThisIsASchemaUpdate) WhereThisIsAnIndexColsIn(args ...string) *ThisIsASchemaUpdate {
	sb := strings.Builder{}
	sb.WriteString("thisIsAnIndexCols IN (")
	for index, v := range args {
		sb.WriteString(fmt.Sprintf("%v", v))
		if index < len(args)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
	u.handler = u.handler.Where(sq.Expr(sb.String()))
	u.whereFlag = true
	return u
}

func (u *ThisIsASchemaUpdate) IgnoreCols(v bool) *ThisIsASchemaUpdate {
	u.handler = u.handler.Set("`ignoreCols`", v)
	return u
}

func (u *ThisIsASchemaUpdate) CreateTime(v time.Time) *ThisIsASchemaUpdate {
	u.handler = u.handler.Set("`createTime`", v)
	return u
}
