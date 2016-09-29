The example test
================

All tests have:

- a README.md to explain what the test is for
- an executable `configuration` which prints to stdout a string
  containing the "configuration". This should contain pertinent stuff
  like which filesystem implementation or network implementation. It
  will be used as a directory name.
- a Makefile with a default target which will set up and run the test.
  Raw log output that we want to keep should go in `logs/`. Gnuplot-format
  data rows should go in `results.dat`. There should also be a `make clean`
  target.

