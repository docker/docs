using System.Windows;
using System.Windows.Controls;
using System.Windows.Media;
using Docker.Core;

namespace Docker.WPF
{
    public partial class NetworkSettings : IActivable
    {
        private readonly IpHelper _ipHelper;
        private readonly IActions _actions;
        private readonly Brush _errorBrush;
        private readonly Brush _originalTextBoxBrush;

        private bool _loaded;

        public NetworkSettings(IpHelper ipHelper, IActions actions)
        {
            _ipHelper = ipHelper;
            _actions = actions;
            _errorBrush = new SolidColorBrush(Colors.Tomato);

            InitializeComponent();
            _originalTextBoxBrush = SwitchMaskTextBox.Background;
        }

        public void Refresh(Settings settings)
        {
            _loaded = false;
            try
            {
                SwitchAddressTextBox.Text = settings.SubnetAddress;
                SwitchMaskTextBox.Text = _ipHelper.GetMaskFromSize(settings.SubnetMaskSize);
                NameServerTextBox.Text = settings.NameServer;
                NameServerAutoButton.IsChecked = settings.UseDnsForwarder;
                NameServerFixedButton.IsChecked = !settings.UseDnsForwarder;
                SwitchAddressDefault.Content = $"default: {new Settings().SubnetAddress}";
                SwitchMaskDefault.Content = $"default: {_ipHelper.GetMaskFromSize(new Settings().SubnetMaskSize)}";

                CheckValues();
            }
            finally
            {
                _loaded = true;
            }
        }

        private void OkButton_Click(object sender, RoutedEventArgs e)
        {
            _actions.RestartVm(_ =>
            {
                _.SubnetAddress = SwitchAddressTextBox.Text;
                _.SubnetMaskSize = _ipHelper.GetMaskSize(SwitchMaskTextBox.Text);
                _.UseDnsForwarder = NameServerAutoButton.IsChecked ?? false;
                _.NameServer = NameServerTextBox.Text;
            });
        }

        private bool IsSwitchAddressValid()
        {
            return IsSwitchMaskValid() && _ipHelper.IsValidAdress(SwitchAddressTextBox.Text) && _ipHelper.IsAddressMatchWithMask(SwitchAddressTextBox.Text, SwitchMaskTextBox.Text);
        }

        private bool IsSwitchMaskValid()
        {
            var maskSize = _ipHelper.GetMaskSize(SwitchMaskTextBox.Text);
            return maskSize > 1 && maskSize < 31;
        }

        private bool IsNameServerValid()
        {
            return (NameServerAutoButton.IsChecked == true) || _ipHelper.IsValidAdress(NameServerTextBox.Text);
        }

        private void SwitchAdressTextBox_TextChanged(object sender, TextChangedEventArgs e)
        {
            if (!_loaded) return;
            CheckValues();
        }

        private void SwitchMaskTextBox_TextChanged(object sender, TextChangedEventArgs e)
        {
            if (!_loaded) return;
            CheckValues();
        }

        private void NameServerTextBox_TextChanged(object sender, TextChangedEventArgs e)
        {
            if (!_loaded) return;
            CheckValues();
        }

        private void NameServerAutoChanged(object sender, RoutedEventArgs e)
        {
            if (!_loaded) return;
            CheckValues();
        }

        private void CheckValues()
        {
            var nameServerIsValid = IsNameServerValid();
            var switchMaskIsValid = IsSwitchMaskValid();
            var switchAddressIsValid = IsSwitchAddressValid();

            NameServerTextBox.Background = nameServerIsValid ? _originalTextBoxBrush : _errorBrush;
            SwitchMaskTextBox.Background = switchMaskIsValid ? _originalTextBoxBrush : _errorBrush;
            SwitchAddressTextBox.Background = switchAddressIsValid ? _originalTextBoxBrush : _errorBrush;
            NameServerTextBox.IsEnabled = IsEnabled && (NameServerFixedButton.IsChecked == true);

            OkButton.IsEnabled = IsEnabled && switchMaskIsValid && switchAddressIsValid && nameServerIsValid;
        }
    }
}