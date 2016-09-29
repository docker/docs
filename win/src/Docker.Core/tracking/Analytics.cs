using System;
using System.Diagnostics;
using System.IO;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Runtime.Serialization;
using System.Text;
using System.Threading.Tasks;
using Newtonsoft.Json;

namespace Docker.Core.Tracking
{
    public interface IAnalytics
    {
        void Track(AnalyticEvent analyticEvent);
    }

    public class SegmentApi : IAnalytics
    {
        private readonly Logger _logger;
        private readonly Channel _channel;
        private readonly Tracking _tracking;
        private readonly IVersion _version;

        public SegmentApi(Channel channel, Tracking tracking, IVersion version)
        {
            _logger = new Logger(GetType());
            _channel = channel;
            _tracking = tracking;
            _version = version;
        }

        public async void Track(AnalyticEvent analyticEvent)
        {
            if (!_tracking.IsEnabled && !analyticEvent.IsCore)
            {
                _logger.Info($"Not tracking: {analyticEvent.Name}");
                return;
            }

            try
            {
                _logger.Info($"Usage statistic: {analyticEvent.Name}");

                if (!await Track(analyticEvent.Name))
                {
                    _logger.Error($"Failed to track event: {analyticEvent.Name}");
                }
            }
            catch (Exception ex)
            {
                _logger.Error($"Failed to track event: {ex.Message}");
            }
        }

        private async Task<bool> Track(string message)
        {
            var data = new SegmentMessage(_tracking.IsEnabled, _channel.Name, _tracking.Id, message, _version).ToJson();
            var content = new StringContent(data, Encoding.UTF8, "application/json");

            using (var client = NewHttpClient(_channel.AnalyticsToken))
            {
                var response = await client.PostAsync("https://api.segment.io/v1/track", content);

                return response.IsSuccessStatusCode;
            }
        }

        private static HttpClient NewHttpClient(string token)
        {
            return new HttpClient
            {
                Timeout = TimeSpan.FromSeconds(3),
                DefaultRequestHeaders = { Authorization = new AuthenticationHeaderValue("Basic", Convert.ToBase64String(Encoding.UTF8.GetBytes($"{token}:"))) }
            };
        }
    }

    [DataContract]
    public class SegmentMessage
    {
        [DataMember(Name = "userId")]
        private string _userId;
        [DataMember(Name = "event")]
        private string _eventIdentifier;
        [DataMember(Name = "properties")]
        private SegmentProperties _properties;

        public SegmentMessage(bool isTrackingEnabled, string channel, string userDistinctId, string eventMessage, IVersion version)
        {
            _userId = userDistinctId;
            _eventIdentifier = eventMessage;
            _properties = new SegmentProperties(version, channel, isTrackingEnabled);
        }

        public string Serialize()
        {
            return ToBase64(ToJson());
        }

        public string ToJson()
        {
            var sw = new StringWriter();
            new JsonSerializer().Serialize(sw, this);
            return sw.ToString();
        }

        private static string ToBase64(string plainText)
        {
            return Convert.ToBase64String(Encoding.UTF8.GetBytes(plainText));
        }
    }

    [DataContract]
    public class SegmentProperties
    {
        [DataMember(Name = "os")]
        public readonly string Os;
        [DataMember(Name = "app major version")]
        public readonly int AppMajorVersion;
        [DataMember(Name = "app minor version")]
        public readonly int AppMinorVersion;
        [DataMember(Name = "app patch version")]
        public readonly int AppPatchVersion;
        [DataMember(Name = "app version name")]
        public readonly string AppVersionName;
        [DataMember(Name = "channel")]
        public readonly string Channel;
        [DataMember(Name = "os major version", EmitDefaultValue = false)]
        public readonly string OsMajorVersion;
        [DataMember(Name = "os minor version", EmitDefaultValue = false)]
        public readonly string OsMinorVersion;
        [DataMember(Name = "os patch version", EmitDefaultValue = false)]
        public readonly string OsPatchVersion;
        [DataMember(Name = "os language", EmitDefaultValue = false)]
        public readonly string OsLanguage;

        public SegmentProperties(IVersion version, string channel, bool addOsInformation)
        {
            Os = "windows";
            AppMajorVersion = version.Major();
            AppMinorVersion = version.Minor();
            AppPatchVersion = version.Build();
            AppVersionName = SafeToLower(version.ToCrossPlatformName());
            Channel = SafeToLower(channel);
            if (addOsInformation)
            {
                OsMajorVersion = SafeToLower(Env.Os.Name);
                OsMinorVersion = Env.Os.ReleaseId;
                OsPatchVersion = Env.Os.BuildNumber;
                OsLanguage = SafeToLower(Env.Os.Language);
            }
        }

        private String SafeToLower(String text)
        {
            return text == null ? "" : text.ToLower();
        }
    }
}