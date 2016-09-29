using System;
using System.Diagnostics;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Navigation;
using Docker.Core;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;

namespace Docker.WPF
{
    public class JsonReturn
    {
        public readonly bool IsValid;
        public readonly string ErrorMessage;

        public JsonReturn(bool isValid, string message)
        {
            IsValid = isValid;
            ErrorMessage = message;
        }
    }

    public partial class DaemonSettings : IActivable
    {
        private readonly ILogger _logger;
        private readonly IActions _actions;

        private bool _loaded;

        public DaemonSettings(IActions actions)
        {
            _logger = new Logger(GetType());
            _actions = actions;

            InitializeComponent();
        }

        private void Apply(object sender, RoutedEventArgs eventArgs)
        {
            _loaded = false;
            try
            {
                var inlinedJson = "";
                if (DaemonOptionsText.Text != "")
                {
                    DaemonOptionsText.Text = JToken.Parse(DaemonOptionsText.Text).ToString();
                    inlinedJson = DaemonOptionsText.Text.Replace("\r", "").Replace("\n", "").Replace("\t", "");
                }

                _actions.RestartVm(_ => _.DaemonOptions = inlinedJson);
            }
            catch (Exception exception)
            {
                ErrorText.Text = exception.Message;
                ErrorText.Visibility = Visibility.Visible;
            }
            finally
            {
                _loaded = true;
            }
        }

        public void Refresh(Settings settings)
        {
            _loaded = false;
            try
            {
                if (settings.DaemonOptions == "")
                {
                    DaemonOptionsText.Text = "";
                }
                else
                {
                    try
                    {
                        var jToken = JToken.Parse(settings.DaemonOptions);
                        DaemonOptionsText.Text = jToken.ToString();
                    }
                    catch (Exception e)
                    {
                        _logger.Warning($"Can't load json from settings: {e.Message}");
                        DaemonOptionsText.Text = "";
                    }
                }
            }
            finally
            {
                _loaded = true;
            }
        }

        private void OnDocumentationClick(object sender, RequestNavigateEventArgs e)
        {
            Process.Start(new ProcessStartInfo(e.Uri.AbsoluteUri));
            e.Handled = true;
        }

        private static JsonReturn IsValidJson(string text)
        {
            try
            {
                JObject.Parse(text);
                return new JsonReturn(true, null);
            }
            catch(JsonReaderException e)
            {
                var near = "";
                if (e.Path != "")
                {
                    near = $"Near '{ e.Path}' ";
                } 
                return new JsonReturn(false, $"{near}Line: {e.LineNumber}, Column: {e.LinePosition}");
            }
        }

        private void OnTextChanged(object sender, TextChangedEventArgs e)
        {
            if (!_loaded) return;

            var error = "";
            var valid = true;

            var json = DaemonOptionsText.Text;
            if (json.Length != 0)
            {
                var parsedJson = IsValidJson(json);
                if (!parsedJson.IsValid)
                {
                    error = $"Invalid json: {parsedJson.ErrorMessage}";
                    valid = false;
                }
            }

            ErrorText.Text = error;
            ErrorText.Visibility = valid ? Visibility.Hidden : Visibility.Visible;
            ApplyButton.IsEnabled = valid;
        }
    }
}
