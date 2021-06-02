# db-crud-generator
The tool which was used to auto generating crud codes with the go struct


## Model struct define

It will scan tag of fields.

Use "primary", "index" to mark whether this field is the primary key, index field.

```go
type ThisIsASchema struct {
	Id int32 `db:"primary;index;shard;not null"`
	ThisIsAnIndexCols string `db:"index;not null"`
}
```

## How to generate

- Generate from command flags
```go
package main

import (
	gen "github.com/Shanghai-Lunara/db-crud-generator"
)

func main() {
	gen.GenerateWithFlagScan()
}

```

run as: 

```shell
go run main.go  -projectName=my_project -scanPath=path/to/model -outputPath=path/to/out 
```

- Or generate from parameter
```go
package main

import (
	gen "github.com/Shanghai-Lunara/db-crud-generator"
)

func main() {
	gen.Generate("my_project", "path/to/model", "path/to/out")
}

```

- Attention

It will generate a go file contains insert, query and update methods.

The where clause will only generate fields marked with primary and index.

## Use after generated

```go
package main

import (
	"context"
	"database/sql"
	"github.com/Shanghai-Lunara/db-crud-generator/example/out"
	"time"
)

func f(tx *sql.Tx, db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// insert, index arg mark value at same row
	err := out.NewThisIsASchemaInsert().
		SetId(1, 1).
		SetThisIsAnIndexCols(1, "emm").
		SetId(2, 2).
		SetId(3, 3).
		SetThisIsAnIndexCols(2, "emm2").
		SetThisIsAnIndexCols(3, "emm3").
		ExecTx(ctx, tx)

	// select
	result, err := out.NewThisIsASchemaSelect().
		SelectThisIsAnIndexCols().
		WhereIdEq(1).
		Query(ctx, db)

	// update
	err := out.NewThisIsASchemaUpdate().
		SetThisIsAnIndexCols("oh").
		WhereIdEq(1).
		ExecTx(tx)
}

```