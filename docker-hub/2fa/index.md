---
description: Enabling two-factor authentication on Docker Hub
keywords: Docker, docker, registry, security, Docker Hub, authentication, two-factor authentication
title: Enable two-factor authentication on Docker Hub
---

Two-factor authentication adds an extra layer of security to your Docker Hub
account by requiring a unique security code with each login.

## About Two-Factor Authentication
Paragraph on two-factor authetication.

When you enable two-factor authentication, you will also be provided a recovery
code. Each recovery code is unique and specific to your account. You can use
this code to recovery your account in case you lose access to your authenticator
app.


## Prerequisites
You need a mobile phone with a time-based on-time password authenticator
application installed. Examples of such apps are YubiKey Authenticator and
Google Authenticator.


## Enable two-factor authentication
To enable two-factor authentication, log in to your Docker Hub account. Click
on your username and select "Account Settings".

![]()

Go to Security and click “Enable Two-Factor Authentication”.

![]()

The next page will remind you to download an authenticator app. Click “Set up
using an app”. You will receive your unique recovery code.

> **Save your recovery code and store it somewhere safe.**
>
> Your recovery code can be used to recover your account in the event you lose
> access to your authenticator app.
{: .important }

![]()

After you have saved your code, click “Next”.

![]()

Open your authenticator app. You can choose between scanning the QR code or
entering a text code into your authenticator app.

![]()

Once you have linked your authenticator app, it will give you a six-digit code
to enter in text field. Click “Enable”.

![]()

You have successfully enabled two-factor authentication.
