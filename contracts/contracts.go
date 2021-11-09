package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	ecommon "github.com/ryanc414/etherscan-api-go/common"
	"github.com/ryanc414/etherscan-api-go/httpapi"
)

// ContractsClient is the client for contracts actions.
type ContractsClient struct {
	API *httpapi.APIClient
}

// GetContractABI returns the contract ABI as a JSON string.
func (c ContractsClient) GetContractABI(
	ctx context.Context, address common.Address,
) (result string, err error) {
	req := struct{ Address common.Address }{address}
	err = c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.ContractsModule,
		Action:  "getabi",
		Request: req,
		Result:  &result,
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
	err = c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.ContractsModule,
		Action:  "getsourcecode",
		Request: req,
		Result:  &result,
	})

	return result, err
}
