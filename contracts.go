package etherscan

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

type ContractsClient struct {
	api *apiClient
}

const contractsModule = "contract"

func (c ContractsClient) GetContractABI(
	ctx context.Context, address common.Address,
) (string, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: contractsModule,
		action: "getabi",
		other:  map[string]string{"address": address.String()},
	})
	if err != nil {
		return "", err
	}

	var abi string
	if err := json.Unmarshal(rspData, &abi); err != nil {
		return "", err
	}

	return abi, nil
}

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

func (c ContractsClient) GetContractSourceCode(
	ctx context.Context, address common.Address,
) ([]ContractInfo, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: contractsModule,
		action: "getsourcecode",
		other:  map[string]string{"address": address.String()},
	})
	if err != nil {
		return nil, err
	}

	var result []ContractInfo
	if err := unmarshalResponse(rspData, &result); err != nil {
		return nil, err
	}

	return result, nil
}
