Clover
---
[![GoDoc](https://godoc.org/plimble/clover?status.svg)](https://godoc.org/github.com/plimble/clover)
[![Build Status](https://travis-ci.org/plimble/clover.svg?branch=master)](https://travis-ci.org/plimble/clover?branch=master)
[![Coverage Status](https://coveralls.io/repos/plimble/clover/badge.svg?branch=master&service=github&foo)](https://coveralls.io/github/plimble/clover?branch=master)
[![Go Report Card](https://goreportcard.com/badge/plimble/clover)](https://goreportcard.com/report/plimble/clover)

Oauth2 server v2.0.0

WIP Reimplement oauth2 server

### Install
```
go get gopkg.in/plimble/clover.v2
```

### Roadmap
- [ ] OAuth2 [RFC6749](https://tools.ietf.org/html/rfc6749)
  - [ ] Authorize Request
    - [ ] Code Response Type
    - [ ] Token Response Type
  - [x] Access Token
    - [x] Authorization Code Grant
    - [x] Resource Owner Password Credentials Grant
    - [x] Client Credentials Grant
    - [x] Extension Grants
  - [x] Verification Access Token Request
- [ ] JSON Web Token (JWT) [RFC7523](https://tools.ietf.org/html/rfc7523)
- [x] Token Revocation [RFC7009](https://tools.ietf.org/html/rfc7009)
- [ ] Security Considerations [RFC6819](https://tools.ietf.org/html/rfc6819)
