using System;
using System.IO;
using System.Net;
using System.Threading.Tasks;
using System.Xml;
using Docker.Core.update;

namespace Docker.Core
{
    public interface IFeedDownloader
    {
        Task<AvailableUpdate> DownloadLatestUpdateInfo();
        string GetAppCastEndPoint();
    }

    public class UpdateFeedDownloader : IFeedDownloader
    {
        private readonly Logger _logger;
        private readonly Channel _channel;
        private readonly bool _useHockeyApp;
        private readonly string _appCastUrl;

        public UpdateFeedDownloader(Channel channel, bool useHockeyApp)
        {
            _logger = new Logger(GetType());
            _channel = channel;
            _useHockeyApp = useHockeyApp;
            _appCastUrl = GetAppCastEndPoint();
        }

        public UpdateFeedDownloader(Channel channel, string appCastUrl)
        {
            _logger = new Logger(GetType());
            _channel = channel;
            _useHockeyApp = false;
            _appCastUrl = appCastUrl;
        }

        public async Task<AvailableUpdate> DownloadLatestUpdateInfo()
        {
            var availableUpdates = new AvailableUpdateList();
            using (var webClient = new WebClient())
            {
                var xmlDocument = new XmlDocument();
                try
                {
                    var bytes = await webClient.DownloadDataTaskAsync(new Uri(_appCastUrl));
                    xmlDocument.Load(new MemoryStream(bytes));
                    availableUpdates.LoadFromXml(xmlDocument);
                    return availableUpdates.LatestVersion();
                }
                catch (Exception ex)
                {
                    _logger.Error(ex.Message);
                    return null;
                }
            }
        }

        public string GetAppCastEndPoint()
        {
            if (_useHockeyApp)
            {
                return $"https://rink.hockeyapp.net/api/2/apps/{_channel.ApplicationId}.rss";
            }
            //Docker AppCast EndPoint
            return $"https://editions-stage-us-east-1-150610-005505.s3.amazonaws.com/build-bucket/win/{_channel.BucketPath}/RELEASE";
        }
    }
}