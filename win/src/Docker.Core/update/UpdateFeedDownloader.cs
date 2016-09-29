using System;
using System.IO;
using System.Net;
using System.Threading.Tasks;
using System.Xml;

namespace Docker.Core.Update
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

        public UpdateFeedDownloader(Channel channel)
        {
            _logger = new Logger(GetType());
            _channel = channel;

        }

        public async Task<AvailableUpdate> DownloadLatestUpdateInfo()
        {
            var availableUpdates = new AvailableUpdateList();
            using (var webClient = new WebClient())
            {
                var xmlDocument = new XmlDocument();
                try
                {
                    var bytes = await webClient.DownloadDataTaskAsync(new Uri(GetAppCastEndPoint()));
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
            return _channel.EndPoint;
        }
    }
}