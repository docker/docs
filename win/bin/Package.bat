echo off
:Loop
if "%1"=="" goto Continue
  goto start
shift
goto Loop
:Continue

:start
REM We hardcode here a few things
set "Manufacturer=Docker Inc."
set "ProductName=Docker"
set "ProductVersion=%1"
set "AssemblyVersion=%4"
set "HumanVersion=%5"
set "ProductId=5F4B755FFFFF"
set "ProductDescription=Docker for Windows"
set "SupportURL=http://docker.com/support"
set "ProductBuildFolder=%2"
set "PackageBuildFolder=%3"

echo "Manufacturer: %Manufacturer%"
echo "Product Name: %ProductName%"
echo "Product Version: %ProductVersion%"
echo "Assembly Version: %AssemblyVersion%"
echo "Human Version: %HumanVersion%"
echo "Product Id: %ProductId%"
echo "Product Description: %ProductDescription%"
echo "Support URL: %SupportURL%"
echo "Product Build Folder: %ProductBuildFolder%"
echo "Package Build Folder: %PackageBuildFolder%"

if ERRORLEVEL 0 goto paraffin else goto error

:paraffin
  echo "Creating %PackageBuildFolder%\%ProductName%"
  md "%PackageBuildFolder%\%ProductName%" 2> NUL
  cd "%~dp0"
  xcopy /y /s "%ProductBuildFolder%" "%PackageBuildFolder%\%ProductName%" > NUL
  del /Q "%PackageBuildFolder%\%ProductName%\*.vshost.exe" 2> NUL
  del /Q "%PackageBuildFolder%\%ProductName%\*.vshost.exe.manifest" 2> NUL
  del /Q "%PackageBuildFolder%\%ProductName%\*.vshost.exe.config" 2> NUL
  paraffin.exe -dir "%PackageBuildFolder%\%ProductName%" -Win64 -x .svn -e BackendCLI -e service -e pdb -e lib -e suo -e ilk -e exp -e obj -e sav -e bak -dirref "%ProductName%RootDir" -g -multiple -c %ProductName% %PackageBuildFolder%\FilesList.wxi
REM  paraffin.exe -dir "%PackageBuildFolder%\%ProductName%" -Win64 -x .svn -e pdb -e iso -e exe -e lib -e ilk -e exp -e obj -e sav -e bak -e vshost.exe -dirref %Manufacturer%RootDir -g -multiple -c %ProductName% %PackageBuildFolder%\FilesList.wxi
  if ERRORLEVEL 0 goto candle else goto error

:candle
REM find install id to use during install/uninstall
  WxiFindId.exe -f "InstallerCli.exe" "%PackageBuildFolder%\FilesList.wxi" > WxiFindDockerInstallerExeId
  set dockerInstallerExeId=
  for /F "delims=" %%i in (WxiFindDockerInstallerExeId) DO set dockerInstallerExeId=%%i
  if NOT DEFINED dockerInstallerExeId goto error
  del /q /f WxiFindDockerInstallerExeId
  echo "InstallerCli.exe File ID: %dockerInstallerExeId%"
REM find executable id to launch after install
  WxiFindId.exe -f "Docker for Windows.exe" "%PackageBuildFolder%\FilesList.wxi" > WxiFindDockerExeId
  set dockerExeId=
  for /F "delims=" %%i in (WxiFindDockerExeId) DO set dockerExeId=%%i
  if NOT DEFINED dockerExeId goto error
  del /q /f WxiFindDockerExeId
  echo "Docker for Windows.exe File ID: %dockerExeId%"
REM find working id 
  FOR /F "delims=" %%i IN ('DirectoryId.exe "%PackageBuildFolder%\FilesList.wxi" "%ProductName%"') DO set workingDirectoryId=%%i
  echo "Working Directory ID: %workingDirectoryId%"
  if "%ProductName%" == "Docker" (set wxs="Docker\Docker.wxs" && set dataFolder="Docker\Data") else (set wxs="ProductTemplate\Product.wxs" && set dataFolder="ProductTemplate\Data")
  echo "Using wxs from %wxs%"
  copy "%wxs%" "%PackageBuildFolder%\%ProductName%.wxs" > NUL
  Replace.exe "%PackageBuildFolder%\%ProductName%.wxs" "--MANUFACTURER--" "%Manufacturer%"
  Replace.exe "%PackageBuildFolder%\%ProductName%.wxs" "--PRODUCTNAME--" "%ProductName%"
  Replace.exe "%PackageBuildFolder%\%ProductName%.wxs" "--PRODUCTVERSION--" "%ProductVersion%"
  Replace.exe "%PackageBuildFolder%\%ProductName%.wxs" "--ASSEMBLYVERSION--" "%AssemblyVersion%"
  Replace.exe "%PackageBuildFolder%\%ProductName%.wxs" "--HUMANVERSION--" %HumanVersion%
  Replace.exe "%PackageBuildFolder%\%ProductName%.wxs" "--PRODUCTDESCRIPTION--" "%ProductDescription%"
  Replace.exe "%PackageBuildFolder%\%ProductName%.wxs" "--SUPPORTURL--" "%SupportURL%"
  Replace.exe "%PackageBuildFolder%\%ProductName%.wxs" "--WORKINGDIRECTORYID--" "%workingDirectoryId%"
  Replace.exe "%PackageBuildFolder%\%ProductName%.wxs" "--WORKINGDIRECTORY--" "%PackageBuildFolder%"
  Replace.exe "%PackageBuildFolder%\%ProductName%.wxs" "--LAUNCHEXECUTABLEID--" "%dockerExeId%"
  Replace.exe "%PackageBuildFolder%\%ProductName%.wxs" "--INSTALLERCMDID--" "%dockerInstallerExeId%"
  for /f "delims=" %%a in ('git rev-parse HEAD') do @set buildSHA1=%%a
  Replace.exe "%PackageBuildFolder%\%ProductName%.wxs" "--BUILDSHA1--" "%buildSHA1%"

  md "%PackageBuildFolder%\Include" 2> NUL
  md "%PackageBuildFolder%\Data" 2> NUL
  xcopy /y /s Include "%PackageBuildFolder%\Include" > NUL
  xcopy /y /s "%dataFolder%" "%PackageBuildFolder%\Data" > NUL
  copy "..\..\resources\Windows\MobyLinux.ps1" "Dependencies\HyperVInstaller.ps1" > NUL
  Replace.exe -U "%PackageBuildFolder%\Docker\resources\License.rtf" "--HUMANVERSION--" %HumanVersion%
  candle.exe -nologo "%PackageBuildFolder%\%ProductName%.wxs" -out "%PackageBuildFolder%\%ProductName%.wixobj"  -ext WixUtilExtension  -ext WixUIExtension
  if ERRORLEVEL 0 goto light else goto error

:light
  light.exe -nologo "%PackageBuildFolder%\%ProductName%.wixobj" -out "%PackageBuildFolder%\Install%ProductName%.msi"  -ext WixUtilExtension  -ext WixUIExtension -dWixUILicenseRtf=./Docker/resources/License.rtf -dWixUIBannerBmp=./Data/WixUIBannerBmp.png -dWixUIDialogBmp=./Data/WixUIDialogBmp.png
  if ERRORLEVEL 0 goto clean else goto error

:usage
  @echo "Usage: %0 Manufacturer ProductName ProductVersion ProductId ProductDescription SupportURL ProductBuildFolder PackageBuildFolder"
  goto end

:error
  @echo "Something Failed!"
  @rmdir /q /s %PackageBuildFolder%\%ProductName%
  @rmdir /q /s %PackageBuildFolder%\Data
  @rmdir /q /s %PackageBuildFolder%\Include
  @del /q /f %PackageBuildFolder%\%ProductName%.wixobj
  @del /q /f %PackageBuildFolder%\Install%ProductName%.wixpdb
  @del /q /f %PackageBuildFolder%\%ProductName%.wxs
  @del /q /f %PackageBuildFolder%\FilesList.wxi
  exit 1
  goto end

:clean
  @echo "Build succeeded!"
  goto end

:end
@rmdir /q /s %PackageBuildFolder%\%ProductName%
@rmdir /q /s %PackageBuildFolder%\Data
@rmdir /q /s %PackageBuildFolder%\Include
@del /q /f %PackageBuildFolder%\%ProductName%.wixobj
@del /q /f %PackageBuildFolder%\Install%ProductName%.wixpdb
@del /q /f %PackageBuildFolder%\%ProductName%.wxs
@del /q /f %PackageBuildFolder%\FilesList.wxi
