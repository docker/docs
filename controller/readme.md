# Orca Controller
This is the core component for Orca

# Setup
The only thing Orca needs to run is RethinkDB.

* Run RethinkDB: `docker run -it -d --name rethinkdb -P dockerorca/rethinkdb`
* Run Orca: `docker run -it --name -P --link rethinkdb:rethinkdb dockerorca/orca`

