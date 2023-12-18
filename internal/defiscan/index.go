package defiscan

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/lowk3v/takeaddr/utils"
	"strings"
)

type DefiYield struct {
	http *resty.Client
}

type DefiYieldScan struct {
	Data struct {
		Scan struct {
			Status  string              `json:"status"`
			Results []DefiYieldScanData `json:"results"`
		} `json:"scan"`
	} `json:"data"`
}

type DefiYieldScanData struct {
	Address   string `json:"address"`
	NetworkId int64  `json:"networkId"`
}

const DefiYieldScanApi = "https://api-scanner.defiyield.app/"

func New() *DefiYield {
	return &DefiYield{
		http: resty.New(),
	}
}

func (d *DefiYield) Scan(address string) ([]DefiYieldScanData, string, error) {
	var resp DefiYieldScan
	post, err := d.http.R().
		SetResult(&resp).
		SetBody(`{"query":"mutation { scan(address: \"`+address+`\") { status results { address networkId } } }"}`).
		SetHeader("Content-Type", "application/json").
		Post(DefiYieldScanApi)
	if err != nil {
		return nil, "", err
	}
	if post.StatusCode() != 200 {
		return nil, "", err
	}
	if resp.Data.Scan.Results == nil || len(resp.Data.Scan.Results) == 0 {
		return nil, "", errors.New(resp.Data.Scan.Status)
	}
	if resp.Data.Scan.Status != "OK" {
		return nil, "", err
	}

	return resp.Data.Scan.Results, beauty(resp.Data.Scan.Results), nil
}

func beauty(data []DefiYieldScanData) string {
	var networkLinkMsg []string
	for _, info := range data {
		chainName := utils.ChainIdToChainName(info.NetworkId, false)
		if len(chainName) > 0 {
			chainName = chainName[:3]
		}
		if len(chainName) == 0 {
			continue
		}
		explorer := utils.ChainIdToExplorer(info.NetworkId)
		networkLinkMsg = append(networkLinkMsg, fmt.Sprintf("\n\t\t[%s](%s/address/%s)", chainName, explorer, info.Address))
	}
	return strings.Join(networkLinkMsg, ", ")
}
