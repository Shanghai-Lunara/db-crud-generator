package input

type ThisIsASchema struct {
	Id int32 `db:"primary;index;shard;not null"`
	ThisIsAnIndexCols string `db:"index;not null"`
}