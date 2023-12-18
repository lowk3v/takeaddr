package utils

import "strings"

type ChainMetadata struct {
	ChainId   int64
	ChainName string
	Explorer  string
}

var ChainMetadataList = []ChainMetadata{
	{
		ChainId:   1,
		ChainName: "ethereum",
		Explorer:  "https://etherscan.io",
	},
	{
		ChainId:   56,
		ChainName: "bsc",
		Explorer:  "https://bscscan.com",
	},
	{
		ChainId:   97,
		ChainName: "testnet bsc",
		Explorer:  "https://testnet.bscscan.com",
	},
	{
		ChainId:   137,
		ChainName: "polygon",
		Explorer:  "https://polygonscan.com",
	},
	{
		ChainId:   80001,
		ChainName: "mumbai",
		Explorer:  "https://mumbai.polygonscan.com",
	},
	{
		ChainId:   250,
		ChainName: "fantom",
		Explorer:  "https://ftmscan.com",
	},
	{
		ChainId:   43114,
		ChainName: "avalanche",
		Explorer:  "https://snowtrace.io/",
	},
	{
		ChainId:   42161,
		ChainName: "arbitrum",
		Explorer:  "https://arbiscan.io",
	},
	{
		ChainId:   1337,
		ChainName: "local",
		Explorer:  "http://localhost:1337",
	},
}

func ChainIdToChainName(chainId int64, isUpper bool) string {
	for _, chainMetadata := range ChainMetadataList {
		if chainMetadata.ChainId == chainId {
			// upper the first letter
			if isUpper {
				return strings.ToUpper(chainMetadata.ChainName)
			}
			return chainMetadata.ChainName
		}
	}
	return ""
}

func ChainIdToExplorer(chainId int64) string {
	for _, chainMetadata := range ChainMetadataList {
		if chainMetadata.ChainId == chainId {
			return chainMetadata.Explorer
		}
	}
	return ""
}
