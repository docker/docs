using System;
using System.Threading.Tasks;
using Docker.Core.update;
using Version = Docker.Core.Version;

namespace Docker.Core
{
    public interface IUpdater
    {
        Task CheckForUpdates(Action startingUpdate, Action upToDate);
    }

    public class NopUpdater : IUpdater
    {
        public Task CheckForUpdates(Action startingUpdate, Action upToDate)
        {
            return Task.FromResult<object>(null);
        }
    }

    public class Updater : IUpdater
    {
        private readonly Version _version;

        private readonly Logger _logger;
        private readonly Channel _channel;
        private readonly IFeedDownloader _feedDownloader;
        private readonly IAskUserToUpdate _installWindow;

        public Updater(Version version, Channel channel, IFeedDownloader feedDownloader, IAskUserToUpdate installWindow)
        {
            _logger = new Logger(GetType());
            _version = version;
            _channel = channel;
            _feedDownloader = feedDownloader;
            _installWindow = installWindow;
        }

        public async Task CheckForUpdates(Action startingUpdate, Action upToDate)
        {

            _logger.Info($"Checking for updates on channel {_channel}...");

            var latestUpdate = await _feedDownloader.DownloadLatestUpdateInfo();

            if (latestUpdate == null)
            {
                _logger.Info("No update available");
                upToDate.Invoke();
                return;
            }

            var remoteVersionBuild = latestUpdate.BuildNumber();
            if (remoteVersionBuild == -1)
            {
                _logger.Warning($"Can't read version number of {latestUpdate}");
                upToDate.Invoke();
                return;
            }

            var localVersionBuild = _version.Revision();
            if (remoteVersionBuild <= localVersionBuild)
            {
                _logger.Info($"Local build {localVersionBuild} is as good as the remote {remoteVersionBuild} on channel {_channel}");
                upToDate.Invoke();
                return;
            }

            _logger.Info($"We got a new version {remoteVersionBuild} which is newer than {localVersionBuild}, asking user.");

            _installWindow.AskUser(latestUpdate, startingUpdate);

            _logger.Info("Check for updates done.");
        }
    }
}