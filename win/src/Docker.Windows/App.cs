using Docker.Core;
using Docker.WPF;
using System;
using System.Diagnostics;
using System.Windows;
using Docker.Core.backend;
using Docker.Core.update;

namespace Docker
{
    public class App
    {
        private readonly Logger _logger;
        private readonly IUpdateCheckTimer _updateCheckTimer;
        private readonly SettingsWindow _settingsWindow;
        private readonly Systray _systray;
        private readonly IActions _actions;
        private readonly AboutBox _aboutBox;
        private readonly INotifications _notifications;
        private readonly KitematicHelper _kitematicHelper;
        private readonly IAppShutdown _appShutdown;
        private readonly ISettingsLoader _settingsLoader;
        private readonly ContainerEngineHelper _containerEngineHelper;

        public App(
            IUpdateCheckTimer updateCheckTimer, SettingsWindow settingsWindow, Systray systray, AboutBox aboutBox,
            IActions actions, INotifications notifications, IBackendExceptionReplier backendExceptionReplier,
            IBackendInstallFeatureExceptionReplier backendInstallFeatureExceptionReplier, KitematicHelper kitematicHelper,
            IAppShutdown appShutdown, ISettingsLoader settingsLoader, ContainerEngineHelper containerEngineHelper)
        {
            _logger = new Logger(GetType());
            _updateCheckTimer = updateCheckTimer;
            _settingsWindow = settingsWindow;
            _systray = systray;
            _aboutBox = aboutBox;
            _actions = actions;
            _notifications = notifications;
            _kitematicHelper = kitematicHelper;
            _appShutdown = appShutdown;
            _settingsLoader = settingsLoader;
            _containerEngineHelper = containerEngineHelper;
            backendExceptionReplier.SetQuitAction(ShutdownAndExit);
            backendInstallFeatureExceptionReplier.SetQuitAction(ShutdownAndExit);
        }

        public void Initialize()
        {
            _systray.Initialize(OpenSystrayMenu);

            _actions.Start();

            if (_settingsLoader.Load().AutoUpdateEnabled)
            {
                _updateCheckTimer.CheckOnce(Shutdown, () => { });
                _updateCheckTimer.Start(Shutdown, () => { });
            }
        }

        internal void OpenSystrayMenu()
        {
            _systray.AddItem("&About Docker...", About_Click);
            _systray.AddItem("&Check for Updates...", CheckForUpdates_Click);
            _systray.AddSeparator();
            _systray.AddItem("&Settings...", Settings_Click);
            if (_containerEngineHelper.CanUseWindowsContainers())
            {
                _systray.AddItem(_containerEngineHelper.UseLinuxContainerEngine
                        ? "Switch to &Windows containers..."
                        : "Switch to &Linux containers...", SwitchDaemon_Click);
            }
            _systray.AddSeparator();
            _systray.AddItem("Documentation...", Documentation_Click);
            _systray.AddItem("Diagnose && &Feedback...", Feedback_Click);
            _systray.AddItem("Open &Kitematic...", Kitematic_Click);
            _systray.AddSeparator();
            _systray.AddItem("E&xit Docker", Quit_Click);
        }

        internal void Shutdown()
        {
            Application.Current?.Dispatcher.InvokeAsync(() =>
            {
                _appShutdown.Shutdown();
            });
        }

        internal void ShutdownAndExit()
        {
            Application.Current?.Dispatcher.InvokeAsync(() =>
            {
                _appShutdown.Shutdown();
                System.Windows.Forms.Application.ExitThread();
            });
        }

        private void About_Click()
        {
            _aboutBox.Show();
        }

        private void CheckForUpdates_Click()
        {
            try
            {
                _updateCheckTimer.CheckOnce(Shutdown, () =>
                {
                    _notifications.Notify("Docker is up-to-date", "You have the latest version of Docker installed", false);
                });
            }
            catch (Exception ex)
            {
                _notifications.NotifyError(ex);
            }
        }

        private void Kitematic_Click()
        {
            _kitematicHelper.Open();
        }

        private void Settings_Click()
        {
            _settingsWindow.Open("General");
        }

        private void Feedback_Click()
        {
            _settingsWindow.Open("Diagnose & Feedback");
        }

        private void Documentation_Click()
        {
            Process.Start(Urls.GettingStartedGuide);
        }

        private void SwitchDaemon_Click()
        {
            _actions.SwitchDaemon();
        }

        private void Quit_Click()
        {
            _logger.Info("User clicked on exit");
            ShutdownAndExit();
        }
    }
}
