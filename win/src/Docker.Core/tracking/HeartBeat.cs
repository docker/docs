using System;
using System.Timers;

namespace Docker.Core.Tracking
{
    public class Heartbeat : IDisposable
    {
        private const double TimerInterval = 3600*1000.0d; // in ms

        private readonly Timer _timer;
        private readonly IAnalytics _analytics;

        public Heartbeat(IAnalytics analytics)
        {
            _analytics = analytics;
            _timer = new Timer(TimerInterval);
            _timer.Elapsed += (s, e) => _analytics.Track(AnalyticEvent.Heartbeat);
        }

        public void Start()
        {
            _analytics.Track(AnalyticEvent.Heartbeat);
            _timer.Start();
        }

        // https://msdn.microsoft.com/library/ms244737.aspx
        ~Heartbeat()
        {
            Dispose(false);
        }

        public void Dispose()
        {
            Dispose(true);
            GC.SuppressFinalize(this);
        }

        protected virtual void Dispose(bool disposing)
        {
            if (disposing)
            {
                _timer.Dispose();
            }
        }
    }
}