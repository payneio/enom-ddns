enom-ddns
=========

Sets the host record at Enom to whatever IP this program resolves to
external machines using Enom's DDNS service. Uses AWS for IP resolution.

Usage:

```
DDNS_DOMAIN=<domain> DDNS_UN=<enom un> DDNS_PW=<enom domain pw> enom-ddns-client
```

This is the domain password, not the Enom user login password. The domain
password needs to be set at Enom.

Build:

- Requires Go
- `go install https://github.com/payneio/enom-ddns`
- Creates a self-contained binary
