# TODOs

## Features

- [ ] Explore [go-attestation](https://github.com/google/go-attestation) as an option for AuthN w/ MSFT?

## Makefile

- [ ] Add `attest` target
    - govulncheck
    - Fuzzing
    - syft
    - sbomgen
    - oss-review-toolkit
    - OSV-Scanner
    - Snyk
    - CycloneDX
    - SPDX
    - trivy
    - zap
    - SigStore OIDC signing
    - SigStore Transparency Log (Trillian)
    - OSS-Fuzz
    - Syzkaller
    - ossf/package-analysis: Open Source Package Analysis
    - Vulnerability Exploitability eXchange (VEX)
    - guacsec/guac
    - cve-search/git-vuln-finder: Finding potential software vulnerabilities from git commit messages
    - chaoss/augur: Python library and web service for Open Source Software Health and Sustainability metrics & data collection. You can find our documentation and new contributor information easily here: https://chaoss.github.io/augur/ and learn more about Augur at our website https://augurlabs.io
    - IBM/CBOM: Cryptography Bill of Materials
    - AppThreat/blint: BLint is a Binary Linter to check the security properties, and capabilities in your executables. It is powered by lief.
    - Contour: A Practical System for Binary Transparency
    - OWASPs SCA tools
    - nexB/scancode-toolkit: ScanCode detects licenses, copyrights, package manifests & dependencies and more by scanning code ... to discover and inventory open source and third-party packages used in your code.
    - GitBOM. It’s not Git or SBOM
    - tern-tools/tern: Tern is a software composition analysis tool and Python library that generates a Software Bill of Materials for container images and Dockerfiles. The SBOM that Tern generates will give you a layer-by-layer view of what's inside your container in a variety of formats including human-readable, JSON, HTML, SPDX and more.
    - DWARF 5 Standard
    - AppThreat/rosa: An experiment that looks very promising so far.
    - eBay/sbom-scorecard: Generate a score for your sbom to understand if it will actually be useful.
    - ossf/scorecard: Security Scorecards - Security health metrics for Open Source
    - Lynis - Security auditing and hardening tool for Linux/Unix
    - anchore/grype: A vulnerability scanner for container images and filesystems
    - Real-time VEX
    - ossillate-inc/packj: The vetting tool 🚀 behind our large-scale security analysis platform to detect malicious/risky open-source packages
    - analysis-tools-dev/static-analysis: ⚙️ A curated list of static analysis (SAST) tools for all programming languages, config files, build tools, and more.
    - https://docs.sigstore.dev/fulcio/overview
    - https://docs.sigstore.dev/cosign/overview
    - https://docs.sigstore.dev/rekor/overview
    - https://landscape.openssf.org/sigstore
    - https://cas.codenotary.com/
    - https://witness.dev/ || https://github.com/testifysec/witness
    - https://github.com/puerco/tejolote
    - https://github.com/marketplace/actions/in-toto-run
    - https://github.com/in-toto/github-action
    - https://github.com/slsa-framework/slsa-github-generator
    - technosophos/helm-gpg: Chart signing and verification with GnuPG for Helm.
    - https://github.com/notaryproject/notary
    - https://github.com/deislabs/ratify
    - https://github.com/aws-solutions/verifiable-controls-evidence-store
    - https://github.com/johnsonshi/image-layer-provenance
    - https://github.com/oras-project/artifacts-spec/
    - https://www.youtube.com/watch?v=UrLdEYVASak
    - https://github.com/transmute-industries/verifiable-actions/tree/main
    <!-- - https://github.com/bureado/awesome-software-supply-chain-security?tab=readme-ov-file#frameworks-and-best-practice-references -->
- [ ] Add `verify` target (cosign)
- [ ] Add `benchmark` target (flamegraphs, traces, etc)
- [ ] Add `coverage` target (codecov)
- [ ] Add `pre-commit` target (pre-commit)
- [ ] Add `codeql` target (pre-commit)

## GitHub Actions

- [ ] Release publishing
- [ ] Release changelog management
- [ ] Release TODO management

<!-- https://blog.ralch.com/articles/golang-conditional-compilation/ -->
<!-- https://opensource.com/article/22/4/go-build-options -->

<!-- Go compilation offers several options to control the build process. These options can be specified using the go build command and its flags.

### Common go build options

- -o file: Specifies the output file name for the compiled binary.
- -v: Enables verbose output, printing the names of compiled packages.
- -work: Prints the path to the temporary work directory and prevents its deletion after the build.
- -x: Prints the commands executed during the build process.
- -buildmode mode: Sets the build mode, which can be exe (default), shared, pie, or plugin.
- -gcflags flags: Passes flags to the Go compiler. For example, -gcflags=-S prints assembly code.
- -ldflags flags: Passes flags to the linker.
- -tags tags: Specifies build tags to include conditional compilation.
- -race: Enables the race detector.
- -a: Forces rebuilding of packages that are already up-to-date.
- -n: Prints the commands that would be executed but does not execute them.
- -p n: Specifies the number of parallel builds.
- -trimpath: Removes file path prefixes from the resulting executable.

### Build modes

- exe: The default mode, producing a standalone executable file.
- shared: Generates a shared library.
- pie: Creates a position-independent executable.
- plugin: Builds a plugin.

### Cross-compilation

Go supports cross-compilation, allowing you to build binaries for different operating systems and architectures. This is achieved by setting the GOOS and GOARCH environment variables before running go build.

### Optimization

The Go compiler performs several optimizations by default. Some techniques to improve performance include: Using inline functions, Avoiding unnecessary memory allocation, Leveraging escape analysis, and Using memory alignment.

### Debugging

#### Go provides several flags for debugging:

- -dwarf: Generates DWARF symbols for debugging.
- -traceprofile file: Writes an execution trace to a file.
- Compiler debugging flags (e.g., -d, -v, -W).

### Additional tools

- go tool compile: Invokes the compiler directly.
- go tool link: Invokes the linker.
- go fix: Updates packages to use newer APIs.
- go env: Displays Go environment variables.
- go mod tidy: Cleans up the go.mod file. -->
