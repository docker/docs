> [!NOTE]
>
> Buildx v0.10 enables support for a minimal [SLSA Provenance](https://slsa.dev/provenance/)
> attestation, which requires support for [OCI-compliant](https://github.com/opencontainers/image-spec)
> multi-platform images. This may introduce issues with registry and runtime support
> (e.g. [Google Cloud Run and AWS Lambda](https://github.com/docker/buildx/issues/1533)).
> You can optionally disable the default provenance attestation functionality
> using `--provenance=false`.
