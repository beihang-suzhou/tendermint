package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

const (
	// FuzzModeDrop is a mode in which we randomly drop reads/writes, connections or sleep
	FuzzModeDrop = iota
	// FuzzModeDelay is a mode in which we randomly sleep
	FuzzModeDelay

	// LogFormatPlain is a format for colored text
	LogFormatPlain = "plain"
	// LogFormatJSON is a format for json output
	LogFormatJSON = "json"
)

// NOTE: Most of the structs & relevant comments + the
// default configuration options were used to manually
// generate the config.toml. Please reflect any changes
// made here in the defaultConfigTemplate constant in
// config/toml.go
// NOTE: libs/cli must know to look in the config dir!
var (
	DefaultTendermintDir = ".tendermint"
	defaultConfigDir     = "config"
	defaultDataDir       = "data"
	//defaultConfigDir     = "G:/root/.bschaind0/config/"
	//defaultDataDir       = "G:/root/.bschaind0/data"
	//defaultConfigDir     = "G:/.bschaind/config/"
	//defaultDataDir       = "G:/.bschaind/data"
	defaultConfigFileName  = "config.toml"
	defaultGenesisJSONName = "genesis.json"

	defaultPrivValKeyName   = "priv_validator_key.json"
	defaultPrivValStateName = "priv_validator_state.json"

	defaultNodeKeyName  = "node_key.json"
	defaultAddrBookName = "addrbook.json"

	defaultConfigFilePath   = filepath.Join(defaultConfigDir, defaultConfigFileName)
	defaultGenesisJSONPath  = filepath.Join(defaultConfigDir, defaultGenesisJSONName)
	defaultPrivValKeyPath   = filepath.Join(defaultConfigDir, defaultPrivValKeyName)
	defaultPrivValStatePath = filepath.Join(defaultDataDir, defaultPrivValStateName)

	defaultNodeKeyPath  = filepath.Join(defaultConfigDir, defaultNodeKeyName)
	defaultAddrBookPath = filepath.Join(defaultConfigDir, defaultAddrBookName)
)

var (
	oldPrivVal     = "priv_validator.json"
	oldPrivValPath = filepath.Join(defaultConfigDir, oldPrivVal)
)

// Config defines the top level configuration for a Tendermint node
type Config struct {
	// Top level options use an anonymous struct
	BaseConfig `mapstructure:",squash"`

	// Options for services
	RPC             *RPCConfig             `mapstructure:"rpc"`
	P2P             *P2PConfig             `mapstructure:"p2p"`
	Mempool         *MempoolConfig         `mapstructure:"mempool"`
	Consensus       *ConsensusConfig       `mapstructure:"consensus"`
	TxIndex         *TxIndexConfig         `mapstructure:"tx_index"`
	Instrumentation *InstrumentationConfig `mapstructure:"instrumentation"`
}

// DefaultConfig returns a default configuration for a Tendermint node
func DefaultConfig() *Config {
	return &Config{
		BaseConfig:      DefaultBaseConfig(),
		RPC:             DefaultRPCConfig(),
		P2P:             DefaultP2PConfig(),
		Mempool:         DefaultMempoolConfig(),
		Consensus:       DefaultConsensusConfig(),
		TxIndex:         DefaultTxIndexConfig(),
		Instrumentation: DefaultInstrumentationConfig(),
	}

    //方法1
	/*var conf *Config
	file,_:=os.Open("tm.toml")
	buf,_:=ioutil.ReadAll(file)
	err :=toml.Unmarshal(buf,&conf)
	if err != nil {
	fmt. Println ( "error:" , err )
	}
	return conf*/

	//方法2
	/*var conf *Config
	conf = new(Config)
	if _, err := toml.DecodeFile("tm.toml", conf); err != nil {
		panic(err)
	}
	spew.Dump(conf)
	return conf
	*/


}
//改为toml方式读取时，此处需要注释掉
// TestConfig returns a configuration that can be used for testing
func TestConfig() *Config {
	return &Config{
		BaseConfig:      TestBaseConfig(),
		RPC:             TestRPCConfig(),
		P2P:             TestP2PConfig(),
		Mempool:         TestMempoolConfig(),
		Consensus:       TestConsensusConfig(),
		TxIndex:         TestTxIndexConfig(),
		Instrumentation: TestInstrumentationConfig(),
	}
}

// SetRoot sets the RootDir for all Config structs
func (cfg *Config) SetRoot(root string) *Config {
	cfg.BaseConfig.RootDir = root
	cfg.RPC.RootDir = root
	cfg.P2P.RootDir = root
	cfg.Mempool.RootDir = root
	cfg.Consensus.RootDir = root
	return cfg
}

// ValidateBasic performs basic validation (checking param bounds, etc.) and
// returns an error if any check fails.
//改为toml方式读取时，此处需要注释掉
func (cfg *Config) ValidateBasic() error {
	if err := cfg.BaseConfig.ValidateBasic(); err != nil {
		return err
	}
	if err := cfg.RPC.ValidateBasic(); err != nil {
		return errors.Wrap(err, "Error in [rpc] section")
	}
	if err := cfg.P2P.ValidateBasic(); err != nil {
		return errors.Wrap(err, "Error in [p2p] section")
	}
	if err := cfg.Mempool.ValidateBasic(); err != nil {
		return errors.Wrap(err, "Error in [mempool] section")
	}
	if err := cfg.Consensus.ValidateBasic(); err != nil {
		return errors.Wrap(err, "Error in [consensus] section")
	}
	return errors.Wrap(
		cfg.Instrumentation.ValidateBasic(),
		"Error in [instrumentation] section",
	)
}

//-----------------------------------------------------------------------------
// BaseConfig

// BaseConfig defines the base configuration for a Tendermint node
type BaseConfig struct {
	// chainID is unexposed and immutable but here for convenience
	chainID string

	// The root directory for all data.
	// This should be set in viper so it can unmarshal into this struct
	RootDir string `mapstructure:"home"`

	// TCP or UNIX socket address of the ABCI application,
	// or the name of an ABCI application compiled in with the Tendermint binary
	ProxyApp string `toml:"proxy_app" mapstructure:"proxy_app"`

	// A custom human readable name for this node
	Moniker string `toml:"moniker" mapstructure:"moniker"`

	// If this node is many blocks behind the tip of the chain, FastSync
	// allows them to catchup quickly by downloading blocks in parallel
	// and verifying their commits
	FastSync bool `toml:"fast_sync" mapstructure:"fast_sync"`

	// Database backend: leveldb | memdb | cleveldb
	DBBackend string `toml:"db_backend" mapstructure:"db_backend"`

	// Database directory
	DBPath string `toml:"db_dir" mapstructure:"db_dir"`

	// Output level for logging
	LogLevel string `toml:"log_level" mapstructure:"log_level"`

	// Output format: 'plain' (colored text) or 'json'
	LogFormat string `toml:"log_format" mapstructure:"log_format"`

	// Path to the JSON file containing the initial validator set and other meta data
	Genesis string `toml:"genesis_file" mapstructure:"genesis_file"`

	// Path to the JSON file containing the private key to use as a validator in the consensus protocol
	PrivValidatorKey string `toml:"priv_validator_key_file" mapstructure:"priv_validator_key_file"`

	// Path to the JSON file containing the last sign state of a validator
	PrivValidatorState string `toml:"priv_validator_state_file" mapstructure:"priv_validator_state_file"`

	// TCP or UNIX socket address for Tendermint to listen on for
	// connections from an external PrivValidator process
	PrivValidatorListenAddr string `toml:"priv_validator_laddr" mapstructure:"priv_validator_laddr"`

	// A JSON file containing the private key to use for p2p authenticated encryption
	NodeKey string `toml:"node_key_file" mapstructure:"node_key_file"`

	// Mechanism to connect to the ABCI application: socket | grpc
	ABCI string `toml:"abci" mapstructure:"abci"`

	// TCP or UNIX socket address for the profiling server to listen on
	ProfListenAddress string `toml:"prof_laddr" mapstructure:"prof_laddr"`

	// If true, query the ABCI app on connecting to a new peer
	// so the app can decide if we should keep the connection or not
	FilterPeers bool `toml:"filter_peers" mapstructure:"filter_peers"` // false

}

// DefaultBaseConfig returns a default base configuration for a Tendermint node
func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		Genesis:            defaultGenesisJSONPath,
		PrivValidatorKey:   defaultPrivValKeyPath,
		PrivValidatorState: defaultPrivValStatePath,
		NodeKey:            defaultNodeKeyPath,
		Moniker:            defaultMoniker,
		ProxyApp:           "tcp://127.0.0.1:26658",
		ABCI:               "socket",
		LogLevel:           DefaultPackageLogLevels(),
		LogFormat:          LogFormatPlain,
		ProfListenAddress:  "",
		FastSync:           true,
		FilterPeers:        false,
		DBBackend:          "leveldb",
		DBPath:             "data",
	}
}

// TestBaseConfig returns a base configuration for testing a Tendermint node
func TestBaseConfig() BaseConfig {
	cfg := DefaultBaseConfig()
	cfg.chainID = "tendermint_test"
	cfg.ProxyApp = "kvstore"
	cfg.FastSync = false
	cfg.DBBackend = "memdb"
	return cfg
}

func (cfg BaseConfig) ChainID() string {
	return cfg.chainID
}

// GenesisFile returns the full path to the genesis.json file
func (cfg BaseConfig) GenesisFile() string {
	return rootify(cfg.Genesis, cfg.RootDir)
}

// PrivValidatorKeyFile returns the full path to the priv_validator_key.json file
func (cfg BaseConfig) PrivValidatorKeyFile() string {
	return rootify(cfg.PrivValidatorKey, cfg.RootDir)
}

// PrivValidatorFile returns the full path to the priv_validator_state.json file
func (cfg BaseConfig) PrivValidatorStateFile() string {
	return rootify(cfg.PrivValidatorState, cfg.RootDir)
}

// OldPrivValidatorFile returns the full path of the priv_validator.json from pre v0.28.0.
// TODO: eventually remove.
func (cfg BaseConfig) OldPrivValidatorFile() string {
	return rootify(oldPrivValPath, cfg.RootDir)
}

// NodeKeyFile returns the full path to the node_key.json file
func (cfg BaseConfig) NodeKeyFile() string {
	return rootify(cfg.NodeKey, cfg.RootDir)
}

// DBDir returns the full path to the database directory
func (cfg BaseConfig) DBDir() string {
	return rootify(cfg.DBPath, cfg.RootDir)
}

// ValidateBasic performs basic validation (checking param bounds, etc.) and
// returns an error if any check fails.
func (cfg BaseConfig) ValidateBasic() error {
	switch cfg.LogFormat {
	case LogFormatPlain, LogFormatJSON:
	default:
		return errors.New("unknown log_format (must be 'plain' or 'json')")
	}
	return nil
}

// DefaultLogLevel returns a default log level of "error"
func DefaultLogLevel() string {
	return "error"
}

// DefaultPackageLogLevels returns a default log level setting so all packages
// log at "error", while the `state` and `main` packages log at "info"
func DefaultPackageLogLevels() string {
	return fmt.Sprintf("%s", DefaultLogLevel()) //main:info,state:info
}

//-----------------------------------------------------------------------------
// RPCConfig

// RPCConfig defines the configuration options for the Tendermint RPC server
type RPCConfig struct {
	RootDir string `toml:"home" mapstructure:"home"`

	// TCP or UNIX socket address for the RPC server to listen on
	ListenAddress string `toml:"laddr" mapstructure:"laddr"`

	// A list of origins a cross-domain request can be executed from.
	// If the special '*' value is present in the list, all origins will be allowed.
	// An origin may contain a wildcard (*) to replace 0 or more characters (i.e.: http://*.domain.com).
	// Only one wildcard can be used per origin.
	CORSAllowedOrigins []string `toml:"cors_allowed_origins" mapstructure:"cors_allowed_origins"`

	// A list of methods the client is allowed to use with cross-domain requests.
	CORSAllowedMethods []string `toml:"cors_allowed_methods" mapstructure:"cors_allowed_methods"`

	// A list of non simple headers the client is allowed to use with cross-domain requests.
	CORSAllowedHeaders []string `toml:"cors_allowed_headers" mapstructure:"cors_allowed_headers"`

	// TCP or UNIX socket address for the gRPC server to listen on
	// NOTE: This server only supports /broadcast_tx_commit
	GRPCListenAddress string `toml:"grpc_laddr" mapstructure:"grpc_laddr"`

	// Maximum number of simultaneous connections.
	// Does not include RPC (HTTP&WebSocket) connections. See max_open_connections
	// If you want to accept a larger number than the default, make sure
	// you increase your OS limits.
	// 0 - unlimited.
	GRPCMaxOpenConnections int `toml:"grpc_max_open_connections" mapstructure:"grpc_max_open_connections"`

	// Activate unsafe RPC commands like /dial_persistent_peers and /unsafe_flush_mempool
	Unsafe bool `toml:"unsafe" mapstructure:"unsafe"`

	// Maximum number of simultaneous connections (including WebSocket).
	// Does not include gRPC connections. See grpc_max_open_connections
	// If you want to accept a larger number than the default, make sure
	// you increase your OS limits.
	// 0 - unlimited.
	// Should be < {ulimit -Sn} - {MaxNumInboundPeers} - {MaxNumOutboundPeers} - {N of wal, db and other open files}
	// 1024 - 40 - 10 - 50 = 924 = ~900
	MaxOpenConnections int `toml:"max_open_connections" mapstructure:"max_open_connections"`
}

// DefaultRPCConfig returns a default configuration for the RPC server
func DefaultRPCConfig() *RPCConfig {
	return &RPCConfig{
		ListenAddress:          "tcp://0.0.0.0:26657",
		CORSAllowedOrigins:     []string{},
		CORSAllowedMethods:     []string{"HEAD", "GET", "POST"},
		CORSAllowedHeaders:     []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "X-Server-Time"},
		GRPCListenAddress:      "",
		GRPCMaxOpenConnections: 900,

		Unsafe:             false,
		MaxOpenConnections: 900,
	}
}

// TestRPCConfig returns a configuration for testing the RPC server
func TestRPCConfig() *RPCConfig {
	cfg := DefaultRPCConfig()
	cfg.ListenAddress = "tcp://0.0.0.0:36657"
	cfg.GRPCListenAddress = "tcp://0.0.0.0:36658"
	cfg.Unsafe = true
	return cfg
}

// ValidateBasic performs basic validation (checking param bounds, etc.) and
// returns an error if any check fails.
func (cfg *RPCConfig) ValidateBasic() error {
	if cfg.GRPCMaxOpenConnections < 0 {
		return errors.New("grpc_max_open_connections can't be negative")
	}
	if cfg.MaxOpenConnections < 0 {
		return errors.New("max_open_connections can't be negative")
	}
	return nil
}

// IsCorsEnabled returns true if cross-origin resource sharing is enabled.
func (cfg *RPCConfig) IsCorsEnabled() bool {
	return len(cfg.CORSAllowedOrigins) != 0
}
type duration struct {
	time.Duration
}
func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
//-----------------------------------------------------------------------------
// P2PConfig

// P2PConfig defines the configuration options for the Tendermint peer-to-peer networking layer
type P2PConfig struct {
	RootDir string `toml:"home" mapstructure:"home"`

	// Address to listen for incoming connections
	ListenAddress string `toml:"laddr" mapstructure:"laddr"`

	// Address to advertise to peers for them to dial
	ExternalAddress string `toml:"external_address" mapstructure:"external_address"`

	// Comma separated list of seed nodes to connect to
	// We only use these if we can’t connect to peers in the addrbook
	Seeds string `toml:"seeds" mapstructure:"seeds"`

	// Comma separated list of nodes to keep persistent connections to
	PersistentPeers string `toml:"persistent_peers" mapstructure:"persistent_peers"`

	// UPNP port forwarding
	UPNP bool `toml:"upnp" mapstructure:"upnp"`

	// Path to address book
	AddrBook string `toml:"addr_book_file" mapstructure:"addr_book_file"`

	// Set true for strict address routability rules
	// Set false for private or local networks
	AddrBookStrict bool `toml:"addr_book_strict" mapstructure:"addr_book_strict"`

	// Maximum number of inbound peers
	MaxNumInboundPeers int `toml:"max_num_inbound_peers" mapstructure:"max_num_inbound_peers"`

	// Maximum number of outbound peers to connect to, excluding persistent peers
	MaxNumOutboundPeers int `toml:"max_num_outbound_peers" mapstructure:"max_num_outbound_peers"`

	// Time to wait before flushing messages out on the connection
	FlushThrottleTimeout time.Duration `toml:"flush_throttle_timeout" mapstructure:"flush_throttle_timeout"`
	//FlushThrottleTimeout duration `toml:"flush_throttle_timeout" mapstructure:"flush_throttle_timeout"`//改为toml方式，需要使用自定义小写duration类型，该类型实现了TextUnmarshaler接口

	// Maximum size of a message packet payload, in bytes
	MaxPacketMsgPayloadSize int `toml:"max_packet_msg_payload_size" mapstructure:"max_packet_msg_payload_size"`

	// Rate at which packets can be sent, in bytes/second
	SendRate int64 `toml:"send_rate" mapstructure:"send_rate"`

	// Rate at which packets can be received, in bytes/second
	RecvRate int64 `toml:"recv_rate" mapstructure:"recv_rate"`

	// Set true to enable the peer-exchange reactor
	PexReactor bool `toml:"pex" mapstructure:"pex"`

	// Seed mode, in which node constantly crawls the network and looks for
	// peers. If another node asks it for addresses, it responds and disconnects.
	//
	// Does not work if the peer-exchange reactor is disabled.
	SeedMode bool `toml:"seed_mode" mapstructure:"seed_mode"`

	// Comma separated list of peer IDs to keep private (will not be gossiped to
	// other peers)
	PrivatePeerIDs string `toml:"private_peer_ids" mapstructure:"private_peer_ids"`

	// Toggle to disable guard against peers connecting from the same ip.
	AllowDuplicateIP bool `toml:"allow_duplicate_ip" mapstructure:"allow_duplicate_ip"`

	// Peer connection configuration.
	HandshakeTimeout time.Duration `mapstructure:"handshake_timeout"`
	DialTimeout      time.Duration `mapstructure:"dial_timeout"`

	//toml方式解析duration
	//HandshakeTimeout duration `toml:"handshake_timeout" mapstructure:"handshake_timeout"`//改为toml方式，需要使用自定义小写duration类型，该类型实现了TextUnmarshaler接口
	//DialTimeout      duration `toml:"dial_timeout" mapstructure:"dial_timeout"`//改为toml方式，需要使用自定义小写duration类型，该类型实现了TextUnmarshaler接口

	// Testing params.
	// Force dial to fail
	TestDialFail bool `toml:"test_dial_fail" mapstructure:"test_dial_fail"`
	// FUzz connection
	TestFuzz       bool            `toml:"test_fuzz" mapstructure:"test_fuzz"`
	TestFuzzConfig *FuzzConnConfig `toml:"test_fuzz_config" mapstructure:"test_fuzz_config"`
}

// DefaultP2PConfig returns a default configuration for the peer-to-peer layer
//改为toml方式读取时，此处需要注释掉
func DefaultP2PConfig() *P2PConfig {
	return &P2PConfig{
		ListenAddress:           "tcp://0.0.0.0:26656",
		ExternalAddress:         "",
		UPNP:                    false,
		AddrBook:                defaultAddrBookPath,
		AddrBookStrict:          true,
		MaxNumInboundPeers:      40,
		MaxNumOutboundPeers:     10,
		FlushThrottleTimeout:    100 * time.Millisecond,
		MaxPacketMsgPayloadSize: 1024,    // 1 kB
		SendRate:                5120000, // 5 mB/s
		RecvRate:                5120000, // 5 mB/s
		PexReactor:              true,
		SeedMode:                false,
		AllowDuplicateIP:        false,
		HandshakeTimeout:        20 * time.Second,
		DialTimeout:             3 * time.Second,
		TestDialFail:            false,
		TestFuzz:                false,
		TestFuzzConfig:          DefaultFuzzConnConfig(),
	}
}

// TestP2PConfig returns a configuration for testing the peer-to-peer layer
//改为toml方式读取时，此处需要注释掉
func TestP2PConfig() *P2PConfig {
	cfg := DefaultP2PConfig()
	cfg.ListenAddress = "tcp://0.0.0.0:36656"
	cfg.FlushThrottleTimeout = 10 * time.Millisecond
	cfg.AllowDuplicateIP = true
	return cfg
}

// AddrBookFile returns the full path to the address book
func (cfg *P2PConfig) AddrBookFile() string {
	return rootify(cfg.AddrBook, cfg.RootDir)
}

// ValidateBasic performs basic validation (checking param bounds, etc.) and
// returns an error if any check fails.
//改为toml方式读取时，此处需要注释掉
func (cfg *P2PConfig) ValidateBasic() error {
	if cfg.MaxNumInboundPeers < 0 {
		return errors.New("max_num_inbound_peers can't be negative")
	}
	if cfg.MaxNumOutboundPeers < 0 {
		return errors.New("max_num_outbound_peers can't be negative")
	}
	if cfg.FlushThrottleTimeout < 0 {
		return errors.New("flush_throttle_timeout can't be negative")
	}
	if cfg.MaxPacketMsgPayloadSize < 0 {
		return errors.New("max_packet_msg_payload_size can't be negative")
	}
	if cfg.SendRate < 0 {
		return errors.New("send_rate can't be negative")
	}
	if cfg.RecvRate < 0 {
		return errors.New("recv_rate can't be negative")
	}
	return nil
}

// FuzzConnConfig is a FuzzedConnection configuration.
type FuzzConnConfig struct {
	Mode         int
	MaxDelay     time.Duration
	ProbDropRW   float64
	ProbDropConn float64
	ProbSleep    float64
}

// DefaultFuzzConnConfig returns the default config.
func DefaultFuzzConnConfig() *FuzzConnConfig {
	return &FuzzConnConfig{
		Mode:         FuzzModeDrop,
		MaxDelay:     3 * time.Second,
		ProbDropRW:   0.2,
		ProbDropConn: 0.00,
		ProbSleep:    0.00,
	}
}

//-----------------------------------------------------------------------------
// MempoolConfig

// MempoolConfig defines the configuration options for the Tendermint mempool
type MempoolConfig struct {
	RootDir   string `toml:"home" mapstructure:"home"`
	Recheck   bool   `toml:"recheck" mapstructure:"recheck"`
	Broadcast bool   `toml:"broadcast" mapstructure:"broadcast"`
	WalPath   string `toml:"wal_dir" mapstructure:"wal_dir"`
	Size      int    `toml:"size" mapstructure:"size"`
	CacheSize int    `toml:"cache_size" mapstructure:"cache_size"`
}

// DefaultMempoolConfig returns a default configuration for the Tendermint mempool
func DefaultMempoolConfig() *MempoolConfig {
	return &MempoolConfig{
		Recheck:   true,
		Broadcast: true,
		WalPath:   "",
		// Each signature verification takes .5ms, size reduced until we implement
		// ABCI Recheck
		Size:      5000,
		CacheSize: 10000,
	}
}

// TestMempoolConfig returns a configuration for testing the Tendermint mempool
func TestMempoolConfig() *MempoolConfig {
	cfg := DefaultMempoolConfig()
	cfg.CacheSize = 1000
	return cfg
}

// WalDir returns the full path to the mempool's write-ahead log
func (cfg *MempoolConfig) WalDir() string {
	return rootify(cfg.WalPath, cfg.RootDir)
}

// WalEnabled returns true if the WAL is enabled.
func (cfg *MempoolConfig) WalEnabled() bool {
	return cfg.WalPath != ""
}

// ValidateBasic performs basic validation (checking param bounds, etc.) and
// returns an error if any check fails.
func (cfg *MempoolConfig) ValidateBasic() error {
	if cfg.Size < 0 {
		return errors.New("size can't be negative")
	}
	if cfg.CacheSize < 0 {
		return errors.New("cache_size can't be negative")
	}
	return nil
}

//-----------------------------------------------------------------------------
// ConsensusConfig

// ConsensusConfig defines the configuration for the Tendermint consensus service,
// including timeouts and details about the WAL and the block structure.
type ConsensusConfig struct {
	RootDir string `toml:"home" mapstructure:"home"`
	WalPath string `toml:"wal_file" mapstructure:"wal_file"`
	walFile string // overrides WalPath if set
	//默认方式解析time.Duration
	TimeoutPropose        time.Duration `mapstructure:"timeout_propose"`
	TimeoutProposeDelta   time.Duration `mapstructure:"timeout_propose_delta"`
	TimeoutPrevote        time.Duration `mapstructure:"timeout_prevote"`
	TimeoutPrevoteDelta   time.Duration `mapstructure:"timeout_prevote_delta"`
	TimeoutPrecommit      time.Duration `mapstructure:"timeout_precommit"`
	TimeoutPrecommitDelta time.Duration `mapstructure:"timeout_precommit_delta"`
	TimeoutCommit         time.Duration `mapstructure:"timeout_commit"`
	//默认方式解析time.Duration
	CreateEmptyBlocksInterval time.Duration `mapstructure:"create_empty_blocks_interval"`
	PeerGossipSleepDuration     time.Duration `mapstructure:"peer_gossip_sleep_duration"`
	//Reactor sleep duration parameters
	PeerQueryMaj23SleepDuration time.Duration `mapstructure:"peer_query_maj23_sleep_duration"`
	//Block time parameters. Corresponds to the minimum time increment between consecutive blocks.
	BlockTimeIota time.Duration `mapstructure:"blocktime_iota"`

	//toml方式解析duration
	//TimeoutPropose        duration `toml:"dial_timeout" mapstructure:"timeout_propose"`//改为toml方式，需要使用自定义小写duration类型，该类型实现了TextUnmarshaler接口
	//TimeoutProposeDelta   duration `toml:"timeout_propose_delta" mapstructure:"timeout_propose_delta"`//改为toml方式，需要使用自定义小写duration类型，该类型实现了TextUnmarshaler接口
	//TimeoutPrevote        duration `toml:"timeout_prevote" mapstructure:"timeout_prevote"`//改为toml方式，需要使用自定义小写duration类型，该类型实现了TextUnmarshaler接口
	//TimeoutPrevoteDelta   duration `toml:"timeout_prevote_delta" mapstructure:"timeout_prevote_delta"`//改为toml方式，需要使用自定义小写duration类型，该类型实现了TextUnmarshaler接口
	//TimeoutPrecommit      duration `toml:"timeout_precommit" mapstructure:"timeout_precommit"`//改为toml方式，需要使用自定义小写duration类型，该类型实现了TextUnmarshaler接口
	//TimeoutPrecommitDelta duration `toml:"timeout_precommit_delta" mapstructure:"timeout_precommit_delta"`//改为toml方式，需要使用自定义小写duration类型，该类型实现了TextUnmarshaler接口
	//TimeoutCommit         duration `toml:"timeout_commit" mapstructure:"timeout_commit"`//改为toml方式，需要使用自定义小写duration类型，该类型实现了TextUnmarshaler接口
	//toml方式解析duration
	//CreateEmptyBlocksInterval duration `toml:"create_empty_blocks_interval" mapstructure:"create_empty_blocks_interval"`//改为toml方式，需要使用自定义小写duration类型，该类型实现了TextUnmarshaler接口
	//// Reactor sleep duration parameters
	//PeerGossipSleepDuration     duration `toml:"peer_gossip_sleep_duration" mapstructure:"peer_gossip_sleep_duration"`//改为toml方式，需要使用自定义小写duration类型，该类型实现了TextUnmarshaler接口
	//PeerQueryMaj23SleepDuration duration `toml:"peer_query_maj23_sleep_duration" mapstructure:"peer_query_maj23_sleep_duration"`
	//// Block time parameters. Corresponds to the minimum time increment between consecutive blocks.
	//BlockTimeIota duration `toml:"blocktime_iota" mapstructure:"blocktime_iota"`//改为toml方式，需要使用自定义小写duration类型，该类型实现了TextUnmarshaler接口

	// Make progress as soon as we have all the precommits (as if TimeoutCommit = 0)
	SkipTimeoutCommit bool `toml:"skip_timeout_commit" mapstructure:"skip_timeout_commit"`
	// EmptyBlocks mode and possible interval between empty blocks
	CreateEmptyBlocks         bool          `toml:"create_empty_blocks" mapstructure:"create_empty_blocks"`
}

// DefaultConsensusConfig returns a default configuration for the consensus service
//改为toml方式，此处需要注释掉
func DefaultConsensusConfig() *ConsensusConfig {
	return &ConsensusConfig{
		WalPath:                     filepath.Join(defaultDataDir, "cs.wal", "wal"),
		TimeoutPropose:              3000 * time.Millisecond,
		TimeoutProposeDelta:         500 * time.Millisecond,
		TimeoutPrevote:              1000 * time.Millisecond,
		TimeoutPrevoteDelta:         500 * time.Millisecond,
		TimeoutPrecommit:            1000 * time.Millisecond,
		TimeoutPrecommitDelta:       500 * time.Millisecond,
		TimeoutCommit:               1000 * time.Millisecond,
		SkipTimeoutCommit:           false,
		CreateEmptyBlocks:           true,
		CreateEmptyBlocksInterval:   0 * time.Second,
		PeerGossipSleepDuration:     100 * time.Millisecond,
		PeerQueryMaj23SleepDuration: 2000 * time.Millisecond,
		BlockTimeIota:               1000 * time.Millisecond,
	}
}

// TestConsensusConfig returns a configuration for testing the consensus service
//改为toml方式，此处需要注释掉
func TestConsensusConfig() *ConsensusConfig {
	cfg := DefaultConsensusConfig()
	cfg.TimeoutPropose = 40 * time.Millisecond
	cfg.TimeoutProposeDelta = 1 * time.Millisecond
	cfg.TimeoutPrevote = 10 * time.Millisecond
	cfg.TimeoutPrevoteDelta = 1 * time.Millisecond
	cfg.TimeoutPrecommit = 10 * time.Millisecond
	cfg.TimeoutPrecommitDelta = 1 * time.Millisecond
	cfg.TimeoutCommit = 10 * time.Millisecond
	cfg.SkipTimeoutCommit = true
	cfg.PeerGossipSleepDuration = 5 * time.Millisecond
	cfg.PeerQueryMaj23SleepDuration = 250 * time.Millisecond
	cfg.BlockTimeIota = 10 * time.Millisecond
	return cfg
}

// MinValidVoteTime returns the minimum acceptable block time.
// See the [BFT time spec](https://godoc.org/github.com/tendermint/tendermint/docs/spec/consensus/bft-time.md).
//改为toml方式，此处需要注释掉
func (cfg *ConsensusConfig) MinValidVoteTime(lastBlockTime time.Time) time.Time {
	return lastBlockTime.Add(cfg.BlockTimeIota)
}

// WaitForTxs returns true if the consensus should wait for transactions before entering the propose step
//改为toml方式，此处需要注释掉
func (cfg *ConsensusConfig) WaitForTxs() bool {
	return !cfg.CreateEmptyBlocks || cfg.CreateEmptyBlocksInterval > 0
}

// Propose returns the amount of time to wait for a proposal
func (cfg *ConsensusConfig) Propose(round int) time.Duration {
	return time.Duration(
		cfg.TimeoutPropose.Nanoseconds()+cfg.TimeoutProposeDelta.Nanoseconds()*int64(round),
	) * time.Nanosecond
}

// Prevote returns the amount of time to wait for straggler votes after receiving any +2/3 prevotes
func (cfg *ConsensusConfig) Prevote(round int) time.Duration {
	return time.Duration(
		cfg.TimeoutPrevote.Nanoseconds()+cfg.TimeoutPrevoteDelta.Nanoseconds()*int64(round),
	) * time.Nanosecond
}

// Precommit returns the amount of time to wait for straggler votes after receiving any +2/3 precommits
func (cfg *ConsensusConfig) Precommit(round int) time.Duration {
	return time.Duration(
		cfg.TimeoutPrecommit.Nanoseconds()+cfg.TimeoutPrecommitDelta.Nanoseconds()*int64(round),
	) * time.Nanosecond
}

// Commit returns the amount of time to wait for straggler votes after receiving +2/3 precommits for a single block (ie. a commit).
//改为toml方式，此处需要注释掉
func (cfg *ConsensusConfig) Commit(t time.Time) time.Time {
	return t.Add(cfg.TimeoutCommit)
}

// WalFile returns the full path to the write-ahead log file
func (cfg *ConsensusConfig) WalFile() string {
	if cfg.walFile != "" {
		return cfg.walFile
	}
	return rootify(cfg.WalPath, cfg.RootDir)
}

// SetWalFile sets the path to the write-ahead log file
func (cfg *ConsensusConfig) SetWalFile(walFile string) {
	cfg.walFile = walFile
}

// ValidateBasic performs basic validation (checking param bounds, etc.) and
// returns an error if any check fails.
//改为toml方式，此处需要注释掉
func (cfg *ConsensusConfig) ValidateBasic() error {
	if cfg.TimeoutPropose < 0 {
		return errors.New("timeout_propose can't be negative")
	}
	if cfg.TimeoutProposeDelta < 0 {
		return errors.New("timeout_propose_delta can't be negative")
	}
	if cfg.TimeoutPrevote < 0 {
		return errors.New("timeout_prevote can't be negative")
	}
	if cfg.TimeoutPrevoteDelta < 0 {
		return errors.New("timeout_prevote_delta can't be negative")
	}
	if cfg.TimeoutPrecommit < 0 {
		return errors.New("timeout_precommit can't be negative")
	}
	if cfg.TimeoutPrecommitDelta < 0 {
		return errors.New("timeout_precommit_delta can't be negative")
	}
	if cfg.TimeoutCommit < 0 {
		return errors.New("timeout_commit can't be negative")
	}
	if cfg.CreateEmptyBlocksInterval < 0 {
		return errors.New("create_empty_blocks_interval can't be negative")
	}
	if cfg.PeerGossipSleepDuration < 0 {
		return errors.New("peer_gossip_sleep_duration can't be negative")
	}
	if cfg.PeerQueryMaj23SleepDuration < 0 {
		return errors.New("peer_query_maj23_sleep_duration can't be negative")
	}
	if cfg.BlockTimeIota < 0 {
		return errors.New("blocktime_iota can't be negative")
	}
	return nil
}

//-----------------------------------------------------------------------------
// TxIndexConfig

// TxIndexConfig defines the configuration for the transaction indexer,
// including tags to index.
type TxIndexConfig struct {
	// What indexer to use for transactions
	//
	// Options:
	//   1) "null"
	//   2) "kv" (default) - the simplest possible indexer, backed by key-value storage (defaults to levelDB; see DBBackend).
	Indexer string `mapstructure:"indexer"`

	// Comma-separated list of tags to index (by default the only tag is "tx.hash")
	//
	// You can also index transactions by height by adding "tx.height" tag here.
	//
	// It's recommended to index only a subset of tags due to possible memory
	// bloat. This is, of course, depends on the indexer's DB and the volume of
	// transactions.
	IndexTags string `mapstructure:"index_tags"`

	// When set to true, tells indexer to index all tags (predefined tags:
	// "tx.hash", "tx.height" and all tags from DeliverTx responses).
	//
	// Note this may be not desirable (see the comment above). IndexTags has a
	// precedence over IndexAllTags (i.e. when given both, IndexTags will be
	// indexed).
	IndexAllTags bool `mapstructure:"index_all_tags"`
}

// DefaultTxIndexConfig returns a default configuration for the transaction indexer.
func DefaultTxIndexConfig() *TxIndexConfig {
	return &TxIndexConfig{
		Indexer:      "kv",
		IndexTags:    "",
		IndexAllTags: false,
	}
}

// TestTxIndexConfig returns a default configuration for the transaction indexer.
func TestTxIndexConfig() *TxIndexConfig {
	return DefaultTxIndexConfig()
}

//-----------------------------------------------------------------------------
// InstrumentationConfig

// InstrumentationConfig defines the configuration for metrics reporting.
type InstrumentationConfig struct {
	// When true, Prometheus metrics are served under /metrics on
	// PrometheusListenAddr.
	// Check out the documentation for the list of available metrics.
	Prometheus bool `toml:"prometheus" mapstructure:"prometheus"`

	// Address to listen for Prometheus collector(s) connections.
	PrometheusListenAddr string `toml:"prometheus_listen_addr" mapstructure:"prometheus_listen_addr"`

	// Maximum number of simultaneous connections.
	// If you want to accept a larger number than the default, make sure
	// you increase your OS limits.
	// 0 - unlimited.
	MaxOpenConnections int `toml:"max_open_connections" mapstructure:"max_open_connections"`

	// Instrumentation namespace.
	Namespace string `toml:"namespace" mapstructure:"namespace"`
}

// DefaultInstrumentationConfig returns a default configuration for metrics
// reporting.
func DefaultInstrumentationConfig() *InstrumentationConfig {
	return &InstrumentationConfig{
		Prometheus:           false,
		PrometheusListenAddr: ":26660",
		MaxOpenConnections:   3,
		Namespace:            "tendermint",
	}
}

// TestInstrumentationConfig returns a default configuration for metrics
// reporting.
func TestInstrumentationConfig() *InstrumentationConfig {
	return DefaultInstrumentationConfig()
}

// ValidateBasic performs basic validation (checking param bounds, etc.) and
// returns an error if any check fails.
func (cfg *InstrumentationConfig) ValidateBasic() error {
	if cfg.MaxOpenConnections < 0 {
		return errors.New("max_open_connections can't be negative")
	}
	return nil
}

//-----------------------------------------------------------------------------
// Utils

// helper function to make config creation independent of root dir
func rootify(path, root string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

//-----------------------------------------------------------------------------
// Moniker

var defaultMoniker = getDefaultMoniker()

// getDefaultMoniker returns a default moniker, which is the host name. If runtime
// fails to get the host name, "anonymous" will be returned.
func getDefaultMoniker() string {
	moniker, err := os.Hostname()
	if err != nil {
		moniker = "anonymous"
	}
	return moniker
}
