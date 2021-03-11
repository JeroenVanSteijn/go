package cmd

import (
	"go/types"

	"github.com/spf13/cobra"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
	"github.com/stellar/go/services/regulated-assets-approval-server/internal/createregulatedassetoffer"
	"github.com/stellar/go/support/config"
)

type CreateRegulatedAssetOffer struct{}

func (c *CreateRegulatedAssetOffer) Command() *cobra.Command {
	opts := createregulatedassetoffer.Options{}
	configOpts := config.ConfigOptions{
		{
			Name:      "asset-code",
			Usage:     "The code of the reguated asset",
			OptType:   types.String,
			ConfigKey: &opts.AssetCode,
			Required:  true,
		},
		{
			Name:      "account-issuer-secret",
			Usage:     "Secret key of the asset issuer's stellar account.",
			OptType:   types.String,
			ConfigKey: &opts.AccountIssuerSecret,
			Required:  true,
		},
		{
			Name:        "horizon-url",
			Usage:       "Horizon URL used for looking up account details",
			OptType:     types.String,
			ConfigKey:   &opts.HorizonURL,
			FlagDefault: horizonclient.DefaultTestNetClient.HorizonURL,
			Required:    true,
		},
		{
			Name:        "network-passphrase",
			Usage:       "Network passphrase of the Stellar network transactions should be signed for",
			OptType:     types.String,
			ConfigKey:   &opts.NetworkPassphrase,
			FlagDefault: network.TestNetworkPassphrase,
			Required:    true,
		},
	}
	cmd := &cobra.Command{
		Use:   "create-regulated-asset-offer",
		Short: "Create a sell offer from the issuing account selling ASSET_CODE for XLM at 1:1.",
		Run: func(_ *cobra.Command, _ []string) {
			configOpts.Require()
			configOpts.SetValues()
			c.Run(opts)
		},
	}
	configOpts.Init(cmd)
	return cmd
}

func (c *CreateRegulatedAssetOffer) Run(opts createregulatedassetoffer.Options) {
	createregulatedassetoffer.Create(opts)
}
