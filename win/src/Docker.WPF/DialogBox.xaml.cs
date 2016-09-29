using System;
using System.IO;
using System.Windows;
using System.Windows.Documents;
using System.Windows.Media.Imaging;
using Docker.Core;

namespace Docker.WPF
{
    public static class NotEnoughtMemoryBox
    {
        public static void ShowOk()
        {
            new DialogBox
            {
                Title = "Failed to start Docker",
                Message = {Text = "Not enough memory to start Docker"},
                Description = "You are trying to start Docker but you don't have enough memory.\nFree some memory or change your settings."
            }.ShowOk();
        }
    }

    public static class ResetToDefaultBox
    {
        public static bool ShowConfirm()
        {
            return new DialogBox
            {
                Title = "Reset to factory defaults",
                Message = {Text = "Reset Docker to factory defaults"},
                Description = "You are about to reset Docker. A factory reset destroys all Docker containers and images local to the machine, and restores the application to its original state, as when it was first installed.\n\nAre you sure you want to continue?",
                OkButton = {Content = "Reset"}
            }.ShowYesNo();
        }
    }

    public static class RestartVmBox
    {
        public static bool ShowConfirm()
        {
            return new DialogBox
            {
                Title = "Restart Docker",
                Message = {Text = "Restart Docker"},
                Description = "You are about to restart Docker. A restart stops all running containers. No data will be lost otherwise.\n\nAre you sure you want to continue?",
                OkButton = {Content = "Restart"}
            }.ShowYesNo();
        }
    }

    public static class DownloadKitematicBox
    {
        public static bool ShowConfirm()
        {
            return new DialogBox
            {
                Title = "Kitematic",
                Message = {Text = "Download Kitematic"},
                Description = $"Kitematic is compatible with Docker for Windows and can be used as a graphical interface to manage your Docker containers. You can download it. Then make sure you install it in {Path.GetDirectoryName(Paths.KitematicPath)}",
                Image = {Source = new BitmapImage(new Uri("Images/kitematic.png", UriKind.Relative))},
                OkButton = {Content = "Download"}
            }.ShowYesNo();
        }
    }

    public static class UpdateKitematicBox
    {
        public static bool ShowConfirm()
        {
            return new DialogBox
            {
                Title = "Kitematic",
                Message = {Text = "Update Kitematic"},
                Description = $"The installed version of Kitematic is outdated. Please download a newer version. Then make sure you install it in {Path.GetDirectoryName(Paths.KitematicPath)}",
                Image = {Source = new BitmapImage(new Uri("Images/kitematic.png", UriKind.Relative))},
                OkButton = {Content = "Download"}
            }.ShowYesNo();
        }
    }

    public static class ImportFromToolBox
    {
        public static bool ShowConfirm()
        {
            return new DialogBox
            {
                Title = "Toolbox migration",
                Message = {Text = "Import Docker Toolbox content"},
                Description = @"This will import the ""default"" Docker Toolbox VirtualBox machine images and containers. The VirtualBox VM will not be deleted",
                Image = {Source = new BitmapImage(new Uri("Images/toolbox.png", UriKind.Relative))},
                OkButton = {Content = "Import"}
            }.ShowYesNo();
        }
    }

    public static class ImportFromToolBoxOnFirstStart
    {
        public static bool ShowConfirm()
        {
            return new DialogBox
            {
                Title = "Toolbox migration",
                Message = {Text = "Toolbox migration assistant"},
                Description = @"A local machine named ""default"" has been found. It is probably the one hosting your Docker images and containers if you previously used Docker Toolbox. Do you want to copy data from that machine?",
                Image = {Source = new BitmapImage(new Uri("Images/toolbox.png", UriKind.Relative))},
                OkButton = {Content = "Copy"},
                CancelButton = {Content = "No"}
            }.ShowYesNo();
        }
    }

    public static class ServiceNotRunningBox
    {
        public static bool ShowConfirm()
        {
            return new DialogBox
            {
                Title = "Service is not running",
                Message = {Text = "Service is not running"},
                Description = @"Docker for Windows service is not running, would you like to start it? Windows will ask you for elevated access.",
                OkButton = {Content = "Start"},
                CancelButton = {Content = "Quit"}
            }.ShowYesNo();
        }
    }

    public static class QuitMessageBox
    {
        public static void ShowOk(string message)
        {
            new DialogBox
            {
                Title = "Error",
                Message = {Text = "An error occured"},
                Description = message,
            }.ShowOk();
        }
    }

    public static class WarnMessageBox
    {
        public static void ShowConfirm(string message)
        {
            new DialogBox
            {
                Title = "Warning",
                Message = {Text = "Warning"},
                Description = message,
                Image = {Source = new BitmapImage(new Uri("Images/WhaleWarning.png", UriKind.Relative))}
            }.ShowOk();
        }
    }

    public static class BackendAskExceptionBox
    {
        public static bool Ask(string message)
        {
            return new DialogBox
            {
                Title = "Docker for Windows",
                Message = {Text = ""},
                Description = message,
                Image = { Source = new BitmapImage(new Uri("Images/WhaleWarning.png", UriKind.Relative)) }
            }.ShowYesNo();
        }
    }

    public partial class DialogBox
    {
        internal DialogBox()
        {
            InitializeComponent();
        }

        public static void Show(string message, string description, string title)
        {
            new DialogBox
            {
                Title = title,
                Message = { Text = message },
                Description = description,
            }.ShowDialog();
        }

        internal void ShowOk()
        {
            OkButton.Content = "OK";
            CancelButton.Width = 0;
            CancelButton.Visibility = Visibility.Hidden;

            ShowDialog();
        }

        internal bool ShowYesNo()
        {
            ShowDialog();

            return DialogResult ?? false;
        }

        internal string Description
        {
            set
            {
                Details.Document.Blocks.Clear();
                var paragraph = new Paragraph();
                paragraph.Inlines.Add(new Run(value));

                Details.Document.Blocks.Add(paragraph);
                Details.IsDocumentEnabled = true;
            }
        }

        private void _okButton_Click(object sender, RoutedEventArgs e)
        {
            DialogResult = true;
        }
    }
}
