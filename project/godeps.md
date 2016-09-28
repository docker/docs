# Notes on Godeps

Bah!  Godeps sucks!  It just isn't well designed for team development
when multiple people work on multiple projects.

Here are some raw notes that might make life a little less painful when
trying to update godeps in the tree.  The basic model is to create a
secondary GOPATH that maps to the **exact** versions of dependencies
that we use within Orca.

* Initial setup of a secondary GOPATH tree
```
export GOPATH=$HOME/go_orca
mkdir -p ${GOPATH}

# Make sure there's a link in there pointing to this orca tree
if [ ! -d ${GOPATH}/src/github.com/docker/orca ] ; then
    ln -s $(pwd) ${GOPATH}/src/github.com/docker/orca
fi

# Initial population
godep restore

```

**WARNING: YMMV** Sometimes `save` seems to work, sometimes `add` works - I still haven't iterated on this enough to nail down the optimal model so you may have to play around with it a bit...

* Periodic work when trying to update something
```
# update what's in our tree if it got out of sync
godep restore

go get stuff...
or go into the GOPATH tree and move the trees to newer commits

rm -rf Godep vendor

godep save ./...

git add Godep vendor
```
