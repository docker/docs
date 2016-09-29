using System;
using System.Diagnostics;
using System.IO;
using System.Text.RegularExpressions;
using System.Threading;
using Docker.Core;

namespace Docker.Backend
{
    public interface IProxyProcess
    {
        void Start(Settings settings);
        void Stop();
    }

    public class ProxyProcess : IProxyProcess
    {
        private const string ProxyExe = "com.docker.proxy.exe";

        private static readonly string[] ProcessNames = {"com.docker.proxy", "com.docker.slirp", "com.docker.db"};

        private readonly Logger _logger;
        private readonly IHyperV _hyperV;
        private readonly IpHelper _ipHelper;

        private bool _isStopping;

        public ProxyProcess(IHyperV hyperV, IpHelper ipHelper)
        {
            _logger = new Logger(GetType());
            _hyperV = hyperV;
            _ipHelper = ipHelper;
        }

        public void Start(Settings settings)
        {
            var fileName = Paths.InResourcesDir(ProxyExe);
            if (!File.Exists(fileName))
            {
                throw new DockerException($"File not found '{Path.GetFullPath(fileName)}'");
            }

            KillExistingProcesses();

            _logger.Info("Starting...");
            StartProcess(fileName, settings);
            _logger.Info("Started");

            _isStopping = false;
        }

        private void StartProcess(string fileName, Settings settings)
        {
            var parentName = AppDomain.CurrentDomain.FriendlyName;
            var vmId = _hyperV.GetId();
            var dbPath = Paths.DatabasePath;
            var upstreamDnsServer = settings.UseDnsForwarder ? "auto" : settings.NameServer;
            var switchIp = _ipHelper.SwitchIp(settings.SubnetAddress, settings.SubnetMaskSize);

            var process = new Process
            {
                StartInfo =
                {
                    FileName = fileName,
                    Arguments = $"-upstreamDnsServer={upstreamDnsServer} -dnsIP={switchIp} -VM={vmId} -db=\"{dbPath}\" -poisonPill=\"{parentName}\"",
                    RedirectStandardOutput = true,
                    RedirectStandardError = true
                },
                EnableRaisingEvents = true
            };

            var subProcessLogger = new SubProcessLogger(_logger);
            process.OutputDataReceived += (sender, args) => subProcessLogger.Log(args.Data);
            process.ErrorDataReceived += (sender, args) => subProcessLogger.Log(args.Data);
            process.StartInfo.UseShellExecute = false;
            process.StartInfo.CreateNoWindow = true;
            process.Exited += (sender, args) =>
            {
                if (_isStopping) return;

                _logger.Info("Restarting...");

                Thread.Sleep(200);
                KillExistingProcesses();
                StartProcess(fileName, settings);

                _logger.Info("Restarted");
            };

            process.Start();
            process.BeginOutputReadLine();
            process.BeginErrorReadLine();
        }

        private void KillExistingProcesses()
        {
            foreach (var processName in ProcessNames)
            {
                KillExistingProcesses(processName);
            }
        }

        private void KillExistingProcesses(string name)
        {
            foreach (var process in Process.GetProcessesByName(name))
            {
                _logger.Info($"Killing existing process with PID {process.Id}");

                try
                {
                    if (!process.HasExited)
                    {
                        process.Kill();
                    }
                }
                catch
                {
                    // Ignore
                }
            }
        }

        public void Stop()
        {
            _isStopping = true;

            KillExistingProcesses();
        }
    }

    public class SubProcessLogger
    {
        private static readonly Regex LogWithDate = new Regex("\\d\\d\\d\\d/\\d\\d/\\d\\d \\d\\d:\\d\\d:\\d\\d .*");

        private readonly ILogger _logger;

        public SubProcessLogger(ILogger logger)
        {
            _logger = logger;
        }

        public void Log(string line)
        {
            if (line == null) return;

            var lineNoDate = LogWithDate.Match(line).Success ? line.Substring(20) : line;

            if (lineNoDate.Contains(": [DEBUG] "))
            {
                _logger.Debug(lineNoDate.Replace(": [DEBUG] ", ": "));
            }
            else if (lineNoDate.Contains(": [WARNING] "))
            {
                _logger.Warning(lineNoDate.Replace(": [WARNING] ", ": "));
            }
            else if (lineNoDate.Contains(": [ERROR] "))
            {
                _logger.Error(lineNoDate.Replace(": [ERROR] ", ": "));
            }
            else if (lineNoDate.Contains(": [INFO] "))
            {
                _logger.Info(lineNoDate.Replace(": [INFO] ", ": "));
            }
            else
            {
                _logger.Info(lineNoDate);
            }
        }
    }
}