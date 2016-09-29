"""Contains client side logic of WinRM SOAP protocol implementation"""
from __future__ import unicode_literals
import base64
import uuid

import xml.etree.ElementTree as ET
import xmltodict

from six import text_type, binary_type

from winrm.transport import Transport
from winrm.exceptions import WinRMError, WinRMOperationTimeoutError

class Protocol(object):
    """This is the main class that does the SOAP request/response logic. There
    are a few helper classes, but pretty much everything comes through here
    first.
    """
    DEFAULT_READ_TIMEOUT_SEC = 30
    DEFAULT_OPERATION_TIMEOUT_SEC = 20
    DEFAULT_MAX_ENV_SIZE = 153600
    DEFAULT_LOCALE = 'en-US'

    def __init__(
            self, endpoint, transport='plaintext', username=None,
            password=None, realm=None, service=None, keytab=None,
            ca_trust_path=None, cert_pem=None, cert_key_pem=None,
            server_cert_validation='validate',
            kerberos_delegation=False,
            read_timeout_sec=DEFAULT_READ_TIMEOUT_SEC,
            operation_timeout_sec=DEFAULT_OPERATION_TIMEOUT_SEC,
            kerberos_hostname_override=None,
        ):
        """
        @param string endpoint: the WinRM webservice endpoint
        @param string transport: transport type, one of 'plaintext' (default), 'kerberos', 'ssl'  # NOQA
        @param string username: username
        @param string password: password
        @param string realm: unused
        @param string service: the service name, default is HTTP
        @param string keytab: the path to a keytab file if you are using one
        @param string ca_trust_path: Certification Authority trust path
        @param string cert_pem: client authentication certificate file path in PEM format  # NOQA
        @param string cert_key_pem: client authentication certificate key file path in PEM format  # NOQA
        @param string server_cert_validation: whether server certificate should be validated on Python versions that suppport it; one of 'validate' (default), 'ignore' #NOQA
        @param bool kerberos_delegation: if True, TGT is sent to target server to allow multiple hops  # NOQA
        @param int read_timeout_sec: maximum seconds to wait before an HTTP connect/read times out (default 30). This value should be slightly higher than operation_timeout_sec, as the server can block *at least* that long. # NOQA
        @param int operation_timeout_sec: maximum allowed time in seconds for any single wsman HTTP operation (default 20). Note that operation timeouts while receiving output (the only wsman operation that should take any significant time, and where these timeouts are expected) will be silently retried indefinitely. # NOQA
        @param string kerberos_hostname_override: the hostname to use for the kerberos exchange (defaults to the hostname in the endpoint URL)
        """

        if operation_timeout_sec >= read_timeout_sec or operation_timeout_sec < 1:
            raise WinRMError("read_timeout_sec must exceed operation_timeout_sec, and both must be non-zero")

        self.read_timeout_sec = read_timeout_sec
        self.operation_timeout_sec = operation_timeout_sec
        self.max_env_sz = Protocol.DEFAULT_MAX_ENV_SIZE
        self.locale = Protocol.DEFAULT_LOCALE

        self.transport = Transport(
            endpoint=endpoint, username=username, password=password,
            realm=realm, service=service, keytab=keytab,
            ca_trust_path=ca_trust_path, cert_pem=cert_pem,
            cert_key_pem=cert_key_pem, read_timeout_sec=self.read_timeout_sec,
            server_cert_validation=server_cert_validation,
            kerberos_delegation=kerberos_delegation,
            kerberos_hostname_override=kerberos_hostname_override,
            auth_method=transport)

        self.username = username
        self.password = password
        self.service = service
        self.keytab = keytab
        self.ca_trust_path = ca_trust_path
        self.server_cert_validation = server_cert_validation
        self.kerberos_delegation = kerberos_delegation
        self.kerberos_hostname_override = kerberos_hostname_override

    def open_shell(self, i_stream='stdin', o_stream='stdout stderr',
                   working_directory=None, env_vars=None, noprofile=False,
                   codepage=437, lifetime=None, idle_timeout=None):
        """
        Create a Shell on the destination host
        @param string i_stream: Which input stream to open. Leave this alone
         unless you know what you're doing (default: stdin)
        @param string o_stream: Which output stream to open. Leave this alone
         unless you know what you're doing (default: stdout stderr)
        @param string working_directory: the directory to create the shell in
        @param dict env_vars: environment variables to set for the shell. For
         instance: {'PATH': '%PATH%;c:/Program Files (x86)/Git/bin/', 'CYGWIN':
          'nontsec codepage:utf8'}
        @returns The ShellId from the SOAP response. This is our open shell
         instance on the remote machine.
        @rtype string
        """
        req = {'env:Envelope': self._get_soap_header(
            resource_uri='http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd',  # NOQA
            action='http://schemas.xmlsoap.org/ws/2004/09/transfer/Create')}
        header = req['env:Envelope']['env:Header']
        header['w:OptionSet'] = {
            'w:Option': [
                {
                    '@Name': 'WINRS_NOPROFILE',
                    '#text': str(noprofile).upper()  # TODO remove str call
                },
                {
                    '@Name': 'WINRS_CODEPAGE',
                    '#text': str(codepage)  # TODO remove str call
                }
            ]
        }

        shell = req['env:Envelope'].setdefault(
            'env:Body', {}).setdefault('rsp:Shell', {})
        shell['rsp:InputStreams'] = i_stream
        shell['rsp:OutputStreams'] = o_stream

        if working_directory:
            # TODO ensure that rsp:WorkingDirectory should be nested within rsp:Shell  # NOQA
            shell['rsp:WorkingDirectory'] = working_directory
            # TODO check Lifetime param: http://msdn.microsoft.com/en-us/library/cc251546(v=PROT.13).aspx  # NOQA
            #if lifetime:
            #    shell['rsp:Lifetime'] = iso8601_duration.sec_to_dur(lifetime)
        # TODO make it so the input is given in milliseconds and converted to xs:duration  # NOQA
        if idle_timeout:
            shell['rsp:IdleTimeOut'] = idle_timeout
        if env_vars:
            env = shell.setdefault('rsp:Environment', {})
            for key, value in env_vars.items():
                env['rsp:Variable'] = {'@Name': key, '#text': value}

        res = self.send_message(xmltodict.unparse(req))
        #res = xmltodict.parse(res)
        #return res['s:Envelope']['s:Body']['x:ResourceCreated']['a:ReferenceParameters']['w:SelectorSet']['w:Selector']['#text']
        root = ET.fromstring(res)
        return next(
            node for node in root.findall('.//*')
            if node.get('Name') == 'ShellId').text

    # Helper method for building SOAP Header
    def _get_soap_header(
            self, action=None, resource_uri=None, shell_id=None,
            message_id=None):
        if not message_id:
            message_id = uuid.uuid4()
        header = {
            '@xmlns:xsd': 'http://www.w3.org/2001/XMLSchema',
            '@xmlns:xsi': 'http://www.w3.org/2001/XMLSchema-instance',
            '@xmlns:env': 'http://www.w3.org/2003/05/soap-envelope',

            '@xmlns:a': 'http://schemas.xmlsoap.org/ws/2004/08/addressing',
            '@xmlns:b': 'http://schemas.dmtf.org/wbem/wsman/1/cimbinding.xsd',
            '@xmlns:n': 'http://schemas.xmlsoap.org/ws/2004/09/enumeration',
            '@xmlns:x': 'http://schemas.xmlsoap.org/ws/2004/09/transfer',
            '@xmlns:w': 'http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd',
            '@xmlns:p': 'http://schemas.microsoft.com/wbem/wsman/1/wsman.xsd',
            '@xmlns:rsp': 'http://schemas.microsoft.com/wbem/wsman/1/windows/shell',  # NOQA
            '@xmlns:cfg': 'http://schemas.microsoft.com/wbem/wsman/1/config',

            'env:Header': {
                'a:To': 'http://windows-host:5985/wsman',
                'a:ReplyTo': {
                    'a:Address': {
                        '@mustUnderstand': 'true',
                        '#text': 'http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous'  # NOQA
                    }
                },
                'w:MaxEnvelopeSize': {
                    '@mustUnderstand': 'true',
                    '#text': '153600'
                },
                'a:MessageID': 'uuid:{0}'.format(message_id),
                'w:Locale': {
                    '@mustUnderstand': 'false',
                    '@xml:lang': 'en-US'
                },
                'p:DataLocale': {
                    '@mustUnderstand': 'false',
                    '@xml:lang': 'en-US'
                },
                # TODO: research this a bit http://msdn.microsoft.com/en-us/library/cc251561(v=PROT.13).aspx  # NOQA
                # 'cfg:MaxTimeoutms': 600
                # Operation timeout in ISO8601 format, see http://msdn.microsoft.com/en-us/library/ee916629(v=PROT.13).aspx  # NOQA
                'w:OperationTimeout': 'PT{0}S'.format(int(self.operation_timeout_sec)),
                'w:ResourceURI': {
                    '@mustUnderstand': 'true',
                    '#text': resource_uri
                },
                'a:Action': {
                    '@mustUnderstand': 'true',
                    '#text': action
                }
            }
        }
        if shell_id:
            header['env:Header']['w:SelectorSet'] = {
                'w:Selector': {
                    '@Name': 'ShellId',
                    '#text': shell_id
                }
            }
        return header

    def send_message(self, message):
        # TODO add message_id vs relates_to checking
        # TODO port error handling code
        return self.transport.send_message(message)

    def close_shell(self, shell_id):
        """
        Close the shell
        @param string shell_id: The shell id on the remote machine.
         See #open_shell
        @returns This should have more error checking but it just returns true
         for now.
        @rtype bool
        """
        message_id = uuid.uuid4()
        req = {'env:Envelope': self._get_soap_header(
            resource_uri='http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd',  # NOQA
            action='http://schemas.xmlsoap.org/ws/2004/09/transfer/Delete',
            shell_id=shell_id,
            message_id=message_id)}

        # SOAP message requires empty env:Body
        req['env:Envelope'].setdefault('env:Body', {})

        res = self.send_message(xmltodict.unparse(req))
        root = ET.fromstring(res)
        relates_to = next(
            node for node in root.findall('.//*')
            if node.tag.endswith('RelatesTo')).text
        # TODO change assert into user-friendly exception
        assert uuid.UUID(relates_to.replace('uuid:', '')) == message_id

    def run_command(
            self, shell_id, command, arguments=(), console_mode_stdin=True,
            skip_cmd_shell=False):
        """
        Run a command on a machine with an open shell
        @param string shell_id: The shell id on the remote machine.
         See #open_shell
        @param string command: The command to run on the remote machine
        @param iterable of string arguments: An array of arguments for this
         command
        @param bool console_mode_stdin: (default: True)
        @param bool skip_cmd_shell: (default: False)
        @return: The CommandId from the SOAP response.
         This is the ID we need to query in order to get output.
        @rtype string
        """
        req = {'env:Envelope': self._get_soap_header(
            resource_uri='http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd',  # NOQA
            action='http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Command',  # NOQA
            shell_id=shell_id)}
        header = req['env:Envelope']['env:Header']
        header['w:OptionSet'] = {
            'w:Option': [
                {
                    '@Name': 'WINRS_CONSOLEMODE_STDIN',
                    '#text': str(console_mode_stdin).upper()
                },
                {
                    '@Name': 'WINRS_SKIP_CMD_SHELL',
                    '#text': str(skip_cmd_shell).upper()
                }
            ]
        }
        cmd_line = req['env:Envelope'].setdefault(
            'env:Body', {}).setdefault('rsp:CommandLine', {})
        cmd_line['rsp:Command'] = {'#text': command}
        if arguments:
            unicode_args = [a if isinstance(a, text_type) else a.decode('utf-8') for a in arguments]
            cmd_line['rsp:Arguments'] = u' '.join(unicode_args)

        res = self.send_message(xmltodict.unparse(req))
        root = ET.fromstring(res)
        command_id = next(
            node for node in root.findall('.//*')
            if node.tag.endswith('CommandId')).text
        return command_id

    def cleanup_command(self, shell_id, command_id):
        """
        Clean-up after a command. @see #run_command
        @param string shell_id: The shell id on the remote machine.
         See #open_shell
        @param string command_id: The command id on the remote machine.
         See #run_command
        @returns: This should have more error checking but it just returns true
         for now.
        @rtype bool
        """
        message_id = uuid.uuid4()
        req = {'env:Envelope': self._get_soap_header(
            resource_uri='http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd',  # NOQA
            action='http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Signal',  # NOQA
            shell_id=shell_id,
            message_id=message_id)}

        # Signal the Command references to terminate (close stdout/stderr)
        signal = req['env:Envelope'].setdefault(
            'env:Body', {}).setdefault('rsp:Signal', {})
        signal['@CommandId'] = command_id
        signal['rsp:Code'] = 'http://schemas.microsoft.com/wbem/wsman/1/windows/shell/signal/terminate'  # NOQA

        res = self.send_message(xmltodict.unparse(req))
        root = ET.fromstring(res)
        relates_to = next(
            node for node in root.findall('.//*')
            if node.tag.endswith('RelatesTo')).text
        # TODO change assert into user-friendly exception
        assert uuid.UUID(relates_to.replace('uuid:', '')) == message_id

    def get_command_output(self, shell_id, command_id):
        """
        Get the Output of the given shell and command
        @param string shell_id: The shell id on the remote machine.
         See #open_shell
        @param string command_id: The command id on the remote machine.
         See #run_command
        #@return [Hash] Returns a Hash with a key :exitcode and :data.
         Data is an Array of Hashes where the cooresponding key
        #   is either :stdout or :stderr.  The reason it is in an Array so so
         we can get the output in the order it ocurrs on
        #   the console.
        """
        stdout_buffer, stderr_buffer = [], []
        command_done = False
        while not command_done:
            try:
                stdout, stderr, return_code, command_done = \
                    self._raw_get_command_output(shell_id, command_id)
                stdout_buffer.append(stdout)
                stderr_buffer.append(stderr)
            except WinRMOperationTimeoutError as e:
                # this is an expected error when waiting for a long-running process, just silently retry
                pass
        return b''.join(stdout_buffer), b''.join(stderr_buffer), return_code

    def _raw_get_command_output(self, shell_id, command_id):
        req = {'env:Envelope': self._get_soap_header(
            resource_uri='http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd',  # NOQA
            action='http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Receive',  # NOQA
            shell_id=shell_id)}

        stream = req['env:Envelope'].setdefault('env:Body', {}).setdefault(
            'rsp:Receive', {}).setdefault('rsp:DesiredStream', {})
        stream['@CommandId'] = command_id
        stream['#text'] = 'stdout stderr'

        res = self.send_message(xmltodict.unparse(req))
        root = ET.fromstring(res)
        stream_nodes = [
            node for node in root.findall('.//*')
            if node.tag.endswith('Stream')]
        stdout = stderr = b''
        return_code = -1
        for stream_node in stream_nodes:
            if not stream_node.text:
                continue
            if stream_node.attrib['Name'] == 'stdout':
                stdout += base64.b64decode(stream_node.text.encode('ascii'))
            elif stream_node.attrib['Name'] == 'stderr':
                stderr += base64.b64decode(stream_node.text.encode('ascii'))

        # We may need to get additional output if the stream has not finished.
        # The CommandState will change from Running to Done like so:
        # @example
        #   from...
        #   <rsp:CommandState CommandId="..." State="http://schemas.microsoft.com/wbem/wsman/1/windows/shell/CommandState/Running"/>
        #   to...
        #   <rsp:CommandState CommandId="..." State="http://schemas.microsoft.com/wbem/wsman/1/windows/shell/CommandState/Done">
        #     <rsp:ExitCode>0</rsp:ExitCode>
        #   </rsp:CommandState>
        command_done = len([
            node for node in root.findall('.//*')
            if node.get('State', '').endswith('CommandState/Done')]) == 1
        if command_done:
            return_code = int(
                next(node for node in root.findall('.//*')
                     if node.tag.endswith('ExitCode')).text)

        return stdout, stderr, return_code, command_done
