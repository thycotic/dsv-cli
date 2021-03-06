package requests

import (
	"context"
	cst "thy/constants"
	"thy/errors"
	"thy/format"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	"github.com/shurcooL/graphql"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . GraphClient
type GraphClient interface {
	DoRequest(uri string, query interface{}, variables map[string]interface{}) ([]byte, *errors.ApiError)
}

type graphClient struct {
}

func NewGraphClient() GraphClient {
	return &graphClient{}
}

func (c *graphClient) DoRequest(uri string, query interface{}, variables map[string]interface{}) ([]byte, *errors.ApiError) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: viper.GetString(cst.NounToken),
		},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := graphql.NewClient(uri, httpClient)
	if err := client.Query(context.Background(), query, variables); err != nil {
		return nil, errors.NewS(err.Error())
	}
	resp, err := format.JsonMarshal(&query)
	if err != nil {
		return nil, MarshalingError
	}

	return resp, nil
}
