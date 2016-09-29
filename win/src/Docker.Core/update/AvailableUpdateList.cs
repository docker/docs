using System.Xml;
using System.Collections.Generic;

namespace Docker.Core.Update
{
    public class AvailableUpdateList
    {
        private readonly List<AvailableUpdate> _versions;

        public AvailableUpdateList()               
        {
            _versions = new List<AvailableUpdate>();
        }

        public void LoadFromXml(XmlDocument xmlDocument)
        {
            var itemNodeList = xmlDocument.DocumentElement?.SelectNodes("/rss/channel/item");
            if (itemNodeList == null)
            {
                return;
            }
            foreach (XmlNode itemNode in itemNodeList)
            {
                var url = itemNode.SelectSingleNode("enclosure")?.Attributes?["url"].Value;
                var version = itemNode.SelectSingleNode("enclosure")?.Attributes?["sparkle:version"].Value;
                var shortVersion = itemNode.SelectSingleNode("enclosure")?.Attributes?["sparkle:shortVersionString"].Value;
                var notes = itemNode.SelectSingleNode("description")?.InnerText;
                _versions.Add(new AvailableUpdate(version, shortVersion, url, notes));
            }
        }

        public AvailableUpdate LatestVersion()
        {
            // HockeyApp latest version is always the first one
            // Once we know more about new FrenchBen RSS feed, we can update this.
            return _versions.Count == 0 ? null : _versions[0];
        }
    }
}