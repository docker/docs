source "https://rubygems.org"

# Update me once in a while: https://github.com/github/pages-gem/releases
# Please ensure, before upgrading, that this version exists as a tag in starefossen/github-pages here:
# https://hub.docker.com/r/starefossen/github-pages/tags/
#
# Remove all gems in Windows:
# ruby -e "`gem list`.split(/$/).each { |line| puts `gem uninstall -Iax #{line.split(' ')[0]}` unless line.empty? }"
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
# ---------------------
# Everything screwed up? To reset environment:
#
# Remove all gems:
# gem uninstall -aIx
#
# Then:
# gem install bundler
# bundle install

gem "github-pages", "147"
gem 'wdm' if Gem.win_platform?
