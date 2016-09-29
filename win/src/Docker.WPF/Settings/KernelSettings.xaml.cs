using System;
using System.Windows;
using Docker.Core;

namespace Docker.WPF
{
    public partial class KernelSettings : IActivable
    {
        private readonly IActions _actions;

        public KernelSettings(IActions actions)
        {
            _actions = actions;

            InitializeComponent();
        }

        public void Refresh(Settings settings)
        {
            SysCtlConfText.Text = settings.SysCtlConf;
        }

        private void Apply(object sender, RoutedEventArgs eventArgs)
        {
            try
            {
                _actions.RestartVm(_ => _.SysCtlConf = SysCtlConfText.Text.Replace("\r", ""));
            }
            catch (Exception exception)
            {
                ErrorText.Text = exception.Message;
                ErrorText.Visibility = Visibility.Visible;
            }
        }
    }
}
