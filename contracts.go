package etherscan

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

// ContractsClient is the client for contracts actions.
type ContractsClient struct {
	api *apiClient
}

const contractsModule = "contract"

// GetContractABI returns the contract ABI as a JSON string.
func (c ContractsClient) GetContractABI(
	ctx context.Context, address common.Address,
) (result string, err error) {
	req := struct{ Address common.Address }{address}
	err = c.api.call(ctx, &callParams{
		module:  contractsModule,
		action:  "getabi",
		request: req,
		result:  &result,
	})

	return result, err
}

// ContractInfo contains information on a contract's source code.
type ContractInfo struct {
	SourceCode           string `etherscan:"SourceCode"`
	ABI                  string `etherscan:"ABI"`
	ContractName         string `etherscan:"ContractName"`
	CompilerVersion      string `etherscan:"CompilerVersion"`
	OptimizationUsed     string `etherscan:"OptimizationUsed"`
	Runs                 uint32 `etherscan:"Runs"`
	ConstructorArguments string `etherscan:"ConstructorArguments"`
	EVMVersion           string `etherscan:"EVMVersion"`
	Library              string `etherscan:"Library"`
	LicenseType          string `etherscan:"LicenseType"`
	Proxy                bool   `etherscan:"Proxy,num"`
	Implementation       string `etherscan:"Implementation"`
	SwarmSource          string `etherscan:"SwarmSource"`
}

// GetContractSourceCode returns the Solidity source code of a verified smart contract.
func (c ContractsClient) GetContractSourceCode(
	ctx context.Context, address common.Address,
) (result []ContractInfo, err error) {
	req := struct{ Address common.Address }{address}
	err = c.api.call(ctx, &callParams{
		module:  contractsModule,
		action:  "getsourcecode",
		request: req,
		result:  &result,
	})

	return result, err
}
