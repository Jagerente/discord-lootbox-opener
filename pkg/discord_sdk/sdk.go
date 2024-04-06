package discord_sdk

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

const (
	baseURL           = "https://discord.com"
	defRef            = "https://discord.com/channels/@me"
	defXSProps        = "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiQ2hyb21lIiwiZGV2aWNlIjoiIiwic3lzdGVtX2xvY2FsZSI6ImVuLVVTIiwiYnJvd3Nlcl91c2VyX2FnZW50IjoiTW96aWxsYS81LjAgKFdpbmRvd3MgTlQgMTAuMDsgV2luNjQ7IHg2NCkgQXBwbGVXZWJLaXQvNTM3LjM2IChLSFRNTCwgbGlrZSBHZWNrbykgQ2hyb21lLzEyMy4wLjAuMCBTYWZhcmkvNTM3LjM2IiwiYnJvd3Nlcl92ZXJzaW9uIjoiMTIzLjAuMC4wIiwib3NfdmVyc2lvbiI6IjEwIiwicmVmZXJyZXIiOiIiLCJyZWZlcnJpbmdfZG9tYWluIjoiIiwicmVmZXJyZXJfY3VycmVudCI6IiIsInJlZmVycmluZ19kb21haW5fY3VycmVudCI6IiIsInJlbGVhc2VfY2hhbm5lbCI6InN0YWJsZSIsImNsaWVudF9idWlsZF9udW1iZXIiOjI4MjA2OCwiY2xpZW50X2V2ZW50X3NvdXJjZSI6bnVsbH0="
	openLootboxMethod = `/api/v9/users/@me/lootboxes/open`
)

type Config struct {
	UserAgent string
	Token     string
}

type SDK struct {
	client           *resty.Client
	xSuperProperties string
}

func New(cfg *Config) *SDK {
	sdk := new(SDK)

	sdk.xSuperProperties = defXSProps

	client := resty.New().
		SetBaseURL(baseURL).
		SetHeader("User-Agent", cfg.UserAgent).
		SetHeader("Authorization", cfg.Token).
		SetHeader("Origin", baseURL).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusTooManyRequests
			},
		)

	sdk.client = client

	return sdk
}

func (s *SDK) UseCustomXSuperProperties(xSuperProperties string) *SDK {
	s.xSuperProperties = xSuperProperties
	return s
}

func (s *SDK) OpenLootbox() (*OpenLootboxResponse, error) {
	res := new(OpenLootboxResponse)

	resp, err := s.client.R().
		SetHeader("X-Super-Properties", defXSProps).
		SetHeader("Accept", "*/*").
		SetHeader("Referer", defRef).
		SetResult(res).
		Post(openLootboxMethod)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to open lootbox: [%v] %s", resp.Status(), string(resp.Body()))
	}

	return res, nil
}
