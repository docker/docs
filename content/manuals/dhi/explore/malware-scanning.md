---
title: Malware scanning
description: Learn how Docker scans Docker Hardened Images for viruses and malware using ClamAV, and how to view and verify the scan attestation.
keywords: docker hardened images, malware scanning, virus scan, clamav, attestation, image security
weight: 50
---

The Docker Hardened Image (DHI) pipeline scans for viruses and malware as part
of the build process. The scan results are embedded as a signed attestation,
which you can independently retrieve and verify.

## How it works

Docker uses [ClamAV](https://www.clamav.net/), an open source antivirus engine,
to scan every layer of each image. The scan runs automatically during the build
process and checks all files in the image, including files inside archives, for
known viruses and malware signatures.

The scan results are published as a signed attestation attached to the image.
The attestation includes the full ClamAV scan report, including the number of
files scanned, the virus signature database version, and whether any infected
files were detected.

## View the malware scan attestation

You can retrieve the malware scan attestation using the Docker Scout CLI.

1. Use the `docker scout attest get` command with the virus scan predicate type:

   ```console
   $ docker scout attest get \
     --predicate-type https://scout.docker.com/virus/v0.1 \
     --predicate \
     dhi.io/<image>:<tag>
   ```

   > [!NOTE]
   >
   > If the image exists locally on your device, you must prefix the image name
   > with `registry://`. For example, use `registry://dhi.io/python` instead of
   > `dhi.io/python`.

   For example:

   ```console
   $ docker scout attest get \
     --predicate-type https://scout.docker.com/virus/v0.1 \
     --predicate \
     dhi.io/python:3.13
   ```

   The output is a JSON object containing the scanner used and the base64-encoded
   scan report:

   ```json
   {
     "scanner": {
       "report": "<base64-encoded ClamAV report>",
       "uri": "clamav/clamav:stable"
     }
   }
   ```

   Decoding the report shows the full ClamAV output, ending with a scan summary:

   ```text
   ----------- SCAN SUMMARY -----------
   Known viruses: 3627833
   Engine version: 1.5.2
   Scanned directories: 4
   Scanned files: 21
   Infected files: 0
   Data scanned: 44.90 MiB
   Data read: 23.88 MiB (ratio 1.88:1)
   Time: 11.473 sec (0 m 11 s)
   Start Date: 2026:04:12 02:36:19
   End Date:   2026:04:12 02:36:30
   ```

2. Verify the attestation signature. To ensure the attestation is authentic and
   signed by Docker, run:

   ```console
   $ docker scout attest get \
     --predicate-type https://scout.docker.com/virus/v0.1 \
     --verify \
     dhi.io/<image>:<tag> --platform <platform>
   ```

   If the attestation is valid, Docker Scout confirms the signature and shows
   the matching `cosign verify` command.

To view other attestations, such as SBOMs or test results, see [Verify
an image](../how-to/verify.md).
