#!/bin/sh
# SUMMARY: bench volumes mounting performance on windows
# LABELS: win,benchmarks
# AUTHOR: simon.ferquel@docker.com

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    rm -Rf /c/users/public/benchfiles
}
trap clean_up EXIT
remove_temp_files_local() {
    rm -Rf /c/users/public/benchfiles/*
}
remove_temp_files_moby() {
    docker run --rm --privileged --pid=host d4w/nsenter /bin/sh -c "rm -Rf /benchfiles/*"
}
expandRelativeToWindowsPath() {
    fullpath="$(pwd)/$1"
    cygpath -w "$fullpath"
}
expandFullToWindowsPath() {
    fullpath="$1"
    cygpath -w "$fullpath"
}


powershell.exe -NoLogo -WindowStyle Hidden -NoProfile -NonInteractive -File "00.BuildPerfTest.ps1"
SETUPENV=$(expandRelativeToWindowsPath 01.setuptestenv.ps1)
RESOURCES=$(expandFullToWindowsPath "${RT_ROOT}/../win/src/resources")
BENCHRESULTS="${RT_RESULTS}/${RT_TEST_NAME}.benchresults.txt"
"${RT_UTILS}/rt-elevate.exe" -wait powershell.exe -File "${SETUPENV}" -networkKind LocalMoby -username "${D4X_USERNAME}" -password "${D4X_PASSWORD}" -resources "${RESOURCES}"
echo "runs LocalMoby large file" > "${BENCHRESULTS}"
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 1 1 8192 32768 /benchfiles &>> "${BENCHRESULTS}"
remove_temp_files_moby
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 1 1 8192 32768 /benchfiles &>> "${BENCHRESULTS}"
remove_temp_files_moby
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 1 1 8192 32768 /benchfiles &>> "${BENCHRESULTS}" 
remove_temp_files_moby
echo "runs LocalMoby many small files" >> "${BENCHRESULTS}"
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 100 10 4096 10 /benchfiles &>> "${BENCHRESULTS}"
remove_temp_files_moby
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 100 10 4096 10 /benchfiles &>> "${BENCHRESULTS}"
remove_temp_files_moby
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 100 10 4096 10 /benchfiles &>> "${BENCHRESULTS}"
"${RT_UTILS}/rt-elevate.exe" -wait powershell.exe -File "${SETUPENV}" -networkKind HyperVSwitch -username "${D4X_USERNAME}" -password "${D4X_PASSWORD}" -resources "${RESOURCES}"
echo "runs HyperVSwitch large file" >> "${BENCHRESULTS}"
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 1 1 8192 32768 /benchfiles &>> "${BENCHRESULTS}"
remove_temp_files_local
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 1 1 8192 32768 /benchfiles &>> "${BENCHRESULTS}"
remove_temp_files_local
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 1 1 8192 32768 /benchfiles &>> "${BENCHRESULTS}"
remove_temp_files_local
echo "runs HyperVSwitch many small files" >> "${BENCHRESULTS}"
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 100 10 4096 10 /benchfiles &>> "${BENCHRESULTS}"
remove_temp_files_local
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 100 10 4096 10 /benchfiles &>> "${BENCHRESULTS}"
remove_temp_files_local
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 100 10 4096 10 /benchfiles &>> "${BENCHRESULTS}"
"${RT_UTILS}/rt-elevate.exe" -wait powershell.exe -File "${SETUPENV}" -networkKind VpnKit -username "${D4X_USERNAME}" -password "${D4X_PASSWORD}" -resources "${RESOURCES}"
echo "runs VpnKit large file" >> "${BENCHRESULTS}"
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 1 1 8192 32768 /benchfiles &>> "${BENCHRESULTS}"
remove_temp_files_local
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 1 1 8192 32768 /benchfiles &>> "${BENCHRESULTS}"
remove_temp_files_local
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 1 1 8192 32768 /benchfiles &>> "${BENCHRESULTS}"
remove_temp_files_local
echo "runs VpnKit many small files" >> "${BENCHRESULTS}"
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 100 10 4096 10 /benchfiles &>> "${BENCHRESULTS}"
remove_temp_files_local
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 100 10 4096 10 /benchfiles &>> "${BENCHRESULTS}"
remove_temp_files_local
docker run --rm -v /benchfiles:/benchfiles perftest /PerfTest 100 10 4096 10 /benchfiles &>> "${BENCHRESULTS}"
# force unmounting and cleaning of the benchfiles folder
"${RT_UTILS}/rt-elevate.exe" -wait powershell.exe -File "${SETUPENV}" -networkKind LocalMoby -username "${D4X_USERNAME}" -password "${D4X_PASSWORD}" -resources "${RESOURCES}"