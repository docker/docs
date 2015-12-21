# Docker open source project template

This directory contains a template for setting up a new open source
project.

## Checklist

- [ ] Update the "Project name" and description of your project in the README.md.
- [ ] Add details of the project's maintainers to the MAINTAINERS file.
- [ ] Setup DCO checks with "leeroy". Instructions can be found [in the README](https://github.com/docker/leeroy).
- [ ] Request the maintainers to be added to the maintainers mailinglist.
- [ ] Request the maintainers to be "voiced" in the #docker-maintainers IRC channel.
- [ ] Add the standard "labels" for triaging issues, as described in [REVIEWING.md](https://github.com/docker/docker/blob/master/project/REVIEWING.md)
      and [ISSUE-TRIAGE.md](https://github.com/docker/docker/blob/master/project/ISSUE-TRIAGE.md)


### After the repository is public

If you're project is ready to go public:

- [ ] Create a pull request to have the project added to the central
      MAINTAINERS file in the docker/opensource repository. Make sure
      to run `make maintainers` to update the MAINTAINERS file with your
      changes. An example pull request can be found here; https://github.com/docker/opensource/pull/45
- [ ] Remove this file from the repository :)
- [ ] Spread the word!


Congratulations, you're done! 
