# Asks nicely to build, test or package Docker for Windows

param(
  [ValidateSet("Package", "Clean", "DestroyVM", "Test", "CopyResources", "AppVeyor", "UpdateVersion", "generateLicenses", "BuildProxy", "BuildDbCli", "Build", "DownloadResources", "DownloadMoby")]
  [string]$target = "Build",
  [switch]$st
)

if ($st) {$singleTarget="-st"}

Push-Location $PSScriptRoot
Try
{
  $fakeExe = ".\bin\FAKE\tools\Fake.exe"
  if (!(Test-Path $fakeExe))
  {
    & bin\nuget.exe install FAKE -OutputDirectory bin -Version 4.39.0 -ExcludeVersion
  }

  $nunitExe = ".\bin\nunit\NUnit.ConsoleRunner\tools\nunit3-console.exe"
  if (!(Test-Path $nunitExe))
  {
    & bin\nuget.exe install NUnit.ConsoleRunner -OutputDirectory bin\nunit -Version 3.4.1 -ExcludeVersion
  }

  $paketExe = ".paket\paket.exe"
  if (!(Test-Path $paketExe))
  {
    & ".paket\paket.bootstrapper.exe"
  }

  & $fakeExe ".\build.fsx" $target $singleTarget

  Exit $LastExitCode
}
Finally
{
  Pop-Location
}

Exit 1
