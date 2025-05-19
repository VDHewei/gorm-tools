# GenTool

Install GEN as a binary tool

## install

```shell
 go install github.com/VDHewei/gorm-tools@latest
```

## usage

```shell
 gorm-tools -h  
 
 Usage of gentool:
  --db string
        input mysql|postgres|sqlite|sqlserver|clickhouse. consult[https://gorm.io/docs/connecting_to_the_database.html] (default "mysql")
  --dsn string
        consult[https://gorm.io/docs/connecting_to_the_database.html]
  --fieldNullable
        generate with pointer when field is nullable
  --fieldCoverable
        generate with pointer when field has default value
  --fieldWithIndexTag
        generate field with gorm index tag
  --fieldWithTypeTag
        generate field with gorm column type tag
  --modelPkgName string
        generated model code's package name
  --outFile string
        query code file name, default: gen.go
  --outPath string
        specify a directory for output (default "./dao/query")
  --tables string
        enter the required data table or leave it blank
  --onlyModel
        only generate models (without query file)
  --withUnitTest
        generate unit test for query code
  --fieldSignable
        detect integer field's unsigned type, adjust generated data type
  -c, --config string
        is path for config yml(eg: gen.yaml)
  --fieldJSONTypeTag
        generate field with gorm json type(jsonb)
  --fieldsTypeMapping []string
        mapping field type mapping ,eg: jsonb:datatypes.JSON
  --importPkgPaths []string
        generate code import package path,eg: github.com/xxx/xxx 
  -d,--defaultYAMLConfigFile string
        generate default yaml config file
  --modelNameSignable
        keep model names and table names consistent, without using plural rewriting
  -v,--version
        print tool version
  -s,--showTables
        show database tables in console
  --showTable
        show table define fields in console
```
#### c
default ""
Is path for gen.yml
Replace the command line with a configuration file
The command line is the highest priority


#### db

default:mysql

input mysql or postgres or sqlite or sqlserver.

consult : https://gorm.io/docs/connecting_to_the_database.html

#### dsn

You can use all gorm's dsn.

consult : https://gorm.io/docs/connecting_to_the_database.html

#### fieldNullable

generate with pointer when field is nullable

#### fieldCoverable

generate with pointer when field has default value

#### fieldWithIndexTag

generate field with gorm index tag

#### fieldWithTypeTag

generate field with gorm column type tag

#### modelPkgName

default table name.

generated model code's package name.

#### outFile

query code file name, default: gen.go

#### outPath

specify a directory for output (default "./dao/query")

#### tables

Value : enter the required data table or leave it blank.

eg :

       --tables="orders" #orders table

       --tables="orders,users" #orders table and users table

       --tables=""          # All data tables in the database.

Generate some tables code.

#### withUnitTest

Value : False / True

Generate unit test.

#### fieldSignable

Value : False / True

detect integer field's unsigned type, adjust generated data type


### fieldJSONTypeTag

Value : False / True

generate field gorm tag with jsonb type(datatypes.JSON)

### showTables

Value : False / True

show database tables in console

```shell
gorm-tools --dsn "user:pwd@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=True&loc=Local" -s
+--------------------------+---------+
|        table_name        | comment |
+--------------------------+---------+
|         user             |         |
|       template_config    |         |
+--------------------------+---------+
```
### showTable

Value: String

show table define fields in console

```shell
gorm-tools --dsn "user:pwd@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=True&loc=Local" --showTable template_config
<template_config>
+----------------+--------------------------+----------+-------+---------+---------+
|     field      |           type           |  null    | pk/uk | default | comment |
+----------------+--------------------------+----------+-------+---------+---------+
|       id       |         integer          |  false   |  pk   |         |         |
|  extra_config  |          jsonb           |  false   |       |   {}    |         |
|   is_hidden    |         boolean          |  false   |       |  false  |         |
|    can_edit    |         boolean          |  false   |       |  true   |         |
|  is_required   |         boolean          |  false   |       |  false  |         |
|     status     |         boolean          |  false   |       |  true   |         |
|   created_at   | timestamp with time zone |  false   |       |         |         |
|   created_by   |           uuid           |  false   |       |         |         |
|   updated_at   | timestamp with time zone |   true   |       |         |         |
|   updated_by   |           uuid           |   true   |       |         |         |
|   deleted_at   | timestamp with time zone |   true   |       |         |         |
|   deleted_by   |           uuid           |   true   |       |         |         |
+----------------+--------------------------+----------+-------+---------+---------+
```


### example

```shell
gorm-tools --dsn "user:pwd@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=True&loc=Local" --tables "orders,doctor"
```
