using System;
using System.Windows;
using Docker.Core;

namespace Docker.WPF
{
    public partial class SettingsWindow
    {
        private readonly Logger _logger;
        private readonly VMSettingsWindow _vmSettingsWindow;
        private readonly SharedVolumesWindow _sharedVolumesWindow;
        private readonly ToolboxMigration _toolboxMigration;

        private Action<int, int> _onApplyVmSettingsChanges;
        private Action _onResetShare;
        private Func<string, bool> _onShare;
        private Action<string> _onUnshare;
        private Action _onResetToDefault;
        private Action _onMigrate;
        private bool _initialized;

        public SettingsWindow(VMSettingsWindow vmSettingsWindow, SharedVolumesWindow sharedVolumesWindow, TaskQueue taskQueue, ToolboxMigration toolboxMigration)
        {
            _logger = new Logger(GetType());
            _vmSettingsWindow = vmSettingsWindow;
            _sharedVolumesWindow = sharedVolumesWindow;
            _toolboxMigration = toolboxMigration;

            taskQueue.Started += DisableActionLinks;
            taskQueue.Ended += EnableActionLinks;

            InitializeComponent();
        }

        public void Open(Action<int, int> onApplyVmSettingsChanges, Action onResetShare, Func<string, bool> onShare, Action<string> onUnshare, Action onResetToDefault, Action onMigrate)
        {
            _onApplyVmSettingsChanges = onApplyVmSettingsChanges;
            _onResetShare = onResetShare;
            _onShare = onShare;
            _onUnshare = onUnshare;
            _onResetToDefault = onResetToDefault;
            _onMigrate = onMigrate;

            var settings = Settings.Load();
            AutoCheckUpdatesBox.IsChecked = settings.AutoUpdateEnabled;
            AutoStartDockerBox.IsChecked = settings.StartAtLogin;

            _initialized = true;
             
            Show();
            Activate();
        }

        private void LaunchDockerBox_Checked(object sender, RoutedEventArgs e)
        {
            SaveSettings();
        }

        private void AutoUpdateBox_Checked(object sender, RoutedEventArgs e)
        {
            SaveSettings();
        }

        private void VMSettingsButton_Click(object sender, RoutedEventArgs e)
        {
            _vmSettingsWindow.OpenModal(_onApplyVmSettingsChanges);
        }

        private void SharedVolumesButton_Click(object sender, RoutedEventArgs e)
        {
            _sharedVolumesWindow.OpenModal(_onResetShare, _onShare, _onUnshare);
        }

        private void MigrateButton_Click(object sender, RoutedEventArgs e)
        {
            if (ImportFromToolBox.ShowImportFromToolBox())
            {
                Hide();
                _onMigrate();
            }
        }

        private void ResetButton_Click(object sender, RoutedEventArgs e)
        {
            if (ResetToDefaultBox.ShowResetToFactoryDefaults())
            {
                _sharedVolumesWindow.Hide();
                _vmSettingsWindow.Hide();
                Hide();
                _onResetToDefault();
            }
        }

        protected override void OnClosing(System.ComponentModel.CancelEventArgs e)
        {
            e.Cancel = true;
            Hide();
            _sharedVolumesWindow.Hide();
            _vmSettingsWindow.Hide();
        }

        private void SaveSettings()
        {
            if (!_initialized) return;

            _logger.Info("Saving settings");

            var settings = Settings.Load();
            settings.AutoUpdateEnabled = AutoCheckUpdatesBox.IsChecked ?? false;
            settings.StartAtLogin = AutoStartDockerBox.IsChecked ?? false;
            settings.Save();
        }

        private void DisableActionLinks()
        {
            Dispatcher.BeginInvoke((Action)(() =>
            {
                ChangeVmSettingsButton.IsEnabled = false;
                ManageSharedVolumesButton.IsEnabled = false;
                ImportToolboxButton.IsEnabled = false;
                ResetFactoryDefaultsButton.IsEnabled = false;
            }));
        }

        private void EnableActionLinks()
        {
            Dispatcher.BeginInvoke((Action)(() =>
            {
                ChangeVmSettingsButton.IsEnabled = true;
                ManageSharedVolumesButton.IsEnabled = true;
                ImportToolboxButton.IsEnabled = _toolboxMigration.IsDefaultMachineExists;
                ResetFactoryDefaultsButton.IsEnabled = true;
            }));
        }
    }
}
