#!/bin/bash

echo "==> Disabling Screensaver"
# Disable Screensaver
defaults write com.apple.screensaver idleTime 0
killall cfprefsd
