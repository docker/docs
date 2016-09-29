using System;
using System.ComponentModel;
using System.Drawing;
using System.Threading;
using System.Windows.Forms;
using System.Windows.Threading;
using Docker.Core;
using Docker.Properties;
using Timer = System.Windows.Forms.Timer;

namespace Docker
{
    public class Systray : ApplicationContext
    {
        private readonly Logger _logger;
        private readonly Icon[] _animationFrames;
        private readonly Timer _animationTimer;

        private IContainer _container;
        private NotifyIcon _notifyIcon;
        private Action _onOpen;

        private int CurrentFrame { get; set; }

        public Systray()
        {
            _logger = new Logger(GetType());
            _animationFrames = new[] {
                Resources.systray_icon_inverted_frame0,
                Resources.systray_icon_inverted_frame1,
                Resources.systray_icon_inverted_frame2,
                Resources.systray_icon_inverted_frame3,
                Resources.systray_icon_inverted_frame4,
                Resources.systray_icon_inverted_frame5,
                Resources.systray_icon_inverted_frame4,
                Resources.systray_icon_inverted_frame3,
                Resources.systray_icon_inverted_frame2,
                Resources.systray_icon_inverted_frame1
            };
            _animationTimer = new Timer { Interval = 150 };
        }

        public void Initialize(Action onOpen)
        {
            _onOpen = onOpen;
            _container = new Container();
            _notifyIcon = new NotifyIcon(_container)
            {
                ContextMenuStrip = new ContextMenuStrip(),
                Icon = Resources.systray_icon_inverted_frame0,
                Text = @"Docker is starting...",
                Visible = true
            };

            _notifyIcon.ContextMenuStrip.Opening += Opening;

            _animationTimer.Tick += AnimationTimer_Tick;
            StartAnimation();
        }

        private void AnimationTimer_Tick(object sender, EventArgs eventArgs)
        {
            try
            {
                CurrentFrame = (CurrentFrame + 1) % _animationFrames.Length;
                _notifyIcon.Icon = _animationFrames[CurrentFrame];
            }
            catch
            {
                // Ignore
            }
        }

        public void StartAnimation()
        {
            new Thread(() =>
            {
                _animationTimer.Start();
                Dispatcher.Run();
            }).Start();
        }

        public void StopAnimation(Icon icon = null)
        {
            new Thread(() =>
            {
                _animationTimer.Stop();
                try
                {
                    _notifyIcon.Icon = icon ?? _animationFrames[0];
                }
                catch
                {
                    // Ignore
                }
                Dispatcher.Run();
            }).Start();
        }

        public void SetStatus(string title, Icon icon, bool? useAnimation)
        {
            try
            {
                switch (useAnimation)
                {
                    case true:
                        StartAnimation();
                        break;
                    case false:
                        StopAnimation(icon);
                        break;
                }
            }
            catch (Exception exception)
            {
                _logger.Error($"Unable to change icon: {exception.Message}");
            }

            try
            {
                if (title.Length >= 64)
                {
                    _notifyIcon.Text = title.Substring(0, 60) + @"...";
                }
                else
                {
                    _notifyIcon.Text = title;
                }
            }
            catch (Exception exception)
            {
                _logger.Error($"Unable to change systray tooltip: {exception.Message}");
            }
        }

        public void NotifyStatus(string title, string details, ToolTipIcon toastIcon, Icon icon, bool? useAnimation)
        {
            SetStatus(title, icon, useAnimation);

            try
            {
                _notifyIcon.ShowBalloonTip(0, title, details, toastIcon);
            }
            catch (Exception exception)
            {
                _logger.Error($"Unable to show balloon tip: {exception.Message}");
            }
        }

        private void Opening(object sender, CancelEventArgs e)
        {
            try
            {
                Cursor.Current = Cursors.WaitCursor;

                _notifyIcon.ContextMenuStrip.Items.Clear();

                _onOpen();

                // The following two lines are the only way I found to make the
                // balloon tooltip disappear when the ContextMenu is displayed
                _notifyIcon.Visible = false;
                _notifyIcon.Visible = true;

                e.Cancel = false;
            }
            catch (Exception ex)
            {
                _logger.Error($"Failed opening Docker menu: {ex.Message}");
            }
            finally
            {
                Cursor.Current = Cursors.Default;
            }
        }

        private static EventHandler CatchError(ToolStripItem item, Action action)
        {
            return delegate {
                try
                {
                    action();
                }
                catch (Exception ex)
                {
                    throw new DockerException($"Failed executing {item.Text}: {ex.Message}");
                }
            };
        }

        public void AddItem(string label, Action action)
        {
            var item = new ToolStripMenuItem(label);

            item.Click += CatchError(item, action);

            _notifyIcon.ContextMenuStrip.Items.Add(item);
        }

        public void AddSeparator()
        {
            _notifyIcon.ContextMenuStrip.Items.Add("-");
        }

        protected override void Dispose(bool disposing)
        {
            if (!disposing) return;

            _container?.Dispose();
            _animationTimer?.Dispose();
        }

        protected override void ExitThreadCore()
        {
            _notifyIcon.Visible = false; // should remove lingering tray icon
            base.ExitThreadCore();
        }
    }
}