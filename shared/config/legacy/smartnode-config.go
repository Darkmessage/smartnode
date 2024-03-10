package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/smartnode/shared"
)

// Constants
const (
	smartnodeTag                       string = "rocketpool/smartnode:v" + shared.RocketPoolVersion
	pruneProvisionerTag                string = "rocketpool/eth1-prune-provision:v0.0.1"
	ecMigratorTag                      string = "rocketpool/ec-migrator:v1.0.0"
	NetworkID                          string = "network"
	ProjectNameID                      string = "projectName"
	SnapshotID                         string = "rocketpool-dao.eth"
	RewardsTreeFilenameFormat          string = "rp-rewards-%s-%d.json"
	MinipoolPerformanceFilenameFormat  string = "rp-minipool-performance-%s-%d.json"
	RewardsTreeIpfsExtension           string = ".zst"
	RewardsTreesFolder                 string = "rewards-trees"
	ChecksumTableFilename              string = "checksums.sha384"
	DaemonDataPath                     string = "/.rocketpool/data"
	WatchtowerFolder                   string = "watchtower"
	WatchtowerStateFile                string = "state.yml"
	RegenerateRewardsTreeRequestSuffix string = ".request"
	RegenerateRewardsTreeRequestFormat string = "%d" + RegenerateRewardsTreeRequestSuffix
	PrimaryRewardsFileUrl              string = "https://%s.ipfs.dweb.link/%s"
	SecondaryRewardsFileUrl            string = "https://ipfs.io/ipfs/%s/%s"
	GithubRewardsFileUrl               string = "https://github.com/rocket-pool/rewards-trees/raw/main/%s/%s"
	FeeRecipientFilename               string = "rp-fee-recipient.txt"
	NativeFeeRecipientFilename         string = "rp-fee-recipient-env.txt"
	SocketFilename                     string = "daemon.sock"
	NodeAddressFilename                string = "address"
	WalletKeystoreFilename             string = "wallet"
	WalletPasswordFilename             string = "password"
)

// Defaults
const (
	defaultProjectName       string = "rocketpool"
	WatchtowerMaxFeeDefault  uint64 = 200
	WatchtowerPrioFeeDefault uint64 = 3
)

// Configuration for the Smartnode
type SmartnodeConfig struct {
	Title string `yaml:"-"`

	// The parent config
	parent *RocketPoolConfig

	////////////////////////////
	// User-editable settings //
	////////////////////////////

	// Docker container prefix
	ProjectName Parameter `yaml:"projectName,omitempty"`

	// The path of the data folder where everything is stored
	DataPath Parameter `yaml:"dataPath,omitempty"`

	// The path of the watchtower's persistent state storage
	WatchtowerStatePath Parameter `yaml:"watchtowerStatePath"`

	// Which network we're on
	Network Parameter `yaml:"network,omitempty"`

	// Manual max fee override
	ManualMaxFee Parameter `yaml:"manualMaxFee,omitempty"`

	// Manual priority fee override
	PriorityFee Parameter `yaml:"priorityFee,omitempty"`

	// Threshold for automatic transactions
	AutoTxGasThreshold Parameter `yaml:"minipoolStakeGasThreshold,omitempty"`

	// The amount of ETH in a minipool's balance before auto-distribute kicks in
	DistributeThreshold Parameter `yaml:"distributeThreshold,omitempty"`

	// Mode for acquiring Merkle rewards trees
	RewardsTreeMode Parameter `yaml:"rewardsTreeMode,omitempty"`

	// Custom URL to download a rewards tree
	RewardsTreeCustomUrl Parameter `yaml:"rewardsTreeCustomUrl,omitempty"`

	// URL for an EC with archive mode, for manual rewards tree generation
	ArchiveECUrl Parameter `yaml:"archiveEcUrl,omitempty"`

	// Manual override for the watchtower's max fee
	WatchtowerMaxFeeOverride Parameter `yaml:"watchtowerMaxFeeOverride,omitempty"`

	// Manual override for the watchtower's priority fee
	WatchtowerPrioFeeOverride Parameter `yaml:"watchtowerPrioFeeOverride,omitempty"`

	// The toggle for rolling records
	UseRollingRecords Parameter `yaml:"useRollingRecords,omitempty"`

	// The rolling record checkpoint interval
	RecordCheckpointInterval Parameter `yaml:"recordCheckpointInterval,omitempty"`

	// The checkpoint retention limit
	CheckpointRetentionLimit Parameter `yaml:"checkpointRetentionLimit,omitempty"`

	// The path of the records folder where snapshots of rolling record info is stored during a rewards interval
	RecordsPath Parameter `yaml:"recordsPath,omitempty"`

	// The toggle for enabling pDAO proposal verification duties
	VerifyProposals Parameter `yaml:"verifyProposals,omitempty"`

	///////////////////////////
	// Non-editable settings //
	///////////////////////////

	// The URL to provide the user so they can follow pending transactions
	txWatchUrl map[Network]string `yaml:"-"`

	// The URL to use for staking rETH
	stakeUrl map[Network]string `yaml:"-"`

	// The map of networks to execution chain IDs
	chainID map[Network]uint `yaml:"-"`

	// The contract address of RocketStorage
	storageAddress map[Network]string `yaml:"-"`

	// The contract address of the RPL token
	rplTokenAddress map[Network]string `yaml:"-"`

	// The contract address of the RPL faucet
	rplFaucetAddress map[Network]string `yaml:"-"`

	// The contract address for Snapshot delegation
	snapshotDelegationAddress map[Network]string `yaml:"-"`

	// The Snapshot API domain
	snapshotApiDomain map[Network]string `yaml:"-"`

	// The contract address of rETH
	rethAddress map[Network]string `yaml:"-"`

	// The contract address of rocketRewardsPool from v1.0.0
	v1_0_0_RewardsPoolAddress map[Network]string `yaml:"-"`

	// The contract address of rocketClaimNode from v1.0.0
	v1_0_0_ClaimNodeAddress map[Network]string `yaml:"-"`

	// The contract address of rocketClaimTrustedNode from v1.0.0
	v1_0_0_ClaimTrustedNodeAddress map[Network]string `yaml:"-"`

	// The contract address of rocketMinipoolManager from v1.0.0
	v1_0_0_MinipoolManagerAddress map[Network]string `yaml:"-"`

	// The contract address of rocketNetworkPrices from v1.1.0
	v1_1_0_NetworkPricesAddress map[Network]string `yaml:"-"`

	// The contract address of rocketNodeStaking from v1.1.0
	v1_1_0_NodeStakingAddress map[Network]string `yaml:"-"`

	// The contract address of rocketNodeDeposit from v1.1.0
	v1_1_0_NodeDepositAddress map[Network]string `yaml:"-"`

	// The contract address of rocketMinipoolQueue from v1.1.0
	v1_1_0_MinipoolQueueAddress map[Network]string `yaml:"-"`

	// The contract address of rocketMinipoolFactory from v1.1.0
	v1_1_0_MinipoolFactoryAddress map[Network]string `yaml:"-"`

	// The contract address of rocketNetworkPrices from v1.2.0
	v1_2_0_NetworkPricesAddress map[Network]string `yaml:"-"`

	// The contract address of rocketNetworkBalances from v1.2.0
	v1_2_0_NetworkBalancesAddress map[Network]string `yaml:"-"`

	// Addresses for RocketRewardsPool that have been upgraded during development
	previousRewardsPoolAddresses map[Network][]common.Address `yaml:"-"`

	// Addresses for RocketNetworkPrices that have been upgraded during development
	previousRocketNetworkPricesAddresses map[Network][]common.Address `yaml:"-"`

	// Addresses for RocketNetworkBalances that have been upgraded during development
	previousRocketNetworkBalancesAddresses map[Network][]common.Address `yaml:"-"`

	// The RocketOvmPriceMessenger Optimism address for each network
	optimismPriceMessengerAddress map[Network]string `yaml:"-"`

	// The RocketPolygonPriceMessenger Polygon address for each network
	polygonPriceMessengerAddress map[Network]string `yaml:"-"`

	// The RocketArbitumPriceMessenger Arbitrum address for each network
	arbitrumPriceMessengerAddress map[Network]string `yaml:"-"`

	// The RocketArbitumPriceMessengerV2 Arbitrum address for each network
	arbitrumPriceMessengerAddressV2 map[Network]string `yaml:"-"`

	// The RocketZkSyncPriceMessenger zkSyncEra address for each network
	zkSyncEraPriceMessengerAddress map[Network]string `yaml:"-"`

	// The RocketBasePriceMessenger Base address for each network
	basePriceMessengerAddress map[Network]string `yaml:"-"`

	// The RocketScrollPriceMessenger Scroll address for each network
	scrollPriceMessengerAddress map[Network]string `yaml:"-"`

	// The Scroll L2 message fee estimator address for each network
	scrollFeeEstimatorAddress map[Network]string `yaml:"-"`

	// The UniswapV3 pool address for each network (used for RPL price TWAP info)
	rplTwapPoolAddress map[Network]string `yaml:"-"`

	// The multicall contract address
	multicallAddress map[Network]string `yaml:"-"`

	// The BalanceChecker contract address
	balancebatcherAddress map[Network]string `yaml:"-"`

	// The FlashBots Protect RPC endpoint
	flashbotsProtectUrl map[Network]string `yaml:"-"`
}

// Generates a new Smartnode configuration
func NewSmartnodeConfig(cfg *RocketPoolConfig) *SmartnodeConfig {

	return &SmartnodeConfig{
		Title:  "Smartnode Settings",
		parent: cfg,

		ProjectName: Parameter{
			ID:                 ProjectNameID,
			Name:               "Project Name",
			Description:        "This is the prefix that will be attached to all of the Docker containers managed by the Smartnode.",
			Type:               ParameterType_String,
			Default:            map[Network]interface{}{Network_All: defaultProjectName},
			AffectsContainers:  []ContainerID{ContainerID_Api, ContainerID_Node, ContainerID_Watchtower, ContainerID_Eth1, ContainerID_Eth2, ContainerID_Validator, ContainerID_Grafana, ContainerID_Prometheus, ContainerID_Exporter},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},

		DataPath: Parameter{
			ID:                 "dataPath",
			Name:               "Data Path",
			Description:        "The absolute path of the `data` folder that contains your node wallet's encrypted file, the password for your node wallet, and all of the validator keys for your minipools. You may use environment variables in this string.",
			Type:               ParameterType_String,
			Default:            map[Network]interface{}{Network_All: getDefaultDataDir(cfg)},
			AffectsContainers:  []ContainerID{ContainerID_Api, ContainerID_Node, ContainerID_Watchtower, ContainerID_Validator},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},

		WatchtowerStatePath: Parameter{
			ID:                 "watchtowerPath",
			Name:               "Watchtower Path",
			Description:        "The absolute path of the watchtower state folder that contains persistent state that is used by the watchtower process on trusted nodes. **Only relevant for trusted nodes.**",
			Type:               ParameterType_String,
			Default:            map[Network]interface{}{Network_All: "$HOME/.rocketpool/watchtower"},
			AffectsContainers:  []ContainerID{ContainerID_Watchtower},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},

		Network: Parameter{
			ID:                 NetworkID,
			Name:               "Network",
			Description:        "The Ethereum network you want to use - select Prater Testnet or Holesky Testnet to practice with fake ETH, or Mainnet to stake on the real network using real ETH.",
			Type:               ParameterType_Choice,
			Default:            map[Network]interface{}{Network_All: Network_Mainnet},
			AffectsContainers:  []ContainerID{ContainerID_Api, ContainerID_Node, ContainerID_Watchtower, ContainerID_Eth1, ContainerID_Eth2, ContainerID_Validator},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
			Options:            getNetworkOptions(),
		},

		ManualMaxFee: Parameter{
			ID:                 "manualMaxFee",
			Name:               "Manual Max Fee",
			Description:        "Set this if you want all of the Smartnode's transactions to use this specific max fee value (in gwei), which is the most you'd be willing to pay (*including the priority fee*).\n\nA value of 0 will show you the current suggested max fee based on the current network conditions and let you specify it each time you do a transaction.\n\nAny other value will ignore the recommended max fee and explicitly use this value instead.\n\nThis applies to automated transactions (such as claiming RPL and staking minipools) as well.",
			Type:               ParameterType_Float,
			Default:            map[Network]interface{}{Network_All: float64(0)},
			AffectsContainers:  []ContainerID{ContainerID_Node, ContainerID_Watchtower},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},

		PriorityFee: Parameter{
			ID:                 "priorityFee",
			Name:               "Priority Fee",
			Description:        "The default value for the priority fee (in gwei) for all of your transactions. This describes how much you're willing to pay *above the network's current base fee* - the higher this is, the more ETH you give to the validators for including your transaction, which generally means it will be included in a block faster (as long as your max fee is sufficiently high to cover the current network conditions).\n\nMust be larger than 0.",
			Type:               ParameterType_Float,
			Default:            map[Network]interface{}{Network_All: float64(2)},
			AffectsContainers:  []ContainerID{ContainerID_Node, ContainerID_Watchtower},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},

		AutoTxGasThreshold: Parameter{
			ID:   "minipoolStakeGasThreshold",
			Name: "Automatic TX Gas Threshold",
			Description: "Occasionally, the Smartnode will attempt to perform some automatic transactions (such as the second `stake` transaction to finish launching a minipool or the `reduce bond` transaction to convert a 16-ETH minipool to an 8-ETH one). During these, your node will use the `Rapid` suggestion from the gas estimator as its max fee.\n\nThis threshold is a limit (in gwei) you can put on that suggestion; your node will not `stake` the new minipool until the suggestion is below this limit.\n\n" +
				"A value of 0 will disable non-essential automatic transactions (such as minipool balance distribution and bond reduction), but essential transactions (such as minipool staking and solo migration promotion) will not be disabled.\n\n" +
				"NOTE: the node will ignore this limit and automatically execute transactions at whatever the suggested fee happens to be once too much time has passed since those transactions were first eligible. You may end up paying more than you wanted to if you set this too low!",
			Type:               ParameterType_Float,
			Default:            map[Network]interface{}{Network_All: float64(150)},
			AffectsContainers:  []ContainerID{ContainerID_Node},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},

		DistributeThreshold: Parameter{
			ID:                 "distributeThreshold",
			Name:               "Auto-Distribute Threshold",
			Description:        "The Smartnode will regularly check the balance of each of your minipools on the Execution Layer (**not** the Beacon Chain).\nIf any of them have a balance greater than this threshold (in ETH), the Smartnode will automatically distribute the balance. This will send your share of the balance to your withdrawal address.\n\nMust be less than 8 ETH.\n\nSet this to 0 to disable automatic distributes.\n[orange]WARNING: if you disable automatic distribution, you **must** ensure you distribute your minipool's balance before it reaches 8 ETH or you will no longer be able to distribute your rewards until you exit the minipool!",
			Type:               ParameterType_Float,
			Default:            map[Network]interface{}{Network_All: float64(1)},
			AffectsContainers:  []ContainerID{ContainerID_Node},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},

		VerifyProposals: Parameter{
			ID:                 "verifyProposals",
			Name:               "Enable PDAO Proposal Checker",
			Description:        "Check this box to opt into the responsibility for verifying Protocol DAO proposals once the Houston upgrade has been activated. Your node will regularly check for new proposals, verify their correctness, and submit challenges to any that do not match the on-chain data (e.g., if someone tampered with voting power and attempted to cheat).\n\nTo learn more about the PDAO proposal checking duty, including requirements and RPL bonding, please see the documentation at <placeholder>.",
			Type:               ParameterType_Bool,
			Default:            map[Network]interface{}{Network_All: false},
			AffectsContainers:  []ContainerID{ContainerID_Node},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},

		RewardsTreeMode: Parameter{
			ID:                 "rewardsTreeMode",
			Name:               "Rewards Tree Mode",
			Description:        "Select how you want to acquire the Merkle Tree files for each rewards interval.",
			Type:               ParameterType_Choice,
			Default:            map[Network]interface{}{Network_All: RewardsMode_Download},
			AffectsContainers:  []ContainerID{ContainerID_Node, ContainerID_Watchtower},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
			Options: []ParameterOption{{
				Name:        "Download",
				Description: "Automatically download the Merkle Tree rewards files that were published by the Oracle DAO after a rewards checkpoint.",
				Value:       RewardsMode_Download,
			}, {
				Name:        "Generate",
				Description: "Use your node to automatically generate the Merkle Tree rewards file once a checkpoint has passed. This option lets you build and verify the file that the Oracle DAO created if you prefer not to trust it and want to generate the tree yourself.\n\n[orange]WARNING: Generating the tree can take a *very long time* if many node operators are opted into the Smoothing Pool, which could impact your attestation performance!",
				Value:       RewardsMode_Generate,
			}},
		},

		RewardsTreeCustomUrl: Parameter{
			ID:                 "rewardsTreeCustomUrl",
			Name:               "Rewards Tree Custom Download URLs",
			Description:        "The Smartnode will automatically download missing rewards tree files from trusted sources like IPFS and Rocket Pool's repository on GitHub. Use this field if you would like to manually specify additional sources that host the rewards tree files, so the Smartnode can download from them as well.\nMultiple URLs can be provided using ';' as separator).\n\nUse '%s' to specify the location of the rewards file name in the URL - for example: `https://my-cool-domain.com/rewards-trees/mainnet/%s`.",
			Type:               ParameterType_String,
			Default:            map[Network]interface{}{Network_All: ""},
			AffectsContainers:  []ContainerID{ContainerID_Watchtower},
			CanBeBlank:         true,
			OverwriteOnUpgrade: false,
		},

		ArchiveECUrl: Parameter{
			ID:                 "archiveECUrl",
			Name:               "Archive-Mode EC URL",
			Description:        "[orange]**For manual Merkle rewards tree generation only.**[white]\n\nGenerating the Merkle rewards tree files for past rewards intervals typically requires an Execution client with Archive mode enabled, which is usually disabled on your primary and fallback Execution clients to save disk space.\nIf you want to generate your own rewards tree files for intervals from a long time ago, you may enter the URL of an Execution client with Archive access here.\n\nFor a free light client with Archive access, you may use https://www.alchemy.com/supernode.",
			Type:               ParameterType_String,
			Default:            map[Network]interface{}{Network_All: ""},
			AffectsContainers:  []ContainerID{ContainerID_Watchtower},
			CanBeBlank:         true,
			OverwriteOnUpgrade: false,
		},

		WatchtowerMaxFeeOverride: Parameter{
			ID:                 "watchtowerMaxFeeOverride",
			Name:               "Watchtower Max Fee Override",
			Description:        fmt.Sprintf("[orange]**For Oracle DAO members only.**\n\n[white]Use this to override the max fee (in gwei) for watchtower transactions. Note that if you set it below %d, the setting will be ignored; it can only be used to set the max fee higher than %d during times of extreme network stress.", WatchtowerMaxFeeDefault, WatchtowerMaxFeeDefault),
			Type:               ParameterType_Float,
			Default:            map[Network]interface{}{Network_All: float64(WatchtowerMaxFeeDefault)},
			AffectsContainers:  []ContainerID{ContainerID_Watchtower},
			CanBeBlank:         false,
			OverwriteOnUpgrade: true,
		},

		WatchtowerPrioFeeOverride: Parameter{
			ID:                 "watchtowerPrioFeeOverride",
			Name:               "Watchtower Priority Fee Override",
			Description:        fmt.Sprintf("[orange]**For Oracle DAO members only.**\n\n[white]Use this to override the priority fee (in gwei) for watchtower transactions. Note that if you set it below %d, the setting will be ignored; it can only be used to set the priority fee higher than %d during times of extreme network stress.", WatchtowerPrioFeeDefault, WatchtowerPrioFeeDefault),
			Type:               ParameterType_Float,
			Default:            map[Network]interface{}{Network_All: float64(WatchtowerPrioFeeDefault)},
			AffectsContainers:  []ContainerID{ContainerID_Watchtower},
			CanBeBlank:         false,
			OverwriteOnUpgrade: true,
		},

		UseRollingRecords: Parameter{
			ID:                 "useRollingRecords",
			Name:               "Use Rolling Records",
			Description:        "[orange]**WARNING: EXPERIMENTAL**\n\n[white]Enable this to use the new rolling records feature, which stores attestation records for the entire Rocket Pool network in real time instead of collecting them all after a rewards period during tree generation.\n\nOnly useful for the Oracle DAO, or if you generate your own rewards trees.",
			Type:               ParameterType_Bool,
			Default:            map[Network]interface{}{Network_All: false},
			AffectsContainers:  []ContainerID{ContainerID_Watchtower},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},

		RecordCheckpointInterval: Parameter{
			ID:                 "recordCheckpointInterval",
			Name:               "Record Checkpoint Interval",
			Description:        "The number of epochs that should pass before saving a new rolling record checkpoint. Used if Rolling Records is enabled.\n\nOnly useful for the Oracle DAO, or if you generate your own rewards trees.",
			Type:               ParameterType_Uint,
			Default:            map[Network]interface{}{Network_All: uint64(45)},
			AffectsContainers:  []ContainerID{ContainerID_Watchtower},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},

		CheckpointRetentionLimit: Parameter{
			ID:                 "checkpointRetentionLimit",
			Name:               "Checkpoint Retention Limit",
			Description:        "The number of checkpoint files to save on-disk before pruning old ones. Used if Rolling Records is enabled.\n\nOnly useful for the Oracle DAO, or if you generate your own rewards trees.",
			Type:               ParameterType_Uint,
			Default:            map[Network]interface{}{Network_All: uint64(200)},
			AffectsContainers:  []ContainerID{ContainerID_Watchtower},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},

		RecordsPath: Parameter{
			ID:                 "recordsPath",
			Name:               "Records Path",
			Description:        "The path of the folder to store rolling record checkpoints in during a rewards interval. Used if Rolling Records is enabled.\n\nOnly useful if you're an Oracle DAO member, or if you generate your own rewards trees.",
			Type:               ParameterType_String,
			Default:            map[Network]interface{}{Network_All: getDefaultRecordsDir(cfg)},
			AffectsContainers:  []ContainerID{ContainerID_Watchtower},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},

		txWatchUrl: map[Network]string{
			Network_Mainnet: "https://etherscan.io/tx",
			Network_Prater:  "https://goerli.etherscan.io/tx",
			Network_Devnet:  "https://holesky.etherscan.io/tx",
			Network_Holesky: "https://holesky.etherscan.io/tx",
		},

		stakeUrl: map[Network]string{
			Network_Mainnet: "https://stake.rocketpool.net",
			Network_Prater:  "https://goerli.rocketpool.net",
			Network_Devnet:  "TBD",
			Network_Holesky: "https://testnet.rocketpool.net",
		},

		chainID: map[Network]uint{
			Network_Mainnet: 1,     // Mainnet
			Network_Prater:  5,     // Goerli
			Network_Devnet:  17000, // Holesky
			Network_Holesky: 17000, // Holesky
		},

		storageAddress: map[Network]string{
			Network_Mainnet: "0x1d8f8f00cfa6758d7bE78336684788Fb0ee0Fa46",
			Network_Prater:  "0xd8Cd47263414aFEca62d6e2a3917d6600abDceB3",
			Network_Devnet:  "0x24503C4f4CEfCfa2d3930Ec7C8CFD44a78c7a25E",
			Network_Holesky: "0x594Fb75D3dc2DFa0150Ad03F99F97817747dd4E1",
		},

		rplTokenAddress: map[Network]string{
			Network_Mainnet: "0xD33526068D116cE69F19A9ee46F0bd304F21A51f",
			Network_Prater:  "0x5e932688e81a182e3de211db6544f98b8e4f89c7",
			Network_Devnet:  "0x60AD78075524eCadA707a2fC7E67024f7F30E3D9",
			Network_Holesky: "0x1Cc9cF5586522c6F483E84A19c3C2B0B6d027bF0",
		},

		rplFaucetAddress: map[Network]string{
			Network_Mainnet: "",
			Network_Prater:  "0x95D6b8E2106E3B30a72fC87e2B56ce15E37853F9",
			Network_Devnet:  "0x218a718A1B23B13737E2F566Dd45730E8DAD451b",
			Network_Holesky: "0xb4565BDe40Cb22282D7287A839c4ce8534674070",
		},

		rethAddress: map[Network]string{
			Network_Mainnet: "0xae78736Cd615f374D3085123A210448E74Fc6393",
			Network_Prater:  "0x178E141a0E3b34152f73Ff610437A7bf9B83267A",
			Network_Devnet:  "0xC932ca658b00DDD1EaA5EeB3BDEc7d7825567255",
			Network_Holesky: "0x7322c24752f79c05FFD1E2a6FCB97020C1C264F1",
		},

		v1_0_0_RewardsPoolAddress: map[Network]string{
			Network_Mainnet: "0xA3a18348e6E2d3897B6f2671bb8c120e36554802",
			Network_Prater:  "0xf9aE18eB0CE4930Bc3d7d1A5E33e4286d4FB0f8B",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		v1_0_0_ClaimNodeAddress: map[Network]string{
			Network_Mainnet: "0x899336A2a86053705E65dB61f52C686dcFaeF548",
			Network_Prater:  "0xc05b7A2a03A6d2736d1D0ebf4d4a0aFE2cc32cE1",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		v1_0_0_ClaimTrustedNodeAddress: map[Network]string{
			Network_Mainnet: "0x6af730deB0463b432433318dC8002C0A4e9315e8",
			Network_Prater:  "0x730982F4439E5AC30292333ff7d0C478907f2219",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		v1_0_0_MinipoolManagerAddress: map[Network]string{
			Network_Mainnet: "0x6293B8abC1F36aFB22406Be5f96D893072A8cF3a",
			Network_Prater:  "0xB815a94430f08dD2ab61143cE1D5739Ac81D3C6d",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		v1_1_0_NetworkPricesAddress: map[Network]string{
			Network_Mainnet: "0xd3f500F550F46e504A4D2153127B47e007e11166",
			Network_Prater:  "0x12f96dC173a806D18d71fAFe3C1BA2149c3E3Dc6",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		v1_1_0_NodeStakingAddress: map[Network]string{
			Network_Mainnet: "0xA73ec45Fe405B5BFCdC0bF4cbc9014Bb32a01cd2",
			Network_Prater:  "0xA73ec45Fe405B5BFCdC0bF4cbc9014Bb32a01cd2",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		v1_1_0_NodeDepositAddress: map[Network]string{
			Network_Mainnet: "0x1Cc9cF5586522c6F483E84A19c3C2B0B6d027bF0",
			Network_Prater:  "0x1Cc9cF5586522c6F483E84A19c3C2B0B6d027bF0",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		v1_1_0_MinipoolQueueAddress: map[Network]string{
			Network_Mainnet: "0x5870dA524635D1310Dc0e6F256Ce331012C9C19E",
			Network_Prater:  "0xEF5EF45bf1CC08D5694f87F8c4023f00CCCB7237",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		v1_1_0_MinipoolFactoryAddress: map[Network]string{
			Network_Mainnet: "0x54705f80D7C51Fcffd9C659ce3f3C9a7dCCf5788",
			Network_Prater:  "0x54705f80D7C51Fcffd9C659ce3f3C9a7dCCf5788",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		snapshotDelegationAddress: map[Network]string{
			Network_Mainnet: "0x469788fE6E9E9681C6ebF3bF78e7Fd26Fc015446",
			Network_Prater:  "0xD0897D68Cd66A710dDCecDe30F7557972181BEDc",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		snapshotApiDomain: map[Network]string{
			Network_Mainnet: "hub.snapshot.org",
			Network_Prater:  "testnet.snapshot.org",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		previousRewardsPoolAddresses: map[Network][]common.Address{
			Network_Mainnet: {
				common.HexToAddress("0x594Fb75D3dc2DFa0150Ad03F99F97817747dd4E1"),
			},
			Network_Prater: {
				common.HexToAddress("0x594Fb75D3dc2DFa0150Ad03F99F97817747dd4E1"),
				common.HexToAddress("0x6e91E3416acf3d015358eeAAF247a0674F6c306f"),
			},
			Network_Devnet: {
				common.HexToAddress("0xd88f4d99e52eb202801628ae37ea1a15e2904d12"),
				common.HexToAddress("0x99856894b360843EB487E252ccED703752cE7C64"),
			},
			Network_Holesky: {},
		},

		previousRocketNetworkPricesAddresses: map[Network][]common.Address{
			Network_Mainnet: {
				common.HexToAddress("0x751826b107672360b764327631cC5764515fFC37"),
			},
			Network_Prater: {},
			Network_Devnet: {
				common.HexToAddress("0xdBe07d7b0B2be5C2EaFe521255245c745422038e"),
			},
			Network_Holesky: {
				common.HexToAddress("0x029d946f28f93399a5b0d09c879fc8c94e596aeb"),
			},
		},

		previousRocketNetworkBalancesAddresses: map[Network][]common.Address{
			Network_Mainnet: {
				common.HexToAddress("0x07FCaBCbe4ff0d80c2b1eb42855C0131b6cba2F4"),
			},
			Network_Prater: {},
			Network_Devnet: {
				common.HexToAddress("0x8c7CbeaAe48642Fe5DCbFC4587276D6ac3e4162A"),
			},
			Network_Holesky: {
				common.HexToAddress("0x9294Fc6F03c64Cc217f5BE8697EA3Ed2De77e2F8"),
			},
		},

		optimismPriceMessengerAddress: map[Network]string{
			Network_Mainnet: "0xdddcf2c25d50ec22e67218e873d46938650d03a7",
			Network_Prater:  "0x87E2deCE7d0A080D579f63cbcD7e1629BEcd7E7d",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		polygonPriceMessengerAddress: map[Network]string{
			Network_Mainnet: "0xb1029Ac2Be4e08516697093e2AFeC435057f3511",
			Network_Prater:  "0x6D736da1dC2562DBeA9998385A0A27d8c2B2793e",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		arbitrumPriceMessengerAddress: map[Network]string{
			Network_Mainnet: "0x05330300f829AD3fC8f33838BC88CFC4093baD53",
			Network_Prater:  "0x2b52479F6ea009907e46fc43e91064D1b92Fdc86",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		arbitrumPriceMessengerAddressV2: map[Network]string{
			Network_Mainnet: "0x312FcFB03eC9B1Ea38CB7BFCd26ee7bC3b505aB1",
			Network_Prater:  "",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		zkSyncEraPriceMessengerAddress: map[Network]string{
			Network_Mainnet: "0x6cf6CB29754aEBf88AF12089224429bD68b0b8c8",
			Network_Prater:  "0x3Fd49431bD05875AeD449Bc8C07352942A7fBA75",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		basePriceMessengerAddress: map[Network]string{
			Network_Mainnet: "0x64A5856869C06B0188C84A5F83d712bbAc03517d",
			Network_Prater:  "",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		scrollPriceMessengerAddress: map[Network]string{
			Network_Mainnet: "0x0f22dc9b9c03757d4676539203d7549c8f22c15c",
			Network_Prater:  "",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		scrollFeeEstimatorAddress: map[Network]string{
			Network_Mainnet: "0x0d7E906BD9cAFa154b048cFa766Cc1E54E39AF9B",
			Network_Prater:  "",
			Network_Devnet:  "",
			Network_Holesky: "",
		},

		rplTwapPoolAddress: map[Network]string{
			Network_Mainnet: "0xe42318ea3b998e8355a3da364eb9d48ec725eb45",
			Network_Prater:  "0x5cE71E603B138F7e65029Cc1918C0566ed0dBD4B",
			Network_Devnet:  "0x7bb10d2a3105ed5cc150c099a06cafe43d8aa15d",
			Network_Holesky: "0x7bb10d2a3105ed5cc150c099a06cafe43d8aa15d",
		},

		multicallAddress: map[Network]string{
			Network_Mainnet: "0x5BA1e12693Dc8F9c48aAD8770482f4739bEeD696",
			Network_Prater:  "0x5BA1e12693Dc8F9c48aAD8770482f4739bEeD696",
			Network_Devnet:  "0x0540b786f03c9491f3a2ab4b0e3ae4ecd4f63ce7",
			Network_Holesky: "0x0540b786f03c9491f3a2ab4b0e3ae4ecd4f63ce7",
		},

		balancebatcherAddress: map[Network]string{
			Network_Mainnet: "0xb1f8e55c7f64d203c1400b9d8555d050f94adf39",
			Network_Prater:  "0x9788C4E93f9002a7ad8e72633b11E8d1ecd51f9b",
			Network_Devnet:  "0xfAa2e7C84eD801dd9D27Ac1ed957274530796140",
			Network_Holesky: "0xfAa2e7C84eD801dd9D27Ac1ed957274530796140",
		},

		flashbotsProtectUrl: map[Network]string{
			Network_Mainnet: "https://rpc.flashbots.net/",
			Network_Prater:  "https://rpc-goerli.flashbots.net/",
			Network_Devnet:  "",
			Network_Holesky: "",
		},
	}

}

// Get the parameters for this config
func (cfg *SmartnodeConfig) GetParameters() []*Parameter {
	return []*Parameter{
		&cfg.Network,
		&cfg.ProjectName,
		&cfg.DataPath,
		&cfg.ManualMaxFee,
		&cfg.PriorityFee,
		&cfg.AutoTxGasThreshold,
		&cfg.DistributeThreshold,
		&cfg.VerifyProposals,
		&cfg.RewardsTreeMode,
		&cfg.RewardsTreeCustomUrl,
		&cfg.ArchiveECUrl,
		&cfg.WatchtowerMaxFeeOverride,
		&cfg.WatchtowerPrioFeeOverride,
		&cfg.UseRollingRecords,
		&cfg.RecordCheckpointInterval,
		&cfg.CheckpointRetentionLimit,
		&cfg.RecordsPath,
	}
}

// Getters for the non-editable parameters

func (cfg *SmartnodeConfig) GetTxWatchUrl() string {
	return cfg.txWatchUrl[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetStakeUrl() string {
	return cfg.stakeUrl[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetChainID() uint {
	return cfg.chainID[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetNodeAddressPath() string {
	if cfg.parent.IsNativeMode {
		return filepath.Join(cfg.DataPath.Value.(string), NodeAddressFilename)
	}

	return filepath.Join(DaemonDataPath, NodeAddressFilename)
}

func (cfg *SmartnodeConfig) GetWalletPath() string {
	if cfg.parent.IsNativeMode {
		return filepath.Join(cfg.DataPath.Value.(string), WalletKeystoreFilename)
	}

	return filepath.Join(DaemonDataPath, WalletKeystoreFilename)
}

func (cfg *SmartnodeConfig) GetPasswordPath() string {
	if cfg.parent.IsNativeMode {
		return filepath.Join(cfg.DataPath.Value.(string), WalletPasswordFilename)
	}

	return filepath.Join(DaemonDataPath, WalletPasswordFilename)
}

func (cfg *SmartnodeConfig) GetValidatorKeychainPath() string {
	if cfg.parent.IsNativeMode {
		return filepath.Join(cfg.DataPath.Value.(string), "validators")
	}

	return filepath.Join(DaemonDataPath, "validators")
}

func (cfg *SmartnodeConfig) GetRecordsPath() string {
	if cfg.parent.IsNativeMode {
		return filepath.Join(cfg.DataPath.Value.(string), "records")
	}

	return filepath.Join(DaemonDataPath, "records")
}

func (cfg *SmartnodeConfig) GeNodeAddressPathInCLI() string {
	return filepath.Join(cfg.DataPath.Value.(string), NodeAddressFilename)
}

func (cfg *SmartnodeConfig) GetVotingPath() string {
	if cfg.parent.IsNativeMode {
		return filepath.Join(cfg.DataPath.Value.(string), "voting", string(cfg.Network.Value.(Network)))
	}

	return filepath.Join(DaemonDataPath, "voting", string(cfg.Network.Value.(Network)))
}

func (cfg *SmartnodeConfig) GetWalletPathInCLI() string {
	return filepath.Join(cfg.DataPath.Value.(string), WalletKeystoreFilename)
}

func (cfg *SmartnodeConfig) GetPasswordPathInCLI() string {
	return filepath.Join(cfg.DataPath.Value.(string), WalletPasswordFilename)
}

func (cfg *SmartnodeConfig) GetValidatorKeychainPathInCLI() string {
	return filepath.Join(cfg.DataPath.Value.(string), "validators")
}

func (cfg *SmartnodeConfig) GetWatchtowerStatePath() string {
	if cfg.parent.IsNativeMode {
		return filepath.Join(cfg.DataPath.Value.(string), WatchtowerFolder, "state.yml")
	}

	return filepath.Join(DaemonDataPath, WatchtowerFolder, "state.yml")
}

func (cfg *SmartnodeConfig) GetCustomKeyPath() string {
	if cfg.parent.IsNativeMode {
		return filepath.Join(cfg.DataPath.Value.(string), "custom-keys")
	}

	return filepath.Join(DaemonDataPath, "custom-keys")
}

func (cfg *SmartnodeConfig) GetCustomKeyPasswordFilePath() string {
	if cfg.parent.IsNativeMode {
		return filepath.Join(cfg.DataPath.Value.(string), "custom-key-passwords")
	}

	return filepath.Join(DaemonDataPath, "custom-key-passwords")
}

func (cfg *SmartnodeConfig) GetStorageAddress() string {
	return cfg.storageAddress[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetRplTokenAddress() string {
	return cfg.rplTokenAddress[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetRplFaucetAddress() string {
	return cfg.rplFaucetAddress[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetSnapshotDelegationAddress() string {
	return cfg.snapshotDelegationAddress[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetSmartnodeContainerTag() string {
	return smartnodeTag
}

func (config *SmartnodeConfig) GetPruneProvisionerContainerTag() string {
	return pruneProvisionerTag
}

func (cfg *SmartnodeConfig) GetEcMigratorContainerTag() string {
	return ecMigratorTag
}

func (cfg *SmartnodeConfig) GetSnapshotApiDomain() string {
	return cfg.snapshotApiDomain[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetVotingSnapshotID() [32]byte {
	// So the contract wants a Keccak'd hash of the voting ID, but Snapshot's service wants ASCII so it can display the ID in plain text; we have to do this to make it play nicely with Snapshot
	buffer := [32]byte{}
	idBytes := []byte(SnapshotID)
	copy(buffer[0:], idBytes)
	return buffer
}

func (config *SmartnodeConfig) GetSnapshotID() string {
	return SnapshotID
}

// The the title for the config
func (cfg *SmartnodeConfig) GetConfigTitle() string {
	return cfg.Title
}

func (cfg *SmartnodeConfig) GetRethAddress() common.Address {
	return common.HexToAddress(cfg.rethAddress[cfg.Network.Value.(Network)])
}

func getDefaultDataDir(cfg *RocketPoolConfig) string {
	return filepath.Join(cfg.RocketPoolDirectory, "data")
}

func getDefaultRecordsDir(cfg *RocketPoolConfig) string {
	return filepath.Join(getDefaultDataDir(cfg), "records")
}

func (cfg *SmartnodeConfig) GetRewardsTreePath(interval uint64, daemon bool) string {
	if daemon && !cfg.parent.IsNativeMode {
		return filepath.Join(DaemonDataPath, RewardsTreesFolder, fmt.Sprintf(RewardsTreeFilenameFormat, string(cfg.Network.Value.(Network)), interval))
	}

	return filepath.Join(cfg.DataPath.Value.(string), RewardsTreesFolder, fmt.Sprintf(RewardsTreeFilenameFormat, string(cfg.Network.Value.(Network)), interval))
}

func (cfg *SmartnodeConfig) GetMinipoolPerformancePath(interval uint64, daemon bool) string {
	if daemon && !cfg.parent.IsNativeMode {
		return filepath.Join(DaemonDataPath, RewardsTreesFolder, fmt.Sprintf(MinipoolPerformanceFilenameFormat, string(cfg.Network.Value.(Network)), interval))
	}

	return filepath.Join(cfg.DataPath.Value.(string), RewardsTreesFolder, fmt.Sprintf(MinipoolPerformanceFilenameFormat, string(cfg.Network.Value.(Network)), interval))
}

func (cfg *SmartnodeConfig) GetRegenerateRewardsTreeRequestPath(interval uint64, daemon bool) string {
	if daemon && !cfg.parent.IsNativeMode {
		return filepath.Join(DaemonDataPath, WatchtowerFolder, fmt.Sprintf(RegenerateRewardsTreeRequestFormat, interval))
	}

	return filepath.Join(cfg.DataPath.Value.(string), WatchtowerFolder, fmt.Sprintf(RegenerateRewardsTreeRequestFormat, interval))
}

func (cfg *SmartnodeConfig) GetWatchtowerFolder(daemon bool) string {
	if daemon && !cfg.parent.IsNativeMode {
		return filepath.Join(DaemonDataPath, WatchtowerFolder)
	}

	return filepath.Join(cfg.DataPath.Value.(string), WatchtowerFolder)
}

func (cfg *SmartnodeConfig) GetFeeRecipientFilePath() string {
	if !cfg.parent.IsNativeMode {
		return filepath.Join(DaemonDataPath, "validators", FeeRecipientFilename)
	}

	return filepath.Join(cfg.DataPath.Value.(string), "validators", NativeFeeRecipientFilename)
}

func (cfg *SmartnodeConfig) GetSocketPath() string {
	if !cfg.parent.IsNativeMode {
		return filepath.Join(DaemonDataPath, SocketFilename)
	}

	return filepath.Join(cfg.DataPath.Value.(string), SocketFilename)
}

func (cfg *SmartnodeConfig) GetSocketPathInCLI() string {
	return filepath.Join(cfg.DataPath.Value.(string), SocketFilename)
}

func (cfg *SmartnodeConfig) GetV100RewardsPoolAddress() common.Address {
	return common.HexToAddress(cfg.v1_0_0_RewardsPoolAddress[cfg.Network.Value.(Network)])
}

func (cfg *SmartnodeConfig) GetV100ClaimNodeAddress() common.Address {
	return common.HexToAddress(cfg.v1_0_0_ClaimNodeAddress[cfg.Network.Value.(Network)])
}

func (cfg *SmartnodeConfig) GetV100ClaimTrustedNodeAddress() common.Address {
	return common.HexToAddress(cfg.v1_0_0_ClaimTrustedNodeAddress[cfg.Network.Value.(Network)])
}

func (cfg *SmartnodeConfig) GetV100MinipoolManagerAddress() common.Address {
	return common.HexToAddress(cfg.v1_0_0_MinipoolManagerAddress[cfg.Network.Value.(Network)])
}

func (cfg *SmartnodeConfig) GetV110NetworkPricesAddress() common.Address {
	return common.HexToAddress(cfg.v1_1_0_NetworkPricesAddress[cfg.Network.Value.(Network)])
}

func (cfg *SmartnodeConfig) GetV120NetworkPricesAddress() common.Address {
	return common.HexToAddress(cfg.v1_2_0_NetworkPricesAddress[cfg.Network.Value.(Network)])
}

func (cfg *SmartnodeConfig) GetV120NetworkBalancesAddress() common.Address {
	return common.HexToAddress(cfg.v1_2_0_NetworkBalancesAddress[cfg.Network.Value.(Network)])
}

func (cfg *SmartnodeConfig) GetV110NodeStakingAddress() common.Address {
	return common.HexToAddress(cfg.v1_1_0_NodeStakingAddress[cfg.Network.Value.(Network)])
}

func (cfg *SmartnodeConfig) GetV110NodeDepositAddress() common.Address {
	return common.HexToAddress(cfg.v1_1_0_NodeDepositAddress[cfg.Network.Value.(Network)])
}

func (cfg *SmartnodeConfig) GetV110MinipoolQueueAddress() common.Address {
	return common.HexToAddress(cfg.v1_1_0_MinipoolQueueAddress[cfg.Network.Value.(Network)])
}

func (cfg *SmartnodeConfig) GetV110MinipoolFactoryAddress() common.Address {
	return common.HexToAddress(cfg.v1_1_0_MinipoolFactoryAddress[cfg.Network.Value.(Network)])
}

func (cfg *SmartnodeConfig) GetPreviousRewardsPoolAddresses() []common.Address {
	return cfg.previousRewardsPoolAddresses[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetPreviousRocketNetworkPricesAddresses() []common.Address {
	return cfg.previousRocketNetworkPricesAddresses[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetPreviousRocketNetworkBalancesAddresses() []common.Address {
	return cfg.previousRocketNetworkBalancesAddresses[cfg.Network.Value.(Network)]
}
func (cfg *SmartnodeConfig) GetOptimismMessengerAddress() string {
	return cfg.optimismPriceMessengerAddress[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetPolygonMessengerAddress() string {
	return cfg.polygonPriceMessengerAddress[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetArbitrumMessengerAddress() string {
	return cfg.arbitrumPriceMessengerAddress[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetArbitrumMessengerAddressV2() string {
	return cfg.arbitrumPriceMessengerAddressV2[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetZkSyncEraMessengerAddress() string {
	return cfg.zkSyncEraPriceMessengerAddress[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetBaseMessengerAddress() string {
	return cfg.basePriceMessengerAddress[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetScrollMessengerAddress() string {
	return cfg.scrollPriceMessengerAddress[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetScrollFeeEstimatorAddress() string {
	return cfg.scrollFeeEstimatorAddress[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetRplTwapPoolAddress() string {
	return cfg.rplTwapPoolAddress[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetMulticallAddress() string {
	return cfg.multicallAddress[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetBalanceBatcherAddress() string {
	return cfg.balancebatcherAddress[cfg.Network.Value.(Network)]
}

func (cfg *SmartnodeConfig) GetFlashbotsProtectUrl() string {
	return cfg.flashbotsProtectUrl[cfg.Network.Value.(Network)]
}

func getNetworkOptions() []ParameterOption {
	options := []ParameterOption{
		{
			Name:        "Ethereum Mainnet",
			Description: "This is the real Ethereum main network, using real ETH and real RPL to make real validators.",
			Value:       Network_Mainnet,
		},
		/*
			{
				Name:        "Prater Testnet",
				Description: "This is the Prater test network, using free fake ETH and free fake RPL to make fake validators.\nUse this if you want to practice running the Smartnode in a free, safe environment before moving to Mainnet.",
				Value:       Network_Prater,
			},
		*/
		{
			Name:        "Holesky Testnet",
			Description: "This is the Holešky (Holešovice) test network, which is the next generation of long-lived testnets for Ethereum. It uses free fake ETH and free fake RPL to make fake validators.\nUse this if you want to practice running the Smartnode in a free, safe environment before moving to Mainnet.",
			Value:       Network_Holesky,
		},
	}

	if strings.HasSuffix(shared.RocketPoolVersion, "-dev") {
		options = append(options, ParameterOption{
			Name:        "Devnet",
			Description: "This is a development network used by Rocket Pool engineers to test new features and contract upgrades before they are promoted to Holesky for staging. You should not use this network unless invited to do so by the developers.",
			Value:       Network_Devnet,
		})
	}

	return options
}