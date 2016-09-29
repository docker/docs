using System;
using System.Diagnostics;
using System.Threading.Tasks;
using System.Windows;
using Docker.Core;
using Docker.WPF.Crash;

namespace Docker.WPF
{
    public partial class FeedbackSettings : IActivable, IAlwaysEnabled
    {
        private readonly ILogger _logger;
        private readonly ICrashReport _crashReport;
        private readonly SettingsWindow _settingsWindow;

        public FeedbackSettings(SettingsWindow settingsWindow, ICrashReport crashReport)
        {
            _logger = new Logger(GetType());
            _settingsWindow = settingsWindow;
            _crashReport = crashReport;

            InitializeComponent();
        }

        public void Refresh(Settings settings)
        {
        }

        private void OpenDocumentationClick(object sender, RoutedEventArgs e)
        {
            _logger.Info("Open documentation");

            Process.Start(Urls.GettingStartedGuide);
        }

        private void OpenIssuesClick(object sender, RoutedEventArgs e)
        {
            _logger.Info("Open github issues");

            Process.Start(Urls.GithubIssues);
        }

        private void OpenLogClick(object sender, RoutedEventArgs e)
        {
            _logger.Info("Open logs");

            new Cmd().RunWindowed(Paths.LogFilename);
        }

        private async void UploadDiagnosticClick(object sender, RoutedEventArgs e)
        {
            _logger.Info("Upload diagnostic");

            var origStatus = _settingsWindow.SetStatus(Status.UploadingDiag);
            try
            {
                DiagnosticId.IsEnabled = false;
                Send.IsEnabled = false;
                DiagnosticId.Text = "";

                var errorId = await Task.Run(() => _crashReport.SendDiagnostic());

                DiagnosticId.Text = $"A diagnostic was uploaded with id: {errorId}";
            }
            catch (Exception ex)
            {
                DiagnosticId.Text = $"Unable to upload a diagnostic: {ex.Message}";
            }
            finally
            {
                _settingsWindow.SetStatus(origStatus);
                DiagnosticId.IsEnabled = true;
                Send.IsEnabled = true;
            }
        }
    }
}