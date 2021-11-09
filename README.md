Etherscan API Client Library - Go
---------------------------------

Client library for the Etherscan library, written in Go. Provides methods
corresponding to each of the Etherscan API endpoints.

Features
========

- Namespaced by API module
- Full context support for cancellation/deadline control
- Allows full configuration of http client object
- Uses standard library and go-ethereum types.

Install
=======

```$ go get github.com/ryanc414/etherscan-api-go```


Basic Usage
===========

```
	client := etherscan.New(&etherscan.Params{
		APIKey: os.Getenv("ETHERSCAN_API_KEY"),
	})

	gas, err := client.Gas.GetGasOracle(ctx)
	if err != nil {
		return errors.Wrap(err, "GetGasOracle")
	}

	log.Print("fast gas price: ", gas.FastGasPrice)
```

See https://pkg.go.dev/github.com/ryanc414/etherscan-api-go for full API
documentation!
