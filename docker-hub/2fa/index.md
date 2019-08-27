---
description: Enabling two-factor authentication on Docker Hub
keywords: Docker, docker, trusted, registry, security, Dockerfile, Docker Hub, webhooks, docs, documentation
title: Enable two-factor authentication on Docker Hub
---

Two-factor authentication adds an extra layer of security to your Docker Hub
account by requiring a unique security code with each login.

## Prerequisites:
A mobile phone with a time-based on-time password (TOTP) application installed.
Examples of such apps are YubiKey Authenticator and Google Authenticator.


## Enable two-factor authentication
To enable two-factor authentication, log in to your Docker Hub account. Click
on your username and select "Account Settings".

Go to Security and click “Enable Two-Factor Authentication”.

The next page will remind you to download a TOTP application. Click “Set up
using an app”. You will receive your unique recovery code.

> **Save your recovery code and store it somewhere safe.**
> Your recovery code can be used to recover your account in the event you lose
> access to your TOTP application.
{: .important }

After you have saved your code, click “Next”.

Open your TOTP application. You can choose between scanning the QR code or
entering a text code into your TOTP application.

Your TOTP application will give you a six-digit code to enter in text field.
Click “Enable”.

You have successfully enabled two-factor authentication.
