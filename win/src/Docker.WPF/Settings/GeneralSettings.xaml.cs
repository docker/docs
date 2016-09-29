using System.Diagnostics;
using System.Windows;
using Docker.Core;
using Docker.Core.Tracking;

namespace Docker.WPF
{
    public partial class GeneralSettings : IActivable
    {
        private readonly ILogger _logger;
        private readonly Channel _channel;
        private readonly Tracking _tracking;
        private readonly ISettingsLoader _settingsLoader;

        private bool _loaded;

        public GeneralSettings(Channel channel, Tracking tracking, ISettingsLoader settingsLoader)
        {
            _logger = new Logger(GetType());
            _channel = channel;
            _tracking = tracking;
            _settingsLoader = settingsLoader;

            InitializeComponent();
        }

        private void LaunchDockerBox_Checked(object sender, RoutedEventArgs e)
        {
            if (!_loaded) return;
            SaveSettings();
        }

        private void AutoUpdateBox_Checked(object sender, RoutedEventArgs e)
        {
            if (!_loaded) return;
            SaveSettings();
        }

        private void TrackingBox_Checked(object sender, RoutedEventArgs e)
        {
            if (!_loaded) return;

            TrackingText.IsEnabled = TrackingBox.IsChecked ?? false;
            SaveSettings();
        }

        private void SaveSettings()
        {
            _logger.Info("Saving settings");

            _settingsLoader.SaveChanges(_ =>
            {
                _.AutoUpdateEnabled = AutoCheckUpdatesBox.IsChecked ?? false;
                _.StartAtLogin = AutoStartDockerBox.IsChecked ?? false;
                _.IsTracking = TrackingBox.IsChecked ?? false;
            });

            _tracking.ChangeTo(TrackingBox.IsChecked ?? false);
        }

        public void Refresh(Settings settings)
        {
            _loaded = false;
            try
            {
                AutoCheckUpdatesBox.IsChecked = settings.AutoUpdateEnabled;
                AutoStartDockerBox.IsChecked = settings.StartAtLogin;

                TrackingBox.IsChecked = settings.IsTracking;
                TrackingBox.IsEnabled = _channel != Channel.Beta;
                TrackingText.IsEnabled = TrackingBox.IsChecked ?? false;

                CurrentChannelText.Text = _channel.Name.ToLower();
                OtherChannelText.Text = _channel == Channel.Beta ? Channel.Stable.Name.ToLower() : Channel.Beta.Name.ToLower();


                _tracking.ChangeTo(settings.IsTracking);
            }
            finally
            {
                _loaded = true;
            }
        }

        private void OpenDocumentation(object sender, RoutedEventArgs e)
        {
            _logger.Info("Open channel link");

            Process.Start(Urls.Channel);
        }
    }
}