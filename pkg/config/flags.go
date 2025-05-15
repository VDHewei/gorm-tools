package config

import (
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

type Options struct {
	YAMLPath              string   `env:"GEN_CONFIG" json:"config" long:"config" short:"c" description:"is path for gen.yml"`
	DSN                   string   `env:"GEN_DSN" json:"dsn" long:"dsn" description:"consult[https://gorm.io/docs/connecting_to_the_database.html]"`
	DB                    string   `env:"GEN_DB" json:"db" long:"db" description:"input mysql|postgres|sqlite|sqlserver|clickhouse. consult[https://gorm.io/docs/connecting_to_the_database.html]"`
	TableList             string   `env:"GEN_TABLES" json:"tables" long:"tables" short:"t" description:"enter the required data table or leave it blank"`
	ExcludeTableList      string   `env:"GEN_EXCLUDE_TABLES" json:"exclude_tables" long:"excludeTables" short:"e" description:"enter the exclude data table or leave it blank"`
	OnlyModel             *bool    `env:"GEN_ONLY_MODEL" json:"onlyModel" long:"onlyModel" description:"only generate models (without query file)"`
	OutPath               string   `env:"GEN_OUT_PATH" json:"outPath" long:"outPath" description:"specify a directory for output"`
	OutFile               string   `env:"GEN_OUTFILE" json:"outFile" long:"outFile" description:"query code file name, default: gen.go" default:"gen.go"`
	Mode                  string   `env:"GEN_MODE" json:"mode" long:"mode" description:"input DefaultQuery|QueryInterface|OutContext. gen mode setting"`
	WithUnitTest          *bool    `env:"GEN_WITH_UNITTEST" json:"withUnitTest" long:"withUnitTest" description:"generate unit test for query code"`
	ModelPkgName          string   `env:"GEN_MODEL_PKG_NAME" json:"modelPkgName" long:"modelPkgName" description:"generated model code's package name"`
	FieldNullable         *bool    `env:"GEN_FIELD_NULLABLE" json:"fieldNullable" long:"fieldNullable" description:"generate with pointer when field is nullable"`
	FieldCoverable        *bool    `env:"GEN_FIELD_COVERABLE" json:"fieldCoverable" long:"fieldCoverable" description:"generate with pointer when field has default value"`
	FieldWithIndexTag     *bool    `env:"GEN_FIELD_WITH_INDEX_TAG" json:"fieldWithIndexTag" long:"fieldWithIndexTag" description:"generate field with gorm index tag"`
	FieldWithTypeTag      *bool    `env:"GEN_FIELD_WITH_TYPE_TAG" json:"fieldWithTypeTag" long:"fieldWithTypeTag" description:"generate field with gorm column type tag"`
	FieldSignable         *bool    `env:"GEN_FIELD_SIGNABLE" json:"fieldSignable" long:"fieldSignable" description:"detect integer field's unsigned type, adjust generated data type"`
	FieldJSONTypeTag      *bool    `env:"GEN_FIELD_JSON_TYPE_TAG" json:"fieldJSONTypeTag" long:"fieldJSONTypeTag" description:"generate field with gorm json type"`
	FieldsTypeMapping     []string `env:"GEN_FIELDS_TYPE_MAPPING" json:"fieldsTypeMapping" long:"fieldsTypeMapping" short:"m" description:"mapping field type mapping ,eg: jsonb:datatypes.JSON"`
	ImportPkgPaths        []string `env:"GEN_IMPORT_PKG_PATHS" json:"importPkgPaths" long:"importPkgPaths" short:"p" description:"generate code import package path,eg: github.com/xxx/xxx"`
	DefaultYAMLConfigFile string   `env:"GEN_Config_File" json:"defaultYAMLConfigFile" long:"defaultYAMLConfigFile" short:"d" description:"generate default yaml config file"`
	V                     *bool    `json:"version" long:"version" short:"v" description:"print tool version"`
	helpMsg               bool
	rowValues             []string
}

func NewOptions() *Options {
	return &Options{}
}

func (f *Options) Parse() (*Options, error) {
	_, err := flags.Parse(f)
	f.rowValues = os.Args[1:]
	if err != nil {
		if errors.Is(err.(*flags.Error).Type, flags.ErrRequired) &&
			f.DefaultYAMLConfigFile != "" {
			return f, nil
		}
		if errors.Is(err.(*flags.Error).Type, flags.ErrHelp) {
			f.helpMsg = true
			return f, nil
		}
		return nil, err
	}
	return f, nil
}

func (f *Options) PrintArgsValues() {
	fmt.Printf("args=: %+v\n", f.rowValues)
}

func (f *Options) GetRawArgs() []string {
	return f.rowValues
}

func (f *Options) GetHelpMsg() bool {
	return f.helpMsg
}
