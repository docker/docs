Agreed upon plan at meeting on July 6 -- 
There are two scenarios: 
Scenario A: GA product
Scenario B: Beta product

There are 3 use cases in each scenario:
1. Unlicensed to Licensed state
2. Licensed to expired state
3. Licensed to out of compliant state (e.g. exceeded node count)

In each use case, we have options of: inform the user or strict enforcement via disabling functionality

We've decided the following table 

Use Case | GA Product| Beta Product
------------ | ------------- | ----------------
Unlicensed to Licensed | Inform at admin login but allow bypass (current UI); Inform throughout GUI | Enforce by refusing GUI login
Licensed to Expired | Inform | Inform
Licensed to Out of Compliant | Inform | Inform

We decided not to shut off any services in the DDC when using "enforce" in any use case. 

In the future, we agree there may be features that should be license enabled, e.g. image scanning only available when using the DDC Pro license by not the DDC Lite license. This is something we'll address in a later release.
