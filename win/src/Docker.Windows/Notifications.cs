using Docker.Core;
using System;
using System.Windows.Forms;
using Docker.Core.Tracking;
using Docker.Properties;
using Docker.WPF;

namespace Docker
{
    public class Notifications : INotifications
    {
        private const string OutOfMemory = "0x8007000E";

        private readonly Logger _logger;
        private readonly Systray _systray;
        private readonly IAnalytics _analytics;
        private readonly ErrorReportWindow _errorReportWindow;

        public Notifications(Systray systray, IAnalytics analytics, ErrorReportWindow errorReportWindow)
        {
            _logger = new Logger(GetType());
            _systray = systray;
            _analytics = analytics;
            _errorReportWindow = errorReportWindow;
        }

        public void NotifyError(string customMessage, Exception exception)
        {
            _logger.Error(exception.Message);
            if (exception.Message.Contains(OutOfMemory))
            {
                _systray.SetStatus("Out of memory", Resources.systray_icon_red, false);
                NotEnoughtMemoryBox.ShowOk();
            }
            else
            {
                _systray.SetStatus(exception.Message, Resources.systray_icon_red, false);
                _errorReportWindow.Show(exception, customMessage);
            }
        }

        public void NotifyError(Exception exception)
        {
            NotifyError(null, exception);
        }

        public void Notify(string message, string details, bool? useAnimation)
        {
            _logger.Info(message);
            _systray.NotifyStatus(message, details, ToolTipIcon.Info, Resources.systray_icon_inverted, useAnimation);
        }

        public void Notify(string message, string details, bool? useAnimation, AnalyticEvent mixpanelEvent)
        {
            _logger.Info(message);
            _analytics.Track(mixpanelEvent);
            _systray.NotifyStatus(message, details, ToolTipIcon.Info, Resources.systray_icon_inverted, useAnimation);
        }

        public void SetStatus(string message, bool? useAnimation, AnalyticEvent mixpanelEvent)
        {
            _logger.Info(message);
            _analytics.Track(mixpanelEvent);
            _systray.SetStatus(message, Resources.systray_icon_inverted, useAnimation);
        }

        public void SetStatus(string message, bool? useAnimation)
        {
            _logger.Info(message);
            _systray.SetStatus(message, Resources.systray_icon_inverted, useAnimation);
        }
    }
}