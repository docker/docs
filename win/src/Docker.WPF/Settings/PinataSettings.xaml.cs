using System;
using System.Diagnostics;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Navigation;
using Docker.Core;
using Microsoft.Win32;

namespace Docker.WPF
{
    public partial class PinataSettings : IActivable, IAlwaysEnabled
    {
        private readonly string[] _channels = {"Stable", "Beta", "Test", "Master"};
        private readonly DevActions _devActions;
        private readonly IAppShutdown _appShutdown;

        private bool _refreshing;

        public PinataSettings(DevActions devActions, IAppShutdown appShutdown)
        {
            _devActions = devActions;
            _appShutdown = appShutdown;

            InitializeComponent();
            ChannelCombo.ItemsSource = _channels;
        }

        public void Refresh(Settings settings)
        {
            _refreshing = true;

            var currentChannel = _devActions.GetChannel().ToString();
            for (var i = 0; i < _channels.Length; i++)
            {
                if (_channels[i].Equals(currentChannel))
                {
                    ChannelCombo.SelectedIndex = i;
                }
            }

            SetupAppCastUrl();

            _refreshing = false;
        }

        private void SetupAppCastUrl()
        {
            if (_devActions.GetAppCastUrl() != "")
            {
                AppCastHyperlink.NavigateUri = new Uri(_devActions.GetAppCastUrl());
            }
            AppCastHyperlinkText.Text = _devActions.GetAppCastUrl();
        }

        private void OnResourceFolder(object sender, RoutedEventArgs e)
        {
            _devActions.OpenResourceFolder();
        }

        private void OnSettingsFolder(object sender, RoutedEventArgs e)
        {
            _devActions.OpenSettingsFolder();
        }

        private void OnOpenAppCast(object sender, RequestNavigateEventArgs e)
        {
            Process.Start(new ProcessStartInfo(e.Uri.AbsoluteUri));
            e.Handled = true;
        }

        private void OnCheckForUpdate(object sender, RoutedEventArgs e)
        {
            _devActions.CheckForUpdate();
        }

        private void OnSendCrashReport(object sender, RoutedEventArgs e)
        {
            _devActions.OpenSendCrashWindow();
        }

        private void ChannelChanged(object sender, SelectionChangedEventArgs e)
        {
            if (_refreshing)
            {
                return;
            }
            var comboBox = sender as ComboBox;
            var value = comboBox?.SelectedItem as string;
            _devActions.ChangeChannelTo(value);
            SetupAppCastUrl();
        }

        private void OnUpdate(object sender, RoutedEventArgs e)
        {
            var openFileDialog = new OpenFileDialog
            {
                Filter = "Installer (*.msi)|*.msi",
                FilterIndex = 1
            };

            if (openFileDialog.ShowDialog() == true)
            {
                InstallUpdateWindow.RunMsi(openFileDialog.FileName, () =>
                {
                    _appShutdown.Shutdown();
                });
            }
        }
    }
}