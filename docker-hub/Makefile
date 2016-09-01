.PHONY: all binary build cross default docs docs-build docs-shell shell test test-unit test-integration test-integration-cli test-docker-py validate

# env vars passed through directly to Docker's build scripts
# to allow things like `make DOCKER_CLIENTONLY=1 binary` easily
# `docs/sources/contributing/devenvironment.md ` and `project/PACKAGERS.md` have some limited documentation of some of these
DOCKER_ENVS := \
	-e BUILDFLAGS \
	-e DOCKER_CLIENTONLY \
	-e DOCKER_EXECDRIVER \
	-e DOCKER_GRAPHDRIVER \
	-e TESTDIRS \
	-e TESTFLAGS \
	-e TIMEOUT
# note: we _cannot_ add "-e DOCKER_BUILDTAGS" here because even if it's unset in the shell, that would shadow the "ENV DOCKER_BUILDTAGS" set in our Dockerfile, which is very important for our official builds

# to allow `make DOCSDIR=docs docs-shell` (to create a bind mount in docs)
DOCS_MOUNT := $(if $(DOCSDIR),-v $(CURDIR):/docs/content/docker-hub/)

# to allow `make DOCSPORT=9000 docs`
DOCSPORT := 8000

# Get the IP ADDRESS
DOCKER_IP=$(shell python -c "import urlparse ; print urlparse.urlparse('$(DOCKER_HOST)').hostname or ''")
HUGO_BASE_URL=$(shell test -z "$(DOCKER_IP)" && echo localhost || echo "$(DOCKER_IP)")
HUGO_BIND_IP=0.0.0.0

GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
DOCKER_IMAGE := docker$(if $(GIT_BRANCH),:$(GIT_BRANCH))
DOCKER_DOCS_IMAGE := docs-base$(if $(GIT_BRANCH),:$(GIT_BRANCH))


DOCKER_RUN_DOCS := docker run --rm -it $(DOCS_MOUNT) -e AWS_S3_BUCKET -e NOCACHE

# for some docs workarounds (see below in "docs-build" target)
GITCOMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)

default: docs

test: docs-build
	$(DOCKER_RUN_DOCS) -p $(if $(DOCSPORT),$(DOCSPORT):)8000 \
		-v $(CURDIR):/docs/content/docker-hub/ \
		-e DOCKERHOST "$(DOCKER_DOCS_IMAGE)" \
		hugo server \
			--log=true --watch=true \
			--port=$(DOCSPORT) --baseUrl=$(HUGO_BASE_URL) --bind=$(HUGO_BIND_IP)


docs: docs-build
	$(DOCKER_RUN_DOCS) -p $(if $(DOCSPORT),$(DOCSPORT):)8000 -e DOCKERHOST "$(DOCKER_DOCS_IMAGE)" hugo server --port=$(DOCSPORT) --baseUrl=$(HUGO_BASE_URL) --bind=$(HUGO_BIND_IP)

docs-draft: docs-build
	$(DOCKER_RUN_DOCS) -p $(if $(DOCSPORT),$(DOCSPORT):)8000 -e DOCKERHOST "$(DOCKER_DOCS_IMAGE)" hugo server --buildDrafts="true" --port=$(DOCSPORT) --baseUrl=$(HUGO_BASE_URL) --bind=$(HUGO_BIND_IP)


docs-shell: docs-build
	$(DOCKER_RUN_DOCS) -p $(if $(DOCSPORT),$(DOCSPORT):)8000 "$(DOCKER_DOCS_IMAGE)" bash


docs-build:
#	( git remote | grep -v upstream ) || git diff --name-status upstream/release..upstream/docs ./ > ./changed-files
#	echo "$(GIT_BRANCH)" > GIT_BRANCH
#	echo "$(AWS_S3_BUCKET)" > AWS_S3_BUCKET
#	echo "$(GITCOMMIT)" > GITCOMMIT
	docker build -t "$(DOCKER_DOCS_IMAGE)" .

# use screenshot container to update screenshots
NOAUTHSCREENSHOT := docker run -it --rm --name screen -v $(CURDIR):/srv --env-file=noauthpasswords.env svendowideit/screenshot
SCREENSHOT := docker run -it --rm --name screen -v $(CURDIR):/srv --env-file=passwords.env svendowideit/screenshot
NOLINKS := docker run -it --rm --name screen -v $(CURDIR):/srv --env-file=passwords.env --env-file=nolinks.env svendowideit/screenshot
GITHUB_DOCSUSER := docker run -it --rm --name screen -v $(CURDIR):/srv --env-file=passwords.env --env-file=githubdocs.env svendowideit/screenshot

# testing
testimage:
	#$(NOAUTHSCREENSHOT) https://hub-beta.docker.com/ hub-images/register-web.png 1280px
	#$(SCREENSHOT) https://hub-beta.docker.com/explore/ hub-images/dashboard.png 1280px
	$(SCREENSHOT) https://hub-beta.docker.com/organizations/ orgs.png 1280px

t2:
	$(NOAUTHSCREENSHOT) https://hub-beta.docker.com/ hub-images/register-web.png 1280px

docs-images:
	# non-authenticated
	$(NOAUTHSCREENSHOT) https://hub-beta.docker.com/ hub-images/register-web.png 1280px
	$(NOAUTHSCREENSHOT) https://hub-beta.docker.com/login/ hub-images/login-web.png 1280px
	# authenticated
	$(SCREENSHOT) https://hub-beta.docker.com/explore/ hub-images/dashboard.png 1280px
	$(SCREENSHOT) https://hub-beta.docker.com/organizations/ hub-images/orgs.png 1280px
	# $(SCREENSHOT) https://hub-beta.docker.com/ hub-images/deploy_key.png 1280px
	$(SCREENSHOT) https://hub-beta.docker.com/r/docsorg/private/~/settings/collaborators/ hub-images/org-repo-collaborators.png 1280px
	$(SCREENSHOT) https://hub-beta.docker.com/ hub-images/repos.png 1280px
	# $(SCREENSHOT) https://hub-beta.docker.com/ hub-images/invite.png 1280px
	# $(SCREENSHOT) https://hub-beta.docker.com/ hub-images/build-trigger.png 1280px
	$(SCREENSHOT) https://hub-beta.docker.com/ hub-images/hub.png 1280px
	$(SCREENSHOT) https://hub-beta.docker.com/u/docsorg/dashboard/teams/boomteam/ hub-images/groups.png 1280px
	$(SCREENSHOT) https://hub-beta.docker.com/r/library/busybox/tags/# hub-images/busybox-image-tags.png 1280px*600px
	# docs/hub-images/getting-started.png (needs a new empty account)

nolinks:
	# bitbucket.md, github.md
	# uses the `nolinks` user: an account that has no accounts linked
	$(NOLINKS) https://hub-beta.docker.com/account/authorized-services/ hub-images/authorized-services.png 1280px
	$(NOLINKS) https://hub-beta.docker.com/account/authorized-services/github-permissions/ add-authorized-github-service.png 1280px

# BROKEN
github:
	$(GITHUB_DOCSUSER) https://github.com/docsuser/private/settings/hooks github-side-hook.png 1280px

# BROKEN, wrong URL and needs hand editing to capture the specific UI elements
broken-github-images:
	$(SCREENSHOT) https://hub-beta.docker.com/ hub-images/gh_settings.png
	$(SCREENSHOT) https://hub-beta.docker.com/ hub-images/gh_menu.png
	$(SCREENSHOT) https://hub-beta.docker.com/ hub-images/gh_add_ssh_user_key.png
	$(SCREENSHOT) https://hub-beta.docker.com/ hub-images/gh_team_members.png
	$(SCREENSHOT) https://hub-beta.docker.com/ hub-images/gh-check-user-org-dh-app-access.png
	$(SCREENSHOT) https://hub-beta.docker.com/ hub-images/gh_service_hook.png
	$(SCREENSHOT) https://hub-beta.docker.com/ hub-images/gh-check-admin-org-dh-app-access.png
	$(SCREENSHOT) https://hub-beta.docker.com/ hub-images/gh_org_members.png
	$(SCREENSHOT) https://hub-beta.docker.com/ hub-images/gh_docker-service.png
	$(SCREENSHOT) https://hub-beta.docker.com/ hub-images/gh_repo_deploy_key.png
