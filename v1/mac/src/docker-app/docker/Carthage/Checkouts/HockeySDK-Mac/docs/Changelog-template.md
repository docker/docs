## Version 4.0.3

- [BUGFIX] Use a directory path for metrics data that's in compliance with Apple's file system guidelines

## Version 4.0.2

- [BUGFIX] Restore compatibility with OS X 10.7 by not accessing `NSURLIsExcludedFromBackupKey` if not available

## Version 4.0.1

- [BUGFIX] Fixes an issue where the whole app's Application Support directory (sandboxed apps) / user’s Application Support directory (non-sandboxed apps) was accidentally excluded from backups.
This SDK release explicitly includes the Application Support directory into backups. If you want to opt-out of this fix and keep the Application Directory's backup flag untouched, add the following line above the SDK setup code:

  	**Objective-C:**
   ```objectivec
   [[NSUserDefaults standardUserDefaults] setBool:YES forKey:@"kBITExcludeApplicationSupportFromBackup"];
   ```
    
	**Swift:**
   ```swift
   NSUserDefaults.standardUserDefaults().setBool(true, forKey: "kBITExcludeApplicationSupportFromBackup")
    ```

## Version 4.0.0

- [IMPROVEMENT] Prefix GZIP category on NSData to prevent symbol collisions

## Version 4.0.0-beta.1

- [NEW] User Metrics including users and sessions data is now in public beta

## Version 4.0.0-alpha.2

- [UPDATE] Add improvements and fixes from 3.2.1

## Version 4.0.0-alpha.1

- [NEW] Add User Metrics support

## Version 3.2.1
- [UPDATE] Some minor refactorings
- [BUGFIX] Fix NSURLSession memory leak in Swift apps
- [BUGFIX] Fix issue preventing attachment from being included when sending non-clean termination report
- [IMPROVEMENT] Anonymize binary path in crash report
- [IMPROVEMENT] Support escaping of additional characters (URL encoding)
- [IMPROVEMENT] Support Bundle Identifiers which contain whitespaces

## Version 3.2.0

- [NEW] Added module definition
- [UPDATE] Added full support for `NSURLSession`
- [UPDATE] Switched to use `@rpath` instead of `@loader_path` (You might need to adjust your `LD_RUNPATH_SEARCH_PATHS` setting to include `@executable_path/../Frameworks`!)
- [UPDATE] Updated PLCrashReporter build
- [BUGFIX] Added missing public header files
- [BUGFIX] Doesn't install the C++ exception when running on 10.7.x
- [BUGFIX] Various other improvements and fixes

## Version 3.1.0

- [NEW] `BITCrashManager`: Added support for unhandled C++ exceptions
- [NEW] `BITCrashManager`: Added process ID to `BITCrashDetails`
- [NEW] `BITCrashManager`: Added `CFBundleShortVersionString` value to crash reports
- [UPDATE] Restructured installation documentation
- [BUGFIX] `BITCrashManager`: Fixed UI not showing centered
- [BUGFIX] `BITCrashManager`: Fixed offline issue showing crash alert over and over again with unsent crash reports
- [BUGFIX] Fixed various compiler warnings & other improvements


## Version 3.0

- [NEW] Requires OS X 10.7 or newer
- [NEW] Converted source code to ARC
- [NEW] Added `BITHockeyAttachment` for more customizable attachments to feedback and crash reports (`content-type`, `filename`)
- [UPDATE] Property `delegate` in all components is now private. Set the delegate on `BITHockeyManager` only!
- [NEW] Added `[BITHockeyManager testIdentifier]` to check if the SDK reaches the server. The result is shown on the HockeyApp website on success.
- [NEW] `BITFeedbackManager`: Updated user interface
- [NEW] `BITFeedbackManager`: Added support for attachments, including preview
- [NEW] `BITCrashManager`: Crash Report UI is not presented modal any longer!
- [NEW] `BITCrashManager`: Added option to define a custom UI flow before sending a crash report, see `setCrashReportUIHandler:`
- [NEW] `BITCrashManager`: Provide details on a crash report (see `lastSessionCrashDetails`)
- [NEW] `BITCrashManager`: Added support for adding a binary attachment to crash reports
- [NEW] `BITCrashManager`: Added the option to define callbacks that will be executed prior to program termination after a crash has occurred. Callback code has to be async-safe! See `setCrashCallbacks`.
- [NEW] `BITCrashManager`: Added `generateTestCrash` method to more quickly test the crash reporting
- [UPDATE] `BITCrashManager`: Updated PLCrashReporter to version 1.2
- [UPDATE] `BITCrashManager`: Mach Exception handler is now enabled by default
- [UPDATE] `BITCrashManager`: Crash reports are now send individually if there are multiple pending
- [BUGFIX] Various bugfixes
<br /><br/>

## Version 2.1

- [NEW] Added Feedback component
- [NEW] Added setter for global `userID`, `userName`, `userEmail`. Can be used instead of the delegates.
- [NEW] Requires 10.6 or newer

## Version 2.0.2

- [BUGFIX] Fix a possible crash before sending the crash report when the selector could not be found

## Version 2.0.1

- [NEW] Crash reports now provide the selector name e.g. for crashes in `objc_msgSend`
- [BUGFIX] Fixed a bug in french localization files that could cause the crash report UI to crash
- [BUGFIX] Enabled `Skip Install` for the Framework target. This fixes a warning when archiving a project and including the SDK as a subproject

## Version 2.0

- General

  - [NEW] Major refactoring of all classes. Please go through the setup or migration guide!
  - [NEW] Added docset with SDK documentation
  - [UPDATE] Updated installation instructions with notes about sandboxing, automatic sending of crash reports and code signing
  - [UPDATE] Improved documentation
  <br /><br/>

- Crash Reporting

  - [NEW] Integrated PLCrashReporter 1.2 beta 3
  - [NEW] Added optional support for Mach exceptions
  - [NEW] Check if additional uncaught exception handlers are installed and print a warning in the console
  - [UPDATE] Replaced `exceptionInterception` property in favor of new `BITCrashReportExceptionApplication` class to catch more exceptions. Check the README file for more details on how to use this.
  - [UPDATE] PLCrashReporter built with `BIT` namespace to avoid collisions
  - [UPDATE] Crash reporting is automatically disabled when the app is invoked with the debugger!
  - [UPDATE] Made all delegates fully optional to simplify setup
  - [UPDATE] Enabled Copy&Paste in crash report UI
  - [UPDATE] Crash report text in the dialog's detail view can now be selected
  - [UPDATE] Adjusted privacy note in the dialog
  <br /><br/>
      
- Beta Updates

  - Added used language to the Sparkle request keys in the BITSystemProfile helper class
  <br /><br/>

## Version 2.0 RC 1

- General

  - [UPDATE] Simplified setup by integrating PLCrashReporter as a static library
  - [UPDATE] Improved documentation
	<br /><br/>

- Crash Reporting

  - [NEW] Integrated PLCrashReporter 1.2 beta 3
  - [UPDATE] Made all delegates fully optional to simplify setup
  - [FIX] Fixed a possible crash when detecting a crash report
  - [FIX] Fixed memory leaks reported by clang on Xcode 5
  - [FIX] Fixed username, email and userid not being sent to the server
  - [FIX] Load previously entered username and email and show them in the dialog
	<br /><br/>
  
## Version 2.0 Beta 1

- General

  - [NEW] Major refactoring of all classes. Please go through the setup or migration guide!
  - [NEW] Added docset with SDK documentation
  - [UPDATE] Improved documentation
	<br /><br/>

- Crash Reporting

  - [NEW] Integrated PLCrashReporter 1.2 beta 1
  - [NEW] Added optional support for Mach exceptions
  - [UPDATE] Replace `exceptionInterception` property in favor of new `BITCrashReportExceptionApplication` class to catch more exceptions. Check the README file for more details on how to use this.
  - [UPDATE] PLCrashReporter build with `BIT` namespace to avoid collisions
  - [UPDATE] Crash reporting is automatically disabled when the app is invoked with the debugger!
	<br /><br/>

## Version 1.1.0 Beta 3

- Crash Reporting

  - [NEW] Check if additional uncaught exception handlers are installed and print a warning in the console
  - [FIX] Fix an exception when the process path is empty
	<br /><br/>

## Version 1.1.0 Beta 2

- Crash Reporting

  - [UPDATE] Enable Copy&Paste in crash report UI
  - [FIX] Fix sending of crash reports in Beta 1 not working
  - [FIX] Change XML format so improper formatted log data can not result in empty crash reports on the server
	<br /><br/>


## Version 1.1.0 Beta 1

- General

  - [UPDATE] Updated installation instructions with notes about sandboxing, automatic sending of crash reports and code signing
  - [FIX] Fix build warnings
	<br /><br/>

- Crash Reporting

  - [NEW] Added beta version of much improved PLCrashReporter 1.2 Beta 1
  - [UPDATE] Crash report text in the dialogs detail view can now be selected
  - [UPDATE] Adjusted privacy note in the dialog
  - [FIX] Fixed converting long executable names in crash reports that broke symbolication
	<br /><br/>
    
- Beta Updates

  - Add used language to the Sparkle request keys in the BITSystemProfile helper class
	<br /><br/>


---

## Version 1.0.3

- Added BITSystemProfile class to send analytics data for beta apps when using Sparkle
- Fixed a few compiler warnings
- Fixed crashes when initializing hockey manager in a different autorelease pool than starting the manager
- Fixed a problem when sending crash reports automatically without user interaction
- Fixed showMainApplicationWindow delegate being invoked multiple times in rare cases
- Improvements to installation and setup instructions

## Version 1.0.2

- Include new PLCrashReporter version, which fixes a crash that can happen when the App/System is unloading images from a process

## Version 1.0.1

- Fixed a App Store rejection cause (only happened if you don't submit with sandbox enabled!): settings data was written into ~/Library/net.hockeyapp.sdk.mac/, and is now written into ~/Library/Caches/<app bundle identifier>/net.hockeyapp.sdk.mac/ next to the queued crash reports
- Fixed an issue writing the queued crashes into the wrong key in the settings file
- Fixed reading the wrong meta data (application log data) for queued crash reports
- Delete crash reports if they can not be processed (only might happen if there is an unknown PLCrashReporter issue)
- Send unique UUID for the crash report to the server (so HockeyApp can identify duplicates in a future version)
- Initialize PLCrashReporter as early as possible instead of waiting until the `startManager` call
- Added new property `didCrashInLastSession`
- Minor code cleanup

## Version 1.0

- Update URL to send crash reports to https://sdk.hockeyapp.net/

## Version 0.9.6 RC 6

- IMPORTANT: Initialization methods and class names changed! Please check the *Setup* section in the *README*. Sorry for that.
- Optimize sending of crash reports. Crash reports will be send synchronously if the app crashes within a customizable time interval (default 5s)
- Added option to ask the user for name and email in the UI
- Removed the delegates to get userid and contact, replaced with username and useremail properties
- Validate app identifier and disable the SDK if it is obviously invalid
- Use proper sandbox safe directories for crash caches and sdk settings
- Fixed symlink error of the framework
- Adjust namespace from CNS (Codenauts) to BIT (Bit Stadium).
- Updated bundle identifiers
- Update copyright information to use Bit Stadium GmbH
- Some internal optimizations

## Version 0.9.5 RC 5

- Add multiple localizations (Finnish, French, Italian, Norwegian, Swedish. Thanks Markus!)
- UI now automatically resizes the buttons to fit the localized strings
- Move CNSCrashReporterManagerDelegate to public headers of the framework

## Version 0.9.4 RC 4

- Update SDK initializer to be less error prone
- Add german localization (Beware, button sizes are not automatically adjusted if you add other languages!)

## Version 0.9.3 RC 3

- Fixed double PLCrashReporter in HockeySDK-Mac framework

## Version 0.9.2 RC 2

- Cleaned up protocols, initialization slightly changed, please check the readme file!
- Fixed company name not appearing in the user interface
- Moved PLCrashReporter framework into the HockeySDK-Mac frameworks folder

## Version 0.9 RC 1

- Fixed memory leak
- Added option to intercept exceptions thrown within the main NSRunLoop before they reach Apple's exception handler
- Send crash report synchronously, so crashes that appear on startup are also safely submitted
- Make sure crash reports are anonymous and don't contain user's home directory name
- Send app binaries UUIDs to the server for server side symbolication improvements

# Version 0.6

- System calls in Last Exception Backtrace are now symbolicated
- Fixed invalid stacktrace and some cases
- Fixed 0x0 appearances in stack traces
- Fixed memory leak

# Version 0.5.1

- Added Mac Sandbox support:
  - Supports 32 and 64 bit Intel X86 architecture
  - Uses brand new PLCrashReporter version instead of crash logs from Libary directories
- Fixed sending crash reports to the HockeyApp servers
