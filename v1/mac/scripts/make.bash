#!/usr/bin/env bash

################################################################################
#                       ____             __                                    #
#                      / __ \____  _____/ /_____  _____                        #
#                     / / / / __ \/ ___/ //_/ _ \/ ___/                        #
#                    / /_/ / /_/ / /__/ ,< /  __/ /                            #
#                   /_____/\____/\___/_/|_|\___/_/                             #
#                                                                              #
#                                                                              #
################################################################################
# make created on 2015-10-09 by michaK

# folder where scripts are located
SCRIPT_FOLDER_PATH=$(cd "$(dirname "$0")"; pwd; cd - > /dev/null)

source "$SCRIPT_FOLDER_PATH/Docker.Config"
 if [ $? -ne 0 ]; then exit -1; fi

source "$SCRIPT_FOLDER_PATH/Docker.Utils"
if [ $? -ne 0 ]; then exit -1; fi

source "$SCRIPT_FOLDER_PATH/Docker.Environment"
if [ $? -ne 0 ]; then LogError "can't source Docker.Environment"; fi

source "$SCRIPT_FOLDER_PATH/Docker.Dependencies"
if [ $? -ne 0 ]; then LogError "can't source Docker.Dependencies"; fi

source "$SCRIPT_FOLDER_PATH/Docker.Build"
if [ $? -ne 0 ]; then LogError "can't source Docker.Build"; fi

source "$SCRIPT_FOLDER_PATH/Docker.Upload"
if [ $? -ne 0 ]; then LogError "can't source Docker.Upload"; fi

source "$SCRIPT_FOLDER_PATH/Docker.Sign"
if [ $? -ne 0 ]; then LogError "can't source Docker.Sign"; fi

################################################################################
# Options Parsing
################################################################################

scriptArgs=''
platforms=''
hasOptions=false

verbose=false                   # print debug output
autoYes=false                   # never prompt the user and aswer yes to all questions
clean=false                     # clean previous build
deps=false                      # install dependencies
build=false                     # build product
upload=false                    # upload to HockeyApp

ParseOptions ()
{
  OPTIND=1
  local opt=':'
  opt="${opt}vy"                # verbose, autoYes
  opt="${opt}d"                 # deps
  opt="${opt}bcu"               # edit, build, clean, run, execute, test, version, sign, package, upload
  opt="${opt}h"                 # help

  while getopts $opt option "$@"; do
    case $option in
      v)  verbose=true;                 LogImportant 'Verbose = true';                      scriptArgs="${scriptArgs} -${option}";;
      y)  autoYes=true;                 LogImportant 'Answer Yes to all questions = true';  scriptArgs="${scriptArgs} -${option}";;

      d)  deps=true;                    LogImportant 'Dependencies = true';                 scriptArgs="${scriptArgs} -${option}";;

      c)  clean=true;                   LogImportant 'Clean = true';                        scriptArgs="${scriptArgs} -${option}";;
      b)  build=true;                   LogImportant 'Build = true';                        scriptArgs="${scriptArgs} -${option}";;
      u)  upload=true;                  LogImportant 'Upload = true';                       scriptArgs="${scriptArgs} -${option}";;

      h)  usage;;
      \?)                               LogError "Unknown option -$OPTARG";                 usage;;
      :)                                LogError "Missing parameter for option -$OPTARG";   usage;;
      *)                                LogError "Unimplemented option -$OPTARG";           usage;;
    esac
    hasOptions=true
  done
  shift $((OPTIND-1))
  return $#
}

usageCommon()
{
  echo "Options:"
  echo "-v, --verbose             increase verbosity"
  echo "-y, --yes                 answer yes to all prompted questions"
  echo ""
  echo "Actions:                  Defaults to build"
  echo "-d  --deps                install dependencies"
  echo "-b, --build               build based on platform specifiers"
  echo "-c, --clean               don't build incrementally and clean target first"
  echo "-u, --upload              upload the built product to HockeyApp"
  echo ""
  echo "Meta:"
  echo "-h, --help                prints this screen"
}

ParseArguments ()
{
  OPTIND=1
  ParseOptions $@
  local argsLeft=$?
  shift $((OPTIND-1))

  if [ $deps = false ] && [ $build = false ] && [ $clean = false ] && [ $upload = false ]; then
    usage 'You need to specify at least one action'
  fi
}

buildStartTime=$(date +%s)
nodeName=$(uname -n)
LogImportant "Build started on ${nodeName%.*} running OSX $scriptPlatformVersion ($(uname -sm)) at $(date +"%T on %m-%d-%Y")"
freeSpace=$(df -h . | tr -d '%' | tail -1)
LogImportant "Local Disks Free Space: $freeSpace"
Log "Logging to ${UWhi}$logFilename${RCol}\n"

rm -f ${logFilename}.previous 2> /dev/null
if [ -e ${logFilename} ]; then
  mv ${logFilename} ${logFilename}.previous 2> /dev/null
fi

LogImportant "Current working directory: $PWD"
LogImportant "Script directory: $SCRIPT_FOLDER_PATH"

ParseArguments $@

Execute "mkdir -p $cacheDirectory"

CheckKeychain

if [ $deps = true ]; then
  EnvironmentCheck
  DependenciesCheck
  Initialize
fi

if [ $clean = true ]; then
  Clean
fi

if [ $build = true ]; then
  EnvironmentCheck
  Initialize
  Build
  if [ -n "$CI" ]; then RestoreLoginKeychain; fi
fi

if [ $upload = true ]; then
  EnvironmentCheck
  Initialize
  Upload
fi

buildEndTime=$(date +%s)
timerDiff=$(( $buildEndTime - $buildStartTime ))
hms $timerDiff
LogCompleted "Build completed successfully in ${hms_str}"
echo ""
