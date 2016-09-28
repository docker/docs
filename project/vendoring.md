# Orca and Vendoring

Godep sucks, but that's what we're using for now.  We might switch to
another tool, but for now, here are some quick tips-and-tricks to try
to make updating things a little less painful.

One area where Godep sucks particularly hard is multi-developer workflow.
If we don't all have the same versions of all the components in our
$GOPATH, all hell breaks lose when you try to update things.  The good
news is there's a little trick that can make this a little less painful.
**Don't use your normal GOPATH location when updating Godep!**

## One time setup
```sh
rm -rf ${HOME}/go_orca/src/github.com/docker
mkdir -p ${HOME}/go_orca/src/github.com/docker
ln -s ${GOPATH}/src/github.com/docker/orca ${HOME}/go_orca/src/github.com/docker/orca
export GOPATH=$HOME/go_orca
godep restore
```

Once you've done that, you now have an alternate GOPATH with the exact versions Orca uses.

# Adding new stuff

**Make sure you've got GOPATH pointing to your alternate first!**

```sh
go get mynewdep
godep update ./...
```

# Changing version
**Make sure you've got GOPATH pointing to your alternate first!**

```sh
pushd $GOPATH/src/XXX
git fetch --all
git checkout SOMENEWVERSION
popd
godep update ./...
make build
# fix problems...
```


# Scorched earth

Sometimes, (often in fact) Godep goes bonkers and needs a clean slate :-(  To add insult to injury
the latest Godep version doesn't like to work in a dirty tree, so this pattern can save your bacon

```sh
rm -rf Godep
git add Godep
git commit -m "nuke godep"
godep save ./...
git add Godep
git commit -m "re-add godep"
git rebase -i HEAD~2
# squash them into one change
```
