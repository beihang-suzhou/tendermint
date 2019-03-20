package config

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"time"
)

func TestDefaultConfig(t *testing.T) {

	assert := assert.New(t)

	// set up some defaults
	cfg := DefaultConfig()
	fmt.Println("---------- BaseConfig ----------")
	fmt.Println(cfg.BaseConfig.ProxyApp)
	fmt.Println(cfg.BaseConfig.Moniker)
	fmt.Println(cfg.BaseConfig.FastSync)
	fmt.Println(cfg.BaseConfig.DBPath)
	fmt.Println(cfg.BaseConfig.NodeKeyFile())
	fmt.Println(cfg.BaseConfig.FilterPeers)
	fmt.Println("---------- RPCConfig ----------")
	fmt.Println(cfg.RPC.CORSAllowedHeaders)
	fmt.Println(cfg.RPC.ListenAddress)
	fmt.Println("---------- P2P ----------")
	fmt.Println(cfg.P2P.MaxPacketMsgPayloadSize)
	fmt.Println(cfg.P2P.TestFuzzConfig)
	fmt.Println("这里为Duration类型： ",cfg.P2P.FlushThrottleTimeout)
	fmt.Println("---------- MempoolConfig ----------")
	fmt.Println(cfg.Mempool.Recheck)
	fmt.Println("---------- ConsensusConfig ----------")
	fmt.Println(cfg.Consensus.CreateEmptyBlocks)
	fmt.Println(cfg.Consensus.CreateEmptyBlocksInterval)
	fmt.Println(cfg.Consensus.SkipTimeoutCommit)
	fmt.Println(cfg.Consensus.TimeoutPrevoteDelta)
	fmt.Println(cfg.Consensus.PeerGossipSleepDuration)
	fmt.Println(cfg.Consensus.TimeoutPrecommit)
	fmt.Println("---------- InstrumentationConfig ----------")
	fmt.Println(cfg.Instrumentation.PrometheusListenAddr)
	fmt.Println("---------- 结束 ----------")
	assert.NotNil(cfg.P2P)
	assert.NotNil(cfg.Mempool)
	assert.NotNil(cfg.Consensus)

	//check the root dir stuff...
	//cfg.SetRoot("/foo")
	//cfg.Genesis = "bar"
	//cfg.DBPath = "/opt/data"
	//cfg.Mempool.WalPath = "wal/mem/"
	//
	//assert.Equal("/foo/bar", cfg.GenesisFile())
	//assert.Equal("/opt/data", cfg.DBDir())
	//assert.Equal("/foo/wal/mem", cfg.Mempool.WalDir())

}

func TestConfigValidateBasic(t *testing.T) {
	cfg := DefaultConfig()
	assert.NoError(t, cfg.ValidateBasic())

	// tamper with timeout_propose
	cfg.Consensus.TimeoutPropose = -10 * time.Second
	assert.Error(t, cfg.ValidateBasic())
}
