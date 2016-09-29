#
# Copyright (C) 2016 Rolf Neugebauer <rolf.neugebauer@docker.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#

"""Main entry point to run the regression test on a (remote) host"""

import abc
import argparse
import base64
import fnmatch
import os
import socket
import stat
import subprocess
import sys
import tarfile
import tempfile
import threading

import rt.httpsvr
import rt.log
import winrm


# Temporary directories on the remote hosts
_WIN_TMP_DIR = "/c/Users/%s/AppData/Local/RT"
_OSX_TMP_DIR = "/tmp/RT"

# globals for ports to use as well as IP address to pass to remotes
_PORT_AUX_HTTP = -1
_PORT_AUX_LOG = -1
_LOCAL_IP = None


def _targz_dir(dst, directory, excludes=None):
    """tar up a directory"""
    tf = tarfile.open(dst, mode="w:gz")
    for root, _, files in os.walk(directory):
        for filename in files:
            fp = os.path.join(root, filename)
            if any(fnmatch.fnmatch(fp, pat) for pat in excludes):
                continue
            tf.add(fp)
    tf.close()


class Unbuffered(object):
    """
    This class is used to make sure stdout is not buffered
    """
    def __init__(self, stream):
        self.stream = stream

    def write(self, data):
        self.stream.write(data)
        self.stream.flush()

    def __getattr__(self, attr):
        return getattr(self.stream, attr)


class _Host:
    """An abstract base class to interact with a host"""
    __metaclass__ = abc.ABCMeta

    def __init__(self, remote, user, password):
        """Initialise the class"""
        self.remote = remote
        self.user = user
        self.password = password
        return

    @abc.abstractmethod
    def run_sh(self, cmd):
        """Run a shell command on Host"""
        pass

    def cp_to(self, lpath, rpath, excludes=None):
        """Copy a file from the local host to the remote host. If @lpath ends
        with '/' copy the directory recursively. @rpath *must* be a
        directory. Optionally supply a list of glob expressions for
        files to ignore."""
        print("COPY TO: %s -> %s" % (lpath, rpath))
        if not lpath.endswith('/'):
            self._cp_to(lpath, rpath)
        # we have a directory to copy.
        # 1. Tar up directory to a temporary file
        fd, fp = tempfile.mkstemp(suffix=".tar.gz")
        _targz_dir(fp, lpath, excludes)
        _, fn = os.path.split(fp)

        # 2. Copy file to the host and remove tmp file
        rc = self._cp_to(fp, rpath)
        os.close(fd)
        os.remove(fp)
        if not rc == 0:
            return rc

        # 3. Untar the file on the host
        rc = self.run_sh("(cd %s; tar xzf %s; rm %s)" %
                         (rpath, fn, fn))
        return rc

    @abc.abstractmethod
    def _cp_to(self, lpath, rpath):
        """ Class specfic implementation of a copy to operation"""
        pass

    @abc.abstractmethod
    def cp_from(self, rpath, lpath):
        """Copy a directory from the remote host to the local host. If @rpath
        ends with a '/' copy the directory recursively. @lpath *must*
        be a directory."""
        pass

    def run_rt_local(self, rpath, args):
        """Run tests on remote system. This is special method as it allows
        different backends to modify the arguments a little if
        needed. @rpath is where 'rt-local' can be found"""

        rc = self.run_sh("cd %s; python ./rt-local %s" % (rpath, args))
        return rc


class _LogThread(threading.Thread):
    """Listen on a logging socket and print stuff till close()"""

    def __init__(self, port):
        threading.Thread.__init__(self)
        self.port = port

    def run(self):
        rt.log.log_svr(self.port)


class WinRMHost(_Host):
    """A class to interact with remote host via Windows Remote Management"""
    # This is fun.
    # - WinRM does not allow one to send file to and fro. So we use
    #   WinRM to execute bash commands on host and use a local web
    #   server to handle GET/POST request to transfer file.  When
    #   transferring directories we tar them before transferring.
    # - WinRM also only gives you the output of a command once it is
    #   done.  So, we log over a socket and have a log server running
    #   locally to print the logs as they come in.
    # - pywinrm chokes when a command runs to long, so never returns.
    #   So we start in the background and exit once the log thread exits.

    def __init__(self, remote, user, password):
        """Initialise the class"""
        super(WinRMHost, self).__init__(remote, user, password)

        self.our_ip = _LOCAL_IP if _LOCAL_IP else self.get_local_ip(remote)

        self.session = winrm.Session("%s" % remote,
                                     auth=(self.user, self.password),
                                     server_cert_validation='ignore')

    @staticmethod
    def get_local_ip(remote):
        """Try to get local IP address used to connect to the remote host"""
        try:
            s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
            s.connect((remote, 53))
            our_ip = s.getsockname()[0]
            s.close()
        except:
            our_ip = socket.gethostbyname(socket.gethostname())
        return our_ip

    def run_sh(self, cmd):
        """Run a shell command on the host"""

        # This is really, really tedious, but I could not get cmd.exe
        # to run bash directly, so we start powershell and run it from
        # there. Further, winrm.run_ps does not work either.
        s = self.session
        print("EXEC: %s" % cmd)
        rs = s.run_cmd("powershell.exe",
                       ['-NonInteractive',
                        '-ExecutionPolicy', 'Unrestricted',
                        '-Command',
                        "\"& 'C:\\Program Files\\Git\\bin\\bash.exe' " +
                        "--login -c '%s' \"" % cmd])
        if len(rs.std_out.strip()):
            print("STDOUT:")
            for l in rs.std_out.split('\n'):
                print(l)
        if len(rs.std_err.strip()):
            print("STDERR:")
            for l in rs.std_err.split('\n'):
                print(l)
        return rs.status_code

    def _cp_to(self, lpath, rpath):
        """Copy a file in @lpath to the @rpath directory on the host."""
        _, fn = os.path.split(lpath)
        # 1. Start web server with file
        s = rt.httpsvr.dl_svr_start(_PORT_AUX_HTTP, lpath)
        # 2. Execute curl (after making sure the dir exist)
        self.run_sh("mkdir -p %s" % rpath)
        rc = self.run_sh("(cd %s; curl -s %s:%d -o %s)" %
                         (rpath, self.our_ip, _PORT_AUX_HTTP, fn))
        # 3. Make sure server is stopped
        rt.httpsvr.dl_svr_wait(s)
        return rc

    def _cp_from(self, rpath, lpath):
        """Copy a file from @rpath to the local @lpath directory."""
        rc = subprocess.call("mkdir -p %s" % lpath, shell=True)
        fp, fn = os.path.split(rpath)
        # 1. Start web server to write a POST to the file
        s = rt.httpsvr.ul_svr_start(_PORT_AUX_HTTP, os.path.join(lpath, fn))
        # 2. Execute curl
        rc = self.run_sh("(cd %s; curl -s -X POST --data-binary @%s %s:%d)" %
                         (fp, fn, self.our_ip, _PORT_AUX_HTTP))
        # 3. Make sure server is stopped
        rt.httpsvr.ul_svr_wait(s)
        return rc

    def cp_from(self, rpath, lpath):
        print("COPY FROM: %s -> %s" % (rpath, lpath))
        if not rpath.endswith('/'):
            self._cp_from(rpath, lpath)

        # local temporary file (only use it for the filename)
        fd, lfp = tempfile.mkstemp(suffix=".tar.gz")
        os.close(fd)
        _, fn = os.path.split(lfp)

        # we have a directory to copy.
        # 1. Tar up directory to a temporary file
        rfn = os.path.join("/c/Users/%s/AppData/Local/" % self.user, fn)
        print(rfn)
        # get parent directory and directory to tar
        rc = self.run_sh("(cd %s; tar czf %s .)" %
                         (rpath, rfn))
        if not rc == 0:
            return rc

        # 2. Copy file to the host and delete the remote file
        rc = self._cp_from(rfn, lpath)
        self.run_sh("rm -f %s" % rfn)
        if not rc == 0:
            return rc

        # 3. Untar and remove local tar file
        rc = subprocess.call("(cd %s; tar xzf %s)" % (lpath, fn), shell=True)
        os.remove(os.path.join(lpath, fn))
        return rc

    def run_rt_local(self, rpath, args):
        # When running tests pipe output to /dev/null. We can only get
        # it after the command was run anyway so it's a bit
        # useless. Instead we run a Logging server in a thread and
        # configure rt-local to log to it.
        bg = False
        if " run" in args:
            thd = _LogThread(_PORT_AUX_LOG)
            thd.start()
            args = "--logger %s:%d " % (self.our_ip, _PORT_AUX_LOG) + args
            bg = True

        # base64 encode the rt-local arguments in case of special chars
        rc = _Host.run_rt_local(self, rpath, "%s %s" %
                                (base64.b64encode(args),
                                 " > /dev/null 2> /dev/null &" if bg else ""))

        if " run" in args:
            print("'rt-local' started in background on remote host")
            thd.join()
        # XXX ToDo: reflect the rt-local return value here
        return rc


class SSHHost(_Host):
    """A class to interact with remote host via SSH"""

    def __init__(self, remote, user, password):
        super(SSHHost, self).__init__(remote, user, password)
        self.connection_string = "%s@%s" % (self.user, self.remote)
        self.key_file = os.path.join(os.getcwd(), "etc", "ssh", "id_rsa")
        self.ignore_host_key = '-oStrictHostKeyChecking=no'
        self.known_hosts_file = '-oUserKnownHostsFile=/dev/null'

        # Check SSH key has correcy permissions
        perms = stat.S_IRUSR & stat.S_IWUSR
        st = os.stat(self.key_file)
        current_perms = stat.S_IMODE(st.st_mode)
        if current_perms != perms:
            os.chmod(self.key_file, 0o600)

    def run_sh(self, cmd):
        """Run a shell command on the host"""
        print("EXEC: %s" % cmd)
        ssh_cmd = [
            "ssh",
            "-q",
            self.ignore_host_key,
            self.known_hosts_file,
            "-i",
            self.key_file,
            self.connection_string,
            '%s' % (cmd)
        ]
        # Use Popen to avoid jamming the stdout buffer
        process = subprocess.Popen(ssh_cmd,
                                   stdin=None,
                                   stdout=subprocess.PIPE,
                                   stderr=subprocess.STDOUT)
        for line in iter(process.stdout.readline, ''):
            sys.stdout.write(line)

        process.communicate()
        return process.returncode

    def cp_from(self, rpath, lpath):
        print("COPY FROM: %s -> %s" % (rpath, lpath))
        if not rpath.endswith('/'):
            rc = self.scp_from(False, rpath, lpath)
        else:
            rc = self.scp_from(True, rpath, lpath)
        return rc

    def _cp_to(self, lpath, rpath):
        print("COPY TO: %s -> %s" % (lpath, rpath))
        scp_cmd = [
            "scp",
            "-q",
            self.ignore_host_key,
            "-i",
            self.key_file,
            lpath,
            "%s:%s" % (self.connection_string, rpath)
        ]
        code = subprocess.call(scp_cmd)
        return code

    def scp_from(self, is_dir, rpath, lpath):
        print("COPY FROM: %s -> %s" % (rpath, lpath))
        if is_dir:
            scp = ["scp", "-q", "-r"]
            if not rpath.endswith("*"):
                rpath = rpath + "*"
        else:
            scp = ["scp", "-q"]

        scp_cmd = scp + [
            self.ignore_host_key,
            "-i",
            self.key_file,
            "%s:%s" % (self.connection_string, rpath),
            lpath
        ]
        code = subprocess.call(scp_cmd)
        return code


def main():
    """Main entry point to run regressions tests on a (remote) host"""

    # Make sure stdout is not buffered
    sys.stdout = Unbuffered(sys.stdout)

    parser = argparse.ArgumentParser()
    parser.description = \
        "Execute rt-local commands on a (remote) host"
    parser.epilog = \
        "Arguments after '--' are passed to rt-local. " + \
        "If no '--' is provided only the tests are copied " + \
        "over to the remote system."

    parser.add_argument('-t', '--type', nargs=1, required=True, dest='rtype',
                        choices=['osx', 'win'],
                        help="Type of host")

    parser.add_argument("-r", "--remote", required=True,
                        help="Remote host IP or name")

    parser.add_argument("--portbase", type=int, default=12000,
                        help="Port(s) to use for auxiliary services")

    parser.add_argument("--local-ip", type=str, dest='localip',
                        help="Override local IP clients connect to.")

    parser.add_argument("-u", "--user", default="docker",
                        help="Username to use")

    parser.add_argument("-p", "--password", default="Pass30rd!",
                        help="Password for user")

    # Split argv at -- and parse
    host_args, _, local_args = " ".join(sys.argv).partition(" -- ")
    args = parser.parse_args(host_args.split()[1:])

    global _PORT_AUX_HTTP
    global _PORT_AUX_LOG
    global _LOCAL_IP
    _PORT_AUX_HTTP = args.portbase
    _PORT_AUX_LOG = args.portbase + 1
    _LOCAL_IP = args.localip

    if args.rtype[0] == "osx":
        h = SSHHost(args.remote, args.user, args.password)
        remote_tmp_dir = _OSX_TMP_DIR
    elif args.rtype[0] == "win":
        h = WinRMHost(args.remote, args.user, args.password)
        remote_tmp_dir = _WIN_TMP_DIR % args.user

    # prepare system
    print("\n=== Clean Remote system")
    h.run_sh("rm -rf %s" % remote_tmp_dir)
    h.run_sh("mkdir %s" % remote_tmp_dir)

    print("\n=== Copy infrastructure and test cases")
    # XXX Should really parse local_args for `-c` and copy them
    # separate if specified.
    h.cp_to("./", remote_tmp_dir, ["./_results/*",
                                   "*.pyc",
                                   "./cases/_tmp/*",
                                   "./etc/*"])

    if len(local_args) == 0:
        return

    print("\n=== Run 'rt-local' on remote host")
    rc = h.run_rt_local(remote_tmp_dir, local_args)

    if " run" in local_args:
        print("\n=== Copy Results")
        t = h.cp_from(remote_tmp_dir + "/_results/", "./_results")
        rc = t if rc == 0 else rc

    return rc
