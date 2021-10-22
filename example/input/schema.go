package input

type ThisIsASchema struct {
	Id int32 `db:"primary;not null"`
	ThisIsAnIndexCols string `db:"index:idx1;not null"`
	IgnoreCols bool `json:"ignoreCols" db:"-"`
}