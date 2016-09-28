# UCP Test Plan
This is a strawman for a UCP test plan.  It is by no means all inclusive and is meant to be a starting point.

# Authentication
- Login with Admin user using Builtin Authenticator
- Login using LDAP Authenticator (configure LDAP in admin settings)
- Logout Admin user (ensure no access to routes)

# Dashboard
- Validate "Applications" link
- Validate "Containers" link
- Validate "Images" link
- Validate "Nodes" link

# User
- View user profile
- Validate change password functionality
- For non-admin, validate team memberships are shown in profile
- Download a client bundle
- Remove a client bundle
- Add the public key of a removed bundle and ensure if it has regained access

# Applications
- Using a client bundle deploy an application
- Confirm application containers are part of the application
- Confirm application stats working properly

# Containers
- Confirm non-UCP containers only by default
- Confirm all containers show when "Show All"
- Validate container actions (restart, stop, remove)
- Validate container info (container details)
- Confirm container logging (container details)
- Confirm container stats reporting (container details)
- COnfirm container console (container details)

# Nodes
- Confirm nodes are shown
- Confirm labels on nodes

# Volumes
- Confirm volumes are listed
- Validate volume creation
- Validate volume removal

# Networks
- Confirm networks are listed
- Validate network creation
- Validate network removal

# Images
- Confirm images are listed
- Validate image pull
- Validate image removal

# Users & Teams (Admin)
- Validate user creation
- Validate team creation
- Validate team management (add / remove users)
- Validate team permissions
    - Create access permission `dev` with permission `Full Control`
    - Create access permission `staging` with permission `Restricted Control`
    - Create access permission `prod` with permission `View Only`

# Settings (Admin)
- Update default log level and confirm
- Configure LDAP (https://docker.atlassian.net/wiki/pages/viewpage.action?pageId=18286175 for configuration)
- Configure DTR
    - Use the first rig from https://docker.atlassian.net/wiki/display/DHE/DTR+CI and follow the instructions at https://docs.docker.com/ucp/dtr-integration/
    - Perform a `docker login`, `docker tag` and `docker push` sequence to DTR from a local docker client.
    - Remove the local image and perform a `docker pull` using a UCP client bundle
    - Switch UCP & DTR configuration to LDAP and repeat


# API / CLI
- Download admin client bundle, load into env, and test docker  commands:
    - `docker ps` (ensure only non-UCP containers are shown)
    - `docker ps -a` (ensure all containers are shown)
    - `docker run --rm alpine echo hello`
    - `docker images`
    - `docker volume ls`
    - `docker network ls`
    - `docker info` (showing correct node/controller information)
    - `docker inspect $(docker run -d --privileged alpine tail -f) | grep Privileged` (preservation of HostConfig in inspect, should be true)

# Access Control
- Create the following users
    - appuser
    - db
    - devuser
- Create the following teams with the following users
    - application
        - appuser
    - db
        - db
    - developers
        - devuser
- Create the following permissions for the following teams
    - application
          - Label: `db` Permission: `View Only`
          - Label: `app` Permission: `Restricted Control`
    - db
          - Label: `db` Permission: `Full Control`
    - developers
          - Label: `dev` Permission: `Full Control`
          - Label: `staging` Permission: `Restricted Control`
          - Label: `prod` Permission: `View Only`

- Deploy the following containers
    - `docker run -ti -d --name app0 --label com.docker.ucp.access.label=app alpine ash`
    - `docker run -ti -d --name app1 --label com.docker.ucp.access.label=app alpine ash`
    - `docker run -ti -d --name db1 --label com.docker.ucp.access.label=db alpine ash`
    - `docker run -ti -d --name db2 --label com.docker.ucp.access.label=db alpine ash`
    - `docker run -ti -d --name dev1 --label com.docker.ucp.access.label=dev alpine ash`
    - `docker run -ti -d --name dev2 --label com.docker.ucp.access.label=dev alpine ash`
    - `docker run -ti -d --name staging1 --label com.docker.ucp.access.label=staging alpine ash`
    - `docker run -ti -d --name staging2 --label com.docker.ucp.access.label=staging alpine ash`
    - `docker run -ti -d --name prod1 --label com.docker.ucp.access.label=prod alpine ash`
    - `docker run -ti -d --name prod2 --label com.docker.ucp.access.label=prod alpine ash`

- User Validation: `appuser`
    - Confirm only `app` and `db` containers are shown
    - Confirm restart, stop for `app`
    - Confirm deny for restart on `db`
    - Confirm deny for console on `app`

- User Validation: `db`
    - Confirm only `db` containers are shown
    - Confirm restart, stop for `app`
    - Confirm deny for restart on `db`

- User Validation: `devuser`
    - Confirm only `dev`, `staging` and `prod` containers are shown
    - Confirm restart, stop for `dev`
    - Confirm console for `dev`
    - Confirm deny for restart on `staging`
    - Confirm deny for console on `staging`
    - Confirm deny for restart on `prod`
    - Confirm deny for console on `prod`
