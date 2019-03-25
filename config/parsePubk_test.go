package config

import (
	"testing"
	"github.com/tendermint/tendermint/privval"
	"fmt"
)

func TestName(t *testing.T) {
	config := DefaultConfig()
	 privVal := privval.LoadFilePV(config.PrivValidatorKeyFile(), config.PrivValidatorStateFile())
     fmt.Println(string(privVal.GetPubKey().Bytes()))
}
