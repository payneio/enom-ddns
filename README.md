# enom-ddns

Sets the host record at Enom to whatever IP this program resolves to
external machines using Enom's DDNS service. Uses AWS for IP resolution.

## Usage

```bash
DDNS_DOMAIN=<domain> ENOM_UN=<enom un> ENOM_PW=<enom domain pw> enom-ddns
```

This is the domain password, not the Enom user login password. The domain
password needs to be set at Enom.

## Running from Docker

```bash
docker run -it -e DDNS_DOMAIN=<domain> -e ENOM_UN=<enom un> -e ENOM_PW=<enom domain pw> quay.io/payneio/enomddns:1
```

# Build

```bash
make build/enom-ddns
```
