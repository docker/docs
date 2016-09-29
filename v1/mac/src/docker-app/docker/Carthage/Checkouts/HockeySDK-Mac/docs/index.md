## Introduction

HockeySDK-Mac implements support for using HockeyApp in your Mac applications.

The following feature is currently supported:

1. **Collect crash reports:** If you app crashes, a crash log with the same format as from the Apple Crash Reporter is written to the device's storage. If the user starts the app again, he is asked to submit the crash report to HockeyApp. This works for both beta and live apps, i.e. those submitted to the App Store!

2. **Feedback:** Collect feedback from your users from within your app and communicate directly with them using the HockeyApp backend.

3. **Add analytics to Sparkle:** If you are using Sparkle to provide app-updates (HockeyApp also supports Sparkle feeds for beta distribution) the SDK contains helpers to add some analytics data to each Sparkle request. 


The main SDK class is `BITHockeyManager`. It initializes all modules and provides access to them, so they can be further adjusted if required. Additionally all modules provide their own protocols.

## Prerequisites

1. Before you integrate HockeySDK into your own app, you should add the app to HockeyApp if you haven't already. Read [this how-to](http://support.hockeyapp.net/kb/how-tos/how-to-create-a-new-app) on how to do it.
2. We also assume that you already have a project in Xcode and that this project is opened in Xcode 6.
3. The SDK supports Mac OS X 10.7 or newer.

## Release Notes

- [Changelog](Changelog)

## Guides

- [Installation & Setup](Guide-Installation-Setup)
- [Migration from previous SDK Versions](Guide-Migration)
- [Mac Desktop Uploader](http://support.hockeyapp.net/kb/services-webhooks-desktop-apps/how-to-upload-to-hockeyapp-on-a-mac)

## HowTos

- [How to do app versioning](HowTo-App-Versioning)
- [How to upload symbols for crash reporting](HowTo-Upload-Symbols)
- [How to add application specific log data](http://support.hockeyapp.net/kb/client-integration-ios-mac-os-x/how-to-add-application-specific-log-data-on-ios-or-osx)

## Troubleshooting

- [Symbolication doesn't work](http://support.hockeyapp.net/kb/client-integration-ios-mac-os-x/how-to-solve-symbolication-problems) (Or the rules of binary UUIDs and dSYMs)
- [Crash Reporting is not working](Troubleshooting-Crash-Reporting-Not-Working)

## Xcode Documentation

This documentation provides integrated help in Xcode for all public APIs and a set of additional tutorials and HowTos.

1. Download the [HockeySDK-Mac documentation](http://hockeyapp.net/releases/).

2. Unzip the file. A new folder `HockeySDK-Mac-documentation` is created.

3. Copy the content into ~`/Library/Developer/Shared/Documentation/DocSet`
