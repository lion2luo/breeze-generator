package core

import (
	"errors"
	"os"
	"strings"
)

//breeze type for generate code.
const (
	Bool = iota
	String
	Byte
	Bytes
	Int16
	Int32
	Int64
	Float32
	Float64
	Map
	Array
	Msg
)

//option keys
const (
	JavaPackage            = "java_package"
	WithPackageDir         = "with_package_dir"
	Alias                  = "alias"
	PackageVersion         = "package_version"
	WithMotanConfiguration = "with_motan_configuration"
)
const (
	MotanPackageName            = "motan_package_name"
	MotanRegistryProtocol       = "motan_registry_protocol"
	MotanRegistryHost           = "motan_registry_host"
	MotanRegistryPort           = "motan_registry_port"
	MotanClientAccesslog        = "motan_client_accesslog"
	MotanClientRequestTimeout   = "motan_client_request_timeout"
	MotanClientCheck            = "motan_client_check"
	MotanClientThrowException   = "motan_client_throw_exception"
	MotanClientRetries          = "motan_client_retries"
	MotanClientMinConnections   = "motan_client_min_connections"
	MotanClientMaxConnections   = "motan_client_max_connections"
	MotanClientLoadbalance      = "motan_client_loadbalance"
	MotanClientHa               = "motan_client_ha"
	MotanClientUsegz            = "motan_client_usegz"
	MotanClientMingzSize        = "motan_client_mingz_size"
	MotanServerAccesslog        = "motan_server_accesslog"
	MotanServerPort             = "motan_server_port"
	MotanServerUsegz            = "motan_server_usegz"
	MotanServerMingzSize        = "motan_server_mingz_size"
	MotanServerMinWorkerThreads = "motan_server_min_worker_threads"
	MotanServerMaxWorkerThreads = "motan_server_max_worker_threads"
	MotanServerMaxConnections   = "motan_server_max_connections"
)

var MotanOptionsDefault = map[string]string{
	MotanPackageName:            "",
	MotanRegistryProtocol:       "",
	MotanRegistryHost:           "",
	MotanRegistryPort:           "",
	MotanClientAccesslog:        "false",
	MotanClientRequestTimeout:   "200",
	MotanClientCheck:            "true",
	MotanClientThrowException:   "true",
	MotanClientRetries:          "0",
	MotanClientMinConnections:   "2",
	MotanClientMaxConnections:   "10",
	MotanClientLoadbalance:      "roundrobin",
	MotanClientHa:               "failover",
	MotanClientUsegz:            "false",
	MotanClientMingzSize:        "1000",
	MotanServerAccesslog:        "false",
	MotanServerPort:             "8100",
	MotanServerUsegz:            "false",
	MotanServerMingzSize:        "1000",
	MotanServerMinWorkerThreads: "20",
	MotanServerMaxWorkerThreads: "200",
	MotanServerMaxConnections:   "100000",
}

const (
	PathSeparator = string(os.PathSeparator)
)

const (
	DefaultNewDirectoryMode   os.FileMode = 0755
	DefaultNewRegularFileMode os.FileMode = 0644
)

//rpc type
const (
	Server          = "server"
	ClientWithAgent = "clientWithAgent" //motan agent client
)

//primitive types
var (
	BoolType    = &Type{Number: Bool, TypeString: "bool"}
	StringType  = &Type{Number: String, TypeString: "string"}
	ByteType    = &Type{Number: Byte, TypeString: "byte"}
	BytesType   = &Type{Number: Bytes, TypeString: "bytes"}
	Int16Type   = &Type{Number: Int16, TypeString: "int16"}
	Int32Type   = &Type{Number: Int32, TypeString: "int32"}
	Int64Type   = &Type{Number: Int64, TypeString: "int64"}
	Float32Type = &Type{Number: Float32, TypeString: "float32"}
	Float64Type = &Type{Number: Float64, TypeString: "float64"}
)

//Parser can parse breeze schema from binary with context.
type Parser interface {
	ParseSchema(content []byte, context *Context) (schema *Schema, err error)
	Name() string
	FileSuffix() string
}

//CodeTemplate is code template for generating code of different languages
type CodeTemplate interface {
	GenerateCode(schema *Schema, context *Context) (contents map[string][]byte, err error)
	PostAllGenerated(context *Context) error
	Name() string
}

//Schema describe a breeze message.
type Schema struct {
	Name       string // file name
	Package    string // file package
	OrgPackage string // schema name package.
	Options    map[string]string
	Messages   map[string]*Message
	Services   map[string]*Service
}

//Message :breeze message. include enum message
type Message struct {
	Name       string
	Alias      string
	Options    map[string]string
	Fields     map[int]*Field
	IsEnum     bool
	EnumValues map[int]string
}

//Field is a breeze message field.
type Field struct {
	Index int
	Name  string
	Type  *Type
}

//Type : message field type
type Type struct {
	Name       string
	Number     int    //get type number. this number is only used for generating code , not for serialize
	KeyType    *Type  //get map key type
	ValueType  *Type  //get map value type or array value type
	TypeString string //get raw type string
}

//Service describe a rpc service, which request and response are breeze messages
type Service struct {
	Name    string
	Options map[string]string
	Methods map[string]*Method
}

//Method : rpc method
type Method struct {
	Name   string
	Params map[int]*Param
	Return string
}

//Param : method param
type Param struct {
	Type string
	Name string
}

//Context : generate context
type Context struct {
	WritePath string
	Parser    Parser
	RPCType   string
	Templates []CodeTemplate
	Schemas   map[string]*Schema
	Messages  map[string]*Message
	Options   map[string]string
}

//GetType : get a Type from type string
func GetType(typeString string, removePackage bool) (*Type, error) {
	typeString = strings.TrimSpace(typeString)
	if typeString == "" {
		return nil, errors.New("type is empty")
	}
	switch typeString {
	case "bool":
		return BoolType, nil
	case "string":
		return StringType, nil
	case "byte":
		return ByteType, nil
	case "bytes":
		return BytesType, nil
	case "int16":
		return Int16Type, nil
	case "int", "int32":
		return Int32Type, nil
	case "int64":
		return Int64Type, nil
	case "float32":
		return Float32Type, nil
	case "float64":
		return Float64Type, nil
	}
	if strings.HasPrefix(typeString, "map<") && strings.HasSuffix(typeString, ">") {
		//only primitive type can be a map key!
		inner := typeString[4 : len(typeString)-1]
		index := strings.Index(inner, ",")
		key := strings.TrimSpace(inner[:index])
		keyType, err := GetType(key, removePackage)
		if err != nil {
			return nil, err
		}
		if keyType.Number > Float64 {
			return nil, errors.New("wrong map key type: " + typeString)
		}
		valueType, err := GetType(strings.TrimSpace(inner[index+1:]), removePackage)
		if err != nil {
			return nil, err
		}
		return &Type{Number: Map, TypeString: typeString, KeyType: keyType, ValueType: valueType}, nil
	}
	if strings.HasPrefix(typeString, "array<") && strings.HasSuffix(typeString, ">") {
		vType, err := GetType(typeString[6:len(typeString)-1], removePackage)
		if err != nil {
			return nil, err
		}
		return &Type{Number: Array, TypeString: typeString, ValueType: vType}, nil
	}
	//message
	if removePackage && strings.Index(typeString, ".") > -1 {
		typeString = typeString[strings.LastIndex(typeString, ".")+1:]
	}
	return &Type{Number: Msg, Name: typeString, TypeString: typeString}, nil
}

//Validate : check schema
func Validate(schema *Schema) error {
	if schema == nil || (len(schema.Messages) == 0 && len(schema.Services) == 0) {
		return errors.New("schema is empty. schema:" + schema.Name)
	}
	return nil
}