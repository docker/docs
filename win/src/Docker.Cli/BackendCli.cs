using System;
using Docker.Core;
using Docker.Core.backend;
using Docker.Core.sharing;
using Docker.WPF;
using Docker.WPF.Crash;
using Version = Docker.Core.Version;

namespace Docker.Cli
{
    internal class BackendCli
    {
        internal const string Secret = "--testftw!928374kasljf039";

        private readonly IActions _actions;
        private readonly ICrashReport _crashReport;
        private readonly ISettingsLoader _settingsLoader;
        private readonly ITaskQueue _taskQueue;
        private readonly IConsoleWriter _console;
        private readonly IShareHelper _shareHelper;
        private readonly IBackend _backend;

        public BackendCli(
            IActions actions, ITaskQueue taskQueue, ICrashReport crashReport,
            ISettingsLoader settingsLoader, IConsoleWriter console, IShareHelper shareHelper,
            IBackend backend)
        {
            _actions = actions;
            _taskQueue = taskQueue;
            _crashReport = crashReport;
            _settingsLoader = settingsLoader;
            _console = console;
            _shareHelper = shareHelper;
            _backend = backend;
        }

        public static void Usage()
        {
            Console.WriteLine(@"Usage: DockerCli.exe [-SwitchDaemon] [-Version]");
            Console.WriteLine(@"  -Version: Show the Docker for Windows version information");
            Console.WriteLine(@"  -SwitchDaemon: Point the Docker CLI to either Linux containers or Windows containers");
            Console.WriteLine(@"  -SharedDrives: List the shared drives");
        }

        public void Run(string[] args)
        {
            try
            {
                var flags = new Flags(args);

                flags.If("-SwitchDaemon", SwitchDaemon);
                flags.If("-Version", Version);
                flags.If("-SharedDrives", SharedDrives);

                // PRIVATE FLAGS
                // !!! THE FOLLOWING ACTIONS ARE FOR TESTS ONLY.
                // SHOULD NOT MADE AVAILABLE TO THE END USERS !!!
                //
                if (!flags.Contains(Secret)) return;

                flags.If("-Start", Start);
                flags.If("-Stop", Stop);
                flags.If("-SendDiagnostic", SendDiagnostic);
                flags.If("-ResetToDefault", ResetToDefault);
                flags.If("-ResetCredential", ResetCredential);
                flags.With("-MigrateVolume=", MigrateVolume);
                flags.With("-Mount=", Mount);
                flags.With("-Unmount=", Unmount);
                flags.With("-SetMemory=", SetMemory);
                flags.With("-SetCpus=", SetCpus);
                flags.With("-SetDNS=", SetDns);
                flags.With("-SetIP=", SetIp);
                flags.With("-SetDaemonJson=", SetDaemonJson);
            }
            finally
            {
                _taskQueue.Shutdown();
            }
        }

        private void SwitchDaemon()
        {
            _actions.SwitchDaemon();
        }

        private void Version()
        {
            var version = new Version();
            var git = new Git();
            var channel = Channel.Load();

            _console.WriteLine("");
            _console.WriteLine(@"Docker for Windows");
            _console.WriteLine($"Version: {version}");
            _console.WriteLine($"Sha1: {git.Sha1()}");
            _console.WriteLine($"Build Number: {version.Revision()}");
            _console.WriteLine($"App Version: {version.ToAppVersionName()}");
            _console.WriteLine($"Human Version: {version.ToHumanStringWithBuildNumber()}");
            _console.WriteLine($"Channel: {channel.Name}");
            _console.WriteLine($"OS Name: {Env.Os.Name}");
            _console.WriteLine($"Windows Edition: {Env.Os.Edition}");
            _console.WriteLine($"Windows Build Number: {Env.Os.BuildNumber}");
        }

        private void Start()
        {
            _actions.Start();
        }

        private void Stop()
        {
            _actions.StopVm();
        }

        private void SendDiagnostic()
        {
            _console.WriteLine(_crashReport.SendDiagnostic());
        }

        private void ResetToDefault()
        {
            _actions.ResetToDefault();
        }

        private void ResetCredential()
        {
            _shareHelper.ResetCredential();
        }

        private void MigrateVolume(string volumePath)
        {
            _actions.MigrateVolume("default", volumePath);
        }

        private void SharedDrives()
        {
            _console.WriteLine(string.Join(",", _backend.SharedDrives()));
        }

        private void Mount(string disk)
        {
            _shareHelper.Mount(disk, _settingsLoader.Load());
        }

        private void Unmount(string disk)
        {
            _shareHelper.Unmount(disk);
        }

        private void SetMemory(string memory)
        {
            _actions.RestartVm(_ => _.VmMemory = int.Parse(memory));
        }

        private void SetCpus(string cpus)
        {
            _actions.RestartVm(_ => _.VmCpus = int.Parse(cpus));
        }

        private void SetDns(string dnsServer)
        {
            if ("automatic".Equals(dnsServer))
            {
                _actions.RestartVm(_ => _.UseDnsForwarder = true);
            }
            else
            {
                _actions.RestartVm(_ =>
                {
                    _.UseDnsForwarder = false;
                    _.NameServer = dnsServer;
                });
            }
        }

        private void SetIp(string ipMask)
        {
            _actions.RestartVm(_ =>
            {
                var parts = ipMask.Split(new[] { '/' }, 2);
                _.SubnetAddress = parts[0];
                _.SubnetMaskSize = new IpHelper().GetMaskSize(parts[1]);
            });
        }

        private void SetDaemonJson(string json)
        {
            _actions.RestartVm(_ => _.DaemonOptions = json);
        }
    }
}