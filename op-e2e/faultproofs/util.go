package faultproofs

import (
	"testing"

	batcherFlags "github.com/ethereum-optimism/optimism/op-batcher/flags"
	op_e2e "github.com/ethereum-optimism/optimism/op-e2e"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

type faultDisputeConfigOpts func(cfg *op_e2e.SystemConfig)

func withBatcherStopped() faultDisputeConfigOpts {
	return func(cfg *op_e2e.SystemConfig) {
		cfg.DisableBatcher = true
	}
}

func withBlobBatches() faultDisputeConfigOpts {
	return func(cfg *op_e2e.SystemConfig) {
		cfg.DataAvailabilityType = batcherFlags.BlobsType

		genesisActivation := hexutil.Uint64(0)
		cfg.DeployConfig.L1CancunTimeOffset = &genesisActivation
		cfg.DeployConfig.L2GenesisDeltaTimeOffset = &genesisActivation
		cfg.DeployConfig.L2GenesisEcotoneTimeOffset = &genesisActivation
	}
}

func startFaultDisputeSystem(t *testing.T, opts ...faultDisputeConfigOpts) (*op_e2e.System, *ethclient.Client) {
	cfg := op_e2e.DefaultSystemConfig(t)
	delete(cfg.Nodes, "verifier")
	for _, opt := range opts {
		opt(&cfg)
	}
	cfg.DeployConfig.SequencerWindowSize = 4
	cfg.DeployConfig.FinalizationPeriodSeconds = 2
	cfg.SupportL1TimeTravel = true
	cfg.DeployConfig.L2OutputOracleSubmissionInterval = 1
	cfg.NonFinalizedProposals = true // Submit output proposals asap
	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")
	return sys, sys.Clients["l1"]
}
