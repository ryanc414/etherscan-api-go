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
	SourceCode           string
	ABI                  string
	ContractName         string
	CompilerVersion      string
	OptimizationUsed     string
	Runs                 uint32
	ConstructorArguments string
	EVMVersion           string
	Library              string
	LicenseType          string
	Proxy                bool
	Implementation       string
	SwarmSource          string
}

type contractResult struct {
	SourceCode       string  `json:"SourceCode"`
	ABI              string  `json:"ABI"`
	CompilerVersion  string  `json:"CompilerVersion"`
	OptimizationUsed string  `json:"OptimizaitonUsed"`
	Runs             uintStr `json:"Runs"`
	ConstructorArgs  string  `json:"ConstructorArguments"`
	EVMVersion       string  `json:"EVMVersion"`
	Library          string  `json:"Library"`
	LicenseType      string  `json:"LicenseType"`
	Proxy            string  `json:"Proxy"`
	Implementation   string  `json:"Implementation"`
	SwarmSource      string  `json:"SwarmSource"`
}

func (c *contractResult) toInfo() *ContractInfo {
	return &ContractInfo{
		SourceCode:           c.SourceCode,
		ABI:                  c.ABI,
		CompilerVersion:      c.CompilerVersion,
		OptimizationUsed:     c.OptimizationUsed,
		Runs:                 uint32(c.Runs.unwrap()),
		ConstructorArguments: c.ConstructorArgs,
		EVMVersion:           c.EVMVersion,
		Library:              c.Library,
		LicenseType:          c.LicenseType,
		Proxy:                c.Proxy != "0",
		Implementation:       c.Implementation,
		SwarmSource:          c.SwarmSource,
	}
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

	var result []contractResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	infos := make([]ContractInfo, len(result))
	for i := range result {
		infos[i] = *result[i].toInfo()
	}

	return infos, nil
}
