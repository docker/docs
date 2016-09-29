using System.Windows;
using Docker.WPF;

namespace Docker
{
    public class AppShutdown : IAppShutdown
    {
        private readonly IActions _actions;
        private readonly Systray _systray;

        public AppShutdown(Systray systray, IActions actions)
        {
            _systray = systray;
            _actions = actions;
        }

        public void Shutdown()
        {
            _systray.StartAnimation();
            foreach (var window in Application.Current.Windows)
            {
                (window as Window)?.Close();
            }
            _actions.ShutdownVm();
        }
    }
}