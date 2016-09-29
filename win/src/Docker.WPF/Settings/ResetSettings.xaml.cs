using System.Windows;
using Docker.Core;

namespace Docker.WPF
{
    public partial class ResetSettings : IActivable
    {
        private readonly IActions _actions;
        private readonly IToolboxMigration _toolboxMigration;

        public ResetSettings(IActions actions, IToolboxMigration toolboxMigration)
        {
            _actions = actions;
            _toolboxMigration = toolboxMigration;

            InitializeComponent();
        }

        public void Refresh(Settings settings)
        {
            ImportToolboxButton.IsEnabled = IsEnabled && DefaultMachineExists();
        }

        private bool DefaultMachineExists()
        {
            return _toolboxMigration.DefaultMachineExists;
        }

        private void ResetButton_Click(object sender, RoutedEventArgs e)
        {
            if (ResetToDefaultBox.ShowConfirm())
            {
                _actions.ResetToDefault();
            }
        }

        private void RestartButton_Click(object sender, RoutedEventArgs e)
        {
            if (RestartVmBox.ShowConfirm())
            {
                _actions.RestartVm();
            }
        }

        private void MigrateButton_Click(object sender, RoutedEventArgs e)
        {
            if (ImportFromToolBox.ShowConfirm())
            {
                var defaultMachinePath = _toolboxMigration.GetMachineVolumePath("default");

                _actions.MigrateVolume("default", defaultMachinePath);
            }
        }
    }
}
