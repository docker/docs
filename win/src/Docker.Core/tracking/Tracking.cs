using System;
using System.IO;

namespace Docker.Core.Tracking
{
    public class Tracking
    {
        private readonly ILogger _logger;
        private readonly Channel _channel;

        public string Id { get; }
        public bool IsEnabled { get; private set; }

        public Tracking(Channel channel) : this(channel, ReadOrCreate(Paths.TrackId), false)
        {
        }

        internal Tracking(Channel channel, string id, bool enabled)
        {
            _logger = new Logger(GetType());
            _channel = channel;
            Id = id;
            IsEnabled = enabled;
        }

        private bool OverrideIsTracking(bool isTracking)
        {
            return _channel == Channel.Beta || isTracking;
        }

        public void ChangeTo(bool newStatus)
        {
            var overridenStatus = OverrideIsTracking(newStatus);
            if (overridenStatus != IsEnabled)
            {
                _logger.Info($"Crash report and usage statistics are {(overridenStatus ? "enabled" : "disabled")}");
            }
            IsEnabled = overridenStatus;
        }

        internal static string ReadOrCreate(string filename)
        {
            if (File.Exists(filename))
            {
                var id = File.ReadAllText(filename);
                if (!string.IsNullOrEmpty(id))
                {
                    return id;
                }
            }

            var newId = Guid.NewGuid().ToString().ToUpper();
            File.WriteAllText(filename, newId);
            return newId;
        }
    }
}