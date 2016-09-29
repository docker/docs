#r "bin/FAKE/tools/FakeLib.dll"
#r "System.Management.Automation"
open Fake
open Fake.AssemblyInfoFile
open Fake.FileUtils
open Fake.Git
open System.Management.Automation
open System.Net
open System.Text.RegularExpressions
open System.Web
open Fake.AppVeyor

let buildDir = "build/win/"
let packageDir = "build/"
let testBuildDir = "tests/win/"
let buildResourceDir = buildDir @@ "resources/"
let nUnitPath = "bin/nunit/NUnit.ConsoleRunner/tools/"
let testDir = "TestResults/"
let toolsBinDir = "bin/"
let srcDir = "src"
let resourceDir = srcDir @@ "Resources/"
let binDir = resourceDir @@ "bin/"
let solutionFile = srcDir @@ "Docker.Windows.sln"
let solutionInfoCs = srcDir @@ "SolutionInfo.cs"
let testDLL = "Docker.Tests.dll"
let defaultSignToolPath = "C:/Program Files (x86)/Windows Kits/10/bin/x86"
let toolsPath = environVarOrDefault "SIGNTOOLPATH" defaultSignToolPath
let isRunningWithAppVeyor = (buildServer = BuildServer.AppVeyor)
let skipTests = (None <> (environVarOrNone "SKIP_TESTS"))
let shouldSign = (None <> (environVarOrNone "PFX"))
let isOnMaster = ("master" = AppVeyorEnvironment.RepoBranch) || (not isRunningWithAppVeyor)
let isOnPR = (None <> (environVarOrNone "APPVEYOR_PULL_REQUEST_NUMBER"))
let isOnTag = AppVeyorEnvironment.RepoTag
let tag = AppVeyorEnvironment.RepoTagName
let sha1 = if isRunningWithAppVeyor then (getCurrentSHA1 ".") else "unknown"
let downloadDir = (environVarOrFail "HOME") @@ ".pinata" @@ "downloads"

let dockerUrl = "https://experimental.docker.com/builds/Windows/x86_64/docker-1.12.2-rc1.zip"
let dockerdUrl = "https://download.docker.com/components/engine/windows-server/cs-1.12/docker.zip"
let dockerMachineUrl = "https://github.com/docker/machine/releases/download/v0.8.2/docker-machine-Windows-x86_64.exe"
let dockerComposeUrl = "https://github.com/docker/compose/releases/download/1.8.1/docker-compose-Windows-x86_64.exe"
let dockerCredentialUrl = "https://github.com/docker/docker-credential-helpers/releases/download/v0.2.0/docker-credential-wincred-v0.2.0-amd64.zip"
let qemuImgUrl = "https://cloudbase.it/downloads/qemu-img-win-x64-2_3_0.zip"

// https://ci.appveyor.com/project/docker/vpnkit/build/1.0.330
// COMMIT 9cb6374ebfd0656961901478e9fc8cf65d000678
let slirpUrl = "https://ci.appveyor.com/api/buildjobs/mywqnrovhqyxvyfo/artifacts/com.docker.slirp.exe"

// https://ci.appveyor.com/project/docker/datakit/build/1.0.292
// COMMIT 09081b419292627b5992ab4a2d5f95cc0257479b
let dbUrl = "https://ci.appveyor.com/api/buildjobs/p6456kmbuixod9pu/artifacts/com.docker.db.exe"
let dbZlibUrl = "https://ci.appveyor.com/api/buildjobs/p6456kmbuixod9pu/artifacts/zlib1.dll"

type Version = {
    Raw : string
    DockerVersion : string
    Build : string
}

let Exec command args =
  let result = Fake.ProcessHelper.Shell.Exec(command, args)
  if result <> 0 then failwithf "%s exited with error %d" command result

let GoBuild dir exe =
  Exec "go" ("build -i -o " + resourceDir + exe + " " + dir)

let GoBuildInTools dir exe =
  Exec "go" ("build -i -o " + toolsBinDir + exe + " " + dir)

let GoTest dir =
  if not skipTests then Exec "go" ("test -race " + dir)

let BatDirectoryName directory =
  let fullPackageDir = FullName directory
  fullPackageDir.[0..fullPackageDir.Length-2]

let updateChannel =
  if isOnTag && tag.StartsWith "win-beta-v" then
    "beta"
  elif isOnTag && tag.StartsWith "win-test-v" then
    "test"
  elif isOnTag && tag.StartsWith "win-v" then
    "stable"
  elif isOnMaster then
    "master"
  else
    "default"

let Version =
  let fullVersion =
    match (GetAttribute "AssemblyInformationalVersion" solutionInfoCs) with
    | Some(v) -> v.Value
    | None -> failwith "Unable to find version"

  let versionArray = split '-' fullVersion
  let dockerVersion = versionArray.[0]
  let build = if isRunningWithAppVeyor then AppVeyorEnvironment.BuildNumber else "0"

  {Raw = fullVersion; DockerVersion = dockerVersion; Build = build}

let GetAssemblyVersion = (fun _ ->
    Version.DockerVersion + "." + Version.Build
)

let GetHumanVersion = (fun _ ->
    let channelQualifier = if updateChannel = "beta" then "" else ("-" + updateChannel)

    Version.Raw + channelQualifier
)

let GetMSIVersion = (fun _ ->
  if isRunningWithAppVeyor then
    "2.0." + AppVeyorEnvironment.BuildNumber
  else
    "2.0.65535" //FIX: to avoid auto-updating and have dual installs because of signing (#1986)
)

let signFile fileToSign name =
  WriteBytesToFile "pfx" (System.Convert.FromBase64String (environVar "PFX"))

  let arguments = sprintf """sign /a /f pfx /p "%s" /d "%s" /du "https://www.docker.com" /t "http://timestamp.verisign.com/scripts/timestamp.dll" "%s" """ (environVar "PFXPASSWORD") name fileToSign

  let result = ExecProcess (fun info ->
                 info.FileName <- toolsPath @@ "signtool.exe"
                 info.Arguments <- arguments) System.TimeSpan.MaxValue
  if result <> 0 then failwithf "Error during sign call"

  PowerShell.Create()
      .AddScript("if ((New-TimeSpan -End $($(get-AuthenticodeSignature \"" + fileToSign + "\").SignerCertificate.NotAfter.DateTime)).Days -lt 30) { throw \"The certificate expires in less than 30 days\" }")
      .Invoke()
      |> Seq.iter (printfn "%O")

  DeleteFiles ["pfx"]

let DownloadFile (url : string) (dest : string) =
  let tmpFile = dest + ".tmp"
  DeleteFile dest
  use client = new WebClient()
  client.DownloadFile(url, tmpFile)
  Rename dest tmpFile

let Sanitize (url : string) =
  url.Replace('/', '-').Replace(':', '-')

let DownloadToCache (url : string) =
  CreateDir downloadDir
  let dest = (downloadDir @@ Sanitize(url))
  if not (TestFile dest) then DownloadFile url dest
  dest

let UnzipSingleFile (dest : string) (name : string) (zipFile : string) = 
  Unzip (downloadDir @@ "tmp") zipFile
  CopyFile dest (downloadDir @@ "tmp" @@ name)

Target "DownloadMoby" (fun _ ->
  let commit = trim (ReadFileAsString (".." @@ "v1" @@ "moby" @@ "COMMIT"))
  logfn "Moby commit: [%s]" commit

  let tarFile = downloadDir @@ "moby-" + commit + ".tgz"
  let tmpFile = tarFile + ".tmp"
  let tarFolder = downloadDir @@ "moby-" + commit
  let isoFile = tarFolder @@ "mobylinux-efi.iso"

  if not (TestFile tarFile) then
    Exec "go" "get -u github.com/justincormack/regextract"
    Exec ((environVarOrFail "GOPATH") @@ "bin" @@ "regextract") ("--output=" + tmpFile + " mobylinux/media:" + commit)
    Rename tarFile tmpFile

  if not (TestFile isoFile) then
    CreateDir tarFolder
    Exec (toolsBinDir @@ "tar" @@ "TarTool.exe") (tarFile  + " " + tarFolder)

  CopyFile (resourceDir @@ "mobylinux.iso") isoFile
)

Target "Clean" (fun _ ->
  CleanDirs [buildDir; packageDir; testDir]
)

Target "BuildProxy" (fun _ ->
  GoTest "../v1/docker_proxy"
  GoBuild "../v1/cmd/com.docker.proxy/" "com.docker.proxy.exe"
)

Target "BuildDbCli" (fun _ ->
  GoBuild "../v1/cmd/9pdb/" "com.docker.9pdb.exe"
)

Target "BuildMobyDiag" (fun _ ->
  GoBuild "../v1/cmd/moby-diag-dl/" "moby-diag-dl.exe"
)

Target "CopyResources" (fun _ ->
  CreateDir buildResourceDir
  cp_r resourceDir buildResourceDir
  DeleteFile (buildResourceDir @@ "bin" @@ ".gitignore" )
  WriteStringToFile false (buildResourceDir @@ "sha1") sha1
  WriteStringToFile false (buildResourceDir @@ "UpdateChannel") updateChannel
)

Target "DownloadResources" (fun _ ->
  Unzip (resourceDir @@ "qemu-img") (DownloadToCache qemuImgUrl)
  CopyFile (resourceDir @@ "com.docker.slirp.exe") (DownloadToCache slirpUrl)
  CopyFile (resourceDir @@ "com.docker.db.exe") (DownloadToCache dbUrl)
  CopyFile (resourceDir @@ "zlib1.dll") (DownloadToCache dbZlibUrl)
  CopyFile (binDir @@ "docker-machine.exe") (DownloadToCache dockerMachineUrl)
  CopyFile (binDir @@ "docker-compose.exe") (DownloadToCache dockerComposeUrl)
  UnzipSingleFile (binDir @@ "docker.exe") "docker/docker.exe" (DownloadToCache dockerUrl)
  UnzipSingleFile (resourceDir @@ "dockerd.exe") "docker/dockerd.exe" (DownloadToCache dockerdUrl)
  UnzipSingleFile (binDir @@ "docker-credential-wincred.exe") "docker-credential-wincred" (DownloadToCache dockerCredentialUrl)
)

Target "CopyIcon" (fun _ ->
  let channel = updateChannel
  if channel.Equals("stable") then
    CopyFile (srcDir @@ "Docker.Windows/icon.ico") (srcDir @@ "Docker.Windows/stable.ico")
)

Target "RenameService" (fun _ ->
  DeleteFile (buildDir @@ "com.docker.service")
  DeleteFile (buildDir @@ "com.docker.service.config")
  Rename (buildDir @@ "com.docker.service") (buildDir @@ "Docker.Service.exe")
  Rename (buildDir @@ "com.docker.service.config") (buildDir @@ "Docker.Service.exe.config")
)

Target "RestorePackages" (fun _ ->
  Fake.Paket.Restore (fun p ->
      {
        p with
          WorkingDir = srcDir
      })
)

Target "BuildSolution" (fun _ ->
  let setParams defaults =
    { defaults with
        Targets = ["Build"]
        Verbosity = Some(Minimal)
        Properties =
          [
            "FakeBuildOutputPath", (directoryInfo ".").FullName @@ buildDir
            "FakeTestOutputPath", (directoryInfo ".").FullName @@ testBuildDir
            "Optimize", "True"
            "DebugSymbols", "True"
            "Configuration", "Release"
            "PreBuildEvent", ""
            "PostBuildEvent", ""
          ]
    }
  build setParams solutionFile |> DoNothing
)

FinalTarget "UploadTestResultsXml" (fun _ ->
  Fake.AppVeyor.UploadTestResultsXml Fake.AppVeyor.TestResultsType.NUnit3 testDir
)

Target "KillNUnit" (fun _ ->
  killProcess "nunit-agent"
)

Target "UnitTest" (fun _ ->
  ActivateFinalTarget "KillNUnit"
  let resourcesFullPath = Fake.FileSystemHelper.FullName buildResourceDir
  Fake.EnvironmentHelper.setEnvironVar "D4W_RESOURCES_PATH" resourcesFullPath
  if isRunningWithAppVeyor then ActivateFinalTarget "UploadTestResultsXml"
  let nunitWhere = if isRunningWithAppVeyor then "cat != CantRunOnAppVeyor" else ""
  [testBuildDir @@ testDLL]
		|> Fake.Testing.NUnit3.NUnit3 (fun p ->
			{p with
        ToolPath = nUnitPath @@ "nunit3-console.exe"
        ResultSpecs = [ testDir @@ "TestResult.xml"]
        Where = nunitWhere
        ShadowCopy = false})
)

Target "UpdateVersion" (fun _ ->
  let version = GetAssemblyVersion()
  let humanVersion = GetHumanVersion()

  UpdateAttributes solutionInfoCs [
    Attribute.Version version
    Attribute.FileVersion version
    Attribute.InformationalVersion humanVersion
  ]
)

Target "BuildInstaller" (fun _ ->
  !! (buildDir @@ "*.xml") |> Seq.iter (fun file -> DeleteFile file)
  Exec "./bin/Package.bat" (GetMSIVersion() + " " + (BatDirectoryName buildDir) + " " + (BatDirectoryName packageDir) + " " + GetAssemblyVersion() + " \"" + GetHumanVersion() + "\"")
)

Target "SignExecutable" (fun _ -> signFile (buildDir @@ "Docker for Windows.exe") "Docker" )

Target "SignInstaller" (fun _ -> signFile (packageDir @@ "InstallDocker.msi") "Docker installer" )

Target "SignCli" (fun _ -> signFile (buildDir @@ "DockerCli.exe") "Docker CLI" )

Target "Upload" (fun _ ->
  let channel = if isOnPR then "pr" else updateChannel
  let version = GetAssemblyVersion()
  let humanVersion = GetHumanVersion()

  let changelog = ReadFile "CHANGELOG" |> Seq.skip 2 |> Seq.takeWhile (fun s -> not <| s.StartsWith("###")) |> Seq.fold (fun r s -> r + s + "\n") ""
  WriteStringToFile false (packageDir @@ "NOTES") changelog

  printfn "====> Upload on %s for version %s <=====" channel version

  GoBuildInTools "../v1/docker-release/" "docker-release.exe"

  let action =
    match channel with
      | "stable" -> "--prod upload"
      | "beta" -> "--prod upload"
      | _ -> sprintf """--human "%s" publish""" humanVersion

  Exec (toolsBinDir @@ "docker-release.exe") ("--channel " + channel + " --arch win --build " + version + " " + action + " " + (packageDir @@ "InstallDocker.msi") + " " + (packageDir @@ "NOTES"))
)

Target "DestroyVM" (fun _ ->
  let script = ReadFileAsString (resourceDir @@ "MobyLinux.ps1")
  PowerShell.Create().AddScript(script).AddParameter("Destroy").Invoke() |> Seq.iter (printfn "%O")
)

Target "GenerateLicenses" (fun _ ->
  GoBuildInTools "./licenses" "generate-license.exe"
  Exec (toolsBinDir @@ "generate-license.exe") "licenses ../v1/opam ../v1/vendor"
  MoveFile resourceDir "OSS-LICENSES.txt"
)

Target "Test" DoNothing         // build and run unit test
Target "Installer" DoNothing    // build and Create an installer
Target "Build" DoNothing        // Clean build ( default target )
Target "Package" DoNothing      // Clean and build and installer
Target "AppVeyor" DoNothing     // Target for AppVeyor CI

"Clean" ?=> "CopyResources"

"BuildProxy"
  ==> "BuildDbCli"
  ==> "BuildMobyDiag"
  ==> "DownloadResources"
  ==> "DownloadMoby"
  ==> "CopyResources"
  ==> "CopyIcon"
  ==> "RestorePackages"
  =?> ("UpdateVersion", isRunningWithAppVeyor)
  ==> "BuildSolution"
  ==> "RenameService"
  =?> ("UnitTest", not skipTests)
  ==> "Test"

"Test"
  =?> ("SignExecutable", shouldSign)
  =?> ("SignCli", shouldSign)
  ==> "BuildInstaller"
  =?> ("SignInstaller", shouldSign)
  ==> "Installer"

"Clean"
   ==> "Installer"
   ==> "Package"

"Clean"
  ==> "Test"
  ==> "Build"

"Package"
  =?> ("Upload", isRunningWithAppVeyor)
  ==> "AppVeyor"

RunTargetOrDefault "Build"
