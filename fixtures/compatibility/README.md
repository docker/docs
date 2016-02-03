This directory contains sample repositories from different versions of Notary client (TUF metadata, trust anchor certificates, and private keys), in order to test backwards compatibility (that newer clients can read old-format repositories).

Notary client makes no guarantees of future-compatibility though (that is, repositories produced by newer clients may not be able to be read by old clients.)

Relevant information for repositories:

- `notary0.1`
	- GUN: `docker.com/notary0.1/samplerepo`
	- key passwords: "randompass"
	- targets:

		```
		   NAME                                  DIGEST                                SIZE (BYTES)
		---------------------------------------------------------------------------------------------
		  LICENSE   9395bac6fccb26bcb55efb083d1b4b0fe72a1c25f959f056c016120b3bb56a62   11309
  		```
  	- It also has a changelist to add a `.gitignore` target, that hasn't been published.
