This is still a work in progress but will shortly be a fully compliant 
Go implementation of [The Update Framework (TUF)](http://theupdateframework.com/).

This implementation was originally forked from [flynn/go-tuf](https://github.com/flynn/go-tuf),
however in attempting to add delegations I found I was making such
significant changes that I could not maintain backwards compatibility
without the code becoming overly convoluted.

This implementation retains the same 3 Clause BSD license present on 
the original flynn implementation.

TODOs:

- [X] Add Targets to existing repo
- [X] Sign metadata files
- [X] Refactor TufRepo to take care of signing ~~and verification~~
- [ ] Ensure consistent capitalization in naming (TUF\_\_\_ vs Tuf\_\_\_)
- [ ] Make caching of metadata files smarter
- [ ] ~~Add configuration for CLI commands. Order of configuration priority from most to least: flags, config file, defaults~~ Notary should be the official CLI
- [ ] Reasses organization of data types. Possibly consolidate a few things into the data package but break up package into a few more distinct files
- [ ] Comprehensive test cases
- [ ] Delete files no longer in use
- [ ] Fix up errors. Some have to be instantiated, others don't, the inconsistency is annoying.
- [X] Bump version numbers in meta files (could probably be done better)
