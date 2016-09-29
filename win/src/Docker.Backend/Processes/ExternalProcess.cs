using System;
using System.Diagnostics;
using System.IO;
using Docker.Core;

namespace Docker.Backend.Processes
{
    internal class ExternalProcess
    {
        private readonly Logger _logger;
        private readonly string _exeName;
        private readonly string _processName;

        private bool _isStopping;

        internal ExternalProcess(Type type, string exeName, string processName)
        {
            _logger = new Logger(type);
            _exeName = exeName;
            _processName = processName;
        }

        public void Start(Func<string> arguments)
        {
            var fileName = Paths.InResourcesDir(_exeName);
            if (!File.Exists(fileName))
            {
                throw new DockerException($"File not found '{Path.GetFullPath(fileName)}'");
            }

            Stop();

            _logger.Info("Starting...");
            StartProcess(fileName, arguments);
            _logger.Info("Started");
        }

        private void StartProcess(string fileName, Func<string> arguments)
        {
            var process = new Process
            {
                StartInfo =
                {
                    FileName = fileName,
                    Arguments = arguments(),
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

                _logger.Error("Process died");
            };

            process.Start();
            process.BeginOutputReadLine();
            process.BeginErrorReadLine();
            _isStopping = false;
        }

        public void Stop()
        {
            _isStopping = true;

            KillExistingProcesses();
        }

        private void KillExistingProcesses()
        {
            foreach (var process in Process.GetProcessesByName(_processName))
            {
                _logger.Info($"Killing existing {_processName} with PID {process.Id}");

                try
                {
                    if (!process.HasExited)
                    {
                        process.Kill();
                        process.WaitForExit(5000);
                    }
                }
                catch
                {
                    // Ignore
                }
            }
        }
    }
}
