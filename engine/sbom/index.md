---
title: Generate the SBOM for Docker images
description: Viewing the Software Bill of Materials (SBOM) for Docker images
keywords: Docker, sbom, Anchore, images, Syft, security
toc_min: 1
toc_max: 2
---

A Software Bill Of Materials (SBOM) is analogous to a packing list for a shipment. It lists all the components that make up the software, or were used to build it. For container images, this includes the operating system packages that are installed (for example, ca-certificates) along with language-specific packages that the software depends on (for example, Log4j). The SBOM could include a subset of this information or even more details, like the versions of components and their source.

> **Warning**
>
> The `docker sbom` command is currently experimental. This may change or be removed from future releases.
{: .warning }

The experimental `docker sbom` command allows you to generate the SBOM of a container image. Today, it does this by scanning the layers of the image using the [Syft project](https://github.com/anchore/syft) but in future it may read the SBOM from the image itself or elsewhere.


## Simple use

To output a tabulated SBOM for an image,  use `docker sbom <image>:<tag>`:

```console
$ docker sbom neo4j:4.4.5
Syft v0.43.0
 ✔ Loaded image            
 ✔ Parsed image            
 ✔ Cataloged packages      [385 packages]

NAME                      VERSION                        TYPE
... 
bsdutils                  1:2.36.1-8+deb11u1             deb
ca-certificates           20210119                       deb
...
log4j-api                 2.17.1                         java-archive  
log4j-core                2.17.1                         java-archive  
...
```

The output includes both system packages and software libraries used by applications in the container image.

## Output formatting and saving outputs

You can view the SBOM output in standard formats like [SPDX](https://spdx.dev){: target="_blank" rel="noopener" class="_"} and [CycloneDX](https://cyclonedx.org){: target="_blank" rel="noopener" class="_"} along with the Syft and GitHub formats using the `--format` option.

```console
$ docker sbom --format spdx-json alpine:3.15
{
 "SPDXID": "SPDXRef-DOCUMENT",
 "name": "alpine-3.15",
 "spdxVersion": "SPDX-2.2",
 "creationInfo": {
  "created": "2022-04-06T21:13:32.035571Z",
  "creators": [
   "Organization: Anchore, Inc",
   "Tool: syft-[not provided]"
  ],
  "licenseListVersion": "3.16"
 },
 "dataLicense": "CC0-1.0",
 "documentNamespace": "https://anchore.com/syft/image/alpine-3.15-4b1b99d8-bbb5-4426-af8e-c510189134ab",
 "packages": [
  {
   "SPDXID": "SPDXRef-1e3f3285636676f3",
   "name": "alpine-baselayout",
   "licenseConcluded": "GPL-2.0-only",
   "description": "Alpine base dir structure and init scripts",
   "downloadLocation": "https://git.alpinelinux.org/cgit/aports/tree/main/alpine-baselayout",
   "externalRefs": [
    {
...
}
```

These outputs are more verbose and contain more information than the default tabulated output.

By default, the command outputs the SBOM to stdout. You can save the output to a file by specifying one with the `--output` flag.

```console
$ docker sbom --format spdx-json --output sbom.json alpine:3.15
Syft v0.43.0
 ✔ Loaded image            
 ✔ Parsed image            
 ✔ Cataloged packages      [14 packages]

$ cat sbom.json
{
 "SPDXID": "SPDXRef-DOCUMENT",
 "name": "alpine-3.15",
 "spdxVersion": "SPDX-2.2",
...
}
```

## Feedback

Thanks for trying the Docker SBOM CLI plugin. We’d love to hear from you. You can provide feedback and report any bugs through the Issues tracker in the [docker/[sbom-cli-plugin](https://github.com/docker/sbom-cli-plugin){: target="_blank" rel="noopener" class="_"} GitHub repository.

