package fsn

import (
	"math/big"
	"strings"
	"time"

	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/anyswap/CrossChain-Bridge/tokens/eth"
	"github.com/anyswap/CrossChain-Bridge/types"
)

const (
	netMainnet = "mainnet"
	netTestnet = "testnet"
	netDevnet  = "devnet"
	netCustom  = "custom"
)

// Bridge fsn bridge inherit from eth bridge
type Bridge struct {
	*eth.Bridge
}

// NewCrossChainBridge new fsn bridge
func NewCrossChainBridge(isSrc bool) *Bridge {
	return &Bridge{Bridge: eth.NewCrossChainBridge(isSrc)}
}

// SetTokenAndGateway set token and gateway config
func (b *Bridge) SetTokenAndGateway(tokenCfg *tokens.TokenConfig, gatewayCfg *tokens.GatewayConfig, check bool) {
	b.CrossChainBridgeBase.SetTokenAndGateway(tokenCfg, gatewayCfg, check)
	b.VerifyChainID()
	b.VerifyConfig()
	b.Init()
}

// VerifyChainID verify chain id
func (b *Bridge) VerifyChainID() {
	tokenCfg := b.TokenConfig
	gatewayCfg := b.GatewayConfig

	networkID := strings.ToLower(tokenCfg.NetID)

	switch networkID {
	case netMainnet, netTestnet, netDevnet:
	case netCustom:
	default:
		log.Fatalf("unsupported fusion network: %v", tokenCfg.NetID)
	}

	var (
		chainID *big.Int
		err     error
	)

	for {
		// call NetworkID instead of ChainID as ChainID may return 0x0 wrongly
		chainID, err = b.NetworkID()
		if err == nil {
			break
		}
		log.Errorf("can not get gateway chainID. %v", err)
		log.Println("retry query gateway", gatewayCfg.APIAddress)
		time.Sleep(3 * time.Second)
	}

	panicMismatchChainID := func() {
		log.Fatalf("gateway chainID %v is not %v", chainID, tokenCfg.NetID)
	}

	switch networkID {
	case netMainnet:
		if chainID.Uint64() != 32659 {
			panicMismatchChainID()
		}
	case netTestnet:
		if chainID.Uint64() != 46688 {
			panicMismatchChainID()
		}
	case netDevnet:
		if chainID.Uint64() != 55555 {
			panicMismatchChainID()
		}
	case netCustom:
	default:
		log.Fatalf("unsupported fusion network %v", networkID)
	}

	b.Signer = types.MakeSigner("EIP155", chainID)

	log.Info("VerifyChainID succeed", "networkID", networkID, "chainID", chainID)
}
