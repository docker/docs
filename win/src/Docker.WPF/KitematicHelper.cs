using System;
using System.Diagnostics;
using System.IO;
using Docker.Core;

namespace Docker.WPF
{
    public class KitematicHelper
    {
        private readonly ILogger _logger;

        public KitematicHelper()
        {
            _logger = new Logger(GetType());
        }

        public void Open()
        {
            var exists = File.Exists(Paths.KitematicPath);

            if (exists && !IsKitematicOutdated())
            {
                Process.Start(Paths.KitematicPath);
            }
            else if ((!exists && DownloadKitematicBox.ShowConfirm()) || (exists && UpdateKitematicBox.ShowConfirm()))
            {
                Process.Start(Urls.KitematicDownload);
            }
        }

        private bool IsKitematicOutdated()
        {
            var version = FileVersionInfo.GetVersionInfo(Paths.KitematicPath).ProductVersion;

            _logger.Info($"Found Kitematic version {version}");

            return IsOutdated(version);
        }

        internal static bool IsOutdated(string version)
        {
            try
            {
                var exeVersion = new System.Version(version);
                var minimumVersion = new System.Version("0.12.0");
                var diff = exeVersion.CompareTo(minimumVersion);

                return diff < 0;
            }
            catch (FormatException)
            {
                return true;
            }
            catch (ArgumentException)
            {
                return true;
            }
        }
    }
}