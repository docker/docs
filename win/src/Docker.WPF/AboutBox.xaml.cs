using System;
using Docker.Core;
using System.Windows;
using System.Windows.Input;
using System.Windows.Media.Imaging;

namespace Docker.WPF
{
    public partial class AboutBox
    {
        private readonly ICmd _cmd;
        private readonly Code _code;
        private readonly DevMode _mode;

        public AboutBox(IVersion version, Git git, ICmd cmd, Code code, DevMode mode, Channel channel)
        {
            _cmd = cmd;
            _code = code;
            _mode = mode;

            InitializeComponent();
            AboutImage.Source = new BitmapImage(new Uri(channel.IsStable() ? "Images/about-docker.png" : "Images/about-docker-beta.png", UriKind.Relative));
            VersionLabel.Content = $"Version {version.ToHumanStringWithBuildNumber()}";
            Sha1Label.Content = git.ShortSha1();
        }

        protected override void OnClosing(System.ComponentModel.CancelEventArgs cancelEvent)
        {
            cancelEvent.Cancel = true;
            Hide();
        }

        private void Acknowledgments_Click(object sender, RoutedEventArgs e)
        {
            OpenWrite(Paths.InResourcesDir("OSS-LICENSES.txt"));
        }

        private void License_Click(object sender, RoutedEventArgs e)
        {
            OpenWrite(Paths.InResourcesDir("LICENSE.rtf"));
        }

        private void OpenWrite(string filename)
        {
            _cmd.RunWindowed("write.exe", $"\"{filename}\"");
        }

        private void OnKeyUp(object sender, KeyEventArgs e)
        {
            _code.OnKeyUp(e, () =>
            {
                _mode.On = !_mode.On;
                var imagePath = _mode.On ? "Images/pinata.png" : "Images/about-docker-beta.png";
                AboutImage.Source = new BitmapImage(new Uri(imagePath, UriKind.Relative));
            });
        }
    }
}