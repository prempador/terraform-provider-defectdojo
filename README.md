# Terraform Provider for Defectdojo

Terraform provider for managing Defectdojo configuration and resources.

This is a community provider and not supported by Hashicorp or Defectdojo.

## Developing the Provider

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20
- [Docker](https://docs.docker.com/engine/install)

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`. This requires Defectdojo to run on your local mashine or you having a Defectdojo instance running somewhere. You will need to set `DEFECTDOJO_HOST` and minimally `DEFECTDOJO_TOKEN` in order for the provider to successfully connect to your Defectdojo instance and to run the Acceptance tests.

```shell
make testacc
```
