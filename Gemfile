source "https://rubygems.org"

# Update me once in a while: https://github.com/github/pages-gem/releases
# Please ensure, before upgrading, that this version exists as a tag in starefossen/github-pages here:
# https://hub.docker.com/r/starefossen/github-pages/tags/
#
# Fresh install?
#
# Windows:
# Install Ruby 2.3.3 x64 and download the Development Kit for 64-bit:
# https://rubyinstaller.org/downloads/
#
# Run this to install devkit after extracting:
# ruby <path_to_devkit>/dk.rb init
# ruby <path_to_devkit>/dk.rb install
#
# then:
# gem install bundler
# bundle install
#
# Mac/Linux:
# Install Ruby 2.3.x and then:
# gem install bundler
# bundle install
#
# ---------------------
# Upgrading? Probably best to reset your environment:
#
# Remove all gems:
# gem uninstall -aIx
#
# (If Windows, do the dk.rb bits above, then go to the next step below)

# Install anew:
# gem install bundler
# bundle install

# This only affects interactive builds (local build, Netlify) and not the
# live site deploy, which uses the Dockerfiles found in the publish-tools
# branch.

gem "github-pages", "198"
gem 'wdm' if Gem.win_platform?
