# Orca Dependency Images

In a perfect world, this directory would be empty.  Alas, we have some
dependencies that don't have official hub images yet.  We have to build
them specifically for our needs so we don't depend on 3rd party images
that could change out from underneath us or disappear.

For any images that mount volumes and have state, version and upgrade
compatibility should be captured so the install/upgrade tool can validate
the users upgrade path is supported.  The images themselves should manage
any data transformations on initial startup.
