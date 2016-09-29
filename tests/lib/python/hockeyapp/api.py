#
# Copyright (C) 2016 Magnus Skjegstad <magnus.skjegstad@docker.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
import urllib2 as ul
import json


class HockeyApp(object):
        """ Basic HockeyApp API functions """
        token = ""

        def __init__(self, token):
            self.token = token

        def _read_api_call(self, call):
            url = "https://rink.hockeyapp.net/%s" % call
            # print url
            req = ul.Request(url,
                             headers={"X-HockeyAppToken": self.token})
            resp = ul.urlopen(req)
            return resp.read()

        def list_apps(self):
            """ Return dict with list of apps available with this token """
            data = self._read_api_call("/api/2/apps")
            return json.loads(data)

        def list_versions(self, app_id, page=1):
            """ List one page of downloads available for this app id.
            Defaults to page 1. """
            data = self._read_api_call("/api/2/apps/%s/app_versions?page=%s" %
                                       (app_id, page))
            return json.loads(data)

        def download_zip(self, app_id, version, filename):
            """ Download an app version as a zip-file """
            # NOTE: This doesn't work with unreleased versions... it
            # will download the closest ID which is released...
            with open(filename, "wb") as f:
                url = ("/api/2/apps/%s/app_versions/%s?format=zip" %
                       (app_id, version))
                f.write(self._read_api_call(url))

        def find_version_str(self, version_str, app_id=None,
                             match_unreleased=False):
            """ Find a specific version string and return the HockeyApp
            (app id, version id, info dict). Returns None if the version_str
            is not found, on other errors an exception may be thrown.

            The IDs can be passed to download_zip() to download the release. If
            app_id is None, all apps seen by the token will be checked for the
            version string. This can be used to download specific build
            number from Pinata HockeyApp, as the version_str == CI build"""

            # If app_id is None, get list of app IDs and call ourself
            if app_id is None:
                apps = self.list_apps()
                for app in apps['apps']:
                    appid = app['public_identifier']
                    f = self.find_version_str(version_str,
                                              app_id=appid,
                                              match_unreleased=match_unreleased)
                    if f is not None:
                        return f
            else:
                # app_id is set, search first page of version strings
                versions = self.list_versions(app_id, page=1)
                for version in versions['app_versions']:
                    if (version['status'] == 2) or match_unreleased:
                        if version['version'] == version_str:
                            return (app_id, version['id'], version)

            # not found, return None
            return None
