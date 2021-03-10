package configureassetissuer

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/log"
	"github.com/stellar/go/txnbuild"
)

type ConfigureAssetIssuerOptions struct {
	HorizonURL          string
	NetworkPassphrase   string
	AccountIssuerSecret string
	AssetCode           string
}

func (opts ConfigureAssetIssuerOptions) horizonClient() horizonclient.ClientInterface {
	var client *horizonclient.Client
	if opts.NetworkPassphrase == network.PublicNetworkPassphrase {
		client = horizonclient.DefaultPublicNetClient
	} else {
		client = horizonclient.DefaultTestNetClient
	}

	client.HorizonURL = opts.HorizonURL

	return client
}

func Configure(opts ConfigureAssetIssuerOptions) {
	err := ConfigureAccountFlags(opts)
	if err != nil {
		log.DefaultLogger.Fatal(errors.Wrap(err, "configuring account flags"))
	}
	err = IssueAssetOffer(opts)
	if err != nil {
		log.DefaultLogger.Fatal(errors.Wrap(err, "issuing asset offer"))
	}
}

// ConfigureAccountFlags will set the account flags needed for regulated assets
// to "auth_required: true" and "auth_revocable: true".
// ref1:https://developers.stellar.org/docs/issuing-assets/control-asset-access/
// ref2:https://github.com/stellar/stellar-protocol/blob/d49e04af8e047474f2c506d9d11bb63b6ad55d2c/ecosystem/sep-0008.md#authorization-flags
func ConfigureAccountFlags(opts ConfigureAssetIssuerOptions) error {
	kp, err := keypair.ParseFull(opts.AccountIssuerSecret)
	if err != nil {
		return errors.Wrap(err, "parsing secret")
	}

	horizonClient := opts.horizonClient()

	log.DefaultLogger.Infof("Account address: %s\n", kp.Address())
	account, err := horizonClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: kp.Address(),
	})
	if err != nil {
		return errors.Wrap(err, "getting account detail")
	}

	if account.Flags.AuthRevocable && account.Flags.AuthRequired {
		log.Info("Account is already auth_revocable:true and auth_required:true.")
		return nil
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &account,
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.SetOptions{
					SetFlags: []txnbuild.AccountFlag{
						txnbuild.AuthRequired,
						txnbuild.AuthRevocable,
					},
				},
			},
			BaseFee:    300,
			Timebounds: txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		return errors.Wrap(err, "creating transaction")
	}

	tx, err = tx.Sign(opts.NetworkPassphrase, kp)
	if err != nil {
		return errors.Wrap(err, "signing transaction")
	}

	_, err = horizonClient.SubmitTransaction(tx)
	if err != nil {
		return errors.Wrap(err, "submitting transaction")
	}

	return nil
}

func IssueAssetOffer(opts ConfigureAssetIssuerOptions) error {
	kp, err := keypair.ParseFull(opts.AccountIssuerSecret)
	if err != nil {
		return errors.Wrap(err, "parsing secret")
	}

	horizonClient := opts.horizonClient()

	account, err := horizonClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: kp.Address(),
	})
	if err != nil {
		return errors.Wrap(err, "getting account detail")
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &account,
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.ManageSellOffer{
					Selling: txnbuild.CreditAsset{
						Code:   "HUE",
						Issuer: kp.Address(),
					},
					Buying: txnbuild.NativeAsset{},
					Amount: "100000",
					Price:  "1",
				},
			},
			BaseFee:    300,
			Timebounds: txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		return errors.Wrap(err, "creating transaction")
	}

	tx, err = tx.Sign(opts.NetworkPassphrase, kp)
	if err != nil {
		return errors.Wrap(err, "signing transaction")
	}

	_, err = horizonClient.SubmitTransaction(tx)
	if err != nil {
		return errors.Wrap(err, "submitting transaction")
	}

	return nil
}
