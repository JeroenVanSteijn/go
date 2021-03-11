package createregulatedassetoffer

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/log"
	"github.com/stellar/go/txnbuild"
)

type Options struct {
	HorizonURL          string
	NetworkPassphrase   string
	AccountIssuerSecret string
	AssetCode           string
}

func (opts Options) horizonClient() horizonclient.ClientInterface {
	var client *horizonclient.Client
	if opts.NetworkPassphrase == network.PublicNetworkPassphrase {
		client = horizonclient.DefaultPublicNetClient
	} else {
		client = horizonclient.DefaultTestNetClient
	}

	client.HorizonURL = opts.HorizonURL

	return client
}

// Create is used to create a sell offer of the regulated asset to create a
// market for this new asset.
func Create(opts Options) {
	kp, err := keypair.ParseFull(opts.AccountIssuerSecret)
	if err != nil {
		log.Fatal(errors.Wrap(err, "parsing secret"))
	}

	horizonClient := opts.horizonClient()

	account, err := horizonClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: kp.Address(),
	})
	if err != nil {
		log.Fatal(errors.Wrap(err, "getting account detail"))
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &account,
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.ManageSellOffer{
					Selling: txnbuild.CreditAsset{
						Code:   opts.AssetCode,
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
		log.Fatal(errors.Wrap(err, "creating transaction"))
	}

	tx, err = tx.Sign(opts.NetworkPassphrase, kp)
	if err != nil {
		log.Fatal(errors.Wrap(err, "signing transaction"))
	}

	log.Info("Will create offer")
	_, err = horizonClient.SubmitTransaction(tx)
	if err != nil {
		log.Fatal(parseHorizonError(err))
	}
	log.Info("Did create offer")

	return
}

func parseHorizonError(err error) error {
	if err == nil {
		return nil
	}

	rootErr := errors.Cause(err)
	if hError := horizonclient.GetError(rootErr); hError != nil {
		resultCode, _ := hError.ResultCodes()
		err = errors.Wrapf(err, "error submitting transaction: %+v, %+v\n", hError.Problem, resultCode)
	} else {
		err = errors.Wrap(err, "error submitting transaction")
	}
	return err
}
