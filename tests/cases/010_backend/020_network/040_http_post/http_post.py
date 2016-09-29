#! /usr/bin/env python

import http.server
import sys

class UploadHandler(http.server.BaseHTTPRequestHandler):
    def do_POST(self):
        length = int(self.headers['Content-Length'])
        print(length)
        data = self.rfile.read(length)
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()

sa = ('', int(sys.argv[1]))
httpd = http.server.HTTPServer(sa, UploadHandler)
httpd.serve_forever()
