Performance testing
===================

TL;DR: we need to crowd-source performance results from developer machines
and do our best to record usable data.

The results of a performance test will vary with

- version of the code under test: this is what we're trying to discover

but will also vary with

- hardware: obviously a faster computer will win
- whether the hardware is on battery or not: low battery has been observed
  to cause 50% network performance loss
- OS version: optimisations in (for example) the hypervisor.framework
  over time may affect results
- free disk space (SSDs are faster when empty)
- other VMs, if the system is virtualised
- other concurrently running software

The scripts in this directory will try to help by:

1. running pre-flight scripts to look for confounding variables
   - a non-idle CPU
   - a system on battery rather than power
   - virtualisation
2. recording test results in a directory per hardware configuration
3. recording as much "uncooked" context data as possible per test

In an ideal world we would use a large set of dedicated customer-representative
software and hardware versions. However this is expensive and hence the
need to "crowd-source" the data instead from developer's machines.

Running tests
-------------

Test results are recorded in a directory per hardware device.
It's up to the caller to 

The following environment variables can be set:

- PINATA_PERF_TESTS -- The list of tests to run.

- PINATA_PERF_VERSION -- The version of Piñata being tested (if not a
  release build etc). Most usefully this can be set to
  "dev.<something>" which indicates an in development build. If set to
  "dev.<int>.<something>" then <int> is the sort order, otherwise the
  strings are sorted lexically (with explicit sort order sorted before
  implicit)

- PINATA_PERF_HWID -- Sets the hardware id. By default this is the
  current machine for running tests and the ID of the benchmarking Mac
  mini ("EC420ECC-EDEE-50AC-9B9C-5984E7F4E23C") when analysing. You
  can use bin/get-hardware-id to get the current machine's ID for
  analysis.

Adding a new test
-----------------

Updating a test
---------------

Cosmetic changes to tests are ok, but be careful if the substance of the
test changes. Consider renaming the test e.g. `test` might become `test.2`
