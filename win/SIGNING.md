# SIGNING EXECTUABLE ON WINDOWS

* Get the original certificate, key and key password and use [keybase.io cli](https://keybase.io/download) tools to decrypt them
* You must install [Micosoft Windows SDK](https://dev.windows.com/en-us/downloads/windows-10-sdk)
* Head over to `C:\Program Files (x86)\Windows Kits\10\bin\x86`
* Choose a password for the `pfx`, we will use `XXXXXX` in this doc
* Convert the `pcs` and `sct` to `pfx`, you'll be prompted for the key password with `pvk2pfx.exe /pvk key.pvk /spc cert.spc /pfx toolbox.pfx /po XXXXXXX`
* You can now sign executable with `signtool.exe sign /a /f toolbox.pfx
 /t http://timestamp.verisign.com/scripts/timestamp.dll /p XXXXXX C:\Users\jeanl\go\src\github.com\docker\pin
ata\win\build\win\Docker.exe`

## Installing the certificate on AppVeyor

The pfx file is a binary file, the only way to get 'secured' data on appveyor is through encrypted environement variable. To make this work we must :
* Base64 encode the pfx file. Done on a mac with `openssl base64 -in cert.pfx -out cert.pfx.b64`
* Copy the base64 encoded cert in your laptop clipboard and head over appveyor settings page for your project
* Click on the environement tab, clear and paste on the `PFX` environement variable
* Check the `PFXPASSWORD` for change if needed

## Links to export the Digicert Certificate from their php website

* Use safari
* Generate and install in the local keystore the certificate
* [Add the intermediary authority](
jeanlaurent [10:16]  
https://www.digicert.com/code-signing/mac-verifying-code-signing-certificate.htm#remove_warning)
* [Export as p12 file](
https://www.digicert.com/code-signing/mac-exporting-code-signing-certificate.htm)
* Rename the p12 as pfx (it's the same, really...)
