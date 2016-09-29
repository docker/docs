using System;
using System.ComponentModel;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Input;
using Docker.Core;
using Docker.Core.backend;
using Docker.WPF.Crash;

namespace Docker.WPF
{
    public partial class ErrorReportWindow
    {
        private readonly Logger _logger;
        private readonly ICrashReport _crashReport;
        private readonly Channel _channel;
        private Exception _exception;

        public ErrorReportWindow(ICrashReport crashReport, Channel channel)
        {
            _logger = new Logger(GetType());
            _crashReport = crashReport;
            _channel = channel;

            InitializeComponent();
        }

        public void Show(Exception exception, string customMessage)
        {
            _exception = exception;
            Dispatcher.Invoke(() =>
            {
                var tempCustomMessage = customMessage == null ? "An error occured" : customMessage.Split('\n')[0];
                ErrorText.Text = tempCustomMessage;
                if (exception is BackendException)
                {
                    ErrorMessage.Text = exception.Message + Environment.NewLine + exception.InnerException.StackTrace;
                }
                else
                {
                    ErrorMessage.Text = exception.Message + Environment.NewLine + exception.StackTrace;
                }
                DiagnosticText.Text = "";
                DiagnosticId.Visibility = Visibility.Hidden;
                if (!_channel.IsStable())
                {
                    UploadButton.Visibility = Visibility.Hidden;
                    UploadDiagnostics();
                }
                else
                {
                    DiagnosticText.Text = "You can send a crash report to help troubleshoot your issue.";
                    UploadButton.Visibility = Visibility.Visible;
                }
                ShowDialog();
                Activate();
            });
        }

        private void OnOkClicked(object sender, RoutedEventArgs e)
        {
            Close();
        }

        private void OpenLogClick(object sender, RoutedEventArgs e)
        {
            _logger.Info("Open logs");

            new Cmd().RunWindowed(Paths.LogFilename);
        }

        private async void UploadDiagnostics()
        {
            Cursor = Cursors.Wait;
            try
            {
                ButtonOk.IsEnabled = false;
                UploadButton.IsEnabled = false;

                DiagnosticText.Text = "A crash report is being collected and uploaded, please wait...";
                var errorId = await Task.Run(() => _crashReport.Send(_exception));
                DiagnosticText.Text = "A crash report has been uploaded, with the diagnostic id:";
                
                DiagnosticId.Text = errorId;
                DiagnosticId.Visibility = Visibility.Visible;
            }
            catch (Exception ex)
            {
                DiagnosticId.Text = $"Unable to upload a diagnostic: {ex.Message}";
            }
            finally
            {
                Cursor = null;
                ButtonOk.IsEnabled = true;
                UploadButton.IsEnabled = true;
            }
        }

        protected override void OnClosing(CancelEventArgs e)
        {
            e.Cancel = true;
            Hide();
        }

        private void UploadDiagnosticClicked(object sender, RoutedEventArgs e)
        {
            UploadDiagnostics();
        }

        private void OnErrorMessageDoubleClick(object sender, MouseButtonEventArgs e)
        {
            ErrorMessage.SelectAll();
        }

        private void OnDiagnosticIdDoubleClick(object sender, MouseButtonEventArgs e)
        {
            DiagnosticId.SelectAll();
        }
    }
}
