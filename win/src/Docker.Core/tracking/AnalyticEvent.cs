namespace Docker.Core.Tracking
{
    public class AnalyticEvent
    {
        public static readonly AnalyticEvent AppLaunched = new AnalyticEvent("appLaunched");
        public static readonly AnalyticEvent InstallShowWelcomePopup = new AnalyticEvent("installShowWelcomePopup");
        public static readonly AnalyticEvent AppRunning = new AnalyticEvent("appRunning");
        public static readonly AnalyticEvent Heartbeat = new AnalyticEvent("heartbeat");

        public string Name { get; }
        public bool IsCore => Heartbeat == this;

        internal AnalyticEvent(string name)
        {
            Name = name;
        }
    }
}
