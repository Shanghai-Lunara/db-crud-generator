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