using System.IO;

namespace Docker.Core
{
    public class Channel
    {
        private const string SegmentBetaToken = "4EFdOmMsMGUinXphTQYnKXB9wXk4vy0T";
        private const string SegmentDevelopmentToken = "qJqnsOy2Bz9n63qFOaSO0mKTv3sIwEls";

        public const string BugSnagPublicApiKey = "1061755028b09100c080d55dfd4012b3";
        public const string BugSnagDevelopmentApiKey = "ffad13852d1a1452af3674f61da6303e";

        // public
        public static readonly Channel Stable = new Channel("Stable", "https://download.docker.com/win/stable/appcast.xml", SegmentBetaToken, BugSnagPublicApiKey);
        public static readonly Channel Beta = new Channel("Beta", "https://download.docker.com/win/beta/appcast.xml", SegmentBetaToken, BugSnagPublicApiKey);

        // staging
        public static readonly Channel Test = new Channel("Test", "https://download-stage.docker.com/win/test/appcast.xml", SegmentDevelopmentToken, BugSnagDevelopmentApiKey);
        public static readonly Channel Master = new Channel("Master", "https://download-stage.docker.com/win/master/appcast.xml", SegmentDevelopmentToken, BugSnagDevelopmentApiKey);

        public string Name;
        public string EndPoint;
        public string AnalyticsToken;
        public string BugsnagApiKey;

        public Channel(string name, string endPoint, string analyticsToken, string bugsnagApiKey)
        {
            Name = name;
            EndPoint = endPoint;
            AnalyticsToken = analyticsToken;
            BugsnagApiKey = bugsnagApiKey;
        }

        public override string ToString()
        {
            return Name;
        }

        public bool IsStaging()
        {
            return Beta != this && Stable != this;
        }

        public bool IsStable()
        {
            return Stable == this;
        }

        public static Channel Load()
        {
            var updateChannelFile = Paths.InResourcesDir("UpdateChannel");
            if (!File.Exists(updateChannelFile))
            {
                return Master;
            }

            var updateChannelFileContent = File.ReadAllLines(updateChannelFile)[0];
            return Parse(updateChannelFileContent);
        }

        public void ChangeTo(string text)
        {
            var channel = Parse(text);
            Name = channel.Name;
            EndPoint = channel.EndPoint;
            AnalyticsToken = channel.AnalyticsToken;
            BugsnagApiKey = channel.BugsnagApiKey;
            File.WriteAllText(Paths.InResourcesDir("UpdateChannel"), Name);
        }

        internal static Channel Parse(string text)
        {
            switch (text?.ToLower())
            {
                case "stable":
                    return Stable;
                case "beta":
                    return Beta;
                case "test":
                    return Test;
                default:
                    return Master;
            }
        }
    }
}
