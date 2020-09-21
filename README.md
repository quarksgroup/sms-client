![Test](https://github.com/quarksgroup/sms-client/workflows/Test/badge.svg)
![Lint](https://github.com/quarksgroup/sms-client/workflows/Lint/badge.svg)
[![codecov](https://codecov.io/gh/quarksgroup/sms-client/branch/master/graph/badge.svg)](https://codecov.io/gh/quarksgroup/sms-client)
# sms-client
The sms client used internally at quarksgroup.


## Goal
Implement several sms API that are available in Rwanda

## Currently supported
* fdi documented at https://fdisms.docs.apiary.io
* more to come...(contributions welcome)

## Example
You can find a full example app using this library at https://github.com/rugwirobaker/helmes

## Contributing
Right now the API is relatively stable and henceforward we would like to add new sms API vendors.
All new vendors or drivers as they are called here must be first declared in the sms/driver.go file. 

```
// Driver identifies source code management driver.
type Driver int

// Drivers(sms gateways) we support
const (
	DriverUnknown Driver = iota
	DriverFdi
    // New driver here 
)

func (d Driver) String() (s string) {
	switch d {
	case DriverFdi:
		return "fdi"
	default:
		return "unknown"
	}
}
```

Then implement the following interfaces:
 ```
 // sms/balance.go

 ...
 // BalanceService handles your sms credit balance.
type BalanceService interface {
	// Current returns the current balance (as of now)
	Current(context.Context) (*Balance, *Response, error)

	// At returns the balance at a given date (eg;2020-09-02).
	At(context.Context, string) (*Balance, *Response, error)
}
 ``` 
 
 
 ```
 // sms/message.go

 ...
// it sends a message returns a delivery report an serializable response
// and if there is a problem returns an error
type SendService interface {
	Send(context.Context, Message) (*Report, *Response, error)
}
``` 
 
 
```
sms/stats.go

...
// StatsService handles information(stats) about your sms sending habits.
StatsService interface {
	// At gets credit balance at any given date (eg;2020-09-02).
	At(ctx context.Context, at string) (*Stats, *Response, error)

	// Current gets the current credit Balance.
	Current(ctx context.Context) (*Stats, *Response, error)
}
``` 

```
sms/token.go

...
// AuthService handles authentication to the underlying API
type AuthService interface {
	// Login the underlying API and get an JWT token
	Login(ctx context.Context, id, secret string) (*Token, *Response, error)

	// Refresh the oauth2 token
	Refresh(ctx context.Context, token *Token, force bool) (*Token, *Response, error)
}

```

For an implementation use the fdi example [sms/driver/fdi](sms/driver/fdi)