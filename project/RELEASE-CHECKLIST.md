This issue tracks progress of the upcoming release of Docker for Mac/Windows.

### Release plan

- [ ] Reschedule tickets that will not be fixed for this release (@djs55, @dgageot and all)

#### Release preparations on Monday

- [ ] Announce release plans in #pinata-dev (@magnuss)
- [ ] Pre-release meeting Monday (@magnuss): [Agenda](https://docs.google.com/a/docker.com/document/d/1b_Mfe3XjW8OOCoUfyy4luY9yZOSakJvXy3Ooz4JjK5A/edit?usp=sharing)
- Mac (@djs55)
  - [ ] Prepare Mac [CHANGELOG](https://github.com/docker/pinata/blob/master/CHANGELOG)
  - [ ] Build Mac UI complete binary (master build is UI complete)
- Windows (@dgageot)
  - [ ] Prepare Windows [CHANGELOG](https://github.com/docker/pinata/blob/master/win/CHANGELOG)
  - [ ] Build Windows UI complete binary
- Docs (@londoncalling) (may be delayed)
  - [ ] Prepare release notes based on [Mac CHANGELOG](https://github.com/docker/pinata/blob/master/CHANGELOG)
  - [ ] Prepare release notes based on [Windows CHANGELOG](https://github.com/docker/pinata/blob/master/win/CHANGELOG)
  - [ ] Update docs based Windows release binary when ready
  - [ ] Update docs based on Mac release binary when read

#### Release preparations on Tuesday

- [ ] Windows binaries (@dgageot)
  - [ ] Build RC and announce in #pinata-dev
  - [ ] Test RC (see TODO for test plan)
- [ ] Mac binaries (@magnuss)
  - [ ] Build RC - announce in #pinata-dev
  - [ ] Test RC (see TODO for test plan)
- [ ] If RCs works fine, tag release binaries (@magnuss, @dgageot)
- [ ] Update License.tar.gz (@frenchben)

#### Synchronised release of binaries

- [ ] Release Win/Mac binaries with release notes (@magnuss)
- [ ] Announce in #ship-it, #pinata and #pinata-dev (@magnuss)

#### Post-release

- [ ] Deploy docs to docs.docker.com (@londoncalling)
- [ ] Bump Mac version (@jeanlaurent)
- [ ] Bump Win version (@dgageot)