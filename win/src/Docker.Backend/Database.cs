using Docker.Core;
using System;
using System.Collections.Generic;
using System.IO;
using System.Security.Cryptography.X509Certificates;
using System.Text.RegularExpressions;
using System.Threading;

namespace Docker.Backend
{
    public interface IDatabase
    {
        void Write(Settings settings);
        void WriteDnsSettings(string dnsSettings);
    }

    public class Database : IDatabase
    {
        private const string CliPath = "com.docker.9pdb.exe";

        private const string KeyNetwork = "com.docker.driver.amd64-linux/network";
        private const string ValueNetworkHybrid = "hybrid";
        private const string KeyDaemonOptions = "com.docker.driver.amd64-linux/etc/docker/daemon.json";
        private const string KeySysCtlConf = "com.docker.driver.amd64-linux/etc/sysctl.conf";
        private const string KeyNetworkConfig = "com.docker.driver.amd64-linux/net/config";
        private const string ValueNetworkConfigStatic = "static";
        private const string KeyNetworkAddress = "com.docker.driver.amd64-linux/net/address";
        private const string KeyNetworkNetmask = "com.docker.driver.amd64-linux/net/netmask";
        private const string KeyNetworkGateway = "com.docker.driver.amd64-linux/net/gateway";
        private const string KeyLocalHostPortForwarding = "com.docker.driver.amd64-linux/native/port-forwarding";
        private const string KeyProxyExclude = "com.docker.driver.amd64-linux/proxy/exclude";
        private const string KeyProxyHttp = "com.docker.driver.amd64-linux/proxy/http";
        private const string KeyProxyHttps = "com.docker.driver.amd64-linux/proxy/https";
        public const string KeyDnsSettings = "com.docker.driver.amd64-linux/slirp/dns";
        private const string KeyCaCertsBundle = "com.docker.driver.amd64-linux/etc/ssl/certs/ca-certificates.crt";

        private readonly object _lock = new object();
        private readonly Logger _logger;
        private readonly ICmd _cmd;
        private readonly IpHelper _ipHelper;
        private readonly WindowsIdentityRecorder _clientIdentityRecorder;

        public Database(ICmd cmd, IpHelper ipHelper, WindowsIdentityRecorder clientIdentityRecorder)
        {
            _logger = new Logger(GetType());
            _cmd = cmd;
            _ipHelper = ipHelper;
            _clientIdentityRecorder = clientIdentityRecorder;
        }

        public void WriteDnsSettings(string dnsSettings)
        {
            if (dnsSettings.Length > 0)
            {
                ModifyDatabase(new List<Tuple<string, string>> { Tuple.Create(KeyDnsSettings, dnsSettings) }, null);
            }
            else
            {
                ModifyDatabase(null, new List<string> { KeyDnsSettings });
            }
        }

        public void Write(Settings settings)
        {
            var subnetMask = _ipHelper.SubnetMask(settings.SubnetAddress, settings.SubnetMaskSize);
            var vmIp = _ipHelper.VmIp(settings.SubnetAddress, settings.SubnetMaskSize);
            var switchIp = _ipHelper.SwitchIp(settings.SubnetAddress, settings.SubnetMaskSize);
            var keyValuesToWrite = new List<Tuple<string, string>>
            {
                Tuple.Create(KeyNetworkConfig, ValueNetworkConfigStatic),
                Tuple.Create(KeyNetworkAddress, vmIp),
                Tuple.Create(KeyNetworkNetmask, subnetMask),
                Tuple.Create(KeyNetworkGateway, switchIp),
                Tuple.Create(KeyNetwork, ValueNetworkHybrid),
                Tuple.Create(KeyDaemonOptions, settings.DaemonOptions),
                Tuple.Create(KeySysCtlConf, settings.SysCtlConf),
                Tuple.Create(KeyLocalHostPortForwarding, "true")
            };

            var keysToRemove = new List<string>();

            if (settings.UseHttpProxy)
            {
                var proxyHttp = ResolveProxyForMoby(settings.ProxyHttp, switchIp);
                var proxyHttps = ResolveProxyForMoby(settings.ProxyHttps, switchIp);
                keyValuesToWrite.Add(Tuple.Create(KeyProxyExclude, settings.ProxyExclude ?? string.Empty));
                keyValuesToWrite.Add(Tuple.Create(KeyProxyHttp, proxyHttp ?? string.Empty));
                keyValuesToWrite.Add(Tuple.Create(KeyProxyHttps, proxyHttps ?? string.Empty));
            }
            else
            {
                keysToRemove.Add(KeyProxyExclude);
                keysToRemove.Add(KeyProxyHttp);
                keysToRemove.Add(KeyProxyHttps);
            }
            ModifyDatabase(keyValuesToWrite, keysToRemove);

            X509Certificate2Collection caCollection;
            using (_clientIdentityRecorder.Impersonate())
            {
                caCollection = FindAllTrustedCa();
            }

            var file = Path.GetTempFileName();
            try
            {
                using (var writer = File.CreateText(file))
                {
                    foreach (var cert in caCollection)
                    {
                        writer.WriteLine("-----BEGIN CERTIFICATE-----");
                        writer.WriteLine(Convert.ToBase64String(cert.Export(X509ContentType.Cert), Base64FormattingOptions.InsertLineBreaks));
                        writer.WriteLine("-----END CERTIFICATE-----");
                    }
                    writer.Flush();
                }
                WriteFile(KeyCaCertsBundle, file);
            }
            finally
            {
                File.Delete(file);
            }
        }

        private static X509Certificate2Collection FindAllTrustedCa()
        {
            var result = new X509Certificate2Collection();
            FillCertificates(result, StoreLocation.CurrentUser, StoreName.AuthRoot);
            FillCertificates(result, StoreLocation.CurrentUser, StoreName.Root);
            FillCertificates(result, StoreLocation.CurrentUser, StoreName.CertificateAuthority);
            FillCertificates(result, StoreLocation.LocalMachine, StoreName.AuthRoot);
            FillCertificates(result, StoreLocation.LocalMachine, StoreName.Root);
            FillCertificates(result, StoreLocation.LocalMachine, StoreName.CertificateAuthority);
            return result;
        }

        private static void FillCertificates(X509Certificate2Collection result, StoreLocation location, StoreName name)
        {
            var store = new X509Store(name, location);
            store.Open(OpenFlags.ReadOnly);
            try
            {
                foreach(var cert in store.Certificates)
                {
                    if (!result.Contains(cert))
                    {
                        result.Add(cert);
                    }
                }
            }
            finally
            {
                store.Close();
            }
        }

        private void ModifyDatabase(ICollection<Tuple<string, string>> keyValuesToWrite, ICollection<string> keysToRemove)
        {
            for (var retry = 5; retry >= 0; retry--)
            {
                try
                {
                    if (keyValuesToWrite != null && keyValuesToWrite.Count > 0)
                        WriteKeys(keyValuesToWrite);
                    if (keysToRemove != null && keysToRemove.Count > 0)
                        RemoveKeys(keysToRemove);
                    return;
                }
                catch (Exception e)
                {
                    if (retry == 1)
                    {
                        throw;
                    }

                    // ReSharper disable once InconsistentlySynchronizedField
                    _logger.Error($"Unable to modify database: {e.Message}. Let's retry in 1s");
                    Thread.Sleep(1000);
                }
            }
        }

        private void RemoveKeys(ICollection<string> keys)
        {
            if (keys.Count == 0)
                return;

            lock (_lock)
            {
                var arguments = "rm";

                foreach (var key in keys)
                {
                    arguments += $" {key}";

                    _logger.Info($"Removing {key}");
                }

                var cli = Paths.InResourcesDir(CliPath);
                var pei = _cmd.RunAsAdministrator(cli, arguments, 0);
                var exitCode = pei.ExitCode;
                if (exitCode != 0)
                {
                    _logger.Error(pei.CombinedOutput);
                    throw new DockerException($"Unable to remove from the database. Exit code: {exitCode}");
                }

                _logger.Info($"Successfully removed {keys.Count} value(s)");
            }
        }

        private void WriteKeys(ICollection<Tuple<string, string>> keyValues)
        {
            if (keyValues.Count == 0)
                return;

            lock (_lock)
            {
                var arguments = "set";

                foreach (var keyValue in keyValues)
                {
                    var key = keyValue.Item1;
                    var value = keyValue.Item2;
                    arguments += $" {key}=\"{value.Replace("\"", "\\\"")}\"";

                    _logger.Info($"Writing {key}={value}");
                }

                var cli = Paths.InResourcesDir(CliPath);
                var pei = _cmd.RunAsAdministrator(cli, arguments, 0);
                var exitCode = pei.ExitCode;
                if (exitCode != 0)
                {
                    _logger.Error(pei.CombinedOutput);
                    throw new DockerException($"Unable to write to the database. Exit code: {exitCode}");
                }

                _logger.Info($"Successfully wrote {keyValues.Count} value(s)");
            }
        }

        private void WriteFile(string key, string file)
        {
            lock (_lock)
            {
                var arguments = "set-file";

                arguments += $" {key}=\"{file.Replace("\"", "\\\"")}\"";

                _logger.Info($"Writing {key}={file}");

                var cli = Paths.InResourcesDir(CliPath);
                var pei = _cmd.RunAsAdministrator(cli, arguments, 0);
                var exitCode = pei.ExitCode;
                if (exitCode != 0)
                {
                    _logger.Error(pei.CombinedOutput);
                    throw new DockerException($"Unable to write to the database. Exit code: {exitCode}");
                }

                _logger.Info("Successfully wrote value");
            }
        }

        private readonly Regex _localHostOr127001Regex = new Regex(@"(localhost|127\.0\.0\.1)", RegexOptions.IgnoreCase);

        private string ResolveProxyForMoby(string localProxy, string hostIp)
        {
            // hostIpResolving
            try
            {
                return _localHostOr127001Regex.Replace(localProxy, hostIp);
            }
            catch
            {
                return localProxy;
            }
        }
    }
}