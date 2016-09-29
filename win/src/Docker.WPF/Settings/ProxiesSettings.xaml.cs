using System.Windows;
using Docker.Core;

namespace Docker.WPF
{
    public partial class ProxiesSettings : IActivable
    {
        private readonly IActions _actions;
        private readonly ISettingsLoader _settingsLoader;

        private HttpProxySettingsViewModel Model
        {
            get { return DataContext as HttpProxySettingsViewModel; }
            set { DataContext = value; }
        }

        public ProxiesSettings(IActions actions, ISettingsLoader settingsLoader)
        {
            _actions = actions;
            _settingsLoader = settingsLoader;

            InitializeComponent();
        }

        public void Refresh(Settings settings)
        {
            Model = HttpProxySettingsViewModel.Load(settings, _actions, _settingsLoader);
        }

        private void OnApply(object sender, RoutedEventArgs e)
        {
            Model?.Apply();
        }
    }
}
