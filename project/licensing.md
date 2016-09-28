# UCP Licensing

This document describes licensing in UCP.  Our general philosophy is
"do no harm" - unlicensed or expired licenses should not prevent the
use of the product.


   | **Trial** | **Enterprise Biz Day** | **Enterprise Biz Critical**
---|-----------|------------------------|----------------------------
Features | Unlimited | Unlimited | Unlimited
Support | E-mail/no-sla| 12x5 | 24x7
License Duration| 30 days free | Monthly or Annual | Monthly or Annual

## UCP Controller

When the system is first installed it is in expired mode.  A license key
must be manually entered into the admin interface.  Only one license key
is used at a time (changing the license key replaces the old with the new
key) In the license UI, the user may select automatic renewal or manual.
(see below for renewal details)

The licenses are paid for based on the number of engines, not the
number of UCP instances.  License keys can be re-used across multiple
UCP instances (they are not tied to a specific UCP instance.)  UCP will
not enforce node counts, but will detect if the node count is over the
license node count, and will report a banner.

## Expired mode:

When a license expires, the following will occur:

* Banner in the UI nagging users their license has expired
    * In the initial unlicensed mode, the message will explain they
      need to download a free trial or purchase a license
* CLI info and or version command will report the expired license
* No upgrades of the controller nodes
* No customer support
    * Support dumps will never be disabled so they can pay up at the time of the incident if they call support


Notable items that will **not** be impacted by licensing (topics that at
one point were considered):

* Adding nodes will be allowed regardless of license state, node count, or current number of nodes
* Disabling phone-home analytics will always be allowed

## Auditing

No auditing will be implemented in this release.  We explored potentially
leveraging our existing Mixpanel integration, but the benefits of
keeping that anonymous outweigh the added license auditing stream.
We instead will rely on users policing themselves by viewing the output
of the /info command as well as the banners that appear in the UI when
they're not in compliance with their license.

## Renewal

In manual renewal mode, the user may at any time browse on hub, and
renew their license and download a new license key.  This updated license
key can then be loaded into the ucp admin UI.  The ucp controller will
never contact the license server directly in this mode, and will instead
rely on the expiration date within the license key.  When the license
is close to expiring the "nag" banner should appear warning the user
that their license is about to expire; they should renew it manually,
or enable automatic renewal.

In automatic renewal mode, when the cert is close to expiring (exact
number of days/hours TBD) the server will start periodically attempting to
contact the license server to get a refreshed license key.  As soon as it is
successful, the old license will be replaced with the new license.  If the
license expires before the server is able to successfully renew, it will
degrade to "expired" mode and continue to attempt to renew the license.
If an error message is returned from the server, the last error reported
will be saved and displayed in the UI (so we can tell users things like
"your credit card expired" etc.)  The user may manually enter a new
license key at any point and if the license key is valid, and not within
the "close to expiring" window the automatic renewal logic will stop.

## License Server

We will leverage the existing DTR license server.

* https://docker.atlassian.net/wiki/display/DHE/DTR+Licensing

Today, when hub makes a request to the license server, it specifies
the license expiration based on what the customer has paid for (months,
years).  When the user downloads a license key, the license server then
generates a key with at most a 24 hour expiration.  This license key
is then manually loaded into DTR, which will call the license server
and renew the license key periodically.  As long as the license hasn't
expired, the license server will renew the key every time DTR requests
a refreshed key.

For UCP, we want to reduce the requirement for continual online renewal,
and will set the license key to expire at the same time as the license the
customer paid for (one year, one month, or possibly multiple months into
the future.)  The downloaded license key will be valid for that duration.

When the user wishes to renew, through hub, the existing license will
be set with a new expiration (same KeyID) and a new download of the license
key will show the new expiration, with the same KeyID.

If the user changes the number of engines, the license will be updated,
but need not be replaced the on the UCP instance unless the user has
exceeded the count and wants to eliminate the warning banner.

In this model, an "online" vs. "offline" setting in the license is no
longer relevant.  In essence, every license is both online and offline,
and it is up to the user to decide if they want the client (UCP) to
automatically attempt renewal, or prefer to manually go to hub, renew,
and download an updated license key, then load that into the UCP admin UI.

At a minimum, the license server will require minor modifications to support
license keys with the same expiration as the underlying license.  Hub will require
branding updates for UCP.

As the enterprise products will be sold as a suite, at this time the
intention is to use the existing (or possibly updated) tiers across
both DTR and UCP.  The UCP implementation will ignore the tiers in this
release (any valid license from the license server will be accepted.)


# Implementation notes: (remove once fleshed out)

Accessing the low-level API:

```bash
export UCP_URL="https://$(echo $DOCKER_HOST | cut -f3 -d/ )"

# Dump out the current license configuration
curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    ${UCP_URL}/api/config/license | jq "."


# Update the license (this will be rejected as malformed)
# Note that the details is read-only and will be ignored if passed in
curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    -XPOST -d '{"auto_refresh":true,"license_config":{"key_id":"bogus","private_key":"privfoo","authorization":"authfoo"}}' \
    ${UCP_URL}/api/config/license | jq "."

# Now load a valid license:
# WARNING - THIS IS A VALID LICENSE - DON'T LET IT OUT IN THE WILD!!!
curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    -XPOST -d '{"auto_refresh":true,"license_config":{"key_id":"4Hg5DGMH78wN5ZjNjbau_agErRqNE5aQ-R3MnUiNYdGg","private_key":"y2dhoSi4jAYKpclaXnmd1R_RJkYy7ySmKCis9e1JpfH6","authorization":"ewogICAicGF5bG9hZCI6ICJleUpsZUhCcGNtRjBhVzl1SWpvaU1qQTBNeTB3TkMweU4xUXhPRG95T0RvMU1sb2lMQ0owYjJ0bGJpSTZJbUpmZVV0VVRFeFVjVUppTmxCdE5ucE5hRlUxT1ZGRlJUQnFaRVYxUVhoWFIzUmZZM1Z0TVZsb1NqQTlJaXdpYldGNFJXNW5hVzVsY3lJNk1UQXNJbXhwWTJWdWMyVlVlWEJsSWpvaVQyWm1iR2x1WlNJc0luUnBaWElpT2lKRmRtRnNkV0YwYVc5dUluMCIsCiAgICJzaWduYXR1cmVzIjogWwogICAgICB7CiAgICAgICAgICJoZWFkZXIiOiB7CiAgICAgICAgICAgICJqd2siOiB7CiAgICAgICAgICAgICAgICJlIjogIkFRQUIiLAogICAgICAgICAgICAgICAia2V5SUQiOiAiSjdMRDo2N1ZSOkw1SFo6VTdCQToyTzRHOjRBTDM6T0YyTjpKSEdCOkVGVEg6NUNWUTpNRkVPOkFFSVQiLAogICAgICAgICAgICAgICAia2lkIjogIko3TEQ6NjdWUjpMNUhaOlU3QkE6Mk80Rzo0QUwzOk9GMk46SkhHQjpFRlRIOjVDVlE6TUZFTzpBRUlUIiwKICAgICAgICAgICAgICAgImt0eSI6ICJSU0EiLAogICAgICAgICAgICAgICAibiI6ICJ5ZEl5LWxVN283UGNlWS00LXMtQ1E1T0VnQ3lGOEN4SWNRSVd1Szg0cElpWmNpWTY3MzB5Q1lud0xTS1Rsdy1VNlVDX1FSZVdSaW9NTk5FNURzNVRZRVhiR0c2b2xtMnFkV2JCd2NDZy0yVVVIX09jQjlXdVA2Z1JQSHBNRk1zeER6V3d2YXk4SlV1SGdZVUxVcG0xSXYtbXE3bHA1blFfUnhyVDBLWlJBUVRZTEVNRWZHd20zaE1PX2dlTFBTLWhnS1B0SUhsa2c2X1djb3hUR29LUDc5ZF93YUhZeEdObDdXaFNuZWlCU3hicGJRQUtrMjFsZzc5OFhiN3ZaeUVBVERNclJSOU1lRTZBZGo1SEpwWTNDb3lSQVBDbWFLR1JDSzR1b1pTb0l1MGhGVmxLVVB5YmJ3MDAwR08td2EyS044VXdnSUltMGk1STF1VzlHa3E0empCeTV6aGdxdVVYYkc5YldQQU9ZcnE1UWE4MUR4R2NCbEp5SFlBcC1ERFBFOVRHZzR6WW1YakpueFpxSEVkdUdxZGV2WjhYTUkwdWtma0dJSTE0d1VPaU1JSUlyWGxFY0JmXzQ2SThnUVdEenh5Y1plX0pHWC1MQXVheVhyeXJVRmVoVk5VZFpVbDl3WE5hSkIta2FDcXo1UXdhUjkzc0d3LVFTZnREME52TGU3Q3lPSC1FNnZnNlN0X05lVHZndjhZbmhDaVhJbFo4SE9mSXdOZTd0RUZfVWN6NU9iUHlrbTN0eWxyTlVqdDBWeUFtdHRhY1ZJMmlHaWhjVVBybWs0bFZJWjdWRF9MU1ctaTd5b1N1cnRwc1BYY2UycEtESW8zMGxKR2hPXzNLVW1sMlNVWkNxekoxeUVtS3B5c0g1SERXOWNzSUZDQTNkZUFqZlpVdk43VSIKICAgICAgICAgICAgfSwKICAgICAgICAgICAgImFsZyI6ICJSUzI1NiIKICAgICAgICAgfSwKICAgICAgICAgInNpZ25hdHVyZSI6ICJMWEtUclBfVTJEUGVlWlBZaFlZdjZJTm1BU1dERWYtMEV5ZDdTb1hwdDdVeDRVOU11VVF2dzlCTFotaDQ0R3JYeXVaeGxYMDdjY2xRc0NRNWZLV1JRdy1XQkphQi04UlRqWXFPaXh3UldCZGlFaDM1c0tNSjRpSzFKbkMxLXpPN1JxTkdycmhscGgtZHM5QUhBT3c0THM5REJRWmZFVURzUzl6X296R1liOHlDR0FmTS1EOWsxTFF5djBoREJJd0ZnMDRmSUFiTkZucmlZckRiRzc4WGl6LWdXQlJieExCdXZxV2lnLTFfOGtab1VkMVRpX1JSREowTThqQkl0WlpDSjRfci0yaGtZdng0SVVMcVU1ajFwV2RjM2pwWE1qUS1lUl85YUxKSUtveEFtME1FYkZSWlBTWU1RTlpTY2pXN2dVZFBJRGxxX3VwenM4LUdRTWhVQW00Q2dVLVlXay1fZXZmSWhVRVh6ZUFtdFQzYWk1R3loYWRUYzNmUDBMOFlUT1ZsRXZWd251WFh6RV9aMjdkQnRpUXlrNk1LWmFjX09mX3ZpT2Rfazd5QktsbzZQVDZSMmJyMUtZTWVReVdpSnBrM3NLNlRacUdnWDEtWThDeFQ2cnowMjAzUmYyN3lwdnM4N1BKN0IyZHpqcHVIbGhFTjZTUDNpTGNYNThUYWphWm5GUnZPWGljTENtVUM5MGFtSFM5RjBxTXZBVFZpN25lUXRjZjlrWVVsdm15ODRaRFc5S0VKcUZrbV93NUtKRURNMGpUcW5DV0x5NVpkR29BaWsyd3J2Wm1vY1Z3OXE2a3YwWjA3dklEUVd2NzBTc0liR0pTZGpyYVh6d0hvZ2NVUTZHU3EyYXpub1cyS2k2RkZQaDlQZHJ5U3lKeEF6Si13TExvQTN1RSIsCiAgICAgICAgICJwcm90ZWN0ZWQiOiAiZXlKbWIzSnRZWFJNWlc1bmRHZ2lPakUxTVN3aVptOXliV0YwVkdGcGJDSTZJbVpSSWl3aWRHbHRaU0k2SWpJd01UVXRNVEl0TVRCVU1UZzZNekE2TVROYUluMCIKICAgICAgfQogICBdCn0="}}' \
    ${UCP_URL}/api/config/license | jq "."

```



* The license will be stored in the kv store, and shared by all nodes
* Add license information into the version string or info
    * It can be displayed on the dashboard
        * Trial should be in the footer - non obtrusive, but there
        * Expired should take over a bunch of the UI and be a bit annoying
