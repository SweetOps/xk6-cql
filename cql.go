package cql

import (
	"errors"
	"time"

	"github.com/gocql/gocql"
	"go.k6.io/k6/js/modules"
)

const ImportPath = "k6/x/cql"

//nolint:gochecknoinits // This is the recommended way to register the module.
func init() {
	modules.Register(ImportPath, new(RootModule))
}

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create `k6/x/cq;` module instances for each VU.
type RootModule struct{}

type CQL struct {
	vu      modules.VU
	session *gocql.Session
}

type Config struct {
	Timeout         string    `json:"timeout"`
	Hosts           []string  `json:"hosts"`
	Username        string    `json:"username"`
	Password        string    `json:"password"`
	Keyspace        string    `json:"keyspace"`
	ProtocolVersion int       `json:"protocol_version"`
	Consistency     string    `json:"consistency"`
	TLS             ConfigTLS `json:"tls"`
}

type ConfigTLS struct {
	CertPath               string `json:"cert_path"`
	KeyPath                string `json:"key_path"`
	CaPath                 string `json:"ca_path"`
	EnableHostVerification bool   `json:"enable_host_verification"`
}

// Ensure the interfaces are implemented correctly.
var (
	_                 modules.Module   = &RootModule{}
	_                 modules.Instance = &CQL{}
	ConnectionTimeout                  = 10 * time.Second
)

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &CQL{vu: vu}
}

// Exports implements the modules.Instance interface and returns the exports
// of the JS module.
func (cql *CQL) Exports() modules.Exports {
	return modules.Exports{Default: cql}
}

func (cql *CQL) Session(cfg Config) error {
	if len(cfg.Hosts) == 0 || cfg.Keyspace == "" {
		return errors.New("hosts and keyspace are required parameters")
	}

	cluster := gocql.NewCluster(cfg.Hosts...)
	cluster.Keyspace = cfg.Keyspace
	cluster.Consistency = cql.resolveConsistency(cfg.Consistency)

	if cfg.ProtocolVersion != 0 {
		cluster.ProtoVersion = cfg.ProtocolVersion
	}

	if cfg.Timeout == "" {
		cluster.Timeout = ConnectionTimeout
		cluster.ConnectTimeout = ConnectionTimeout
		cluster.WriteTimeout = ConnectionTimeout
	} else {
		timeout, err := time.ParseDuration(cfg.Timeout)
		if err != nil {
			return errors.New("invalid timeout value: " + err.Error())
		}
		cluster.Timeout = timeout
		cluster.ConnectTimeout = timeout
	}

	if cfg.Username != "" && cfg.Password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: cfg.Username,
			Password: cfg.Password,
		}
	}

	if cfg.TLS.CertPath != "" && cfg.TLS.KeyPath != "" {
		cluster.SslOpts = &gocql.SslOptions{
			CertPath: cfg.TLS.CertPath,
			KeyPath:  cfg.TLS.KeyPath,
		}
	}

	if cfg.TLS.CaPath != "" {
		cluster.SslOpts = &gocql.SslOptions{
			CaPath: cfg.TLS.CaPath,
		}
	}

	if cfg.TLS.EnableHostVerification {
		cluster.SslOpts = &gocql.SslOptions{
			EnableHostVerification: cfg.TLS.EnableHostVerification,
		}
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return errors.New("failed to create session: " + err.Error())
	}
	cql.session = session
	return nil
}

func (cql *CQL) Exec(query string) error {
	if err := cql.isConnected(); err != nil {
		return err
	}

	return cql.session.Query(query).Exec()
}

func (cql *CQL) Close() {
	if cql.session != nil {
		cql.session.Close()
	}
}

func (cql *CQL) Batch(batchType string, queries []string) error {
	if err := cql.isConnected(); err != nil {
		return err
	}

	bType := cql.resolveBatchType(batchType)
	b := cql.session.NewBatch(bType)
	for _, query := range queries {
		b.Query(query)
	}

	return cql.session.ExecuteBatch(b)
}

func (cql *CQL) resolveBatchType(batchType string) gocql.BatchType {
	switch batchType {
	case "unlogged":
		return gocql.UnloggedBatch
	case "counter":
		return gocql.CounterBatch
	default:
		return gocql.LoggedBatch
	}
}

func (cql *CQL) resolveConsistency(consistency string) gocql.Consistency {
	switch consistency {
	case "all":
		return gocql.All
	case "any":
		return gocql.Any
	case "one":
		return gocql.One
	case "two":
		return gocql.Two
	case "three":
		return gocql.Three
	case "each_quorum":
		return gocql.EachQuorum
	case "quorum":
		return gocql.Quorum
	case "local_one":
		return gocql.LocalOne
	case "local_quorum":
		return gocql.LocalQuorum
	default:
		return gocql.Quorum
	}
}

func (cql *CQL) isConnected() error {
	if cql.session == nil {
		return errors.New("not connected to a cluster")
	}
	return nil
}
