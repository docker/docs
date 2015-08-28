:: Batch file that installer can run. Simply runs bash script with the same name

cd %~dp0
"C:\Program Files\Git\bin\sh.exe" --login -i migrate.sh
