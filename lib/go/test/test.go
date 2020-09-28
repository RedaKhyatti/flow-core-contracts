package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/cadence"
	sdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"

	emulator "github.com/dapperlabs/flow-emulator"
	"github.com/dapperlabs/flow-emulator/server/backend"
)

var (
	fungibleTokenAddress = "notset"
	flowTokenAddress     = "notset"
	feesAddress          = "notset"
)

func newBlockchain() backend.Emulator {
	networkType := os.Getenv("FLOW_NETWORK_TYPE")
	if networkType == "testnet" {
		return newNetwork()
	}
	fmt.Println("Using Emulator")
	return newEmulator()
}

// newEmulator returns a emulator object for testing.
func newEmulator() *emulator.Blockchain {
	b, err := emulator.NewBlockchain()
	if err != nil {
		panic(err)
	}

	generator := sdk.NewAddressGenerator(sdk.Emulator)
	// Skip Service account addr
	generator.Next()
	fungibleTokenAddress = generator.NextAddress().Hex()
	flowTokenAddress = generator.NextAddress().Hex()
	feesAddress = generator.NextAddress().Hex()

	return b
}

func newNetwork() backend.Emulator {
	b, err := emulator.NewNetwork(os.Getenv("FLOW_ADDRESS"), os.Getenv("FLOW_SERVICE_ACCOUNT_PRIVATE_KEY"))

	if err != nil {
		panic(err)
	}

	generator := sdk.NewAddressGenerator(sdk.Testnet)
	// Skip Service account addr
	generator.Next()
	fungibleTokenAddress = generator.NextAddress().Hex()
	flowTokenAddress = generator.NextAddress().Hex()
	feesAddress = generator.NextAddress().Hex()

	return b
}

// signAndSubmit signs a transaction with an array of signers and adds their signatures to the transaction
// before submitting it to the emulator.
//
// If the private keys do not match up with the addresses, the transaction will not succeed.
//
// The shouldRevert parameter indicates whether the transaction should fail or not.
//
// This function asserts the correct result and commits the block if it passed.
func signAndSubmit(
	t *testing.T,
	b backend.Emulator,
	tx *sdk.Transaction,
	signerAddresses []sdk.Address,
	signers []crypto.Signer,
	shouldRevert bool,
) {
	latestBlockID, err := b.GetLatestBlockID()
	assert.NoError(t, err)

	tx = tx.SetReferenceBlockID(latestBlockID)
	// sign transaction with each signer
	for i := len(signerAddresses) - 1; i >= 0; i-- {
		signerAddress := signerAddresses[i]
		signer := signers[i]

		if i == 0 {
			err := tx.SignEnvelope(signerAddress, 0, signer)
			assert.NoError(t, err)
		} else {
			err := tx.SignPayload(signerAddress, 0, signer)
			assert.NoError(t, err)
		}
	}

	Submit(t, b, tx, shouldRevert)
}

// Submit submits a transaction and checks if it fails or not.
func Submit(
	t *testing.T,
	b backend.Emulator,
	tx *sdk.Transaction,
	shouldRevert bool,
) {
	// submit the signed transaction
	err := b.AddTransaction(*tx)
	require.NoError(t, err)

	result, err := b.ExecuteNextTransaction()
	require.NoError(t, err)

	if shouldRevert {
		assert.True(t, result.Reverted())
	} else {
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
			//cmd.PrettyPrintError(result.Error, "", map[string]string{"": ""})
		}
	}

	_, err = b.CommitBlock()
	assert.NoError(t, err)
}

// executeScriptAndCheck executes a script and checks to make sure that it succeeded.
func executeScriptAndCheck(t *testing.T, b backend.Emulator, script []byte) cadence.Value {
	result, err := b.ExecuteScript(script, nil)
	require.NoError(t, err)
	if !assert.True(t, result.Succeeded()) {
		t.Log(result.Error.Error())
	}

	return result.Value
}

func readFile(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return contents
}

// CadenceUFix64 returns a UFix64 value
func CadenceUFix64(value string) cadence.Value {
	newValue, err := cadence.NewUFix64(value)

	if err != nil {
		panic(err)
	}

	return newValue
}

func bytesToCadenceArray(b []byte) cadence.Array {
	values := make([]cadence.Value, len(b))

	for i, v := range b {
		values[i] = cadence.NewUInt8(v)
	}

	return cadence.NewArray(values)
}
