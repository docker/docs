Docker Toolbox Testing
======================

Testing is mostly manual for now. Below you'll find a checklist of cases that have been covered by each release.

Things to try:
- Old version of VirtualBox running with old version of Boot2Docker

## Toolbox 1.8.0 RC7

### Mac OS X 10.11

- [ ] Install + Quick Start
- [ ] Install + Migrating Boot2Docker
- [ ] Install + Start Kitematic 

### Mac OS X 10.10

- [ ] Install + Quick Start
- [ ] Install + Migrating Boot2Docker
- [ ] Install + Start Kitematic 

### Mac OS X 10.9

- [ ] Install + Quick Start
- [ ] Install + Migrating Boot2Docker
- [ ] Install + Start Kitematic 

### Windows 10

- Currently unsupported by Toolbox as VirtualBox 5 does not yet support Windows 10.

### Windows 8 (8.1) 

- [x] Install + Quick Start
- [failed] Install + Migrating Boot2Docker
- [x] Install + Start Kitematic 
*User may encourter Docker Machine errors. 
** Boot2docker installer failed to work on the test PC 

Clean install: 
1. Open VirtualBox (if installed), and remove all the VMs by choosing Delete all files. 
2. Open Control Panel -> Program and Features -> Uninstall Virtualbox 
3. Uninstall Docker Toolbox (if installed) 
4. Open "File Explorer" -> C:\Users\YOUR_USERNAME\ and delete .docker folder 
5. Run toolbox installer again

### Windows 7

- [x] Install + Quick Start
- [x] Install + Migrating Boot2Docker
- [x] Uninstalling
*User may encourter Docker Machine errors. 
