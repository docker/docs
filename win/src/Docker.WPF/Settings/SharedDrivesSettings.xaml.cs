using Docker.Core;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using Docker.Core.sharing;

namespace Docker.WPF
{
    public partial class SharedDrivesSettings : IActivable
    {
        private class SharedDriveView
        {
            public string DriveLetter { get; }
            public bool ShareEnabled { get; set; }
            public SharedDriveView(string driveLetter)
            {
                DriveLetter = driveLetter;
            }
        }

        private readonly ILogger _logger;
        private readonly IShareHelper _shareHelper;
        private readonly ISettingsLoader _settingsLoader;
        private readonly SettingsWindow _settingsWindow;

        private List<SharedDriveView> _driveList;

        public SharedDrivesSettings(SystemChangesNotifier notifier, ISettingsLoader settingsLoader, IShareHelper shareHelper, SettingsWindow settingsWindow)
        {
            _logger = new Logger(GetType());
            _settingsLoader = settingsLoader;
            _shareHelper = shareHelper;
            _settingsWindow = settingsWindow;

            InitializeComponent();

            notifier.OnDiskMount += DriveModificationDetected;
            notifier.OnDiskUnmount += DriveModificationDetected;
        }

        public void Refresh(Settings settings)
        {
            _driveList = new List<SharedDriveView>();

            foreach (var drive in DriveInfo.GetDrives())
            {
                if (!drive.IsReady) continue;
                if (drive.DriveType == DriveType.Network) continue;

                var driveLetter = drive.Name.Substring(0, 1);

                var newDriveView = new SharedDriveView(driveLetter);

                bool shared;
                if (settings.SharedDrives.TryGetValue(driveLetter, out shared))
                {
                    newDriveView.ShareEnabled = shared;
                }

                _driveList.Add(newDriveView);
            }

            SharedVolumeListView.ItemsSource = _driveList;
        }

        private async void ResetButton_Click(object sender, RoutedEventArgs e)
        {
            var origStatus = _settingsWindow.SetStatus(Status.UpdatingDrives);
            try
            {
                await Task.Run(() =>
                {
                    _shareHelper.ResetCredential();
                });
            }
            finally
            {
                _settingsWindow.SetStatus(origStatus);
            }
        }

        private async void OkButton_Click(object sender, RoutedEventArgs e)
        {
            _logger.Info("Apply shared drive settings");

            var origStatus = _settingsWindow.SetStatus(Status.UpdatingDrives);
            try
            {
                await Task.Run(() =>
                {
                    foreach (var drive in _driveList)
                    {
                        if (drive.ShareEnabled)
                        {
                            if (!_shareHelper.Mount(drive.DriveLetter, _settingsLoader.Load()))
                            {
                                // User pressed cancel on the Credential window
                                return;
                            }
                        }
                        else
                        {
                            _shareHelper.Unmount(drive.DriveLetter);
                        }
                    }
                });
            }
            finally
            {
                _settingsWindow.SetStatus(origStatus);
            }
        }

        private void PromptText_OnMouseDoubleClick(object sender, MouseButtonEventArgs e)
        {
            (sender as TextBox)?.SelectAll();
        }

        private void DriveModificationDetected()
        {
            Dispatcher.Invoke(() =>
            {
                _driveList = new List<SharedDriveView>();

                foreach (var drive in DriveInfo.GetDrives())
                {
                    if (!drive.IsReady) continue;
                    if (drive.DriveType == DriveType.Network) continue;

                    var driveLetter = drive.Name.Substring(0, 1);

                    var newDriveView = new SharedDriveView(driveLetter);
                    var foundDriveView = (SharedVolumeListView.ItemsSource as List<SharedDriveView>)?.FirstOrDefault(s => s.DriveLetter == driveLetter);
                    if (foundDriveView != null)
                    {
                        newDriveView.ShareEnabled = foundDriveView.ShareEnabled;
                    }

                    _driveList.Add(newDriveView);
                }

                SharedVolumeListView.ItemsSource = _driveList;
            });
        }
    }
}
