package input

type ThisIsASchema struct {
	Id int32 `db:"primary;shard;not null"`
	ThisIsAnIndexCols string `db:"index:idx1;not null"`
}

//
type ThisIsAShardSchema struct {
	Id int32 `db:"primary;shard;not null"`
	ThisIsAnIndexCols string `db:"index:idx1;not null"`
}