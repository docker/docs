---
title: Integrate with multiple registries
description: Integrate UCP with multiple registries
keywords: trust, registry, integrate, UCP, DTR
---

Universal Control Plane can pull and run images from any image registry,
including Docker Trusted Registry and Docker Store.

If your registry uses globally-trusted TLS certificates, everything works
out of the box, and you don't need to configure anything.
If your registries use self-signed certificates or use your own Certificate
Authority to sign the public key certificates, you need to configure UCP
to trust those repositories.

## Trust Docker Trusted Registry

To configure UCP to trust a DTR deployment, you need to update the
[UCP system configuration]() to include one entry for each DTR deployment:

```
[[registries]]
  host_address = "dtr.example.org"
  ca_bundle = """
-----BEGIN CERTIFICATE-----
...
-----END CERTIFICATE-----"""

[[registries]]
  host_address = "internal-dtr.example.org:444"
  ca_bundle = """
-----BEGIN CERTIFICATE-----
...
-----END CERTIFICATE-----"""
```

You only need to include the port section, if your DTR deployment is running
on a port other than 443.

You can customize and use the script below to generate a file named
`trust-dtr.toml` with the configuration needed for your DTR deployment.

```
# Replace this url by your DTR deployment url
DTR_URL=https://ddc-prod-dtr.testing.dckr.io:443
DTR_CA_URL=$DTR_URL/ca

# Create the registry configuration and save it it
cat <<EOL > trust-dtr.toml

[[registries]]
  # host address should not contain protocol or port if using 443
  host_address = "$(echo $DTR_URL | sed 's/^https\{0,1\}:\/\///' | sed 's/:443$//')"
  ca_bundle = """
$(curl -sk $DTR_CA_URL)"""
EOL
```

You can then append the content of `trust-dtr.toml` to your current UCP
configuration to make UCP trust this DTR deployment. [Learn more]()
