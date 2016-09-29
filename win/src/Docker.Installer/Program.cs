using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Globalization;
using System.IO;
using System.Linq;
using System.Threading;
using Docker.Core;
using Docker.Core.backend;
using Docker.Core.di;
using Version = Docker.Core.Version;

namespace Docker.Installer
{
    internal class InstallerModule : Module
    {
        protected override void Configure()
        {
            Bind<IBackend, Backend.Backend>();
            Bind<IToolboxMigration, ToolboxMigration>();
        }
    }

    internal class Program
    {
        private static Logger _logger;

        private readonly IBackend _backend;
        private readonly IToolboxMigration _toolboxMigration;

        public Program(IBackend backend, IToolboxMigration toolboxMigration)
        {
            _backend = backend;
            _toolboxMigration = toolboxMigration;
        }

        [STAThread]
        private static void Main(string[] args)
        {
            try
            {
                Thread.CurrentThread.CurrentCulture = new CultureInfo("en-US");

                Logger.Initialize(Path.Combine(Paths.LocalApplicationData, "install-log.txt"));

                _logger = new Logger("Installer");
                _logger.Info("Version: " + new Version().ToHumanString());
                _logger.Info("Starting on: " + DateTime.Now);
                _logger.Info("Resources: " + Paths.ResourcesPath);
                _logger.Info("OS: " + Env.Os.Name);
                _logger.Info("Edition: " + Env.Os.Edition);
                _logger.Info("Id: " + Env.Os.ReleaseId);
                _logger.Info("Build: " + Env.Os.BuildNumber);
                _logger.Info("BuildLabName: " + Env.Os.BuildLabName);
                _logger.Info("Sha1: " + new Git().Sha1());

                using (var singletons = new Singletons(new InstallerModule()))
                {
                    var program = singletons.Get<Program>();
                    program.Run(args);
                }
            }
            catch (Exception ex)
            {
                _logger?.Error(ex.Message);
            }

            Environment.Exit(0); // Always exit 0
        }

        private void Run(IReadOnlyCollection<string> args)
        {
            if (args.Count == 0)
            {
                _logger.Info("Usage: Docker.Installer.exe -[k|m|p|u]");
                _logger.Info(" -k Kill processes");
                _logger.Info(" -m Migrate toolbox");
                _logger.Info(" -p Prepare upgrade");
                _logger.Info(" -u Uninstall");
                Environment.Exit(2);
            }

            if (args.Contains("-k"))
            {
                KillProcesses();
            }
            if (args.Contains("-m"))
            {
                Migrate();
            }
            if (args.Contains("-p"))
            {
                PrepareUpgrade();
            }
            if (args.Contains("-u"))
            {
                Uninstall();
            }

            WindowsMessage.BroadcastSettingsChange();
        }

        private static void KillProcesses()
        {
            _logger.Info("Killing processes...");

            string[] processNames = { "com.docker.service", "Docker for Windows", "com.docker.proxy", "com.docker.db", "com.docker.vpn", "dockerd"};
            foreach (var name in processNames)
            {
                foreach (var process in Process.GetProcessesByName(name))
                {
                    _logger.Info($"Killing existing {name} process with PID {process.Id}");
                    process.Kill();
                    Thread.Sleep(200);
                }
            }

            try
            {
                File.Delete(Paths.DaemonSocketPath);
            }
            catch
            {
                // ignored
            }
        }

        private void Migrate()
        {
            try
            {
                if (!_toolboxMigration.IsToolboxInstalled)
                return;

                _logger.Info("Migrating Toolbox user...");

                _toolboxMigration.MigrateUser();

                _logger.Info("Toolbox user migration completed successfully");
            }
            catch (Exception ex)
            {
                _logger.Info(ex.Message);
            }
        }

        private void Uninstall()
        {
            _logger.Info("Uninstalling...");

            try
            {
                Directory.Delete(Paths.LocalRoamingApplicationData, true);
            }
            catch
            {
                // ignored
            }

            _backend.Destroy(false);

            try
            {
                File.Delete(Paths.DaemonSocketPath);
            }
            catch
            {
                // ignored
            }

            _logger.Info("Uninstall completed successfully");
        }

        private void PrepareUpgrade()
        {
            _logger.Info("Preparing upgrade...");

            try
            {
                File.Delete(Paths.DaemonSocketPath);
            }
            catch
            {
                // ignored
            }
            _backend.Destroy(true);

            _logger.Info("Prepare upgrade completed successfully");
        }
    }
}
