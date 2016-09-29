using System;
using System.Diagnostics;
using System.IO;
using System.Management.Automation;

namespace Docker.Core
{
    public interface ICmd
    {
        int Run(string filename, string arguments);
        Cmd.ProcessExecutionInfo Run(string filename, string arguments, int timeout);
        Cmd.ProcessExecutionInfo RunAsAdministrator(string filename, string arguments, int timeout);
        void RunWindowed(string filename, string arguments = null);
        int AskForElevated(string filename, string arguments);
    }

    public class Cmd : ICmd
    {
        private readonly Logger _logger;

        public Cmd()
        {
            _logger = new Logger(GetType());
        }

        public int Run(string filename, string arguments)
        {
            var process = CreateProcess(filename, arguments);

            process.Start();
            process.OutputDataReceived += (sender, args) =>
            {
                if (args.Data == null) return;
                _logger.Info(args.Data);
            };
            process.ErrorDataReceived += (sender, args) =>
            {
                if (args.Data == null) return;
                _logger.Info(args.Data);
            };
            process.BeginOutputReadLine();
            process.BeginErrorReadLine();
            process.WaitForExit();

            return process.ExitCode;
        }

        public class ProcessExecutionInfo
        {
            public int ExitCode;
            public string StandardOutput;
            public string ErrorOutput;
            public string CombinedOutput;
            public bool TimedOut;
        }

        public ProcessExecutionInfo Run(string filename, string arguments, int timeout)
        {
            var process = CreateProcess(filename, arguments);
            return StartProcessAndCaptureOutput(timeout, process);
        }

        public ProcessExecutionInfo RunAsAdministrator(string filename, string arguments, int timeout)
        {
            var process = CreateAdministratorProcess(filename, arguments);
            return StartProcessAndCaptureOutput(timeout, process);
        }

        public int AskForElevated(string filename, string arguments)
        {
            var process = CreateAskElevatedProcess(filename, arguments);
            process.Start();
            process.WaitForExit();
            return process.ExitCode;
        }

        private static ProcessExecutionInfo StartProcessAndCaptureOutput(int timeout, Process process)
        {
            process.Start();

            var standardOutput = "";
            var errorOutput = "";
            var combinedOutput = "";
            var combinedOutputLock = new object();

            process.OutputDataReceived += (sender, args) =>
            {
                var line = args.Data;
                if (line == null) return;

                standardOutput += line + Environment.NewLine;
                lock (combinedOutputLock)
                {
                    combinedOutput += line + Environment.NewLine;
                }
            };

            process.ErrorDataReceived += (sender, args) =>
            {
                var line = args.Data;
                if (line == null) return;

                errorOutput += line + Environment.NewLine;
                lock (combinedOutputLock)
                {
                    combinedOutput += line + Environment.NewLine;
                }
            };

            process.BeginOutputReadLine();
            process.BeginErrorReadLine();

            var timedOut = false;
            if (timeout == 0)
                process.WaitForExit();
            else if (!process.WaitForExit(timeout)) // process timed out
            {
                try { process.Kill(); }
                catch
                {
                    // ignored
                }
                timedOut = true;
            }

            return new ProcessExecutionInfo
            {
                ExitCode = process.ExitCode,
                StandardOutput = standardOutput,
                ErrorOutput = errorOutput,
                CombinedOutput = combinedOutput,
                TimedOut = timedOut
            };
        }

        private static Process CreateProcess(string filename, string arguments)
        {
            return new Process
            {
                StartInfo =
                {
                    CreateNoWindow = true,
                    UseShellExecute = false,
                    RedirectStandardOutput = true,
                    RedirectStandardError = true,
                    FileName = filename,
                    Arguments = arguments
                }
            };
        }

        private static Process CreateAdministratorProcess(string filename, string arguments)
        {
            return new Process
            {
                StartInfo =
                {
                    CreateNoWindow = true,
                    UseShellExecute = false,
                    RedirectStandardOutput = true,
                    RedirectStandardError = true,
                    FileName = filename,
                    Arguments = arguments,
                    Verb = "runas"
                }
            };
        }

        private static Process CreateAskElevatedProcess(string filename, string arguments)
        {
            return new Process
            {
                StartInfo =
                {
                    UseShellExecute = true,
                    FileName = filename,
                    Arguments = arguments,
                    Verb = "runas"
                }
            };
        }

        public void RunWindowed(string filename, string arguments = null)
        {
            var workingDirectory = Path.GetDirectoryName(filename);
            if (workingDirectory == null)
            {
                throw new RuntimeException($"Working directory not found for {filename}");
            }

            Process.Start(new ProcessStartInfo(filename, arguments)
            {
                WorkingDirectory = workingDirectory
            });
        }
    }
}