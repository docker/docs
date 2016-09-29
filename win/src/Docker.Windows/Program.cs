using System;
using System.Collections.Generic;
using System.Globalization;
using System.ServiceProcess;
using System.Threading;
using System.Threading.Tasks;
using Docker.Core;
using Docker.Core.Backend;
using Docker.WPF;
using Docker.Backend;
using Docker.Core.backend;
using Docker.Core.di;
using Docker.Core.Pipe;
using Docker.Core.Tracking;
using Docker.Core.Update;
using Docker.WPF.Crash;
using Docker.WPF.Credentials;
using Logger = Docker.Core.Logger;

namespace Docker
{
    internal class DockerModule : Module
    {
        private readonly Logger _logger;
        private readonly IReadOnlyCollection<string> _args;

        internal DockerModule(Logger logger, IReadOnlyCollection<string> args)
        {
            _logger = logger;
            _args = args;
        }

        protected override void Configure()
        {
            var commandLineOptions = new CommandLineOptions(_args);
            Bind(commandLineOptions);
            Bind(_logger);
            Bind(Channel.Load());
            Bind<IBackend, BackendClient>();
            Bind<ICrashReport, CrashReport>();
            Bind<INotifications, Notifications>();
            Bind<IAnalytics, SegmentApi>();
            Bind<IInstallUpdateWindow, InstallUpdateWindow>();
            Bind<ISharedDrives, SettingsBasedSharedDrives>();

            if (commandLineOptions.DisableCheckForUpdates)
            {
                Bind<IUpdater, NopUpdater>();
            }
            else
            {
                Bind<IUpdater, Updater>();
            }

            if (commandLineOptions.DisableToolboxMigration)
            {
                Bind<IToolboxMigration, NopToolboxMigration>();
            }
            else
            {
                Bind<IToolboxMigration, ToolboxMigration>();
            }

            if (commandLineOptions.Username != null)
            {
                Bind<ICredentialAsker>(new FixedCredentialAsker(commandLineOptions.Username, commandLineOptions.Password));
            }
            else
            {
                Bind<ICredentialAsker, CredentialAsker>();
            }

            if (commandLineOptions.DisableWelcomeWhale)
            {
                Bind<IWelcomeShower, NopWelcomeShower>();
            }
            else
            {
                Bind<IWelcomeShower, WelcomeWindow>();
            }
        }
    }

    internal static class Program
    {
        [STAThread]
        private static void Main(string[] args)
        {
            Environment.Exit(Run(args) ? 0 : 1);
        }

        private static bool Run(IReadOnlyCollection<string> args)
        {
            Logger logger = null;
            ErrorReportWindow errorWindow = null;
            CrashReport crashReport = null;
            Tracking tracking = null;
            try
            {
                Thread.CurrentThread.CurrentCulture = new CultureInfo("en-US");

                Logger.Initialize(Paths.LogFilename);

                logger = new Logger("Program");
                var version = new Core.Version();
                logger.Info("Version: " + version.ToHumanString());
                logger.Info("Starting on: " + DateTime.Now);
                logger.Info("Resources: " + Paths.ResourcesPath);
                logger.Info("OS: " + Env.Os.Name);
                logger.Info("Edition: " + Env.Os.Edition);
                logger.Info("Id: " + Env.Os.ReleaseId);
                logger.Info("Build: " + Env.Os.BuildNumber);
                logger.Info("BuildLabName: " + Env.Os.BuildLabName);
                var git = new Git();
                logger.Info("Sha1: " + git.Sha1());
                logger.Info($"You can send feedback, including this log file, at {Urls.GithubIssues}");

                // ReSharper disable once UnusedVariable
                var app = new System.Windows.Application
                {
                    ShutdownMode = System.Windows.ShutdownMode.OnExplicitShutdown
                };
                System.Windows.Forms.Application.EnableVisualStyles();
                System.Windows.Forms.Application.SetCompatibleTextRenderingDefault(false);

                if (!SingleInstance.Start())
                {
                    logger.Error("An instance is already running. Exiting.");
                    return false;
                }

                var singletons = new Singletons(new DockerModule(logger, args));

                // First initialize crash report tools
                errorWindow = singletons.Get<ErrorReportWindow>();
                crashReport = singletons.Get<CrashReport>();
                var settingsLoader = singletons.Get<ISettingsLoader>();
                tracking = singletons.Get<Tracking>();
                tracking.ChangeTo(settingsLoader.Load().IsTracking);

                var analytics = singletons.Get<IAnalytics>();
                analytics.Track(AnalyticEvent.AppLaunched);

                // After AppLaunched
                var startPumpingBackendLogs = true;
                if (!IsServiceRunning())
                {
                    if (IsDebugMode())
                    {
                        startPumpingBackendLogs = false;
                        new Singletons(new BackendServerModule()).Get<BackendServer>().Run();
                    }
                    else if (!ServiceNotRunningBox.ShowConfirm())
                    {
                        return false;
                    }
                    else
                    {
                        var cmd = singletons.Get<ICmd>();
                        var exitCode = cmd.AskForElevated("net", "start \"com.docker.service\"");
                        if (exitCode != 0)
                        {
                            throw new Exception($"Unable to start Docker for Windows service: {exitCode}");
                        }
                    }
                }

                var backendClient = singletons.Get<IBackend>();
                var remoteVersion = backendClient.Version();
                if (remoteVersion != version.ToHumanString())
                {
                    throw new Exception($"A wrong version of Docker for Windows service is running.{Environment.NewLine}Expected: {version.ToHumanString()}.{Environment.NewLine}Got: {remoteVersion}");
                }

                if (startPumpingBackendLogs)
                {
                    var logPipeClient = new LogPipeClient("dockerLogs", new Logger("Service"));
                    Task.Run(() => logPipeClient.Run());
                }

                var heartbeat = singletons.Get<Heartbeat>();
                heartbeat.Start();

                singletons.Get<Actions>().Settings = singletons.Get<SettingsWindow>(); // TODO: inject by constructor

                var application = singletons.Get<App>();
                application.Initialize();

                var sysTray = singletons.Get<Systray>();
                System.Windows.Forms.Application.Run(sysTray);

                return true;
            }
            catch (Exception exception)
            {
                if (errorWindow != null)
                {
                    errorWindow.Show(exception, "Error during initialization");
                    return false;
                }

                // we are in the danger zone here, we may lack a few classes.
                // Anything can be null.
                logger?.Error(exception.Message);
                if (tracking?.IsEnabled != null)
                {
                    crashReport?.SendOrFailSilently(exception);
                }
                DialogBox.Show("Docker failed to initialize", exception.Message, "Docker failed to initialize");

                return false;
            }
            finally
            {
                SingleInstance.Release();
            }
        }

        private static bool IsDebugMode()
        {
#if DEBUG
            return true;
#else
            return false;
#endif
        }

        private static bool IsServiceRunning()
        {
            try
            {
                var sc = new ServiceController { ServiceName = "com.docker.service" };
                if (sc.Status == ServiceControllerStatus.Running)
                {
                    return true;
                }
            }
            catch
            {
                // Service is not installed
            }

            return false;
        }
    }
}
