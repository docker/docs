using System.Security.Principal;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using System.Windows.Media;
using Docker.Core;

namespace Docker.WPF.Credentials
{
    public partial class CredentialWindow
    {
        private readonly Brush _transparent;
        private readonly Brush _white;

        public CredentialWindow()
        {
            _transparent = new SolidColorBrush(Colors.Transparent);
            _white = new SolidColorBrush(Colors.White);

            InitializeComponent();

            Error.Content = "";
            Username.Text = WindowsIdentity.GetCurrent().Name;
            Password.Focus();

            PreviewKeyUp += (sender, e) =>
            {
                if (e.Key == Key.Escape)
                {
                    e.Handled = true;
                    Close();
                }
            };
        }

        public Credential Ask()
        {
            ShowDialog();

            return DialogResult == true
                ? new Credential(Username.Text, Password.Password)
                : null;
        }

        private void OkButton_Click(object sender, RoutedEventArgs e)
        {
            DialogResult = true;
        }

        private void CancelButton_Click(object sender, RoutedEventArgs e)
        {
            DialogResult = false;
        }

        private void UpdateButtons()
        {
            OkButton.IsEnabled = (Password.Password.Length > 0) && (Username.Text.Length > 0);
        }

        private void UsernameChanged(object sender, TextChangedEventArgs e)
        {
            Username.Background = Username.Text.Length > 0 ? _white : _transparent;
            UpdateButtons();
        }

        private void PasswordChanged(object sender, RoutedEventArgs e)
        {
            Password.Background = Password.Password.Length > 0 ? _white : _transparent;
            UpdateButtons();
        }
    }
}