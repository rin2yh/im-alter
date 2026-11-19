# im-alter

Select an enabled macOS input source by ID. Built with Go and purego.

## Install

Download the archive for your Mac from [GitHub Releases](https://github.com/rin2yh/im-alter/releases).

- Apple Silicon uses the `darwin_arm64` archive.
- Intel uses the `darwin_amd64` archive.

Extract the archive and move the `im-alter` executable to a directory on your
`PATH`.

## Usage

```sh
im-alter com.apple.keylayout.ABC
```

Selection succeeds silently. An unavailable source exits with status 1. Calling
the command without exactly one argument prints usage and exits with status 2.

## Build from source

Run this from the repository root on macOS.

```sh
CGO_ENABLED=0 go build -o im-alter .
```

## Releases

tagpr opens a release pull request on pushes to main. Merging it tags the
release and runs GoReleaser to publish macOS arm64 and amd64 archives and
checksums to GitHub Releases.

## License

MIT. See [LICENSE](./LICENSE).
