package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/primevprotocol/mev-oracle/pkg/keysigner"
	"github.com/primevprotocol/mev-oracle/pkg/node"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

const (
	defaultHTTPPort = 8080
)

var (
	optionConfig = &cli.StringFlag{
		Name:     "config",
		Usage:    "path to config file",
		Required: true,
		EnvVars:  []string{"MEV_ORACLE_CONFIG"},
	}
)

func main() {
	app := &cli.App{
		Name:  "mev-oracle",
		Usage: "Entry point for mev-oracle",
		Commands: []*cli.Command{
			{
				Name:  "start",
				Usage: "Start the mev-oracle node",
				Flags: []cli.Flag{
					optionConfig,
				},
				Action: func(c *cli.Context) error {
					return start(c)
				},
			},
		}}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(app.Writer, "exited with error: %v\n", err)
	}
}

type config struct {
	PrivKeyFile         string   `yaml:"priv_key_file" json:"priv_key_file"`
	KeystorePath        string   `yaml:"keystore_path" json:"keystore_path"`
	KeystorePassword    string   `yaml:"keystore_password" json:"keystore_password"`
	HTTPPort            int      `yaml:"http_port" json:"http_port"`
	LogLevel            string   `yaml:"log_level" json:"log_level"`
	L1RPCUrl            string   `yaml:"l1_rpc_url" json:"l1_rpc_url"`
	SettlementRPCUrl    string   `yaml:"settlement_rpc_url" json:"settlement_rpc_url"`
	OracleContractAddr  string   `yaml:"oracle_contract_addr" json:"oracle_contract_addr"`
	PreconfContractAddr string   `yaml:"preconf_contract_addr" json:"preconf_contract_addr"`
	PgHost              string   `yaml:"pg_host" json:"pg_host"`
	PgPort              int      `yaml:"pg_port" json:"pg_port"`
	PgUser              string   `yaml:"pg_user" json:"pg_user"`
	PgPassword          string   `yaml:"pg_password" json:"pg_password"`
	PgDbname            string   `yaml:"pg_dbname" json:"pg_dbname"`
	LaggerdMode         int      `yaml:"laggerd_mode" json:"laggerd_mode"`
	OverrideWinners     []string `yaml:"override_winners" json:"override_winners"`
}

func checkConfig(cfg *config) error {
	if cfg.PrivKeyFile == "" && (cfg.KeystorePath == "" || cfg.KeystorePassword == "") {
		return fmt.Errorf("priv_key_file or keystore_path and keystore_password are required")
	}

	if cfg.HTTPPort == 0 {
		cfg.HTTPPort = defaultHTTPPort
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	if cfg.L1RPCUrl == "" {
		return fmt.Errorf("l1_rpc_url is required")
	}

	if cfg.SettlementRPCUrl == "" {
		return fmt.Errorf("settlement_rpc_url is required")
	}

	if cfg.OracleContractAddr == "" {
		return fmt.Errorf("oracle_contract_addr is required")
	}

	if cfg.PreconfContractAddr == "" {
		return fmt.Errorf("preconf_contract_addr is required")
	}

	if cfg.PgHost == "" || cfg.PgPort == 0 || cfg.PgUser == "" || cfg.PgPassword == "" || cfg.PgDbname == "" {
		return fmt.Errorf("pg_host, pg_port, pg_user, pg_password, pg_dbname are required")
	}

	return nil
}

func start(c *cli.Context) error {
	configFile := c.String(optionConfig.Name)
	fmt.Fprintf(c.App.Writer, "starting mev-oracle with config file: %s\n", configFile)

	var cfg config
	buf, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file at '%s': %w", configFile, err)
	}

	if err := yaml.Unmarshal(buf, &cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config file at '%s': %w", configFile, err)
	}

	if err := checkConfig(&cfg); err != nil {
		return fmt.Errorf("invalid config file at '%s': %w", configFile, err)
	}

	lvl, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to parse log level '%s': %w", cfg.LogLevel, err)
	}

	zerolog.SetGlobalLevel(lvl)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(os.Stdout).With().Caller().Logger()

	keySigner, err := setupKeySigner(&cfg)
	if err != nil {
		return fmt.Errorf("failed to setup key signer: %w", err)
	}

	common.HexToAddress(cfg.OracleContractAddr)

	nd, err := node.NewNode(&node.Options{
		KeySigner:           keySigner,
		HTTPPort:            cfg.HTTPPort,
		L1RPCUrl:            cfg.L1RPCUrl,
		SettlementRPCUrl:    cfg.SettlementRPCUrl,
		OracleContractAddr:  common.HexToAddress(cfg.OracleContractAddr),
		PreconfContractAddr: common.HexToAddress(cfg.PreconfContractAddr),
		PgHost:              cfg.PgHost,
		PgPort:              cfg.PgPort,
		PgUser:              cfg.PgUser,
		PgPassword:          cfg.PgPassword,
		PgDbname:            cfg.PgDbname,
		LaggerdMode:         cfg.LaggerdMode,
		OverrideWinners:     cfg.OverrideWinners,
	})
	if err != nil {
		return fmt.Errorf("failed starting node: %w", err)
	}

	<-c.Done()
	fmt.Fprintf(c.App.Writer, "shutting down...\n")
	closed := make(chan struct{})

	go func() {
		defer close(closed)

		err := nd.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close node")
		}
	}()

	select {
	case <-closed:
	case <-time.After(5 * time.Second):
		log.Error().Msg("failed to close node in time")
	}

	return nil
}

func setupKeySigner(cfg *config) (keysigner.KeySigner, error) {
	if cfg.PrivKeyFile == "" {
		return setupKeystoreSigner(cfg)
	}
	return setupPrivateKeySigner(cfg)
}

func setupKeystoreSigner(cfg *config) (keysigner.KeySigner, error) {
	// Load the keystore file
	ks := keystore.NewKeyStore(cfg.KeystorePath, keystore.LightScryptN, keystore.LightScryptP)
	accounts := ks.Accounts()
	if len(accounts) == 0 {
		return nil, fmt.Errorf("no accounts found in keystore, path: %s", cfg.KeystorePath)
	}

	account := accounts[0]
	return keysigner.NewKeystoreSigner(ks, cfg.KeystorePassword, account), nil
}

func setupPrivateKeySigner(cfg *config) (keysigner.KeySigner, error) {
	privKeyFile := cfg.PrivKeyFile
	if strings.HasPrefix(privKeyFile, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}

		privKeyFile = filepath.Join(homeDir, privKeyFile[2:])
	}

	privKey, err := crypto.LoadECDSA(privKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key from file '%s': %w", cfg.PrivKeyFile, err)
	}
	return keysigner.NewPrivateKeySigner(privKey), nil
}
