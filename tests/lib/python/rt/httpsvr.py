#
# Copyright (C) 2016 Rolf Neugebauer <rolf.neugebauer@docker.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#

"""For some systems we start a webserver to allow downloading and
uploading of files. Here is the code for it"""

import os
import shutil
import sys
import threading


if sys.version_info < (3, 0):
    import BaseHTTPServer as httpserver
else:
    import http.server as httpserver


# This is a bit weird. This creates a custom RequestHandler class for
# the file we want the client to download.
class _DownloadHandler(httpserver.BaseHTTPRequestHandler):
    def __init__(self, filepath, *args):
        self.filepath = filepath
        httpserver.BaseHTTPRequestHandler.__init__(self, *args)
        return

    def do_GET(self):
        with open(self.filepath, 'rb') as f:
            self.send_response(200)
            self.send_header("Content-Type", 'application/octet-stream')
            self.send_header("Content-Disposition",
                             'attachment; filename="%s"' %
                             os.path.basename(self.filepath))
            fs = os.fstat(f.fileno())
            self.send_header("Content-Length", str(fs.st_size))
            self.end_headers()
            shutil.copyfileobj(f, self.wfile)
            return


def _make_dl_handler(filepath):
    return lambda *args: _DownloadHandler(filepath, *args)


class _DownloadThread(threading.Thread):
    """Create a thread with a webserver, serving one GET request"""
    def __init__(self, port, filepath):
        threading.Thread.__init__(self)
        self.port = port
        self.filepath = filepath

    def run(self):
        sa = ('', self.port)
        handler = _make_dl_handler(self.filepath)
        httpd = httpserver.HTTPServer(sa, handler)
        httpd.handle_request()


def dl_svr_start(port, filepath):
    """Start a thread with a Webserver to download a file from"""
    thd = _DownloadThread(port, filepath)
    thd.start()
    return thd


def dl_svr_wait(thd):
    """Wait for the server thread to exit"""
    thd.join()


class _UploadHandler(httpserver.BaseHTTPRequestHandler):
    def __init__(self, filepath, *args):
        self.filepath = filepath
        httpserver.BaseHTTPRequestHandler.__init__(self, *args)
        return

    def do_POST(self):
        with open(self.filepath, 'wb') as f:
            length = int(self.headers['Content-Length'])
            print("POST: Length=%d" % length)
            data = self.rfile.read(length)
            f.write(data)
        self.send_response(200)


def _make_ul_handler(filepath):
    return lambda *args: _UploadHandler(filepath, *args)


class _UploadThread(threading.Thread):
    """Create a thread with a webserver, serving one GET request"""
    def __init__(self, port, filepath):
        threading.Thread.__init__(self)
        self.port = port
        self.filepath = filepath

    def run(self):
        sa = ('', self.port)
        handler = _make_ul_handler(self.filepath)
        httpd = httpserver.HTTPServer(sa, handler)
        httpd.handle_request()


def ul_svr_start(port, filepath):
    """Start a thread with a Webserver to upload a file to"""
    thd = _UploadThread(port, filepath)
    thd.start()
    return thd


def ul_svr_wait(thd):
    """Wait for the server thread to exit"""
    thd.join()
