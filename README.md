# WhereAmI

`whereami` is a simple **command line utility that displays your current global/public IP address**; works on macOS, Linux and Windows.

Useful for finding out the ephemeral (current external) IPv4 address.

```shellsession
$ whereami
123.234.123.124
```

```shellsession
$ whereami -help
Usage of whereami:
  -verbose
        prints detailed information if any. such as IPv6 and etc.
```

- Note:
  - This command only displays IPv4 addresses. However, **some service providers will return IPv6 addresses and more detailed information**. In these cases, the `--verbose` option can be used to view the details of the provider's response.
  - To avoid a large number of API requests to the service providers, **this application sleeps for one second** after printing the obtained global/public IP address.

## Install

- Manual download and install:
  - [Latest Releases Page](https://github.com/KEINOS/whereami/releases/latest)
    - **macOS** (x86_64/M1), **Windows** (x86_64/ARM64), **Linux** (x86_64/ARM64/ARM v5, 6, 7)
    - Download the archive of your OS and architecture then extract it. Place the extracted binary in your PATH with executable permission.
    - Public Key of the signature: [https://github.com/KEINOS.gpg](https://github.com/KEINOS.gpg)

- Install via [Homebrew](https://brew.sh/):
  - macOS, Linux and Windows WSL2. (x86_64/ARM64, M1)

    ```bash
    brew install KEINOS/apps/whereami
    ```

- Install via `go install`:
  - Go v1.16 or above.

    ```bash
    go install github.com/KEINOS/whereami/cmd/whereami@latest
    ```

- Run via Docker:
  - Multiarch build for x86_64 (Intel/AMD) and ARM64/M1 architectures.

    ```bash
    # The image is around 5.5MB in size
    docker pull keinos/whereami:latest
    docker run --rm keinos/whereami
    ```

## Statuses

[![Unit Test (Versions)](https://github.com/KEINOS/whereami/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/KEINOS/whereami/actions/workflows/unit-tests.yml)
[![Unit Tests (Platform)](https://github.com/KEINOS/whereami/actions/workflows/platform-test.yml/badge.svg)](https://github.com/KEINOS/whereami/actions/workflows/platform-test.yml)
[![golangci-lint](https://github.com/KEINOS/whereami/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/KEINOS/whereami/actions/workflows/golangci-lint.yml)
[![CodeQL](https://github.com/KEINOS/whereami/actions/workflows/codeQL-analysis.yml/badge.svg)](https://github.com/KEINOS/whereami/actions/workflows/codeQL-analysis.yml)

[![codecov](https://codecov.io/gh/KEINOS/whereami/branch/main/graph/badge.svg?token=wwZpJLfm0l)](https://codecov.io/gh/KEINOS/whereami)
[![Go Report Card](https://goreportcard.com/badge/github.com/KEINOS/dev-go)](https://goreportcard.com/report/github.com/KEINOS/dev-go)

## Contribute

[![go1.16+](https://img.shields.io/badge/Go-1.16+-blue?logo=go)](https://github.com/KEINOS/whereami/actions/workflows/go-versions.yml "Supported versions")
[![Go Reference](https://pkg.go.dev/badge/github.com/KEINOS/whereami.svg)](https://pkg.go.dev/github.com/KEINOS/whereami/ "View document")

[![Opened Issues](https://img.shields.io/github/issues/KEINOS/whereami?color=lightblue&logo=github)](https://github.com/KEINOS/whereami/issues "opened issues")
[![PR](https://img.shields.io/github/issues-pr/KEINOS/whereami?color=lightblue&logo=github)](https://github.com/KEINOS/whereami/pulls "Pull Requests")

- [GolangCI Lint](https://golangci-lint.run/) rules: [.golangci-lint.yml](https://github.com/KEINOS/whereami/blob/main/.golangci.yml)
- To run tests in a container:
  - `docker-compose --file ./.github/docker-compose.yml run v1_17`
  - This will run the below on Go 1.17:
    - `go test -cover -race ./...`
    - `golangci-lint run`
    - `golint ./...`
- Branch to PR:
  - `main`
  - ( It is recommended that [DraftPR](https://github.blog/2019-02-14-introducing-draft-pull-requests/) be done first to avoid duplication of work )

## License

- [MIT](https://github.com/KEINOS/whereami/blob/main/LICENSE). Copyright: (c) 2021 [KEINOS and the WhereAmI contributors](https://github.com/KEINOS/whereami/graphs/contributors).

## Acknowledgment

### Service Providers Used

We would like to thank the following service providers that this command uses to discover global/public IP addresses.

- [https://ipinfo.io/](https://ipinfo.io/)
- [https://inet-ip.info/](https://inet-ip.info/)
- [http://inetclue.com/](http://inetclue.com/)
- [https://toolpage.org/](https://en.toolpage.org/tool/ip-address)
- [https://ipinfo.io/](https://ipinfo.io/)
<!-- Disabled due to the issue #2 // - [https://whatismyip.com/](https://www.whatismyip.com/) -->

> **This command requests these providers in random order and returns the first IP address with the same response**. As soon as 3 of the same IP address are returned, the command stops and prints that IP address.
> If you notice that a provider is not working or not responding properly, please [report an issue](https://github.com/KEINOS/whereami/issues).
