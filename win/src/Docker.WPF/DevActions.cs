using System;
using Docker.Core;
using Docker.Core.Update;

namespace Docker.WPF
{
    public class DevActions
    {
        private readonly ICmd _cmd;
        private readonly Channel _channel;
        private readonly IUpdater _updater;
        private readonly ErrorReportWindow _errorReportWindow;

        public DevActions(ICmd cmd, Channel channel, IUpdater updater, ErrorReportWindow errorReportWindow)
        {
            _cmd = cmd;
            _channel = channel;
            _updater = updater;
            _errorReportWindow = errorReportWindow;
        }

        public void OpenResourceFolder()
        {
            _cmd.RunWindowed(Paths.ResourcesPath);
        }

        public void OpenSettingsFolder()
        {
            _cmd.RunWindowed(Paths.LocalRoamingApplicationData);
        }

        public Channel GetChannel()
        {
            return _channel;
        }

        public string GetAppCastUrl()
        {
            return _updater.GetAppCastEndPoint();
        }

        public void CheckForUpdate()
        {
            _updater.CheckForUpdates(() => { DialogBox.Show("No Update", "No Update available at this time", "No Update");}, () => { });
        }

        public void OpenSendCrashWindow()
        {
            _errorReportWindow.Show(new Exception("Foobar"), null);
        }

        public void ChangeChannelTo(string value)
        {
            _channel.ChangeTo(value);
            CheckForUpdate();
        }
    }
}