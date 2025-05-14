package core

import (
	"fmt"
	"github.com/VDHewei/gorm-tools/pkg/config"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gen"
	"gorm.io/gorm"
	"log"
	"os"
)

type (
	GenTools struct {
		db     *gorm.DB
		models []interface{}
		g      *gen.Generator
		params *config.CmdParams
	}
	Option func(*GenTools)
)

// ConnectDB choose db type for connection to database
func ConnectDB(t config.DBType, dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn cannot be empty")
	}
	switch t {
	case config.DbMySQL:
		return gorm.Open(mysql.Open(dsn))
	case config.DbPostgres:
		return gorm.Open(postgres.Open(dsn))
	case config.DbSQLite:
		return gorm.Open(sqlite.Open(dsn))
	case config.DbSQLServer:
		return gorm.Open(sqlserver.Open(dsn))
	case config.DbClickHouse:
		return gorm.Open(clickhouse.Open(dsn))
	default:
		return nil, fmt.Errorf("unknow db %q (support mysql || postgres || sqlite || sqlserver for now)", t)
	}
}

// GenModels is gorm/gen generated models
func GenModels(g *gen.Generator, db *gorm.DB, tables []string, opts ...gen.ModelOpt) (models []interface{}, err error) {
	if len(tables) == 0 {
		// Execute tasks for all tables in the database
		tables, err = db.Migrator().GetTables()
		if err != nil {
			return nil, fmt.Errorf("GORM migrator get all tables fail: %w", err)
		}
	}

	// Execute some data table tasks
	models = make([]interface{}, len(tables))
	for i, tableName := range tables {
		meta := g.GenerateModel(tableName, opts...)
		models[i] = meta
	}
	return models, nil
}

func (g *GenTools) OpenDB() error {
	var err error
	g.params.Revise()
	if g.db, err = ConnectDB(g.params.GetDBType(), g.params.DSN); err != nil {
		return err
	}
	return nil
}

func (g *GenTools) GetDB() *gorm.DB {
	if g.db == nil {
		if err := g.OpenDB(); err != nil {
			log.Fatalln("open db fail:", err)
			return nil
		}
	}
	return g.db
}

func (g *GenTools) GenModels() error {
	var (
		err    error
		tables = g.params.Tables
		opts   = g.params.GetModelOptions()
	)
	if g.models, err = GenModels(g.g, g.GetDB(), tables, opts...); err != nil {
		return err
	}
	return nil
}

func (g *GenTools) GetModels() []interface{} {
	if g.models == nil {
		if err := g.GenModels(); err != nil {
			log.Fatalln("get models fail:", err)
			return nil
		}
	}
	return g.models
}

func (g *GenTools) RegisterModels(models ...interface{}) {
	for _, m := range models {
		g.models = append(g.models, m)
	}
}

func (g *GenTools) Execute() {
	g.g.UseDB(g.GetDB())
	if g.params.OnlyModel {
		if err := g.GenModels(); err != nil {
			log.Fatalln("gen models fail:", err)
			os.Exit(-1)
		}
	} else {
		g.g.ApplyBasic(g.GetModels()...)
	}
	g.g.Execute()
}

func (g *GenTools) LoadConfig() gen.Config {
	var c = gen.Config{
		OutPath:           g.params.OutPath,
		OutFile:           g.params.OutFile,
		ModelPkgPath:      g.params.ModelPkgName,
		Mode:              g.params.GetMode(),
		WithUnitTest:      g.params.WithUnitTest,
		FieldNullable:     g.params.FieldNullable,
		FieldCoverable:    g.params.FieldCoverable,
		FieldSignable:     g.params.FieldSignable,
		FieldWithIndexTag: g.params.FieldWithIndexTag,
		FieldWithTypeTag:  g.params.FieldWithTypeTag,
	}
	// mappings
	if m := g.params.GetTypeMappings(); len(m) > 0 {
		c.WithDataTypeMap(m)
	}
	// "gorm.io/datatypes"
	if paths := g.params.GetImportPkgPaths(); len(paths) > 0 {
		c.WithImportPkgPath(paths...)
	}
	return c
}

func WithDB(db *gorm.DB) Option {
	return func(g *GenTools) {
		g.g.UseDB(db)
	}
}

func WithGen(g *gen.Generator) Option {
	return func(tools *GenTools) {
		tools.g = g
	}
}

func WithConfig(c *config.CmdParams) Option {
	return func(tools *GenTools) {
		tools.params = c
	}
}

func New(opts ...Option) *GenTools {
	var ins = &GenTools{}
	for _, o := range opts {
		o(ins)
	}
	if ins.params == nil {
		ins.params = config.New().Parse()
	}
	if ins.g == nil {
		ins.g = gen.NewGenerator(ins.LoadConfig())
	}
	return ins
}
