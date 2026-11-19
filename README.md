# im-alter

Select an enabled macOS input source by ID. Built with Go and purego.

```sh
CGO_ENABLED=0 go build -o im-alter .
./im-alter com.apple.keylayout.ABC
```

Selection succeeds silently. An unavailable source exits with status 1. Calling
the command without exactly one argument prints usage and exits with status 2.

## Releases

tagpr opens a release pull request on pushes to main. Merging it tags the
release and runs GoReleaser to publish macOS arm64 and amd64 archives and
checksums to GitHub Releases.
