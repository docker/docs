{{< tab name="Arch Linux" >}}

Arch Linux now defaults to the `iptables-nft` backend instead of the legacy
`iptables` backend. This change may cause networking issues when running
Docker in rootless mode.

If you encounter networking problems, consider switching to the legacy
iptables backend or ensuring compatibility with nftables.

See the [Arch Linux announcement](https://archlinux.org/news/iptables-now-defaults-to-the-nft-backend/).

{{< /tab >}}

