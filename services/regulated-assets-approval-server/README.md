# regulated-assets-approval-server

```
Varsion: v0.0.1
Status: unreleased
```

This is a [SEP-8] Approval Server reference implementation based on SEP-8 v1.6.1
intended for **testing only**. It is being concieved to:

1. Be used as an example of how regulated assets transactions can be validated
   and revised by an anchor.
2. Serve as a demo server where wallets can test and validate their SEP-8
   implementation.

## Usage

```sh
$ go install
$ regulated-assets-approval-server --help
SEP-8 Approval Server

Usage:
  regulated-assets-approval-server [command] [flags]
  regulated-assets-approval-server [command]

Available Commands:
  configure-asset-issuer       Configure asset issuer to use SEP-8 regulated assets.
  create-regulated-asset-offer Create a sell offer from the issuing account selling ASSET_CODE for XLM at 1:1.
  serve                        Serve the SEP-8 Approval Server

Use "regulated-assets-approval-server [command] --help" for more information about a command.
```

### Usage: `configure-asset-issuer`

```sh
$ go install
$ regulated-assets-approval-server serve --help
Configure asset issuer to use SEP-8 regulated assets.

Usage:
  regulated-assets-approval-server configure-asset-issuer [flags]

Flags:
      --account-issuer-secret string   Secret key of the asset issuer's stellar account. (ACCOUNT_ISSUER_SECRET)
      --horizon-url string             Horizon URL used for looking up account details (HORIZON_URL) (default "https://horizon-testnet.stellar.org/")
      --network-passphrase string      Network passphrase of the Stellar network transactions should be signed for (NETWORK_PASSPHRASE) (default "Test SDF Network ; September 2015")
```

### Usage: `create-regulated-asset-offer`

```sh
$ go install
$ regulated-assets-approval-server serve --help
Create a sell offer from the issuing account selling ASSET_CODE for XLM at 1:1.

Usage:
  regulated-assets-approval-server create-regulated-asset-offer [flags]

Flags:
      --account-issuer-secret string   Secret key of the asset issuer's stellar account. (ACCOUNT_ISSUER_SECRET)
      --asset-code string              The code of the reguated asset (ASSET_CODE)
      --horizon-url string             Horizon URL used for looking up account details (HORIZON_URL) (default "https://horizon-testnet.stellar.org/")
      --network-passphrase string      Network passphrase of the Stellar network transactions should be signed for (NETWORK_PASSPHRASE) (default "Test SDF Network ; September 2015")
```

### Usage: `serve`

```sh
$ go install
$ regulated-assets-approval-server serve --help
Serve the SEP-8 Approval Server

Usage:
  regulated-assets-approval-server serve [flags]

Flags:
      --horizon-url string          Horizon URL used for looking up account details (HORIZON_URL) (default "https://horizon-testnet.stellar.org/")
      --network-passphrase string   Network passphrase of the Stellar network transactions should be signed for (NETWORK_PASSPHRASE) (default "Test SDF Network ; September 2015")
      --port int                    Port to listen and serve on (PORT) (default 8000)
```

[SEP-8]: https://github.com/stellar/stellar-protocol/blob/7c795bb9abc606cd1e34764c4ba07900d58fe26e/ecosystem/sep-0008.md
