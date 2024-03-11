package config

import (
	"path/filepath"

	"github.com/rocket-pool/node-manager-core/config"
	"github.com/rocket-pool/smartnode/shared/config/ids"
)

// Configuration for the Validator Client
type ValidatorClientConfig struct {
	// The command for restarting the validator container in native mode
	NativeValidatorRestartCommand config.Parameter[string]

	// The command for stopping the validator container in native mode
	NativeValidatorStopCommand config.Parameter[string]

	// Subconfigs
	VcCommon   *config.ValidatorClientCommonConfig
	Lighthouse *config.LighthouseVcConfig
	Lodestar   *config.LodestarVcConfig
	Nimbus     *config.NimbusVcConfig
	Prysm      *config.PrysmVcConfig
	Teku       *config.TekuVcConfig
}

// Generates a new Validator Client config
func NewValidatorClientConfig(rocketPoolDirectory string) *ValidatorClientConfig {
	cfg := &ValidatorClientConfig{
		NativeValidatorRestartCommand: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.NativeValidatorRestartCommandID,
				Name:               "VC Restart Script",
				Description:        "The absolute path to a custom script that will be invoked when the Smart Node needs to restart your validator client to load the new key after a minipool is staked.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_Daemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]string{
				config.Network_All: getDefaultValidatorRestartCommand(rocketPoolDirectory),
			},
		},

		NativeValidatorStopCommand: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.NativeValidatorStopCommandID,
				Name:               "Validator Stop Command",
				Description:        "The absolute path to a custom script that will be invoked when the Smart Node needs to stop your validator client in case of emergency.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_Daemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]string{
				config.Network_All: getDefaultValidatorStopCommand(rocketPoolDirectory),
			},
		},

		VcCommon:   config.NewValidatorClientCommonConfig(),
		Lighthouse: config.NewLighthouseVcConfig(),
		Lodestar:   config.NewLodestarVcConfig(),
		Nimbus:     config.NewNimbusVcConfig(),
		Prysm:      config.NewPrysmVcConfig(),
		Teku:       config.NewTekuVcConfig(),
	}

	cfg.Lighthouse.ContainerTag.Default[Network_Devnet] = cfg.Lighthouse.ContainerTag.Default[config.Network_Holesky]
	cfg.Lodestar.ContainerTag.Default[Network_Devnet] = cfg.Lodestar.ContainerTag.Default[config.Network_Holesky]
	cfg.Nimbus.ContainerTag.Default[Network_Devnet] = cfg.Nimbus.ContainerTag.Default[config.Network_Holesky]
	cfg.Prysm.ContainerTag.Default[Network_Devnet] = cfg.Prysm.ContainerTag.Default[config.Network_Holesky]
	cfg.Teku.ContainerTag.Default[Network_Devnet] = cfg.Teku.ContainerTag.Default[config.Network_Holesky]

	return cfg
}

// The title for the config
func (cfg *ValidatorClientConfig) GetTitle() string {
	return "Validator Client"
}

// Get the parameters for this config
func (cfg *ValidatorClientConfig) GetParameters() []config.IParameter {
	return []config.IParameter{
		&cfg.NativeValidatorRestartCommand,
		&cfg.NativeValidatorStopCommand,
	}
}

// Get the sections underneath this one
func (cfg *ValidatorClientConfig) GetSubconfigs() map[string]config.IConfigSection {
	return map[string]config.IConfigSection{
		ids.VcCommonID:   cfg.VcCommon,
		ids.LighthouseID: cfg.Lighthouse,
		ids.LodestarID:   cfg.Lodestar,
		ids.NimbusID:     cfg.Nimbus,
		ids.PrysmID:      cfg.Prysm,
		ids.TekuID:       cfg.Teku,
	}
}

func getDefaultValidatorRestartCommand(rocketPoolDirectory string) string {
	return filepath.Join(rocketPoolDirectory, "restart-vc.sh")
}

func getDefaultValidatorStopCommand(rocketPoolDirectory string) string {
	return filepath.Join(rocketPoolDirectory, "stop-validator.sh")
}
