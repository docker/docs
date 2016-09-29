using System;
using static System.Int32;

namespace Docker.Core.Update
{
    public class AvailableUpdate
    {
        public readonly string Version;
        public readonly string ShortVersion;
        public readonly string DownloadUrl;
        public readonly string Notes;

        public AvailableUpdate(string version, string shortVersion, string downloadUrl, string notes)
        {
            Version = version;
            ShortVersion = shortVersion;
            DownloadUrl = downloadUrl;
            Notes = notes;
        }

        public override string ToString()
        {
            return $"Version: {Version}, Notes: {Notes}, ShortVersion: {ShortVersion}, DownloadUrl: {DownloadUrl}";
        }

        public int BuildNumber()
        {
            var buildString = Version;
            if (buildString == null)
            {
                return -1;
            }
            var indexDot = buildString.LastIndexOf(".", StringComparison.Ordinal);
            if (indexDot != -1)
            {
                buildString = buildString.Substring(indexDot + 1);
            }
            int remoteVersionBuild;
            if (TryParse(buildString, out remoteVersionBuild))
            {
                return remoteVersionBuild;
            }
            return -1;
        }

        public string ToHumanStringWithBuildNumber()
        {
            return $"{ShortVersion} (build: {Version})";
        }
    }
}