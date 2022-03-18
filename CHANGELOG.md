# 18 Mar 2022
- Updated to 8.1.0.
- Initial setup become easier a bit and now can be fully automated.
- Now CA using PEM cert instead of PKCS#12 to get rid of useless certificate doubles.
  Config files still are available.
- Now Elasticsearch using only one PKCS#12 certificate for HTTP and TCP.
- Now passwords are stored in .env file.
- Now service configs are stored as environment variables.
- Fixed healthcheck args.
- Some small fixes.

## 24 Dec 2021
- Init.
