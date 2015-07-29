# Preparing development environment (for maintainers)

1. Install Inno Setup 5 non-unicode: http://www.jrsoftware.org/isdl.php (`isetup-5.x.x-unicode.exe`).

2. Install `kSignCMD`: http://codesigning.ksoftware.net/ (click "Download kSign"
and then "FREE DOWNLOAD" or "Click Here to Download kSign", which will likely be
a link to http://cdn1.ksoftware.net/ksign_installer.exe)

3. Install the `docker-code-signing.pfx` certificate somewhere (the instructions
below assume `Z:\sven\src\docker\windows-installer\docker-code-signing.pfx`);
you will need the password (`d_get_from_core` below).

4. Open `Toolbox.iss` in the Inno Setup Compiler.  It has a few constants at
the top that are important to make note of (especially `MyAppVersion` and the
path variables `dockerCli`, `dockerMachineCli`, `kitematicSetup`, `b2dIso`, `msysGit`, `virtualBoxMsi`, and `virtualBoxCommon`).

5. Launch **Inno Setup Compiler** by opening `Toolbox.iss` and add code signing
by applying the following steps:

- Click "Tools" --> "Configure Sign Tools" > "Add"
- "Name of the Sign Tool:" `ksign`
- "Command of the Sign Tool:" `"C:\Program Files (x86)\kSign\kSignCMD.exe" /f Z:\sven\src\docker\windows-installer\docker-code-signing.pfx /p d_get_from_core $p`

# Releasing a new version

Update the versions of the dependencies in `bundle.sh`.

Update `#define MyAppVersion` line in `Toolbox.iss`.

# Downloading bundle dependencies

Open a git bash window in this directory and run script:

    ./bundle.sh

This should be downloading dependencies with their correct versions to `bundle\`
folder where the Inno Setup Compiler can pick up from.

# Compiling the installer

After configuring, open `Toolbox.iss` with Inno Setup Compiler and hit
**'Build'**. The results will be in the `Output` folder.

This can be done through commandline as welll. Launch a cmd.exe shell from this
directory and run:

    "c:\Program Files (x86)\Inno Setup 5\ISCC.exe" Toolbox.iss
