# route53-registerer
Simple static binary to update a route53 dns record

I'm 90% sure I already wrote basically this code, so if you see another copy I
wrote lying around somewhere let me know... but in the meantime it was easier
to just write it again.

This is meant to be best used as a simple static binary for cases where that's convenient to have.


## Usage

It takes the usual means of configuring AWS stuff (environment variables preferred though) in addition to the `REGISTER_IP=$ip`, `REGISTER_DOMAIN=foo.my.domain`, and `DNS_TYPE=<A|SRV|AAAA> (default A)` environment variables.

# License

MIT
