using System;
using System.ComponentModel;
using System.Linq;
using System.Windows;
using System.Windows.Input;
using System.Windows.Media.Animation;

namespace Docker.WPF
{
    public interface IWelcomeShower
    {
        void Show(ISettings settings);
    }

    public class NopWelcomeShower : IWelcomeShower
    {
        public void Show(ISettings settings)
        {
        }
    }

    public partial class WelcomeWindow : IWelcomeShower
    {
        private ISettings _settings;

        public WelcomeWindow()
        {
            InitializeComponent();

            string[] texts = { "docker ps", "docker images", "docker version", "docker info", "docker run hello-world" };
            const int initialDelay = 5000;
            const int keystrokeFrameLength = 150;
            const int commandFrameLength = 3000;
            var animationLength = texts.Sum(text => text.Length * keystrokeFrameLength + commandFrameLength);
            var timeline = new StringAnimationUsingKeyFrames
            {
                BeginTime = TimeSpan.FromMilliseconds(initialDelay),
                RepeatBehavior = RepeatBehavior.Forever,
                Duration = new Duration(TimeSpan.FromMilliseconds(animationLength))
            };

            var elapsed = 0;
            foreach (var text in texts)
            {
                var displayedString = "";
                foreach (var character in text)
                {
                    displayedString += character;
                    timeline.KeyFrames.Add(new DiscreteStringKeyFrame("> " + displayedString, KeyTime.FromTimeSpan(TimeSpan.FromMilliseconds(elapsed))));
                    elapsed += keystrokeFrameLength;
                }
                timeline.KeyFrames.Add(new DiscreteStringKeyFrame("> " + text, KeyTime.FromTimeSpan(TimeSpan.FromMilliseconds(elapsed)))); // this call can be avoided by checking in the loop above if this is the last char.
                elapsed += commandFrameLength;
            }
            PromptText.BeginAnimation(ContentProperty, timeline);

            LaunchWaveAnim();
            PreviewKeyUp += (sender, e) =>
            {
                if (e.Key == Key.Escape)
                {
                    e.Handled = true;
                    Hide();
                }
            };
        }

        private void LaunchWaveAnim()
        {
            var sb = FindResource("WavesSb") as Storyboard;
            sb?.Begin();
        }

        protected override void OnClosing(CancelEventArgs e)
        {
            e.Cancel = true;
            Hide();
        }

        private void GotIt_Click(object sender, RoutedEventArgs e)
        {
            Hide();
        }

        private void OnPrivacyCheck(object sender, RoutedEventArgs routedEventArgs)
        {
            Hide();
            _settings.Open("General");
        }

        public void Show(ISettings settings)
        {
            _settings = settings;

            Dispatcher.Invoke(Show);
        }
    }
}
