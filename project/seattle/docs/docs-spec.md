# Docs for Seattle

For CSE 1.12, UCP 1.2 and DTR 2.1, the plan is for us to do an internal
workshop with SAs and SEs to collect feedback. After that, we'll start a closed
beta on those products.

In terms of documentation, the plan is to have documentation ready for the
workshop. Then we'll continuously deploy documentation throughout Beta and GA.

## Documentation drafts

To make it productive for everyone and ensure we have these deliverables on
time, it's easier if the dev team creates drafts that Joao can then use
to create the public-facing docs.

A draft has:
* The user story
* Explanation of important concepts the user might not be familiar with
* Steps and commands the user needs to execute the story

A draft doesn't assume that the user:
* Knows the internals of the product
* Knows the product inside out
* Has read every single piece of documentation

[See this example draft to learn more](draft-example.md).

Below you can find a prioritized list of the documentation we need to have ready
for the workshops.
If we don't have time and have to cut tasks from this backlog, we'll cut the
tasks from bottom to top.

Tasks marked with `*` are tasks that we need a team member to create a draft for.

## CS Docker Engine

* Release notes
* Install
* Upgrade

## UCP

* Release notes
* System requirements (Alex)
  * What ports need to be open and why?
* Install UCP for production (Alex)
  * Assuming that you have 5 fresh linux machines, what do you need to do to
  create a UCP cluster with 3 managers and 2 workers
* Docker Editions (Daniel)
  * Same as above, but on the cloud
* Upgrade UCP (Daniel)
  * Assuming you already have a production-grade installation, how do you
  upgrade to this version?
* Uninstall UCP (Adrian)
* Add and remove nodes (Alex)
* Deploy apps (Arunan)
  * From the UI
  * From the CLI
* UCP architecture (Daniel)
* Permission model
* docker/ucp image reference
* Backup and disaster recovery (probably not in Beta)
  * Do we have support for this?

## DTR

* Release notes
* System requirements
  * What ports need to be open and why?
* DTR architecture*
* docker/ucp image reference

## Notary

* Notary integration with UCP and DTR*
  * Why should you care
  * How to enable it / what are the consequences
  * What is the development flow
  * What is the promotion flow
