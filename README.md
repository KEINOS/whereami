# WhereAmI

"`whereami`" is a command line utility that displays the current global/public IP address.

Useful for finding out the ephemeral (current external) IP address; works on macOS, Linux and Windows.

```shellsession
$ whereami
123.234.123.124
```

```shellsession
$ whereami -help
Usage of whereami:
  -verbose
        prints detailed information if any
```

## Install

- Via [Homebrew](https://brew.sh/) for macOS, Linux and Windows WSL2. (Intel, AMD64, ARM64, M1)

    ```bash
    brew install KEINOS/apps/whereami
    ```

- For manual download or other architectures see [releases page](https://github.com/KEINOS/whereami/releases/latest).
- Note:
    - To avoid a large number of API requests to the service providers, **this application sleeps for one second** after obtaining a global/public IP address.

## Statuses

[![go1.14+](https://github.com/KEINOS/whereami/actions/workflows/go-versions.yml/badge.svg)](https://github.com/KEINOS/whereami/actions/workflows/go-versions.yml)
[![golangci-lint](https://github.com/KEINOS/whereami/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/KEINOS/whereami/actions/workflows/golangci-lint.yml)
[![codecov](https://codecov.io/gh/KEINOS/whereami/branch/main/graph/badge.svg?token=wwZpJLfm0l)](https://codecov.io/gh/KEINOS/whereami)
[![Go Report Card](https://goreportcard.com/badge/github.com/KEINOS/dev-go)](https://goreportcard.com/report/github.com/KEINOS/dev-go)
[![CodeQL](https://github.com/KEINOS/whereami/actions/workflows/codeQL-analysis.yml/badge.svg)](https://github.com/KEINOS/whereami/actions/workflows/codeQL-analysis.yml)

## Contribute

[![go1.16+](https://img.shields.io/badge/Go-1.16+-blue?logo=go)](https://github.com/KEINOS/whereami/actions/workflows/go-versions.yml "Supported versions")
[![Go Reference](https://pkg.go.dev/badge/github.com/KEINOS/whereami.svg)](https://pkg.go.dev/github.com/KEINOS/whereami/)

## License

- [MIT](https://github.com/KEINOS/whereami/blob/main/LICENSE)
- Copyright: (c) 2021 [KEINOS and the WhoAmI contributors](https://github.com/KEINOS/whereami/graphs/contributors).

## Acknowledgment

### Service Providers Used

We would like to thank the following service providers that this command uses to discover global/public IP addresses.
**This command uses the first IP address with the same response from these providers in random order**. (Max match: 3)

- [https://ipinfo.io/](https://ipinfo.io/)
- [https://inet-ip.info/](https://inet-ip.info/)
- [http://inetclue.com/](http://inetclue.com/)
