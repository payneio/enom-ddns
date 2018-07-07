enom-ddns
=========

Sets the host record at Enom to whatever IP this program resolves to external machines.
Uses AWS for IP resolution.

Usage:

```
DDNS_DOMAIN=<domain> DDNS_UN=<enom un> DDNS_PW=<enom pw> enom-ddns-client
```

Build:

- Requires Go
- `go install https://github.com/payneio/enom-ddns`
- Creates a self-contained binary

