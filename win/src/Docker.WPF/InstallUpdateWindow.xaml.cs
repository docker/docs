using System;
using System.ComponentModel;
using System.Diagnostics;
using System.IO;
using System.Net;
using System.Windows;
using System.Windows.Input;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;
using Docker.Core;
using Docker.Core.Update;

namespace Docker.WPF
{
    public partial class InstallUpdateWindow : IInstallUpdateWindow
    {
        private readonly ILogger _logger;

        private AvailableUpdate _latestUpdate;
        private Action _startingUpdateAction;
        private bool _canceled;

        public InstallUpdateWindow(Channel channel)
        {
            _logger = new Logger(GetType());

            InitializeComponent();
            Whale.Source = new BitmapImage(new Uri(channel.IsStable() ? "Images/about-docker.png" : "Images/about-docker-beta.png", UriKind.Relative));

            ReleaseNotes.Navigating += ReleaseNotes_Navigating;
            PreviewKeyUp += (sender, e) =>
            {
                if (e.Key == Key.Escape)
                {
                    e.Handled = true;
                    _canceled = true;
                    Close();
                }
            };
        }

        public void Open(AvailableUpdate latestUpdate, Action startingUpdateAction)
        {
            Dispatcher.Invoke(() => {
                _latestUpdate = latestUpdate;
                _startingUpdateAction = startingUpdateAction;
                var version = new Core.Version();
                VersionLabel.Text = $"Docker {latestUpdate.ToHumanStringWithBuildNumber()} is now available.{Environment.NewLine}You have {version.ToHumanStringWithBuildNumber()}.{Environment.NewLine}Would you like to download it now?";
                ReleaseNotes.NavigateToString(ToHtml(latestUpdate.Notes));
                if (!IsVisible)
                {
                    DownloadProgressBar.Visibility = Visibility.Hidden;
                    DownloadProgressBar.Value = 0;
                    CancelButton.IsEnabled = true;
                    InstallButton.IsEnabled = true;
                    _canceled = false;
                }
                Show();
            });
        }

        private static void ReleaseNotes_Navigating(object sender, NavigatingCancelEventArgs e)
        {
            try
            {
                var uri = e.Uri;
                if (uri == null)
                    return;

                var scheme = uri.Scheme;
                if (scheme == "file")
                    return;

                e.Cancel = true;

                if (scheme == "http" || scheme == "https")
                {
                    Process.Start(uri.AbsoluteUri);
                }
            }
            catch
            {
                // ignored
            }
        }

        protected override void OnClosing(CancelEventArgs e)
        {
            e.Cancel = true;
            _canceled = true;
            Hide();
        }

        private static string ToHtml(string releaseNotes)
        {
            return @"<html><style>body{font-family:arial; font-size:.9em;} p{margin:1em,0,.5em,0;} ul{padding:0; margin:0,0,0,30px;}</style><body>" + releaseNotes + "</body></html>";
        }

        private void btnInstall_Click(object sender, RoutedEventArgs e)
        {
            DownloadProgressBar.Visibility = Visibility.Visible;
            CancelButton.IsEnabled = false;
            InstallButton.IsEnabled = false;

            DownloadAndInstallMsi(_latestUpdate, progressInfo => {
                Dispatcher.BeginInvoke((Action)(() =>
                {
                    DownloadProgressBar.Value = progressInfo.ProgressPercentage;
                }));
                return _canceled;
            }, () => Dispatcher.BeginInvoke((Action)Close));
        }

        private void btnCancel_Click(object sender, RoutedEventArgs e)
        {
            Close();
        }

        private class DownloadProgressInformation
        {
            internal DownloadProgressInformation(int progressPercentage)
            {
                ProgressPercentage = progressPercentage;
            }
            public int ProgressPercentage { get; }
        }

        private async void DownloadAndInstallMsi(AvailableUpdate update, Func<DownloadProgressInformation, bool> progress, Action completed)
        {
            try
            {
                var tmpFilename = Path.GetTempFileName();
                var tmpMsiFilename = Path.ChangeExtension(tmpFilename, ".msi");

                _logger.Info($"Download msi to {tmpMsiFilename}");

                var wc = new WebClient();
                wc.DownloadFileCompleted += (a, b) => completed();
                wc.DownloadProgressChanged += (a, b) =>
                {
                    if (progress(new DownloadProgressInformation(b.ProgressPercentage)))
                    {
                        wc.CancelAsync();
                    }
                };

                var task = wc.DownloadFileTaskAsync(new Uri(update.DownloadUrl), tmpMsiFilename);
                await task;

                if (_canceled) return;

                Close();

                RunMsi(tmpMsiFilename, _startingUpdateAction);
            }
            catch
            {
                // Cancelled
            }
        }

        public static void RunMsi(string msiPath, Action startingUpdateAction)
        {
            startingUpdateAction();

            Process.Start(new ProcessStartInfo
            {
                FileName = "msiexec",
                Arguments = $"/passive /qb /norestart /i \"{msiPath}\" LAUNCH_APPLICATION_ON=\"1\"",
                Verb = "runas"
            });
            System.Windows.Forms.Application.ExitThread();
        }
    }
}