// Copyright 2019 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlserver

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/dolthub/go-mysql-server/sql"

	"github.com/dolthub/dolt/go/cmd/dolt/commands/engine"
	"github.com/dolthub/dolt/go/libraries/doltcore/sqle/cluster"
	"github.com/dolthub/dolt/go/libraries/utils/argparser"
	"github.com/dolthub/dolt/go/libraries/utils/filesys"
)

// LogLevel defines the available levels of logging for the server.
type LogLevel string

const (
	LogLevel_Trace   LogLevel = "trace"
	LogLevel_Debug   LogLevel = "debug"
	LogLevel_Info    LogLevel = "info"
	LogLevel_Warning LogLevel = "warning"
	LogLevel_Error   LogLevel = "error"
	LogLevel_Fatal   LogLevel = "fatal"
)

const (
	defaultHost                    = "localhost"
	defaultPort                    = 3306
	defaultUser                    = "root"
	defaultPass                    = ""
	defaultTimeout                 = 8 * 60 * 60 * 1000 // 8 hours, same as MySQL
	defaultReadOnly                = false
	defaultLogLevel                = LogLevel_Info
	defaultAutoCommit              = true
	defaultDoltTransactionCommit   = false
	defaultMaxConnections          = 100
	defaultQueryParallelism        = 0
	defaultPersistenceBahavior     = loadPerisistentGlobals
	defaultDataDir                 = "."
	defaultCfgDir                  = ".doltcfg"
	defaultPrivilegeFilePath       = "privileges.db"
	defaultBranchControlFilePath   = "branch_control.db"
	defaultMetricsHost             = ""
	defaultMetricsPort             = -1
	defaultAllowCleartextPasswords = false
	defaultUnixSocketFilePath      = "/tmp/mysql.sock"
	defaultMaxLoggedQueryLen       = 0
	defaultEncodeLoggedQuery       = false
)

const (
	ignorePeristentGlobals = "ignore"
	loadPerisistentGlobals = "load"
)

// String returns the string representation of the log level.
func (level LogLevel) String() string {
	switch level {
	case LogLevel_Trace:
		fallthrough
	case LogLevel_Debug:
		fallthrough
	case LogLevel_Info:
		fallthrough
	case LogLevel_Warning:
		fallthrough
	case LogLevel_Error:
		fallthrough
	case LogLevel_Fatal:
		return string(level)
	default:
		return "unknown"
	}
}

// ServerConfig contains all of the configurable options for the MySQL-compatible server.
type ServerConfig interface {
	// Host returns the domain that the server will run on. Accepts an IPv4 or IPv6 address, in addition to localhost.
	Host() string
	// Port returns the port that the server will run on. The valid range is [1024, 65535].
	Port() int
	// User returns the username that connecting clients must use.
	User() string
	// Password returns the password that connecting clients must use.
	Password() string
	// ReadTimeout returns the read timeout in milliseconds
	ReadTimeout() uint64
	// WriteTimeout returns the write timeout in milliseconds
	WriteTimeout() uint64
	// ReadOnly returns whether the server will only accept read statements or all statements.
	ReadOnly() bool
	// LogLevel returns the level of logging that the server will use.
	LogLevel() LogLevel
	// Autocommit defines the value of the @@autocommit session variable used on every connection
	AutoCommit() bool
	// DoltTransactionCommit defines the value of the @@dolt_transaction_commit session variable that enables Dolt
	// commits to be automatically created when a SQL transaction is committed.
	DoltTransactionCommit() bool
	// DataDir is the path to a directory to use as the data dir, both to create new databases and locate existing ones.
	DataDir() string
	// CfgDir is the path to a directory to use to store the dolt configuration files.
	CfgDir() string
	// MaxConnections returns the maximum number of simultaneous connections the server will allow.  The default is 1
	MaxConnections() uint64
	// QueryParallelism returns the parallelism that should be used by the go-mysql-server analyzer
	QueryParallelism() int
	// TLSKey returns a path to the servers PEM-encoded private TLS key. "" if there is none.
	TLSKey() string
	// TLSCert returns a path to the servers PEM-encoded TLS certificate chain. "" if there is none.
	TLSCert() string
	// RequireSecureTransport is true if the server should reject non-TLS connections.
	RequireSecureTransport() bool
	// MaxLoggedQueryLen is the max length of queries written to the logs.  Queries longer than this number are truncated.
	// If this value is 0 then the query is not truncated and will be written to the logs in its entirety.  If the value
	// is less than 0 then the queries will be omitted from the logs completely
	MaxLoggedQueryLen() int
	// ShouldEncodeLoggedQuery determines if logged queries are base64 encoded.
	// If true, queries will be logged as base64 encoded strings.
	// If false (default behavior), queries will be logged as strings, but newlines and tabs will be replaced with spaces.
	ShouldEncodeLoggedQuery() bool
	// PersistenceBehavior is "load" if we include persisted system globals on server init
	PersistenceBehavior() string
	// DisableClientMultiStatements is true if we want the server to not
	// process incoming ComQuery packets as if they had multiple queries in
	// them, even if the client advertises support for MULTI_STATEMENTS.
	DisableClientMultiStatements() bool
	// MetricsLabels returns labels that are applied to all prometheus metrics
	MetricsLabels() map[string]string
	MetricsHost() string
	MetricsPort() int
	// PrivilegeFilePath returns the path to the file which contains all needed privilege information in the form of a
	// JSON string.
	PrivilegeFilePath() string
	// BranchControlFilePath returns the path to the file which contains the branch control permissions.
	BranchControlFilePath() string
	// UserVars is an array containing user specific session variables
	UserVars() []UserSessionVars
	// SystemVars is a map setting global SQL system variables. For example, `secure_file_priv`.
	SystemVars() engine.SystemVariables
	// JwksConfig is an array containing jwks config
	JwksConfig() []engine.JwksConfig
	// AllowCleartextPasswords is true if the server should accept cleartext passwords.
	AllowCleartextPasswords() bool
	// Socket is a path to the unix socket file
	Socket() string
	// RemotesapiPort is the port to use for serving a remotesapi interface with this sql-server instance.
	// A remotesapi interface will allow this sql-server process to be used
	// as a dolt remote for things like `clone`, `fetch` and read
	// replication.
	RemotesapiPort() *int
	// RemotesapiReadOnly is true if the remotesapi interface should be read only.
	RemotesapiReadOnly() *bool
	// ClusterConfig is the configuration for clustering in this sql-server.
	ClusterConfig() cluster.Config
	// EventSchedulerStatus is the configuration for enabling or disabling the event scheduler in this server.
	EventSchedulerStatus() string
	// ValueSet returns whether the value string provided was explicitly set in the config
	ValueSet(value string) bool
}

// WritableServerConfig is a ServerConfig that support overwriting certain values.
type WritableServerConfig interface {
	ServerConfig
	// SetUserName sets the username for servers with no other auth established
	SetUserName(string)
	// SetPassword sets the password for servers with no other auth established
	SetPassword(string)
}

type validatingServerConfig interface {
	ServerConfig
	// goldenMysqlConnectionString returns a connection string for a mysql
	// instance that can be used to validate query results
	goldenMysqlConnectionString() string
}

// ServerConfigReader is an interface for reading a ServerConfig from a file or command line arguments.
type ServerConfigReader interface {
	// ReadConfigFile reads a config file and returns a ServerConfig for it
	ReadConfigFile(cwdFS filesys.Filesys, file string) (ServerConfig, error)
	// ReadConfigArgs reads command line arguments and returns a ServerConfig for them
	ReadConfigArgs(args *argparser.ArgParseResults) (ServerConfig, error)
}

// ValidateConfig returns an `error` if any field is not valid.
func ValidateConfig(config ServerConfig) error {
	if config.Host() != "localhost" {
		ip := net.ParseIP(config.Host())
		if ip == nil {
			return fmt.Errorf("address is not a valid IP: %v", config.Host())
		}
	}
	if config.Port() < 1024 || config.Port() > 65535 {
		return fmt.Errorf("port is not in the range between 1024-65535: %v\n", config.Port())
	}
	if config.LogLevel().String() == "unknown" {
		return fmt.Errorf("loglevel is invalid: %v\n", string(config.LogLevel()))
	}
	if config.RequireSecureTransport() && config.TLSCert() == "" && config.TLSKey() == "" {
		return fmt.Errorf("require_secure_transport can only be `true` when a tls_key and tls_cert are provided.")
	}
	return ValidateClusterConfig(config.ClusterConfig())
}

const (
	maxConnectionsKey = "max_connections"
	readTimeoutKey    = "net_read_timeout"
	writeTimeoutKey   = "net_write_timeout"
	eventSchedulerKey = "event_scheduler"
)

// ApplySystemVariables sets the global system variables based on the given `ServerConfig`.
func ApplySystemVariables(cfg ServerConfig) error {
	if cfg.ValueSet(maxConnectionsKey) {
		err := sql.SystemVariables.SetGlobal("max_connections", cfg.MaxConnections())
		if err != nil {
			return err
		}
	}

	if cfg.ValueSet(readTimeoutKey) {
		err := sql.SystemVariables.SetGlobal("net_read_timeout", cfg.ReadTimeout())
		if err != nil {
			return err
		}
	}

	if cfg.ValueSet(writeTimeoutKey) {
		err := sql.SystemVariables.SetGlobal("net_write_timeout", cfg.WriteTimeout())
		if err != nil {
			return err
		}
	}

	if cfg.ValueSet(eventSchedulerKey) {
		err := sql.SystemVariables.SetGlobal("event_scheduler", cfg.EventSchedulerStatus())
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateClusterConfig(config cluster.Config) error {
	if config == nil {
		return nil
	}
	remotes := config.StandbyRemotes()
	if len(remotes) == 0 {
		return errors.New("cluster config: must supply standby_remotes when supplying cluster configuration.")
	}
	for i := range remotes {
		if remotes[i].Name() == "" {
			return fmt.Errorf("cluster: standby_remotes[%d]: name: Cannot be empty", i)
		}
		if strings.Index(remotes[i].RemoteURLTemplate(), "{database}") == -1 {
			return fmt.Errorf("cluster: standby_remotes[%d]: remote_url_template: is \"%s\" but must include the {database} template parameter", i, remotes[i].RemoteURLTemplate())
		}
	}
	if config.BootstrapRole() != "" && config.BootstrapRole() != "primary" && config.BootstrapRole() != "standby" {
		return fmt.Errorf("cluster: boostrap_role: is \"%s\" but must be \"primary\" or \"standby\"", config.BootstrapRole())
	}
	if config.BootstrapEpoch() < 0 {
		return fmt.Errorf("cluster: boostrap_epoch: is %d but must be >= 0", config.BootstrapEpoch())
	}
	if config.RemotesAPIConfig().Port() < 0 || config.RemotesAPIConfig().Port() > 65535 {
		return fmt.Errorf("cluster: remotesapi: port: is not in range 0-65535: %d", config.RemotesAPIConfig().Port())
	}
	if config.RemotesAPIConfig().TLSKey() == "" && config.RemotesAPIConfig().TLSCert() != "" {
		return fmt.Errorf("cluster: remotesapi: tls_key: must supply a tls_key if you supply a tls_cert")
	}
	if config.RemotesAPIConfig().TLSKey() != "" && config.RemotesAPIConfig().TLSCert() == "" {
		return fmt.Errorf("cluster: remotesapi: tls_cert: must supply a tls_cert if you supply a tls_key")
	}
	return nil
}

// ConnectionString returns a Data Source Name (DSN) to be used by go clients for connecting to a running server.
// If unix socket file path is defined in ServerConfig, then `unix` DSN will be returned.
func ConnectionString(config ServerConfig, database string) string {
	user := config.User()
	if user == "" {
		user = "root"
	}
	var dsn string
	if config.Socket() != "" {
		dsn = fmt.Sprintf("%v:%v@unix(%v)/%v", user, config.Password(), config.Socket(), database)
	} else {
		dsn = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", user, config.Password(), config.Host(), config.Port(), database)
	}
	if config.AllowCleartextPasswords() {
		dsn += "?allowCleartextPasswords=1"
	}
	return dsn
}

// ConfigInfo returns a summary of some of the config which contains some of the more important information
func ConfigInfo(config ServerConfig) string {
	socket := ""
	sock, useSock, err := checkForUnixSocket(config)
	if err != nil {
		panic(err)
	}
	if useSock {
		socket = fmt.Sprintf(`|S="%v"`, sock)
	}
	return fmt.Sprintf(`HP="%v:%v"|T="%v"|R="%v"|L="%v"%s`, config.Host(), config.Port(),
		config.ReadTimeout(), config.ReadOnly(), config.LogLevel(), socket)
}

// LoadTLSConfig loads the certificate chain from config.TLSKey() and config.TLSCert() and returns
// a *tls.Config configured for its use. Returns `nil` if key and cert are `""`.
func LoadTLSConfig(cfg ServerConfig) (*tls.Config, error) {
	if cfg.TLSKey() == "" && cfg.TLSCert() == "" {
		return nil, nil
	}
	c, err := tls.LoadX509KeyPair(cfg.TLSCert(), cfg.TLSKey())
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{
			c,
		},
	}, nil
}

// DoltServerConfigReader is the default implementation of ServerConfigReader suitable for parsing Dolt config files
// and command line options.
type DoltServerConfigReader struct{}

var _ ServerConfigReader = DoltServerConfigReader{}

func (d DoltServerConfigReader) ReadConfigFile(cwdFS filesys.Filesys, file string) (ServerConfig, error) {
	return YamlConfigFromFile(cwdFS, file)
}

func (d DoltServerConfigReader) ReadConfigArgs(args *argparser.ArgParseResults) (ServerConfig, error) {
	return NewCommandLineConfig(nil, args)
}
