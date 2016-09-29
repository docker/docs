using System;
using Docker.Core;
using Docker.Core.Backend;
using Docker.WPF;
using Docker.Core.Tracking;
using System.IO;
using Docker.Core.backend;
using Docker.Core.sharing;

namespace Docker
{
    public class Actions : IActions
    {
        private readonly ITaskQueue _taskQueue;
        private readonly IBackend _backend;
        private readonly IShareHelper _shareHelper;
        private readonly INotifications _notifications;
        private readonly IWelcomeShower _welcomeWindow;
        private readonly IAnalytics _analytics;
        private readonly IToolboxMigration _toolboxMigration;
        private readonly ISettingsLoader _settingsLoader;

        public Actions(ITaskQueue taskQueue, BackendClient backend, IShareHelper shareHelper,
            INotifications notifications, IWelcomeShower welcomeWindow, IAnalytics analytics,
            IToolboxMigration toolboxMigration, ISettingsLoader settingsLoader)
        {
            _taskQueue = taskQueue;
            _backend = backend;
            _shareHelper = shareHelper;
            _notifications = notifications;
            _welcomeWindow = welcomeWindow;
            _analytics = analytics;
            _toolboxMigration = toolboxMigration;
            _settingsLoader = settingsLoader;
        }

        public ISettings Settings { private get; set; }

        public void Start()
        {
            // Do that before SettingsLoader.Load()
            var isFirstLaunchOrHasExecutableChanged = _settingsLoader.IsFirstLaunchOrHasExecutableChanged();

            var settings = _settingsLoader.Load();

            var migrate = false;

            if (_toolboxMigration.DefaultMachineExists
                && _toolboxMigration.HyperVDiskFolderExists
                && !File.Exists(Paths.MobyDiskPath))
            {
                if (ImportFromToolBoxOnFirstStart.ShowConfirm())
                {
                    migrate = true;
                }
            }

            _taskQueue.Queue(() =>
            {
                if (migrate)
                {
                    _notifications.Notify("Docker is migrating default machine...", "This may take some time", true);
                    _backend.MigrateVolume(_toolboxMigration.GetMachineVolumePath("default"));
                }

                var showWelcomeWhale = migrate || isFirstLaunchOrHasExecutableChanged;
                DoStart(settings, showWelcomeWhale);
            });
        }

        // Must be executed in a _taskQueue.Queue. Didn't find a way to enforce that.
        private void DoStart(Settings settings, bool showWelcomeWindow)
        {
            // send a notification before starting the Backend only in the case
            // we will show the welcome popup later
            if (showWelcomeWindow)
                _notifications.Notify("Docker is starting...", "This will only take a few seconds", true);

            _backend.Start(settings);
            _shareHelper.UpdateMounts(settings);

            if (showWelcomeWindow)
            {
                // app is showing the welcome popup
                _notifications.SetStatus("Docker is running", false, AnalyticEvent.InstallShowWelcomePopup);
                _analytics.Track(AnalyticEvent.AppRunning);
                _welcomeWindow.Show(Settings);
            }
            else
            {
                // app is NOT showing the welcome popup
                _notifications.Notify("Docker is running", "Open PowerShell and start hacking with docker or compose", false, AnalyticEvent.AppRunning);
            }
        }

        public void ResetToDefault()
        {
            _taskQueue.QueueWithWaitMessage("Docker will reset to default...", () =>
            {
                _notifications.Notify("Docker is resetting to default...", "This may take some time", true);

                _backend.Stop();
                _backend.Destroy(false);
                _shareHelper.ResetCredential();
                try
                {
                    Directory.Delete(Paths.LocalRoamingApplicationData, true);
                }
                catch
                {
                    // ignored
                }

                DoStart(_settingsLoader.Load(), true);
            });
        }

        public void StopVm()
        {
            _taskQueue.QueueWithWaitMessage("Docker will stop", () =>
            {
                _notifications.Notify("Stopping", "Docker is stopping...", true);

                _backend.Stop();

                _notifications.Notify("Stopped", "Docker is stopped...", false);
            });
        }

        public void ShutdownVm()
        {
            _taskQueue.Shutdown("Docker will shut down", () => _backend.Stop());
        }

        public void MigrateVolume(string name, string volumePath)
        {
            var settings = _settingsLoader.Load();

            _taskQueue.QueueWithWaitMessage($"Docker will migrate {name} machine...", () =>
            {
                _notifications.Notify($"Docker is migrating {name} machine...", "This may take some time", true);

                _backend.MigrateVolume(volumePath);
                _backend.Start(settings);
                _shareHelper.UpdateMounts(settings);

                _notifications.SetStatus("Docker is running", false);
            });
        }

        public void RestartVm(Action<Settings> changes = null)
        {
            var recreate = false;

            if (changes != null)
            {
                var before = _settingsLoader.Load();
                _settingsLoader.SaveChanges(changes);
                var after = _settingsLoader.Load();

                if (!after.SubnetAddress.Equals(before.SubnetAddress) ||
                    !after.SubnetMaskSize.Equals(before.SubnetMaskSize))
                {
                    // The hyperv switch needs to be recreated
                    recreate = true;
                }
            }

            RestartVm(recreate);
        }

        public void RecreateVm()
        {
            RestartVm(true);
        }

        private void RestartVm(bool mustDestroyVm)
        {
            var settings = _settingsLoader.Load();

            _taskQueue.QueueWithWaitMessage("Docker will restart...", () =>
            {
                _notifications.Notify("Docker is restarting...", "This may take some time", true);

                _backend.Stop();
                if (mustDestroyVm) _backend.Destroy(true);
                _backend.Start(settings);
                _shareHelper.UpdateMounts(settings);

                _notifications.SetStatus("Docker is running", false);
            });
        }

        public void SwitchDaemon()
        {
            var settings = _settingsLoader.Load();

            _taskQueue.QueueWithWaitMessage("Docker will switch...", () =>
            {
                _notifications.Notify("Docker is switching...", "This may take some time", true);

                try
                {
                    _backend.SwitchDaemon(settings);
                }
                finally
                {
                    _notifications.SetStatus("Docker is running", false);
                }
            });
        }
    }
}