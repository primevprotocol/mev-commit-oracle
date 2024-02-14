package main

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	contracts "github.com/primevprotocol/contracts-abi/config"
	"github.com/primevprotocol/mev-oracle/pkg/keysigner"
	"github.com/primevprotocol/mev-oracle/pkg/node"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"golang.org/x/crypto/sha3"
)

const (
	defaultHTTPPort  = 8080
	defaultConfigDir = "~/.mev-commit-oracle"
	defaultKeyFile   = "key"
	defaultKeystore  = "keystore"
)

var (
	optionConfig = &cli.StringFlag{
		Name:    "config",
		Usage:   "path to config file",
		EnvVars: []string{"MEV_ORACLE_CONFIG"},
	}

	optionPrivKeyFile = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "priv-key-file",
		Usage:   "path to private key file",
		EnvVars: []string{"MEV_ORACLE_PRIV_KEY_FILE"},
		Value:   filepath.Join(defaultConfigDir, defaultKeyFile),
	})

	optionHTTPPort = altsrc.NewIntFlag(&cli.IntFlag{
		Name:    "http-port",
		Usage:   "port to listen on for HTTP requests",
		EnvVars: []string{"MEV_ORACLE_HTTP_PORT"},
		Value:   defaultHTTPPort,
		Action: func(c *cli.Context, p int) error {
			if p < 0 || p > 65535 {
				return fmt.Errorf("invalid port number: %d", p)
			}
			return nil
		},
	})

	optionLogLevel = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-level",
		Usage:   "log level",
		EnvVars: []string{"MEV_ORACLE_LOG_LEVEL"},
		Value:   "info",
		Action: func(c *cli.Context, l string) error {
			_, err := zerolog.ParseLevel(l)
			return err
		},
	})

	optionL1RPCUrl = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "l1-rpc-url",
		Usage:   "URL for L1 RPC",
		EnvVars: []string{"MEV_ORACLE_L1_RPC_URL"},
	})

	optionSettlementRPCUrl = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "settlement-rpc-url",
		Usage:   "URL for settlement RPC",
		EnvVars: []string{"MEV_ORACLE_SETTLEMENT_RPC_URL"},
		Value:   "http://localhost:8545",
	})

	optionOracleContractAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "oracle-contract-addr",
		Usage:   "address of the oracle contract",
		EnvVars: []string{"MEV_ORACLE_ORACLE_CONTRACT_ADDR"},
		Value:   contracts.TestnetContracts.Oracle,
	})

	optionPreconfContractAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "preconf-contract-addr",
		Usage:   "address of the preconf contract",
		EnvVars: []string{"MEV_ORACLE_PRECONF_CONTRACT_ADDR"},
		Value:   contracts.TestnetContracts.PreconfCommitmentStore,
	})

	optionPgHost = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "pg-host",
		Usage:   "PostgreSQL host",
		EnvVars: []string{"MEV_ORACLE_PG_HOST"},
		Value:   "localhost",
	})

	optionPgPort = altsrc.NewIntFlag(&cli.IntFlag{
		Name:    "pg-port",
		Usage:   "PostgreSQL port",
		EnvVars: []string{"MEV_ORACLE_PG_PORT"},
		Value:   5432,
	})

	optionPgUser = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "pg-user",
		Usage:   "PostgreSQL user",
		EnvVars: []string{"MEV_ORACLE_PG_USER"},
		Value:   "postgres",
	})

	optionPgPassword = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "pg-password",
		Usage:   "PostgreSQL password",
		EnvVars: []string{"MEV_ORACLE_PG_PASSWORD"},
		Value:   "postgres",
	})

	optionPgDbname = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "pg-dbname",
		Usage:   "PostgreSQL database name",
		EnvVars: []string{"MEV_ORACLE_PG_DBNAME"},
		Value:   "mev_oracle",
	})

	optionLaggerdMode = altsrc.NewIntFlag(&cli.IntFlag{
		Name:    "laggerd-mode",
		Usage:   "No of blocks to lag behind for L1 chain",
		EnvVars: []string{"MEV_ORACLE_LAGGERD_MODE"},
		Value:   0,
	})

	optionOverrideWinners = altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
		Name:    "override-winners",
		Usage:   "Override winners for testing",
		EnvVars: []string{"MEV_ORACLE_OVERRIDE_WINNERS"},
	})

	optionKeystorePassword = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "keystore-password",
		Usage:   "use to access keystore",
		EnvVars: []string{"MEV_ORACLE_KEYSTORE_PASSWORD"},
	})

	optionKeystorePath = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "keystore-path",
		Usage:   "path to keystore location",
		EnvVars: []string{"MEV_ORACLE_KEYSTORE_PATH"},
		Value:   filepath.Join(defaultConfigDir, defaultKeystore),
	})
)

func main() {
	flags := []cli.Flag{
		optionConfig,
		optionPrivKeyFile,
		optionHTTPPort,
		optionLogLevel,
		optionL1RPCUrl,
		optionSettlementRPCUrl,
		optionOracleContractAddr,
		optionPreconfContractAddr,
		optionPgHost,
		optionPgPort,
		optionPgUser,
		optionPgPassword,
		optionPgDbname,
		optionLaggerdMode,
		optionOverrideWinners,
		optionKeystorePath,
		optionKeystorePassword,
	}
	app := &cli.App{
		Name:  "mev-oracle",
		Usage: "Entry point for mev-oracle",
		Commands: []*cli.Command{
			{
				Name:   "start",
				Usage:  "Start the mev-oracle node",
				Flags:  flags,
				Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc(optionConfig.Name)),
				Action: func(c *cli.Context) error {
					return initializeApplication(c)
				},
			},
		}}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(app.Writer, "exited with error: %v\n", err)
	}
}

func createKeyIfNotExists(c *cli.Context, path string) error {
	// check if key already exists
	if _, err := os.Stat(path); err == nil {
		fmt.Fprintf(c.App.Writer, "Using existing private key: %s\n", path)
		return nil
	}

	fmt.Fprintf(c.App.Writer, "Creating new private key: %s\n", path)

	// check if parent directory exists
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		// create parent directory
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return err
		}
	}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	if err := crypto.SaveECDSA(path, privKey); err != nil {
		return err
	}

	wallet := getEthAddressFromPubKey(&privKey.PublicKey)

	fmt.Fprintf(c.App.Writer, "Private key saved to file: %s\n", path)
	fmt.Fprintf(c.App.Writer, "Wallet address: %s\n", wallet.Hex())
	return nil
}

func resolveFilePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}

	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		return filepath.Join(home, path[1:]), nil
	}

	return path, nil
}

func getEthAddressFromPubKey(key *ecdsa.PublicKey) common.Address {
	pbBytes := crypto.FromECDSAPub(key)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pbBytes[1:])
	address := hash.Sum(nil)[12:]

	return common.BytesToAddress(address)
}

func initializeApplication(c *cli.Context) error {
	if err := verifyKeystorePasswordPresence(c); err != nil {
		return err
	}
	if err := launchOracleWithConfig(c); err != nil {
		return err
	}
	return nil
}

// verifyKeystorePasswordPresence checks for the presence of a keystore password.
// it returns error, if keystore path is set and keystore password is not
func verifyKeystorePasswordPresence(c *cli.Context) error {
	if c.IsSet(optionKeystorePath.Name) && !c.IsSet(optionKeystorePassword.Name) {
		return cli.Exit("Password for encrypted keystore is missing", 1)
	}
	return nil
}

// launchOracleWithConfig configures and starts the oracle based on the CLI context or config.yaml file.
func launchOracleWithConfig(c *cli.Context) error {
	keySigner, err := setupKeySigner(c)
	if err != nil {
		return fmt.Errorf("failed to setup key signer: %w", err)
	}

	lvl, _ := zerolog.ParseLevel(c.String(optionLogLevel.Name))

	zerolog.SetGlobalLevel(lvl)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(os.Stdout).With().Caller().Logger()

	nd, err := node.NewNode(&node.Options{
		KeySigner:           keySigner,
		HTTPPort:            c.Int(optionHTTPPort.Name),
		L1RPCUrl:            c.String(optionL1RPCUrl.Name),
		SettlementRPCUrl:    c.String(optionSettlementRPCUrl.Name),
		OracleContractAddr:  common.HexToAddress(c.String(optionOracleContractAddr.Name)),
		PreconfContractAddr: common.HexToAddress(c.String(optionPreconfContractAddr.Name)),
		PgHost:              c.String(optionPgHost.Name),
		PgPort:              c.Int(optionPgPort.Name),
		PgUser:              c.String(optionPgUser.Name),
		PgPassword:          c.String(optionPgPassword.Name),
		PgDbname:            c.String(optionPgDbname.Name),
		LaggerdMode:         c.Int(optionLaggerdMode.Name),
		OverrideWinners:     c.StringSlice(optionOverrideWinners.Name),
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

func setupKeySigner(c *cli.Context) (keysigner.KeySigner, error) {
	if c.IsSet(optionKeystorePath.Name) {
		return setupKeystoreSigner(c)
	}
	return setupPrivateKeySigner(c)
}

func setupKeystoreSigner(c *cli.Context) (keysigner.KeySigner, error) {
	// Load the keystore file
	ks := keystore.NewKeyStore(c.String(optionKeystorePath.Name), keystore.LightScryptN, keystore.LightScryptP)
	password := c.String(optionKeystorePassword.Name)
	ksAccounts := ks.Accounts()

	var account accounts.Account
	if len(ksAccounts) == 0 {
		var err error
		account, err = ks.NewAccount(password)
		if err != nil {
			return nil, fmt.Errorf("failed to create account: %w", err)
		}
	} else {
		account = ksAccounts[0]
	}

	fmt.Fprintf(c.App.Writer, "Public address of the key: %s\n", account.Address.Hex())
	fmt.Fprintf(c.App.Writer, "Path of the secret key file: %s\n", account.URL.Path)

	return keysigner.NewKeystoreSigner(ks, password, account), nil
}

func setupPrivateKeySigner(c *cli.Context) (keysigner.KeySigner, error) {
	privKeyFile, err := resolveFilePath(c.String(optionPrivKeyFile.Name))
	if err != nil {
		return nil, fmt.Errorf("failed to get private key file path: %w", err)
	}

	if err := createKeyIfNotExists(c, privKeyFile); err != nil {
		return nil, fmt.Errorf("failed to create private key: %w", err)
	}

	privKey, err := crypto.LoadECDSA(privKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key from file '%s': %w", privKeyFile, err)
	}

	return keysigner.NewPrivateKeySigner(privKey), nil
}
