#
# Copyright (C) 2016 Magnus Skjegstad <magnus.skjegstad@docker.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#

""" Main for hockeyapp pinata CLI """

import argparse


def main():
    """ Main for hockeyapp pinata CLI """

    parser = argparse.ArgumentParser()
    parser.description = \
        "List or download Pinata builds from Hockeyapp"

    parser.add_argument("-t", "--token", required=True,
                        help="HockeyApp token - required for access")

    parser.add_argument("--list-apps", action='store_true',
                        help="List available app IDs")

    parser.add_argument("--list-downloads", type=str,
                        metavar="APP_ID",
                        help="List available downloads for one app ID")

    parser.add_argument("--download", type=str, nargs=3,
                        metavar=("APP_ID", "VERSION", "FILENAME"),
                        help="Download zip'ed app via the specified URL.")

    parser.add_argument("--download-latest", type=str, nargs=2,
                        metavar=("APP_ID", "FILENAME"),
                        help="Download latest version " +
                             "of the specified app_id.")

    parser.add_argument("--download-build", type=str, nargs=2,
                        metavar=("BUILD", "FILENAME"),
                        help="Download a specific Pinata build by checking " +
                             "the 'shortversion' field. The tool will try " +
                             "to find this build on the first page of the " +
                             "download results for all apps the token has " +
                             "access to.")

    parser.add_argument("--match-unreleased", action='store_true',
                        default=False,
                        help="Include unreleased versions when searching " +
                             "for latest downloads and specific builds.")

    args = parser.parse_args()

    from hockeyapp.api import HockeyApp
    ha_api = HockeyApp(args.token)

    if args.list_apps:
        apps = ha_api.list_apps()
        for app in apps['apps']:
            print("%s,%s,%s,%s,%s,%s" % (app['public_identifier'],
                                         app['custom_release_type'],
                                         app['title'],
                                         app['platform'],
                                         app['updated_at'],
                                         app['status']))
        return 0

    if args.list_downloads is not None:
        vers = ha_api.list_versions(args.list_downloads)
        for ver in vers['app_versions']:
            if ver['status'] == 2 or args.match_unreleased:
                if ver['status'] == 2:
                    status = "released"
                    url = ver['download_url']
                else:
                    status = "unreleased"
                    # no url for unreleased versions
                    url = ""
                print("%s,%s,%s,%s,%s,%s,%s" % (ver['id'],
                                                ver['shortversion'],
                                                ver['version'],
                                                url,
                                                ver['updated_at'],
                                                ver['timestamp'],
                                                status))
        return 0

    if args.download is not None:
        appid, version, filename = args.download
        ha_api.download_zip(appid, version, filename)
        return 0

    if args.download_latest is not None:
        appid, filename = args.download_latest
        # get highest id
        vers = ha_api.list_versions(appid)
        maxid = 0
        txt = "unknown"
        for ver in vers['app_versions']:
            if ver['id'] > maxid and \
              (ver['status'] == 2 or args.match_unreleased):
                maxid = ver['id']
                if ver['status'] == 2:
                    status = ""
                else:
                    status = " (unreleased)"
                    print("Unreleased binaries can't be downloaded " +
                          "with this version - aborting...")
                    return 1
                txt = "%s, build %s%s" % (ver['shortversion'], ver['version'],
                                          status)
        if maxid != 0:
            print("Downloading version %s (%s)" % (maxid, txt))
            ha_api.download_zip(appid, maxid, filename)
        else:
            print("No downloads found.")

    if args.download_build is not None:
        version_str, filename = args.download_build
        print("Looking for build %s" % version_str)
        if args.match_unreleased:
            print("(Unreleased versions included in result)")
        else:
            print("(Unreleased versions NOT included in result)")
        res = ha_api.find_version_str(version_str,
                                      match_unreleased=args.match_unreleased)
        if res is None:
            print("Build not found.")
            return 1

        app_id, dlid, info = res
        if info['status'] == 2:
            status = "released"
        else:
            status = "unreleased"
            # print("WARNING: Downloading unreleased binary!")
            print("Unreleased binaries can't be downloaded with this version" +
                  " - aborting...")
            return 1

        txt = "%s build %s (status=%s)" % (info['shortversion'],
                                           info['version'],
                                           status)
        print("Downloading version %s (%s)" % (dlid, txt))
        ha_api.download_zip(app_id, dlid, filename)

    return 0
