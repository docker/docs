#define MyAppName "Docker Toolbox"
#define MyAppPublisher "Docker"
#define MyAppURL "https://docker.com"
#define MyAppContact "https://docs.docker.com"

#define b2dIsoPath "..\bundle\boot2docker.iso"
#define dockerCli "..\bundle\docker.exe"
#define dockerMachineCli "..\bundle\docker-machine.exe"
#define kitematic "..\bundle\kitematic"
#define git "..\bundle\Git.exe"
#define virtualBoxCommon "..\bundle\common.cab"
#define virtualBoxMsi "..\bundle\VirtualBox_amd64.msi"

[Setup]
AppCopyright={#MyAppPublisher}
AppId={{FC4417F0-D7F3-48DB-BCE1-F5ED5BAFFD91}
AppContact={#MyAppContact}
AppComments={#MyAppURL}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
ArchitecturesAllowed=x64
ArchitecturesInstallIn64BitMode=x64
DefaultDirName={pf}\{#MyAppName}
DefaultGroupName=Docker
DisableProgramGroupPage=yes
OutputBaseFilename=DockerToolbox
Compression=lzma
SolidCompression=yes
WizardImageFile=windows-installer-side.bmp
WizardSmallImageFile=windows-installer-logo.bmp
WizardImageStretch=no
WizardImageBackColor=$22EBB8
UninstallDisplayIcon={app}\unins000.exe
SetupIconFile=toolbox.ico
ChangesEnvironment=true

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Types]
Name: "full"; Description: "Full installation"
Name: "custom"; Description: "Custom installation"; Flags: iscustom

[Run]
Filename: "{win}\explorer.exe"; Parameters: "{userprograms}\Docker\"; Flags: postinstall skipifsilent; Description: "View Shortcuts in File Explorer"

[Tasks]
Name: desktopicon; Description: "{cm:CreateDesktopIcon}"
Name: modifypath; Description: "Add docker.exe & docker-machine.exe to &PATH"

[Components]
Name: "Docker"; Description: "Docker Client for Windows" ; Types: full custom; Flags: fixed
Name: "DockerMachine"; Description: "Docker Machine for Windows" ; Types: full custom; Flags: fixed
Name: "VirtualBox"; Description: "VirtualBox"; Types: full custom; Flags: disablenouninstallwarning
Name: "Kitematic"; Description: "Kitematic for Windows (Alpha)" ; Types: full custom
Name: "Git"; Description: "Git for Windows"; Types: full custom; Flags: disablenouninstallwarning

[Files]
Source: ".\docker-quickstart-terminal.ico"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#dockerCli}"; DestDir: "{app}"; Flags: ignoreversion; Components: "Docker"
Source: ".\start.sh"; DestDir: "{app}"; Flags: ignoreversion; Components: "Docker"
Source: ".\delete.sh"; DestDir: "{app}"; Flags: ignoreversion; Components: "Docker"
Source: "{#dockerMachineCli}"; DestDir: "{app}"; Flags: ignoreversion; Components: "DockerMachine"
Source: ".\migrate.sh"; DestDir: "{app}"; Flags: ignoreversion; Components: "DockerMachine"
Source: ".\migrate.bat"; DestDir: "{app}"; Flags: ignoreversion; Components: "DockerMachine"
Source: "{#kitematic}\*"; DestDir: "{app}\kitematic"; Flags: ignoreversion recursesubdirs; Components: "Kitematic"
Source: "{#b2dIsoPath}"; DestDir: "{app}"; Flags: ignoreversion; Components: "DockerMachine"; AfterInstall: CopyBoot2DockerISO()
Source: "{#git}"; DestDir: "{app}\installers\git"; DestName: "git.exe"; AfterInstall: RunInstallGit();  Components: "Git"
Source: "{#virtualBoxCommon}"; DestDir: "{app}\installers\virtualbox"; Components: "VirtualBox"
Source: "{#virtualBoxMsi}"; DestDir: "{app}\installers\virtualbox"; DestName: "virtualbox.msi"; AfterInstall: RunInstallVirtualBox(); Components: "VirtualBox"

[Icons]
Name: "{userprograms}\Docker\Kitematic (Alpha)"; WorkingDir: "{app}"; Filename: "{app}\kitematic\Kitematic.exe"; Components: "Kitematic"
Name: "{commondesktop}\Kitematic (Alpha)"; WorkingDir: "{app}"; Filename: "{app}\kitematic\Kitematic.exe"; Tasks: desktopicon; Components: "Kitematic"
Name: "{userprograms}\Docker\Docker Quickstart Terminal"; WorkingDir: "{app}"; Filename: "{pf64}\Git\git-bash.exe"; Parameters: """{app}\start.sh"""; IconFilename: "{app}/docker-quickstart-terminal.ico"; Components: "Docker"
Name: "{commondesktop}\Docker Quickstart Terminal"; WorkingDir: "{app}"; Filename: "{pf64}\Git\git-bash.exe"; Parameters: """{app}\start.sh"""; IconFilename: "{app}/docker-quickstart-terminal.ico"; Tasks: desktopicon; Components: "Docker"

[UninstallRun]
Filename: "{app}\delete.sh"

[UninstallDelete]
Type: filesandordirs; Name: "{localappdata}\..\Roaming\Kitematic"

[Code]
#include "base64.iss"
#include "guid.iss"

var
	restart: boolean;
  TrackingDisabled: Boolean;
	DockerInstallDocs: TLabel;
  TrackingCheckBox: TNewCheckBox;

function uuid(): String;
var
  dirpath: String;
  filepath: String;
  ansiresult: AnsiString;
begin
  dirpath := ExpandConstant('{userappdata}\DockerToolbox');
  filepath := dirpath + '\id.txt';
  ForceDirectories(dirpath);

  Result := '';
  if FileExists(filepath) then
    LoadStringFromFile(filepath, ansiresult);
    Result := String(ansiresult)

  if Length(Result) = 0 then
    Result := GetGuid('');
    StringChangeEx(Result, '{', '', True);
    StringChangeEx(Result, '}', '', True);
    SaveStringToFile(filepath, AnsiString(Result), False);
end;

function WindowsVersionString(): String;
var
	ResultCode: Integer;
	lines : TArrayOfString;
begin
	if not Exec(ExpandConstant('{cmd}'), ExpandConstant('/c wmic os get caption | more +1 > C:\windows-version.txt'), '', SW_HIDE, ewWaitUntilTerminated, ResultCode) then begin
		Result := 'N/A';
		exit;
	end;

	if LoadStringsFromFile(ExpandConstant('C:\windows-version.txt'), lines) then begin
		Result := lines[0];
	end else begin
		Result := 'N/A'
	end;
end;

procedure TrackEvent(name: String);
var
  payload: String;
  WinHttpReq: Variant;
begin
  if TrackingDisabled or WizardSilent() then
    exit;
  try
    payload := Encode64(Format(ExpandConstant('{{"event": "%s", "properties": {{"token": "{#MixpanelToken}", "distinct_id": "%s", "os": "win32", "os version": "%s", "version": "{#MyAppVersion}"}}'), [name, uuid(), WindowsVersionString()]));
    WinHttpReq := CreateOleObject('WinHttp.WinHttpRequest.5.1');
    WinHttpReq.Open('POST', 'https://api.mixpanel.com/track/?data=' + payload, false);
    WinHttpReq.SetRequestHeader('Content-Type', 'application/json');
    WinHttpReq.Send('');
  except
  end;
end;

function NeedRestart(): Boolean;
begin
	Result := restart;
end;

function NeedToInstallVirtualBox(): Boolean;
begin
	// TODO: Also compare versions
	Result := (
		(GetEnv('VBOX_INSTALL_PATH') = '')
		and
		(GetEnv('VBOX_MSI_INSTALL_PATH') = '')
	);
end;

function NeedToInstallGit(): Boolean;
begin
	// TODO: Find a better way to see if Git is installed
	Result := not DirExists('C:\Program Files\Git') or not FileExists('C:\Program Files\Git\git-bash.exe')
end;

procedure DocLinkClick(Sender: TObject);
var
	ErrorCode: Integer;
begin
	ShellExec('', 'https://docs.docker.com/installation/windows/', '', '', SW_SHOW, ewNoWait, ErrorCode);
end;

procedure TrackingCheckBoxClicked(Sender: TObject);
begin
  if TrackingCheckBox.Checked then begin
		TrackingDisabled := False;
    TrackEvent('Enabled Tracking');
  end else begin
    TrackEvent('Disabled Tracking');
    TrackingDisabled := True;
  end;
end;

procedure InitializeWizard;
var
  WelcomePage: TWizardPage;
	TrackingLabel: TLabel;
begin
	DockerInstallDocs := TLabel.Create(WizardForm);
	DockerInstallDocs.Parent := WizardForm;
	DockerInstallDocs.Left := 8;
	DockerInstallDocs.Top := WizardForm.ClientHeight - DockerInstallDocs.ClientHeight - 8;
	DockerInstallDocs.Cursor := crHand;
	DockerInstallDocs.Font.Color := clBlue;
	DockerInstallDocs.Font.Style := [fsUnderline];
	DockerInstallDocs.Caption := '{#MyAppName} installation documentation';
	DockerInstallDocs.OnClick := @DocLinkClick;
	DockerInstallDocs.Visible := True;

  WelcomePage := PageFromID(wpWelcome)

	WizardForm.WelcomeLabel2.AutoSize := True;

  TrackingCheckBox := TNewCheckBox.Create(WizardForm);
  TrackingCheckBox.Top := 168;
  TrackingCheckBox.Left := WizardForm.WelcomeLabel2.Left;
  TrackingCheckBox.Width := WizardForm.WelcomeLabel2.Width;
  TrackingCheckBox.Height := 40;
  TrackingCheckBox.Caption := 'Send one-time, anonymous diagnostics during install.';
  TrackingCheckBox.Checked := True;
  TrackingCheckBox.Parent := WelcomePage.Surface;
  TrackingCheckBox.OnClick := @TrackingCheckboxClicked;

	TrackingLabel := TLabel.Create(WizardForm);
	TrackingLabel.Parent := WelcomePage.Surface;
	TrackingLabel.Font := WizardForm.WelcomeLabel2.Font;
	TrackingLabel.Font.Color := clGray;
  TrackingLabel.Caption := 'This data helps us detect problems and improve the installation experience. We only use it for aggregate statistics and will never share it with third parties.';
	TrackingLabel.WordWrap := True;
	TrackingLabel.Visible := True;
	TrackingLabel.Left := WizardForm.WelcomeLabel2.Left;
	TrackingLabel.Width := WizardForm.WelcomeLabel2.Width;
	TrackingLabel.Top := 200;
	TrackingLabel.Height := 100;

		// Don't do this until we can compare versions
		// Wizardform.ComponentsList.Checked[2] := NeedToInstallVirtualBox();
		Wizardform.ComponentsList.ItemEnabled[2] := not NeedToInstallVirtualBox();
		Wizardform.ComponentsList.Checked[4] := NeedToInstallGit();
		Wizardform.ComponentsList.ItemEnabled[4] := not NeedToInstallGit();
end;

function InitializeSetup(): boolean;
begin
  TrackEvent('Installer Started');
  Result := True;
end;

function NextButtonClick(CurPageID: Integer): Boolean;
begin
  if CurPageID = wpWelcome then begin
      TrackEvent('Continued from Overview');
  end;
  Result := True
end;

procedure RunInstallVirtualBox();
var
	ResultCode: Integer;
begin
	WizardForm.FilenameLabel.Caption := 'installing VirtualBox'
	if not Exec(ExpandConstant('msiexec'), ExpandConstant('/qn /i "{app}\installers\virtualbox\virtualbox.msi" /norestart'), '', SW_HIDE, ewWaitUntilTerminated, ResultCode) then
		MsgBox('virtualbox install failure', mbInformation, MB_OK);
end;

procedure RunInstallGit();
var
	ResultCode: Integer;
begin
	WizardForm.FilenameLabel.Caption := 'installing Git for Windows'
	if Exec(ExpandConstant('{app}\installers\git\git.exe'), '/sp- /verysilent /norestart', '', SW_HIDE, ewWaitUntilTerminated, ResultCode) then
	begin
		// handle success if necessary; ResultCode contains the exit code
		//MsgBox('git installed OK', mbInformation, MB_OK);
	end
	else begin
		// handle failure if necessary; ResultCode contains the error code
		MsgBox('git install failure', mbInformation, MB_OK);
	end;
end;

procedure CopyBoot2DockerISO();
begin
  WizardForm.FilenameLabel.Caption := 'copying boot2docker iso'
  if not ForceDirectories(ExpandConstant('{userdocs}\..\.docker\machine\cache')) then
      MsgBox('Failed to create docker machine cache dir', mbError, MB_OK);
  if not FileCopy(ExpandConstant('{app}\boot2docker.iso'), ExpandConstant('{userdocs}\..\.docker\machine\cache\boot2docker.iso'), false) then
      MsgBox('File moving failed!', mbError, MB_OK);
end;

function MigrateVM() : Boolean;
var
  ResultCode: Integer;
begin
  if NeedToInstallGit() or NeedToInstallVirtualBox() or not FileExists(ExpandConstant('{app}\docker-machine.exe')) then begin
    Result := true
    exit
  end;

  ExecAsOriginalUser('C:\Program Files\Oracle\VirtualBox\VBoxManage.exe', 'showvminfo default', '', SW_HIDE, ewWaitUntilTerminated, ResultCode)
  if ResultCode <> 1 then begin
    Result := true
    exit
  end;

  ExecAsOriginalUser('C:\Program Files\Oracle\VirtualBox\VBoxManage.exe', 'showvminfo boot2docker-vm', '', SW_HIDE, ewWaitUntilTerminated, ResultCode)
  if ResultCode <> 0 then begin
    Result := true
    exit
  end;

  if MsgBox('Migrate your existing Boot2Docker VM to work with the Docker Toolbox? Your existing Boot2Docker VM will not be affected. This should take about a minute.', mbConfirmation, MB_YESNO) = IDYES then
  begin
		TrackEvent('Boot2Docker Migration Started');
    WizardForm.StatusLabel.Caption := 'Migrating Boot2Docker VM...'
    WizardForm.FilenameLabel.Caption := 'This will take a minute...'
    ExecAsOriginalUser(ExpandConstant('{app}\docker-machine.exe'), ExpandConstant('rm -f default > nul 2>&1'), '', SW_HIDE, ewWaitUntilTerminated, ResultCode)
    DelTree(ExpandConstant('{userdocs}\..\.docker\machine\machines\default'), True, True, True);
    ExecAsOriginalUser(ExpandConstant('{app}\migrate.bat'), ExpandConstant('> {localappdata}\Temp\toolbox-migration-logs.txt 2>&1'), '', SW_HIDE, ewWaitUntilTerminated, ResultCode)
    if ResultCode = 0 then
    begin
      TrackEvent('Boot2Docker Migration Succeeded');
      MsgBox('Succcessfully migrated Boot2Docker VM to a Docker Machine VM named "default"', mbInformation, MB_OK);
    end
    else begin
      TrackEvent('Boot2Docker Migration Failed');
      MsgBox('Migration of Boot2Docker VM failed. Please file an issue with the migration logs at https://github.com/docker/toolbox/issues/new.', mbCriticalError, MB_OK);
      Exec(ExpandConstant('{win}\notepad.exe'), ExpandConstant('{localappdata}\Temp\toolbox-migration-logs.txt'), '', SW_SHOW, ewNoWait, ResultCode)
      Result := false
			WizardForm.Close;
      exit;
    end;
  end
  else
  begin
    TrackEvent('Boot2Docker Migration Skipped');
  end;
  Result := true
end;

const
	ModPathName = 'modifypath';
	ModPathType = 'user';

function ModPathDir(): TArrayOfString;
begin
	setArrayLength(Result, 1);
	Result[0] := ExpandConstant('{app}');
end;
#include "modpath.iss"

procedure CurPageChanged(CurPageID: Integer);
begin
	if (CurPageID = wpFinished) then
	begin
      WizardForm.FinishedLabel.AutoSize := True;
			WizardForm.FinishedLabel.Font.Style := [fsBold];
			WizardForm.FinishedLabel.Caption :=
				'Upgrade your default Docker VM using the command:' + \
				#13#10 + \
				#13#10 + \
				'         `docker-machine upgrade default`';
  end;
end;

procedure CurStepChanged(CurStep: TSetupStep);
var
	Success: Boolean;
begin
	Success := True;
	if CurStep = ssPostInstall then
  begin
    if IsTaskSelected(ModPathName) then
			ModPath();
    if not WizardSilent() then
      Success := MigrateVM();
		if Success then
			trackEvent('Installer Finished');
  end;
end;
