; Inno Setup script for the Retro Shooter Windows installer.
;
; Build with scripts/build-installer.sh (WSL/Linux) or in CI. Directly:
;   ISCC.exe /DAppVersion=0.1.0 installer\retro-shooter.iss
; Expects the game at dist\RetroShooter.exe (built with -H windowsgui).

#ifndef AppVersion
  #define AppVersion "0.0.0-dev"
#endif

#define AppName      "Retro Shooter"
#define AppExe       "RetroShooter.exe"
#define AppPublisher "grump606"
#define AppURL       "https://github.com/grump606/retro-shooter"

[Setup]
; AppId identifies the app for upgrades and uninstall. Never change it.
AppId={{6F1C2A9E-4B7D-4E35-9C1A-8D2F5E7B3A41}
AppName={#AppName}
AppVersion={#AppVersion}
AppVerName={#AppName} {#AppVersion}
AppPublisher={#AppPublisher}
AppPublisherURL={#AppURL}
AppSupportURL={#AppURL}/issues

; Per-user install: no admin prompt. {autopf} becomes
; %LOCALAPPDATA%\Programs, and shortcuts go to the user's desktop/Start menu.
PrivilegesRequired=lowest
DefaultDirName={autopf}\{#AppName}
DisableProgramGroupPage=yes
DisableDirPage=auto

ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible
MinVersion=10.0

OutputDir=..\dist
OutputBaseFilename=RetroShooter-Setup-{#AppVersion}
SetupIconFile=..\build\windows\icon.ico
UninstallDisplayIcon={app}\{#AppExe}
UninstallDisplayName={#AppName}
WizardStyle=modern
Compression=lzma2/max
SolidCompression=yes

; Close a running copy of the game before upgrading or uninstalling.
CloseApplications=yes

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
; Ticked by default; the user can untick it in the wizard.
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"

[Files]
Source: "..\dist\{#AppExe}"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\README.md"; DestDir: "{app}"; DestName: "README.txt"; Flags: ignoreversion
Source: "..\CREDITS.md"; DestDir: "{app}"; DestName: "CREDITS.txt"; Flags: ignoreversion

[Icons]
Name: "{autoprograms}\{#AppName}"; Filename: "{app}\{#AppExe}"; WorkingDir: "{app}"
Name: "{autodesktop}\{#AppName}"; Filename: "{app}\{#AppExe}"; WorkingDir: "{app}"; Tasks: desktopicon

[Run]
Filename: "{app}\{#AppExe}"; Description: "{cm:LaunchProgram,{#AppName}}"; Flags: nowait postinstall skipifsilent
