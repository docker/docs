using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Windows;
using System.Windows.Controls;
using Docker.Core;
using System.Windows.Forms.Integration;
using System.Windows.Input;
using System.Windows.Media.Imaging;
using Docker.Core.di;
using UserControl = System.Windows.Controls.UserControl;

namespace Docker.WPF
{
    public interface ISettings
    {
        void Open(string defaultPanelName = "");
    }

    public class SettingPanel
    {
        internal readonly string Name;
        internal readonly UserControl UserControl;

        public SettingPanel(string name, UserControl userControl)
        {
            Name = name;
            UserControl = userControl;
        }

        public override string ToString()
        {
            return Name;
        }
    }

    public class Status
    {
        public static readonly Status Running = new Status(true, "Docker is running", new Uri("../Images/green.png", UriKind.Relative));
        public static readonly Status Starting = new Status(false, "Docker is starting...", new Uri("../Images/orange.png", UriKind.Relative));
        public static readonly Status Failed = new Status(true, "Failed to start", new Uri("../Images/red.png", UriKind.Relative));
        public static readonly Status UploadingDiag = new Status(false, "Uploading Diagnostic...", new Uri("../Images/orange.png", UriKind.Relative));
        public static readonly Status UpdatingDrives = new Status(false, "Updating Drives...", new Uri("../Images/orange.png", UriKind.Relative));

        public readonly bool IsEnabled;
        public readonly string Text;
        public readonly Uri Image;

        private Status(bool isEnabled, string text, Uri image)
        {
            IsEnabled = isEnabled;
            Text = text;
            Image = image;
        }
    }

    public partial class SettingsWindow : ISettings
    {
        private readonly DevMode _devMode;
        private readonly ISettingsLoader _settingsLoader;
        private readonly Singletons _singletons;
        private readonly List<SettingPanel> _tabs;
        private readonly SettingPanel _devPanel;

        private Status _status;

        public SettingsWindow(
            ITaskQueue taskQueue, DevMode devMode, ISettingsLoader settingsLoader,
            DevActions devActions, IAppShutdown appShutdown, Singletons singletons)
        {
            _devMode = devMode;
            _settingsLoader = settingsLoader;
            _singletons = singletons;
            _tabs = new List<SettingPanel>();
            _devPanel = new SettingPanel("Dev", new PinataSettings(devActions, appShutdown));
            _status = Status.Starting;

            InitializeComponent();

            taskQueue.Started += () => SetStatus(Status.Starting);
            taskQueue.Ended += exceptionOccured => SetStatus(exceptionOccured ? Status.Failed : Status.Running);
        }

        public void Open(string defaultPanelName = "")
        {
            if (_tabs.Count == 0)
            {
                _tabs.Add(new SettingPanel("General", _singletons.Get<GeneralSettings>()));
                _tabs.Add(new SettingPanel("Shared Drives", _singletons.Get<SharedDrivesSettings>()));
                _tabs.Add(new SettingPanel("Advanced", _singletons.Get<AdvancedSettings>()));
                _tabs.Add(new SettingPanel("Network", _singletons.Get<NetworkSettings>()));
                _tabs.Add(new SettingPanel("Proxies", _singletons.Get<ProxiesSettings>()));
                _tabs.Add(new SettingPanel("Docker Daemon", _singletons.Get<DaemonSettings>()));
                // _tabs.Add(new SettingPanel("Kernel", _singletons.Get<KernelSettings>()));
                _tabs.Add(new SettingPanel("Diagnose & Feedback", _singletons.Get<FeedbackSettings>()));
                _tabs.Add(new SettingPanel("Reset", _singletons.Get<ResetSettings>()));

                foreach (var tab in _tabs)
                {
                    PanelList.Items.Add(tab);
                }

                PanelList.SelectionChanged += PanelList_OnSelectionChanged;
                PanelList.SelectedIndex = 0;
            }

            ElementHost.EnableModelessKeyboardInterop(this);
            if (_devMode.On)
            {
                if (!_tabs.Contains(_devPanel))
                {
                    _tabs.Add(_devPanel);
                    PanelList.Items.Add(_devPanel);
                }
            }
            else
            {
                PanelList.Items.Remove(_devPanel);
                _tabs.RemoveAll(panel => ReferenceEquals(panel, _devPanel));
            }

            var settings = _settingsLoader.Load();
            foreach (var settingPanel in _tabs)
            {
                RefreshTab(settingPanel.UserControl, settings, true);
            }

            if (!"".Equals(defaultPanelName))
            {
                for (var i = 0; i < _tabs.Count; i++)
                {
                    if (_tabs[i].Name.Equals(defaultPanelName))
                    {
                        PanelList.SelectedIndex = i;
                    }
                }
            }

            Show();
            Activate();
        }

        protected override void OnClosing(CancelEventArgs e)
        {
            e.Cancel = true;
            Hide();
        }

        internal Status SetStatus(Status status)
        {
            var previousStatus = _status;
            _status = status;

            Dispatcher.BeginInvoke((Action)(() =>
            {
                Cursor = status.IsEnabled ? null : Cursors.Wait;

                DockerLoading.Text = status.Text;
                StatusImage.Source = new BitmapImage(status.Image);

                var settings = _settingsLoader.Load();
                foreach (var settingPanel in _tabs)
                {
                    RefreshTab(settingPanel.UserControl, settings, status.IsEnabled);
                }
            }));

            return previousStatus;
        }

        private void PanelList_OnSelectionChanged(object sender, SelectionChangedEventArgs e)
        {
            var selectedPanel = ((SettingPanel)e.AddedItems[0]).UserControl;
            if (!PanelGrid.Children.Contains(selectedPanel))
            {
                PanelGrid.Children.Add(selectedPanel);
            }

            var settings = _settingsLoader.Load();
            foreach (var child in PanelGrid.Children)
            {
                var panel = (UserControl)child;

                panel.Visibility = ReferenceEquals(panel, selectedPanel) ? Visibility.Visible : Visibility.Hidden;
                RefreshTab(panel, settings, _status.IsEnabled);
            }
        }

        private void RefreshTab(UIElement control, Settings settings, bool isEnabled)
        {
            control.IsEnabled = control is IAlwaysEnabled || isEnabled;

            if (isEnabled)
            {
                (control as IActivable)?.Refresh(settings);
            }
        }
    }
}
