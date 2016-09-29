using System;
using System.Windows;
using Docker.Core;

namespace Docker.WPF
{
    public partial class AdvancedSettings : IActivable
    {
        private readonly IActions _actions;

        public AdvancedSettings(IActions actions)
        {
            _actions = actions;

            InitializeComponent();
        }

        public void Refresh(Settings settings)
        {
            MemorySlider.Minimum = 1024;
            MemorySlider.Value = settings.VmMemory;
            MemorySlider.Maximum = MaxMemory();

            CpuSlider.Value = settings.VmCpus;
            CpuSlider.Maximum = MaxCpus();
        }

        private int MaxMemory()
        {
            var maxMemory = Env.MaxMemory - 2048;
            var vmMaxMemory = maxMemory - maxMemory % 256;

            return Math.Max(vmMaxMemory, 2048);
        }

        private int MaxCpus()
        {
            return Environment.ProcessorCount;
        }

        private void OkButton_Click(object sender, RoutedEventArgs e)
        {
            _actions.RestartVm(_ =>
            {
                _.VmCpus = Convert.ToInt32(CpuSlider.Value);
                _.VmMemory = Convert.ToInt32(MemorySlider.Value);
            });
        }

        private void CPUSlider_ValueChanged(object sender, RoutedPropertyChangedEventArgs<double> e)
        {
            CpuTextBlock.Content = $"CPUs: {CpuSlider.Value}";
        }

        private void MemorySlider_ValueChanged(object sender, RoutedPropertyChangedEventArgs<double> e)
        {
            MemoryTextBlock.Content = $"Memory: {MemorySlider.Value} MB";
        }
    }
}
