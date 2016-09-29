using System;
using System.IO;
using Bugsnag;
using Bugsnag.Clients;
using Docker.Core;
using Newtonsoft.Json;
using Version = Docker.Core.Version;
using System.Text.RegularExpressions;
using Docker.Core.backend;
using Docker.Core.Tracking;

namespace Docker.WPF.Crash
{
    public class BugSnag
    {
        private readonly Logger _logger;
        private readonly Version _version;
        private readonly DebugInfo _debugInfo;
        private readonly Channel _channel;
        private readonly Tracking _tracking;
        private readonly Git _git;
        private readonly ISettingsLoader _settingsLoader;
        private readonly string _apikey;

        private Metadata _metadata;

        public BugSnag(Logger logger, Version version, DebugInfo debugInfo, Channel channel, Tracking tracking, Git git, ISettingsLoader settingsLoader)
        {
            _logger = logger;
            _version = version;
            _debugInfo = debugInfo;
            _channel = channel;
            _tracking = tracking;
            _git = git;
            _settingsLoader = settingsLoader;
            _apikey = channel.BugsnagApiKey;
        }

        private string GenerateHashingGroup(Exception ex)
        {
            try
            {
                var hashingGroup = ex.Message;
                hashingGroup = hashingGroup.Replace(Environment.MachineName, "@1");
                hashingGroup = hashingGroup.Replace(Environment.UserName, "@2");
                hashingGroup = hashingGroup.Replace("The running command stopped because the preference variable \"ErrorActionPreference\" or common parameter is set to Stop:", "");
                var guids = Regex.Matches(hashingGroup, @"(\{){0,1}[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}(\}){0,1}");
                for (var i = 0; i < guids.Count; i++)
                {
                    hashingGroup = hashingGroup.Replace(Environment.UserName, $"@3{i}");
                }
                hashingGroup = hashingGroup.Substring(0, Math.Min(hashingGroup.Length, 80));
                return hashingGroup;
            }
            catch(Exception exception)
            {
                _logger.Error($"Failed to generate hashing group: {exception.Message}, use default message");
                return ex.Message.Substring(0, Math.Min(ex.Message.Length, 80));
            }
        }

        public void Notify(Exception exception)
        {
            if (_metadata == null)
            {
                _metadata = BuildMetadata(exception);
            }

            try
            {
                var errorId = Guid.NewGuid().ToString();
                _logger.Info($"Sending Bugsnag report {errorId}...");

                var bugsnag = new BaseClient(_apikey);

                bugsnag.Config.SetUser(_tracking.Id, "", errorId);
                bugsnag.Config.AppVersion = _version.ToString();
                bugsnag.Config.AutoNotify = false;

                if (exception.Message.Length > 0)
                {
                    bugsnag.Config.BeforeNotify(error =>
                    {
                        error.GroupingHash = GenerateHashingGroup(error.Exception);
                        return true;
                    });
                }

                bugsnag.Notify(exception, Severity.Error, _metadata);

                _logger.Info($"Bugsnag report {errorId} sent");
            }
            catch (Exception bugsnagException)
            {
                _logger.Error(bugsnagException.Message);
                _logger.Error(bugsnagException.StackTrace);
                if (bugsnagException.InnerException != null)
                {
                    _logger.Error(bugsnagException.InnerException.Message);
                    _logger.Error(bugsnagException.InnerException.StackTrace);
                }
            }
        }

        public Metadata BuildMetadata(Exception exception)
        {
            var metadata = new Metadata();
            if (exception is BackendException)
            {
                metadata.AddToTab("SERVICE", "stacktrace", exception.InnerException.StackTrace);
            }
            metadata.AddToTab("APP", "sha1", _git.ShortSha1());
            metadata.AddToTab("APP", "channel", _channel.ToString());
            metadata.AddToTab("OS", "osEdition", Env.Os.Edition);
            metadata.AddToTab("OS", "osReleaseId", Env.Os.ReleaseId);
            metadata.AddToTab("OS", "osBuild", Env.Os.BuildNumber);
            metadata.AddToTab("OS", "osBuildLabName", Env.Os.BuildLabName);
            metadata.AddToTab("OS", "osLanguage", Env.Os.Language);
            metadata.AddToTab("LOG", "logfile", Logger.AllLogs());
            metadata.AddToTab("SETTINGS", "settings", ToJson(_settingsLoader.Load()));
            metadata.AddToTab("DEBUG", "debug", _debugInfo.GetDebugInfo());
            return metadata;
        }

        private static string ToJson(Settings settings)
        {
            try
            {
                using (var writer = new StringWriter())
                {
                    new JsonSerializer().Serialize(writer, settings);

                    return writer.ToString();
                }
            }
            catch
            {
                return "UNKNOWN";
            }
        }
    }
}