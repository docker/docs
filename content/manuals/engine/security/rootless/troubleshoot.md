{{< tab name="Arch Linux" >}}

Arch Linux defaults to the `iptables-nft` backend. This may cause networking
issues when running Docker in rootless mode.

If you encounter issues, consider switching to the legacy iptables backend or
ensure nftables compatibility.

See the Arch Linux announcement:
https://archlinux.org/news/iptables-now-defaults-to-the-nft-backend/

{{< /tab >}}
