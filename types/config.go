package types

import (
	"fmt"
	"os"

	"github.com/irisnet/core-sdk-go/common/crypto"
	"github.com/irisnet/core-sdk-go/types/store"
)

const (
	defaultGas           = 200000
	defaultFees          = "4iris"
	defaultTimeout       = 5
	defaultLevel         = "info"
	defaultMaxTxsBytes   = 1073741824
	defaultAlgo          = "secp256k1"
	defaultMode          = Sync
	defaultPath          = "$HOME/irishub-sdk-go/leveldb"
	defaultGasAdjustment = 1.0
	defaultTxSizeLimit   = 1048576
	BIP44Prefix          = "44'/118'/"
	PartialPath          = "0'/0/0"
	FullPath             = BIP44Prefix + PartialPath
)

type ClientConfig struct {
	// irishub node rpc address
	NodeURI string

	// irishub grpc address
	GRPCAddr string

	// irishub chain-id
	ChainID string

	// max gas limit
	Gas uint64

	// Fee amount of point
	Fee DecCoins

	// PrivKeyArmor DAO Implements
	KeyDAO store.KeyDAO

	// Private key generation algorithm(sm2,secp256k1)
	Algo string

	// Transaction broadcast Mode
	Mode BroadcastMode

	// Transaction broadcast timeout(seconds)
	Timeout uint

	// log level(trace|debug|info|warn|error|fatal|panic)
	Level string

	// maximum bytes of a transaction
	MaxTxBytes uint64

	// adjustment factor to be multiplied against the estimate returned by the tx simulation;
	GasAdjustment float64

	// whether to enable caching
	Cached bool

	TokenManager TokenManager

	KeyManager crypto.KeyManager

	TxSizeLimit uint64

	// bech32 Address Prefix
	Bech32AddressPrefix AddrPrefixCfg

	// BIP44 path
	BIP44Path string

	// BSN ProjectId ProjectKey ChainAccountAddress
	BSNProject BSNProjectInfo
}

type BSNProjectInfo struct {
	ProjectId           string
	ProjectKey          string
	ChainAccountAddress string
}

func NewClientConfig(uri, grpcAddr, chainID string, options ...Option) (ClientConfig, error) {
	cfg := ClientConfig{
		NodeURI:  uri,
		ChainID:  chainID,
		GRPCAddr: grpcAddr,
	}
	for _, optionFn := range options {
		if err := optionFn(&cfg); err != nil {
			return ClientConfig{}, err
		}
	}

	if err := cfg.checkAndSetDefault(); err != nil {
		return ClientConfig{}, err
	}
	return cfg, nil
}

func (cfg *ClientConfig) checkAndSetDefault() error {
	if len(cfg.NodeURI) == 0 {
		return fmt.Errorf("nodeURI is required")
	}

	if len(cfg.ChainID) == 0 {
		return fmt.Errorf("chainID is required")
	}

	if err := GasOption(cfg.Gas)(cfg); err != nil {
		return err
	}

	if err := FeeOption(cfg.Fee)(cfg); err != nil {
		return err
	}

	if err := AlgoOption(cfg.Algo)(cfg); err != nil {
		return err
	}

	if err := KeyDAOOption(cfg.KeyDAO)(cfg); err != nil {
		return err
	}

	if err := ModeOption(cfg.Mode)(cfg); err != nil {
		return err
	}

	if err := TimeoutOption(cfg.Timeout)(cfg); err != nil {
		return err
	}

	if err := LevelOption(cfg.Level)(cfg); err != nil {
		return err
	}

	if err := MaxTxBytesOption(cfg.MaxTxBytes)(cfg); err != nil {
		return err
	}

	if err := TokenManagerOption(cfg.TokenManager)(cfg); err != nil {
		return err
	}

	if err := TxSizeLimitOption(cfg.TxSizeLimit)(cfg); err != nil {
		return err
	}

	if err := Bech32AddressPrefixOption(cfg.Bech32AddressPrefix)(cfg); err != nil {
		return err
	}

	if err := BIP44PathOption(cfg.BIP44Path)(cfg); err != nil {
		return err
	}
	if err := BSNProjectInfoOption(cfg.BSNProject)(cfg); err != nil {
		return err
	}
	return GasAdjustmentOption(cfg.GasAdjustment)(cfg)
}

type Option func(cfg *ClientConfig) error

func FeeOption(fee DecCoins) Option {
	return func(cfg *ClientConfig) error {
		if fee == nil || fee.Empty() || !fee.IsValid() {
			fees, _ := ParseDecCoins(defaultFees)
			fee = fees
		}
		cfg.Fee = fee
		return nil
	}
}

func KeyDAOOption(dao store.KeyDAO) Option {
	return func(cfg *ClientConfig) error {
		if dao == nil {
			defaultPath := os.ExpandEnv(defaultPath)
			levelDB, err := store.NewLevelDB(defaultPath, nil)
			if err != nil {
				return err
			}
			dao = levelDB
		}
		cfg.KeyDAO = dao
		return nil
	}
}

func GasOption(gas uint64) Option {
	return func(cfg *ClientConfig) error {
		if gas <= 0 {
			gas = defaultGas
		}
		cfg.Gas = gas
		return nil
	}
}

func AlgoOption(algo string) Option {
	return func(cfg *ClientConfig) error {
		if algo == "" {
			algo = defaultAlgo
		}
		cfg.Algo = algo
		return nil
	}
}

func ModeOption(mode BroadcastMode) Option {
	return func(cfg *ClientConfig) error {
		if mode == "" {
			mode = defaultMode
		}
		cfg.Mode = mode
		return nil
	}
}

func TimeoutOption(timeout uint) Option {
	return func(cfg *ClientConfig) error {
		if timeout <= 0 {
			timeout = defaultTimeout
		}
		cfg.Timeout = timeout
		return nil
	}
}

func LevelOption(level string) Option {
	return func(cfg *ClientConfig) error {
		if level == "" {
			level = defaultLevel
		}
		cfg.Level = level
		return nil
	}
}

func MaxTxBytesOption(maxTxBytes uint64) Option {
	return func(cfg *ClientConfig) error {
		if maxTxBytes <= 0 {
			maxTxBytes = defaultMaxTxsBytes
		}
		cfg.MaxTxBytes = maxTxBytes
		return nil
	}
}

func GasAdjustmentOption(gasAdjustment float64) Option {
	return func(cfg *ClientConfig) error {
		if gasAdjustment <= 0 {
			gasAdjustment = defaultGasAdjustment
		}
		cfg.GasAdjustment = gasAdjustment
		return nil
	}
}

func CachedOption(enabled bool) Option {
	return func(cfg *ClientConfig) error {
		cfg.Cached = enabled
		return nil
	}
}

func TokenManagerOption(tokenManager TokenManager) Option {
	return func(cfg *ClientConfig) error {
		if tokenManager == nil {
			tokenManager = DefaultTokenManager{}
		}
		cfg.TokenManager = tokenManager
		return nil
	}
}

func TxSizeLimitOption(txSizeLimit uint64) Option {
	return func(cfg *ClientConfig) error {
		if txSizeLimit <= 0 {
			txSizeLimit = defaultTxSizeLimit
		}
		cfg.TxSizeLimit = txSizeLimit
		return nil
	}
}

func KeyManagerOption(keyManager crypto.KeyManager) Option {
	return func(cfg *ClientConfig) error {
		cfg.KeyManager = keyManager
		return nil
	}
}

func Bech32AddressPrefixOption(bech32AddressPrefix AddrPrefixCfg) Option {
	return func(cfg *ClientConfig) error {

		if bech32AddressPrefix.AccountAddr == "" || bech32AddressPrefix.ValidatorAddr == "" || bech32AddressPrefix.ConsensusAddr == "" {
			bech32AddressPrefix = *PrefixCfg
		}
		if bech32AddressPrefix.AccountPub == "" || bech32AddressPrefix.ValidatorPub == "" || bech32AddressPrefix.ConsensusPub == "" {
			bech32AddressPrefix = *PrefixCfg
		}
		cfg.Bech32AddressPrefix = bech32AddressPrefix
		return nil
	}
}

func BIP44PathOption(bIP44Path string) Option {
	return func(cfg *ClientConfig) error {
		if bIP44Path == "" {
			bIP44Path = FullPath
		}
		cfg.BIP44Path = bIP44Path
		return nil
	}
}

func BSNProjectInfoOption(info BSNProjectInfo) Option {
	return func(cfg *ClientConfig) error {
		cfg.BSNProject.ProjectId = info.ProjectId
		cfg.BSNProject.ProjectKey = info.ProjectKey
		cfg.BSNProject.ChainAccountAddress = info.ChainAccountAddress
		return nil
	}

}
