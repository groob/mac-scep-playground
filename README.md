# Overview

Testing SCEP and related profile payloads is hard. You need a SCEP Server. You need a profile with the right keys. You need a CA cert. Your profile likely needs an Identiy Preference. Once you've installed the profile you probably want to test it against a running server. So you gotta generate some server certs, set up nginx and configure it to require TLS Auth. By the time you're done you forgot what you even wanted to test. This repo is for you(and me in a few months when I look into this functionality again).

This repo comes with the following resources:

- A custom Go server with TLS config that requires a cert signed by the CA in `ca/depot/ca.pem`.
- A SCEP profile with the SCEP CA above embedded and a custom Identiy preference for `*.corp.acme.co`.

# Requirements & Usage

- Install [Go](https://golang.org/dl/)
- Install [MicroMDM/SCEP](https://github.com/micromdm/scep), both `scepserver` and `scepclient` though you'll likely only need the server. We're testing Profiles.

- Run the SCEP Server.
    ```
    # Note: The CA pass is the password of the CA private key in ca/depot/ca.pem.
    # The -challenge is the SCEP challeng you'll be prompted for when installing a profile.
    # Keep the allowrenew at 0 otherwise you wont be able to renew the SCEP cert for two weeks.

    scepserver -port 9001 -challenge=secret -allowrenew=0 -capass=secret
    ```

- Run the server. 
    `go run server.go`

- Install the profile in `client/profile.mobileconfig`.

- Edit `/etc/hosts` to point `foo.corp.acme.co` to `127.0.0.1`

- Visit Safari at `https://foo.corp.acme.co:9000`. Did you get prompted for the cert from the profile?

# Extra

## Creating random self signed certs (like the server cert)

The Go stdlib comes with a handy utility to generate self signed certs you can use for testing. You can also use `openssl`.
Anything in the `--host=` flag is what the SAN of the cert will be.

```
go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host=localhost:9000,127.0.0.1:9000 --ecdsa-curve=P256 --ca=true

# read the cert info
openssl x509 -in cert.pem -text
```

## Initialize a new CA with SCEP

```
scepserver ca -init -organization groob-io -key-password=secret
```

## Use `scepclient` to get a cert

Run this instead of the profile if you need to get a client cert.

```
scepclient -server-url=http://localhost:9001/scep -private-key=./client/key.pem -challenge=secret
```

## Profile documentation

- https://mosen.github.io/profiledocs/
- [Apple Docs](https://developer.apple.com/library/content/featuredarticles/iPhoneConfigurationProfileRef/Introduction/Introduction.html#//apple_ref/doc/uid/TP40010206-CH1-SW18)
- You can't use the client cert in your own app because Keychain imposes a strict ACL on a SCEP installed cert:
[via frogor](https://github.com/munki/munki/issues/662#issuecomment-250538851)
