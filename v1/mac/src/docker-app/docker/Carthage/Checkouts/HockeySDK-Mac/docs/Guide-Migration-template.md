## Introduction

This guide will help you migrate from HockeySDK-Mac version 1.x to the latest release of the unified HockeySDK for Mac.

First of all we will cleanup the obsolete installation files and then convert your existing code to the new API calls.

## Cleanup

First of all you should remove all files from prior versions of either HockeySDK-Mac. If you not sure which files you added, here are a few easy steps.

### HockeySDK-Mac v1.x

In Xcode open the `Project Navigator` (⌘+1). In the search field at the bottom enter `HockeySDK.framework`. If search returns any results you have the first release of our unified SDK added to your project. Even if you added it as a git submodule we would suggest you remove it first.

### Final Steps

Search again in the `Project Navigator` (⌘+1) for "CrashReporter.framework". You shouldn't get any results now. If not, remove the CrashReporter.framework from your project.

## Installation

Follow the steps in our installation guide [Installation & Setup](Guide-Installation-Setup).

After you finished the steps for either of the installation procedures, we have to migrate your existing code.

## Setup

### HockeySDK-Mac 1.x

There might be minor to the SDK setup code required. Some delegates methods are deprecated and have to be replaced.

- The protocol `BITCrashReportManagerDelegate` has been replaced by `BITCrashManagerDelegate`.
- A new protocol `BITHockeyManagerDelegate` which also implements `BITCrashManagerDelegate` has been introduced and should be used in the appDelegate
- The class `BITCrashReportManager` has been replaced by `BITCrashManager` and is no singleton any longer
- The properties `userName` and `userEmail` of `BITCrashReportManager` are now delegates of `BITHockeyManagerDelegate`
  - `- (NSString *)userNameForHockeyManager:(BITHockeyManager *)hockeyManager componentManager:(BITHockeyBaseManager *)componentManager;`
  - `- (NSString *)userEmailForHockeyManager:(BITHockeyManager *)hockeyManager componentManager:(BITHockeyBaseManager *)componentManager;`
- The required delegate `crashReportApplicationLog` is replaced by `-(NSString *)applicationLogForCrashManager:(id)crashManager`
- The property `loggingEnabled` in `BITHockeyManager` has been replaced by the property `debugLogEnabled`

### HockeySDK-Mac 2.x

- The call `[BITHockeyManager configureWithIdentifier:companyName:delegate:]` has been deprecated. Use either `[BITHockeyManager configureWithIdentifier:delegate:]` or `[BITHockeyManager configureWithIdentifier:]`

- The delegate `[BITCrashManager showMainApplicationWindowForCrashManager:]` has been deprecated.

- The property `BITCrashManager.enableMachExceptionHandler` is now deprecated since Mach Exception Handler is no enabled by default. Use `BITCrashManager.disableMachExceptionHandler` to disable it.

- The crash report window is not presented modal any longer! If you are presenting a window and give it focus, this might hide the crash report UI.

### Troubleshooting

Error message:

    dyld: Library not loaded: @rpath/HockeySDK.framework/Versions/A/HockeySDK
      Referenced from: /Users/USER/Library/Developer/Xcode/DerivedData/HockeyMac/Build/Products/Debug/APPNAME.app/Contents/MacOS/APPNAME
      Reason: image not found
  
Solution: Add the following entry to your `Runpath Search Paths` in the targets build settings

    @loader_path/../Frameworks