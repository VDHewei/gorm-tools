package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	gen "gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type (
	CmdParams struct {
		args                  *Options `yaml:"-" json:"-"`
		DSN                   string   `yaml:"dsn"`               // consult[https://gorm.io/docs/connecting_to_the_database.html]"
		DB                    string   `yaml:"db"`                // input mysql or postgres or sqlite or sqlserver. consult[https://gorm.io/docs/connecting_to_the_database.html]
		Tables                []string `yaml:"tables"`            // enter the required data table or leave it blank
		ExcludeTableList      []string `yaml:"exclude_tables"`    // enter the exclude data table or leave it blank
		OnlyModel             bool     `yaml:"onlyModel"`         // only generate model
		OutPath               string   `yaml:"outPath"`           // specify a directory for output
		OutFile               string   `yaml:"outFile"`           // query code file name, default: gen.go
		WithUnitTest          bool     `yaml:"withUnitTest"`      // generate unit test for query code
		ModelPkgName          string   `yaml:"modelPkgName"`      // generated model code's package name
		FieldNullable         bool     `yaml:"fieldNullable"`     // generate with pointer when field is nullable
		FieldCoverable        bool     `yaml:"fieldCoverable"`    // generate with pointer when field has default value
		FieldWithIndexTag     bool     `yaml:"fieldWithIndexTag"` // generate field with gorm index tag
		FieldWithTypeTag      bool     `yaml:"fieldWithTypeTag"`  // generate field with gorm column type tag
		FieldSignable         bool     `yaml:"fieldSignable"`     // detect integer field's unsigned type, adjust generated data type
		FieldJSONTypeTag      bool     `yaml:"fieldJSONTypeTag"`  // generate field with gorm json type
		FieldsTypeMapping     []string `yaml:"fieldsTypeMapping"` // generate table field with gorm type
		ImportPkgPaths        []string `yaml:"importPkgPaths"`    // generate code import package path
		Mode                  string   `yaml:"mode"`              // generate mode (input DefaultQuery|QueryInterface|OutContext)
		defaultYAMLConfigFile string   `json:"-" yaml:"-"`        // generate default yaml config file
	}
	// YamlConfig is yaml config struct
	YamlConfig struct {
		Version  string     `yaml:"version"`
		Database *CmdParams `yaml:"database"`
	}
	// DBType database type
	DBType string
)

func (d DBType) String() string {
	return string(d)
}

// NewFromYaml parse cmd param from yaml
func NewFromYaml(path string) *CmdParams {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("parseCmdFromYaml fail %s", err.Error())
		return nil
	}
	defer file.Close() // nolint
	var yamlConfig YamlConfig
	if err = yaml.NewDecoder(file).Decode(&yamlConfig); err != nil {
		log.Fatalf("parseCmdFromYaml fail %s", err.Error())
		return nil
	}
	return yamlConfig.Database
}

func New() *CmdParams {
	var c = &CmdParams{}
	return c
}

func (c *CmdParams) Revise() *CmdParams {
	if c == nil {
		return c
	}
	if c.DB == "" {
		if c.DSN == "" {
			c.DB = DbMySQL.String()
		} else {
			c.DB = extractDSNDBType(c.DSN)
		}
	}
	if c.DB == DbPostgres.String() {
		c.DSN = schemaPostgresToValues(c.DSN)
	} else {
		if c.DB == DbSQLite.String() {
			c.DSN = strings.TrimPrefix(c.DSN, fmt.Sprintf("%s://", c.DB))
			c.DSN = fmt.Sprintf("%s:%s", "file", c.DSN)
		} else if c.DB != DbSQLServer.String() {
			c.DSN = strings.TrimPrefix(c.DSN, fmt.Sprintf("%s://", c.DB))
		}
	}
	if c.Mode == "" {
		c.Mode = "DefaultQuery|Context|QueryInterface"
	}
	if c.OutPath == "" {
		c.OutPath = DefaultQueryPath
	}
	var excludes = make(map[string]struct{})
	if len(c.ExcludeTableList) > 0 {
		tableList := make([]string, 0, len(c.ExcludeTableList))
		for _, tableName := range c.ExcludeTableList {
			_tableName := strings.TrimSpace(tableName) // trim leading and trailing space in tableName
			if _tableName == "" {                      // skip empty tableName
				continue
			}
			if _, ok := excludes[_tableName]; ok {
				continue
			}
			excludes[_tableName] = struct{}{}
			tableList = append(tableList, _tableName)
		}
		c.ExcludeTableList = tableList
	}
	if len(c.Tables) > 0 {
		indexes := make(map[string]struct{})
		tableList := make([]string, 0, len(c.Tables))
		for _, tableName := range c.Tables {
			_tableName := strings.TrimSpace(tableName) // trim leading and trailing space in tableName
			if _tableName == "" {                      // skip empty tableName
				continue
			}
			if _, ok := indexes[_tableName]; ok {
				continue
			}
			indexes[_tableName] = struct{}{}
			tableList = append(tableList, _tableName)
		}
		c.Tables = tableList
	}
	return c
}

func (c *CmdParams) LoadArgs() (*Options, error) {
	// choose is file or flag
	var options = NewOptions()
	return options.Parse()
}

// Parse is parser for cmd
func (c *CmdParams) Parse() *CmdParams {
	args, err := c.LoadArgs()
	if err != nil {
		log.Fatalf("parse cli args fail %s", err.Error())
		return nil
	}
	if args.YAMLPath == "" {
		c.args = args
		return c.argsParse(args)
	}
	//use yml config
	if cx := NewFromYaml(args.YAMLPath); cx != nil {
		cx.args = args
		return cx
	}
	return c
}

func (c *CmdParams) argsParse(args *Options) *CmdParams {
	// cmd first
	if args.DSN != "" {
		c.DSN = args.DSN
	}
	if args.DB != "" {
		c.DB = args.DB
	}
	if args.TableList != "" {
		c.Tables = strings.Split(args.TableList, ",")
	}
	if args.ExcludeTableList != "" {
		c.ExcludeTableList = strings.Split(args.ExcludeTableList, ",")
	}
	if args.OnlyModel {
		c.OnlyModel = args.OnlyModel
	}
	if args.OutPath != "" {
		c.OutPath = args.OutPath
	}
	if args.OutFile != "" {
		c.OutFile = args.OutFile
	}
	if args.WithUnitTest {
		c.WithUnitTest = args.WithUnitTest
	}
	if args.ModelPkgName != "" {
		c.ModelPkgName = args.ModelPkgName
	}
	if args.FieldNullable {
		c.FieldNullable = args.FieldNullable
	}
	if args.FieldCoverable {
		c.FieldCoverable = args.FieldCoverable
	}
	if args.FieldWithIndexTag {
		c.FieldWithIndexTag = args.FieldWithIndexTag
	}
	if args.FieldSignable {
		c.FieldSignable = args.FieldSignable
	}
	if args.FieldWithTypeTag {
		c.FieldWithTypeTag = args.FieldWithTypeTag
	}
	if args.FieldJSONTypeTag {
		c.FieldJSONTypeTag = args.FieldJSONTypeTag
	}
	if args.DefaultYAMLConfigFile != "" {
		c.defaultYAMLConfigFile = args.DefaultYAMLConfigFile
	}
	if len(args.FieldsTypeMapping) > 0 {
		c.FieldsTypeMapping = args.FieldsTypeMapping
	}
	if len(args.ImportPkgPaths) > 0 {
		c.ImportPkgPaths = args.ImportPkgPaths
	}
	return c
}

func (c *CmdParams) GetDBType() DBType {
	if c.DB == "" {
		return DbMySQL
	}
	return DBType(c.DB)
}

func (c *CmdParams) GetMode() gen.GenerateMode {
	if ms := strings.Split(c.Mode, "|"); len(ms) > 0 {
		var v gen.GenerateMode = 0
		for _, m := range ms {
			switch m {
			case "DefaultQuery":
				v |= gen.WithDefaultQuery
			case "QueryInterface":
				v |= gen.WithQueryInterface
			case "OutContext":
				v |= gen.WithoutContext
			}
		}
		return v
	}
	return gen.WithoutContext | gen.WithQueryInterface | gen.WithDefaultQuery
}

func (c *CmdParams) GetImportPkgPaths() []string {
	if len(c.ImportPkgPaths) > 0 {
		return c.ImportPkgPaths
	}
	return []string{}
}

func (c *CmdParams) GetTypeMappings() map[string]func(columnType gorm.ColumnType) (dataType string) {
	var mappings = make(map[string]func(columnType gorm.ColumnType) (dataType string))
	if c.FieldsTypeMapping != nil {
		for _, v := range c.FieldsTypeMapping {
			offset := strings.Index(v, ":")
			if offset <= 0 {
				continue
			}
			name := v[:offset]
			value := v[offset+1:]
			if name == "" || value == "" {
				continue
			}
			if mappings[name] != nil {
				log.Printf("duplicate type mapping name %s \n", name)
				continue
			}
			mappings[name] = createTypeMapping(name, value)
		}
	}
	if c.FieldJSONTypeTag {
		mappings["jsonb"] = createTypeMapping("jsonb", "datatypes.JSON")
	}
	return mappings
}

func (c *CmdParams) GetModelOptions() []gen.ModelOpt {
	return []gen.ModelOpt{
		gen.FieldGORMTagReg("*", nullFieldForGo),
		gen.FieldRegexCommentReplace(`\{\{*\}\}`, replaceComment),
	}
}

func (c *CmdParams) IsHelp() bool {
	return c.args != nil && c.args.GetHelpMsg()
}

func (c *CmdParams) GetGenDefaultYAMLFile() string {
	return c.defaultYAMLConfigFile
}

func (c *CmdParams) withDefault() *CmdParams {
	if c.Mode == "" {
		c.Mode = "DefaultQuery|QueryInterface|OutContext"
	}
	if len(c.ImportPkgPaths) <= 0 {
		c.ImportPkgPaths = []string{"gorm.io/datatypes"}
	}
	if c.DSN == "" {
		if c.DB == "" {
			c.DSN = defaultMysqlDSN
		} else {
			switch strings.ToLower(c.DB) {
			case DbMySQL.String():
				c.DSN = defaultMysqlDSN
			case DbPostgres.String():
				c.DSN = defaultPostgresDSN
			case DbSQLite.String():
				c.DSN = defaultSQLiteDSN
			case DbSQLServer.String():
				c.DSN = defaultSQLServerDSN
			case DbClickHouse.String():
				c.DSN = defaultClickHouseDSN
			default:
				c.DSN = defaultMysqlDSN
			}
		}
	}
	if c.DB == "" {
		c.DB = extractDSNDBType(c.DSN)
	}
	if c.OutPath == "" {
		c.OutPath = "./models"
	}
	if len(c.ExcludeTableList) <= 0 {
		c.ExcludeTableList = []string{"ignore_table_name1", "ignore_table_name2"}
	}
	if c.ModelPkgName == "" {
		c.ModelPkgName = "models"
	}
	if len(c.Tables) <= 0 {
		c.Tables = []string{"test", "user"}
	}
	// c.args.PrintArgsValues()
	if len(c.FieldsTypeMapping) <= 0 {
		c.FieldsTypeMapping = []string{"jsonb:datatypes.JSON"}
	}
	return c
}

func createTypeMapping(typ, value string) func(columnType gorm.ColumnType) (dataType string) {
	return func(columnType gorm.ColumnType) (dataType string) {
		vs, _ := columnType.ColumnType()
		if columnType.DatabaseTypeName() == typ || vs == typ {
			return value
		}
		return "string"
	}
}

func extractDSNDBType(dsn string) string {
	if strings.Contains(dsn, "://") {
		schema := strings.ToLower(strings.Split(dsn, "://")[0])
		switch schema {
		case DbMySQL.String():
			return DbMySQL.String()
		case DbPostgres.String():
			return DbPostgres.String()
		case DbSQLite.String(), `file`:
			return DbSQLite.String()
		case DbSQLServer.String():
			return DbSQLServer.String()
		case DbClickHouse.String(), `tcp`:
			return DbClickHouse.String()
		default:
			return schema
		}
	}
	if strings.HasPrefix(dsn, "file:") {
		return DbSQLite.String()
	}
	if strings.Contains(dsn, "search_path=") {
		return DbPostgres.String()
	}
	if strings.Contains(dsn, " ") && strings.Contains(dsn, "=") {
		return DbPostgres.String()
	}
	return DbMySQL.String()
}

func nullFieldForGo(tag field.GormTag) field.GormTag {
	for key, values := range tag {
		log.Printf("tag=%s,vlaues=%+v", key, values)
	}
	return tag
}

func replaceComment(comment string) string {
	if strings.Contains(comment, "{{") && strings.Contains(comment, "}}") {
		return strings.ReplaceAll(strings.ReplaceAll(comment, "{{", "{"), "}}", "}")
	}
	return comment
}

func SaveYAMLConfigFile(params *CmdParams, saveFile string) (string, error) {
	var data, err = yaml.Marshal(params.withDefault().Revise())
	if err != nil {
		return saveFile, err
	}
	ext := strings.ToLower(filepath.Ext(saveFile))
	if ext == "" {
		var state os.FileInfo
		if state, err = os.Stat(saveFile); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return saveFile, err
			}
			_ = os.MkdirAll(saveFile, os.ModePerm)
		}
		if state != nil && state.IsDir() {
			ext = ".yaml"
			saveFile = fmt.Sprintf("%s/config.yaml", saveFile)
		}
	}
	if ext != ".yaml" && ext != ".yml" {
		return saveFile, errors.New("save file must be  yaml or yml")
	}
	if err = os.WriteFile(saveFile, data, 0644); err != nil {
		return saveFile, err
	}
	return saveFile, nil
}

func schemaPostgresToValues(dsn string) string {
	var (
		values  []string
		indexes = map[string]struct{}{}
	)
	if strings.Contains(dsn, "://") {
		uri, err := url.Parse(dsn)
		if err != nil {
			log.Printf("parse dsn fail:%s", err.Error())
			return dsn
		}
		values = append(values, fmt.Sprintf("host=%s", uri.Hostname()))
		indexes["host"] = struct{}{}
		if p := uri.Port(); p != "" {
			values = append(values, fmt.Sprintf("port=%s", p))
		} else {
			values = append(values, "port=5432")
		}
		indexes["port"] = struct{}{}
		if uri.User != nil {
			values = append(values, fmt.Sprintf("user=%s", uri.User.Username()))
			indexes["user"] = struct{}{}
			if p, ok := uri.User.Password(); ok {
				values = append(values, fmt.Sprintf("password=%s", p))
				indexes["password"] = struct{}{}
			}
		}
		values = append(values, fmt.Sprintf("dbname=%s", uri.Path[1:]))
		indexes["dbname"] = struct{}{}
		for k, v := range uri.Query() {
			if _, ok := indexes[k]; ok {
				continue
			}
			indexes[k] = struct{}{}
			values = append(values, fmt.Sprintf("%s=%s", k, v[0]))
		}
	}
	if len(values) > 0 {
		return strings.Join(values, " ")
	}
	return dsn
}
