package client

import (
	"github.com/guowenshuai/ieth/modules/repo"
	"github.com/urfave/cli/v2"

)

func NewCli(c *cli.Context) (*APIClient, error) {
	r, err := repo.NewRepo(c.String("repo"))
	if err != nil {
		return nil, err
	}
	apibyte, err := r.GetApi()
	if err != nil {
		return nil, err
	}
	cmd := NewAPIClient(&Configuration{
		BasePath:      string(apibyte),
		Host:          "",
		Scheme:        "",
		DefaultHeader: nil,
		UserAgent:     "Swagger-Codegen/1.0.0/go",
		HTTPClient:    nil,
	})
	return cmd, nil
}

