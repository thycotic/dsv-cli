package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"thy/constants"
	cst "thy/constants"
	"thy/errors"
	"thy/format"
	"thy/paths"
	preds "thy/predictors"
	"thy/requests"
	"thy/utils"

	"github.com/posener/complete"
	"github.com/spf13/viper"
	"github.com/thycotic-rd/cli"
)

type client struct {
	request   requests.Client
	outClient format.OutClient
}

func GetNoDataOpClientWrappers(targetEntity string) cli.PredictorWrappers {
	return cli.PredictorWrappers{
		preds.LongFlag(cst.ClientID): cli.PredictorWrapper{complete.PredictAnything, preds.NewFlagValue(preds.Params{Name: cst.ClientID, Usage: fmt.Sprintf("ID of the %s ", targetEntity)}), false},
	}
}

func GetClientCmd() (cli.Command, error) {
	return NewCommand(CommandArgs{
		Path: []string{cst.NounClient},
		RunFunc: func(args []string) int {
			name := viper.GetString(cst.DataName)
			if name == "" && len(args) > 0 {
				name = args[0]
			}
			if name == "" {
				return cli.RunResultHelp
			}
			return client{requests.NewHttpClient(), nil}.handleClientReadCmd(args)
		},
		SynopsisText: "client (<client-id> | --client-id)",
		HelpText: fmt.Sprintf(`Execute an action on a %[1]s in %[2]s

Usage:
   • %[1]s %[3]s
   • %[1]s --client-id %[3]s
		`, cst.NounClient, cst.ProductName, cst.ExampleClientID),
		FlagsPredictor: GetNoDataOpClientWrappers(cst.NounClient),
		MinNumberArgs:  1,
	})
}

func GetClientReadCmd() (cli.Command, error) {
	return NewCommand(CommandArgs{
		Path:         []string{cst.NounClient, cst.Read},
		RunFunc:      client{requests.NewHttpClient(), nil}.handleClientReadCmd,
		SynopsisText: fmt.Sprintf("%s %s (<client-id> | --client-id)", cst.NounClient, cst.Read),
		HelpText: fmt.Sprintf(`Execute an action on a %[1]s in %[2]s

Usage:
   • %[1]s %[4]s %[3]s
   • %[1]s %[4]s --client-id %[3]s
		`, cst.NounClient, cst.ProductName, cst.ExampleClientID, cst.Read),
		FlagsPredictor: GetNoDataOpClientWrappers(cst.NounClient),
		MinNumberArgs:  1,
	})
}

func GetClientDeleteCmd() (cli.Command, error) {
	return NewCommand(CommandArgs{
		Path:         []string{cst.NounClient, cst.Delete},
		RunFunc:      client{requests.NewHttpClient(), nil}.handleClientDeleteCmd,
		SynopsisText: fmt.Sprintf("%s %s (<client-id> | --client-id)", cst.NounClient, cst.Delete),
		HelpText: fmt.Sprintf(`Delete a %[1]s from %[2]s

Usage:
   • %[1]s %[4]s %[3]s
   • %[1]s %[4]s --client-id %[3]s --force
		`, cst.NounClient, cst.ProductName, cst.ExampleClientID, cst.Delete),
		FlagsPredictor: cli.PredictorWrappers{
			preds.LongFlag(cst.ClientID): cli.PredictorWrapper{complete.PredictAnything, preds.NewFlagValue(preds.Params{Name: cst.ClientID, Usage: fmt.Sprintf("ID of the %s ", cst.NounClient)}), false},
			preds.LongFlag(cst.Force):    cli.PredictorWrapper{complete.PredictAnything, preds.NewFlagValue(preds.Params{Name: cst.Force, Shorthand: "", Usage: fmt.Sprintf("Immediately delete %s", cst.NounClient), Global: false, ValueType: "bool"}), false},
		},
		MinNumberArgs: 1,
	})
}

func GetClientRestoreCmd() (cli.Command, error) {
	return NewCommand(CommandArgs{
		Path:         []string{cst.NounClient, cst.Restore},
		RunFunc:      client{requests.NewHttpClient(), nil}.handleClientRestoreCmd,
		SynopsisText: fmt.Sprintf("%s %s (<client-id> | --client-id)", cst.NounClient, cst.Restore),
		HelpText: fmt.Sprintf(`Restore a deleted %[1]s in %[2]s

Usage:
   • %[1]s %[4]s %[3]s
		`, cst.NounClient, cst.ProductName, cst.ExampleClientID, cst.Restore),
		FlagsPredictor: GetNoDataOpClientWrappers(cst.NounClient),
		MinNumberArgs:  1,
	})
}

func GetClientCreateCmd() (cli.Command, error) {
	return NewCommand(CommandArgs{
		Path:         []string{cst.NounClient, cst.Create},
		RunFunc:      client{requests.NewHttpClient(), nil}.handleClientUpsertCmd,
		SynopsisText: fmt.Sprintf("%s %s (<role> | --role) |(<uses> | --uses)|(<desc> | --desc)|(<ttl> | --ttl)| (<url> | --url) | ( <urlTTL> | --urlTTL)", cst.NounClient, cst.Create),
		HelpText: fmt.Sprintf(`%[4]s a %[1]s in %[2]s

Usage:
   • %[1]s %[4]s %[3]s
   • %[1]s %[4]s --role %[3]s
		`, cst.NounClient, cst.ProductName, cst.ExampleRoleName, cst.Create),
		FlagsPredictor: cli.PredictorWrappers{
			preds.LongFlag(cst.NounRole):            cli.PredictorWrapper{complete.PredictAnything, preds.NewFlagValue(preds.Params{Name: cst.NounRole, Usage: fmt.Sprintf("Name of the %s ", cst.NounRole)}), false},
			preds.LongFlag(cst.NounBootstrapUrl):    cli.PredictorWrapper{complete.PredictAnything, preds.NewFlagValue(preds.Params{Name: cst.NounBootstrapUrl, Usage: fmt.Sprint("Generate one time use url instead of secret (optional)"), ValueType: "bool"}), false},
			preds.LongFlag(cst.NounBootstrapurlTTL): cli.PredictorWrapper{complete.PredictNothing, preds.NewFlagValue(preds.Params{Name: cst.NounBootstrapurlTTL, Usage: fmt.Sprint("UrlTTL for generated url (optional)")}), false},
			preds.LongFlag(cst.NounClientUses):      cli.PredictorWrapper{complete.PredictNothing, preds.NewFlagValue(preds.Params{Name: cst.NounClientUses, Usage: fmt.Sprint("The number of times the client credential can be read.  if set to 0, it can be used infinitely.  default is 0 (optional)")}), false},
			preds.LongFlag(cst.NounClientDesc):      cli.PredictorWrapper{complete.PredictNothing, preds.NewFlagValue(preds.Params{Name: cst.NounClientDesc, Usage: fmt.Sprint("Client description (optional)")}), false},
			preds.LongFlag(cst.NounClientTTL):       cli.PredictorWrapper{complete.PredictNothing, preds.NewFlagValue(preds.Params{Name: cst.NounClientTTL, Usage: fmt.Sprint("How long until the client credential expires. if set to 0, it can be used indefinitely.  default is 0.\n (optional)")}), false},
		},
		MinNumberArgs: 1,
	})
}

func GetClientSearchCmd() (cli.Command, error) {
	return NewCommand(CommandArgs{
		Path:         []string{cst.NounClient, cst.Search},
		RunFunc:      client{requests.NewHttpClient(), nil}.handleClientSearchCmd,
		SynopsisText: fmt.Sprintf("%s (<role> | --role)", cst.Search),
		HelpText: fmt.Sprintf(`Search for %[1]ss attached to a given %[5]s in %[2]s

Usage:
   • %[1]s %[4]s %[3]s
   • %[1]s %[4]s --role %[3]s
		`, cst.NounClient, cst.ProductName, cst.ExampleRoleName, cst.Search, cst.NounRole),
		FlagsPredictor: cli.PredictorWrappers{
			preds.LongFlag(cst.NounRole): cli.PredictorWrapper{complete.PredictAnything, preds.NewFlagValue(preds.Params{Name: cst.NounRole, Usage: "Role that has attached clients (required)"}), false},
			preds.LongFlag(cst.Limit):    cli.PredictorWrapper{complete.PredictAnything, preds.NewFlagValue(preds.Params{Name: cst.Limit, Shorthand: "l", Usage: fmt.Sprint("Maximum number of results per cursor (optional)")}), false},
			preds.LongFlag(cst.Cursor):   cli.PredictorWrapper{complete.PredictAnything, preds.NewFlagValue(preds.Params{Name: cst.Cursor, Usage: constants.CursorHelpMessage}), false},
		},
		MinNumberArgs: 1,
	})
}

func (c client) handleClientReadCmd(args []string) int {
	var err *errors.ApiError
	var data []byte
	clientID := viper.GetString(cst.ClientID)
	if clientID == "" && len(args) > 0 {
		clientID = args[0]
	}
	if clientID == "" {
		err = errors.NewS("error: must specify " + cst.ClientID)
	} else {
		uri := paths.CreateResourceURI(cst.NounClient, clientID, "", true, nil, true)
		data, err = c.request.DoRequest("GET", uri, nil)
	}

	outClient := c.outClient
	if outClient == nil {
		outClient = format.NewDefaultOutClient()
	}
	outClient.WriteResponse(data, err)
	return utils.GetExecStatus(err)
}

func (c client) handleClientDeleteCmd(args []string) int {
	var err *errors.ApiError
	var data []byte
	force := viper.GetBool(cst.Force)
	clientID := viper.GetString(cst.ClientID)
	if clientID == "" && len(args) > 0 {
		clientID = args[0]
	}
	if clientID == "" {
		err = errors.NewS("error: must specify " + cst.ClientID)
	} else {
		query := map[string]string{"force": strconv.FormatBool(force)}
		uri := paths.CreateResourceURI(cst.NounClient, clientID, "", true, query, true)
		data, err = c.request.DoRequest("DELETE", uri, nil)
	}
	if c.outClient == nil {
		c.outClient = format.NewDefaultOutClient()
	}

	c.outClient.WriteResponse(data, err)
	return utils.GetExecStatus(err)
}

func (c client) handleClientRestoreCmd(args []string) int {
	var err *errors.ApiError
	var data []byte
	if c.outClient == nil {
		c.outClient = format.NewDefaultOutClient()
	}
	clientID := viper.GetString(cst.ClientID)
	if clientID == "" && len(args) > 0 {
		clientID = args[0]
	}
	if clientID == "" {
		err = errors.NewS("error: must specify " + cst.ClientID)
	} else {
		uri := paths.CreateResourceURI(cst.NounClient, clientID, "/restore", true, nil, true)
		data, err = c.request.DoRequest("PUT", uri, nil)
	}
	c.outClient.WriteResponse(data, err)
	return utils.GetExecStatus(err)
}

func (c client) handleClientUpsertCmd(args []string) int {
	var err *errors.ApiError
	var data []byte
	roleName := viper.GetString(cst.NounRole)
	url := viper.GetBool(cst.NounBootstrapUrl)
	urlTTL := viper.GetInt64(cst.NounBootstrapurlTTL)
	ttl := viper.GetInt64(cst.NounClientTTL)
	desc := viper.GetString(cst.NounClientDesc)
	uses := viper.GetInt(cst.NounClientUses)
	if roleName == "" && len(args) > 0 {
		roleName = args[0]
	}
	if roleName == "" {
		return cli.RunResultHelp
	}
	client := Client{
		Role:        roleName,
		Url:         url,
		UrlTTL:      urlTTL,
		Uses:        uses,
		Description: desc,
		TTL:         ttl,
	}
	reqMethod := strings.ToLower(viper.GetString(cst.LastCommandKey))
	var uri string
	if reqMethod == cst.Create {
		reqMethod = "POST"
		uri = paths.CreateResourceURI(cst.NounClient, "", "", true, nil, true)
	} else {
		reqMethod = "PUT"
		uri = paths.CreateResourceURI(cst.NounClient, viper.GetString(cst.ClientID), "", true, nil, true)
	}
	data, err = c.request.DoRequest(reqMethod, uri, &client)
	outClient := c.outClient
	if outClient == nil {
		outClient = format.NewDefaultOutClient()
	}
	outClient.WriteResponse(data, err)
	return utils.GetExecStatus(err)
}

func (c client) handleClientSearchCmd(args []string) int {
	var err *errors.ApiError
	var data []byte
	role := viper.GetString(cst.NounRole)
	limit := viper.GetString(cst.Limit)
	cursor := viper.GetString(cst.Cursor)
	if role == "" && len(args) > 0 {
		role = args[0]
	}
	if role == "" {
		err = errors.NewS("error: must specify " + cst.NounRole)
	} else {
		queryParams := map[string]string{
			cst.NounRole: role,
			cst.Limit:    limit,
			cst.Cursor:   cursor,
		}
		uri := paths.CreateResourceURI(cst.NounClient, "", "", false, queryParams, true)
		data, err = c.request.DoRequest("GET", uri, nil)
	}

	outClient := c.outClient
	if outClient == nil {
		outClient = format.NewDefaultOutClient()
	}
	outClient.WriteResponse(data, err)
	return utils.GetExecStatus(err)
}

type Client struct {
	Role        string `json:"role"`
	Url         bool   `json:"url,omitempty"`
	UrlTTL      int64  `json:"urlTTL,omitempty"`
	TTL         int64  `json:"ttl,omitempty"`
	Uses        int    `json:"usesLimit,omitempty"`
	Description string `json:"description,omitempty"`
}
