---
title: Integrate with multiple registries
description: Integrate UCP with multiple registries
keywords: trust, registry, integrate, UCP, DTR
redirect_from:
  - /datacenter/ucp/3.0/guides/admin/configure/integrate-with-multiple-registries/
---

Universal Control Plane can pull and run images from any image registry,
including Docker Trusted Registry and Docker Store.

If your registry uses globally-trusted TLS certificates, everything works
out of the box, and you don't need to configure anything. But if your registries
use self-signed certificates or certificates issues by your own Certificate
Authority, you need to configure UCP to trust those registries.

## Trust Docker Trusted Registry

To configure UCP to trust a DTR deployment, you need to update the
[UCP system configuration](ucp-configuration-file.md) to include one entry for
each DTR deployment:

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

You only need to include the port section if your DTR deployment is running
on a port other than 443.

You can customize and use the script below to generate a file named
`trust-dtr.toml` with the configuration needed for your DTR deployment.

```
# Replace this url by your DTR deployment url and port
DTR_URL=https://dtr.example.org
DTR_PORT=443

dtr_full_url=${DTR_URL}:${DTR_PORT}
dtr_ca_url=${dtr_full_url}/ca

# Strip protocol and default https port
dtr_host_address=${dtr_full_url#"https://"}
dtr_host_address=${dtr_host_address%":443"}

# Create the registry configuration and save it it
cat <<EOL > trust-dtr.toml

[[registries]]
  # host address should not contain protocol or port if using 443
  host_address = $dtr_host_address
  ca_bundle = """
$(curl -sk $dtr_ca_url)"""
EOL
```

You can then append the content of `trust-dtr.toml` to your current UCP
configuration to make UCP trust this DTR deployment.

## Where to go next

- [Integrate with LDAP by using a configuration file](external-auth/enable-ldap-config-file.md)
