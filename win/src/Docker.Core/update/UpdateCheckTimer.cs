using System;
using System.Timers;
using Docker.Core.Update;

namespace Docker.Core.update
{
    public interface IUpdateCheckTimer
    {
        void Start(Action startingUpdate, Action upToDate);
        void CheckOnce(Action startingUpdate, Action upToDate);
    }

    public class UpdateCheckTimer : IUpdateCheckTimer, IDisposable
    {
        private const double TimerInterval = 6 * 3600 * 1000.0d; // in ms

        private readonly Timer _timer;
        private readonly IUpdater _updater;

        public UpdateCheckTimer(IUpdater updater)
        {
            _updater = updater;
            _timer = new Timer(TimerInterval);
        }

        public void CheckOnce(Action startingUpdate, Action upToDate)
        {
            _updater.CheckForUpdates(startingUpdate, upToDate);
        }

        public void Start(Action startingUpdate, Action upToDate)
        {
            _timer.Elapsed += (s, e) => CheckOnce(startingUpdate, upToDate);
            _timer.Start();
        }

        // https://msdn.microsoft.com/library/ms244737.aspx
        ~UpdateCheckTimer()
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
