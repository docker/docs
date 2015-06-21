# Notary is still a work in progress and we invite contributions and reviews from the security community. It will need to go through a formal security review process before it should be used in production.


# Notary

notary is a tool for publishing and consuming trusted collections of content. Publishers can digitally sign collections and consumers can verify integrity and origin of content. This ability is built on a straightforward key management and signing interface to create signed collections and configure trusted publishers.

notary is based on [The Update Framework](http://theupdateframework.com/), a secure general design for the problem of software distribution and updates. By using TUF, notary achieves a number of key advantages:


* **Survivable Key Compromise**: Content publishers must manage keys in order to sign their content. Signing keys may be compromised or lost so systems must be designed in order to be flexible and recoverable in the case of key compromise. TUF's notion of key roles is utilized to separate responsibilities across a hierarchy of keys such that loss of any particular key (except the root role) by itself is not fatal to the security of the system.
* **Freshness Guarantees**: Replay attacks are a common problem in designing secure systems, where previously valid payloads are replayed to trick another system. The same problem exists in the software update systems, where old signed can be presented as the most recent. notary makes use of timestamping on publishing so that consumers can know that they are receiving the most up to date content. This is particularly important when dealing with software update where old vulnerable versions could be used to attack users.
* **Configurable Trust Thresholds**: Oftentimes there are a large number of publishers that are allowed to publish a particular piece of content. For example, open source projects where there are a number of core maintainers. Trust thresholds can be used so that content consumers require a configurable number of signatures on a piece of content in order to trust it. Using thresholds increases security so that loss of individual signing keys doesn't allow publishing of malicious content.
* **Signing Delegation**: To allow for flexible publishing of trusted collections, a content publisher can delegate part of their collection to another signer. This delegation is represented as signed metadata so that a consumer of the content can verify both the content and the delegation.
* **Use of Existing Distribution**: notary's trust guarantees are not tied at all to particular distribution channels from which content is delivered. Therefore, trust can be added to any existing content delivery mechanism.
* **Untrusted Mirrors and Transport**: All of the notary metadata can be mirrored and distributed via arbitrary channels.


# Using Notary
Lets try using notary.

First, lets initiate a notary collection called `example.com/scripts`

```sh
notary init example.com/scripts
```

Now, look at the keys you created as a result of initialization
```sh
notary keys
```

Cool, now add a local file `install.sh` and call it `v1`
```sh
notary add example.com/scripts v1 install.sh
```

Wouldn't it be nice if others could know that you've signed this content? Use `publish` to publish your collection to your default notary-server
```sh
notary publish example.com/scripts
```

Now, others can pull your trusted collection
```sh
notary list example.com/scripts
```

More importantly, they can verify the content of your script by using `notary verify`:
```sh
curl example.com/install.sh | notary verify example.com/scripts v1 | sh
```
