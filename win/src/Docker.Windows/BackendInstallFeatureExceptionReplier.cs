using Docker.Core;
using Docker.Core.Features;
using Docker.WPF;
using Microsoft.Win32;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using Docker.Core.backend;

namespace Docker
{
    internal class BackendInstallFeatureExceptionReplier : IBackendInstallFeatureExceptionReplier
    {
        private const string VirtualBoxRegistryPath = @"SOFTWARE\Oracle\VirtualBox";
        private const string VirtualBoxRegistryKey = @"InstallDir";

        private readonly IBackend _backend;
        private readonly ICmd _cmd;
        private readonly IToolboxMigration _toolboxMigration;

        private Action _quitAction;

        public BackendInstallFeatureExceptionReplier(IBackend backend, ICmd cmd, IToolboxMigration toolboxMigration)
        {
            _backend = backend;
            _cmd = cmd;
            _toolboxMigration = toolboxMigration;
        }

        public void SetQuitAction(Action action)
        {
            _quitAction = action;
        }

        private string FormatMessage(IReadOnlyCollection<Feature> features)
        {
            var descriptions = string.Join(" and ", features.Select(feature => feature.Description));

            string message;

            if (features.Count == 1)
            {
                message = $"{descriptions} feature is not enabled." + Environment.NewLine +
                    "Do you want to enable it for Docker to be able to work properly?" + Environment.NewLine +
                    "Your computer will restart automatically.";
            }
            else
            {
                message = $"{descriptions} features are not enabled." + Environment.NewLine +
                    "Do you want to enable them for Docker to be able to work properly?" + Environment.NewLine +
                    "Your computer will restart automatically.";
            }

            if (features.Contains(Feature.HyperV))
            {
                if (_toolboxMigration.IsToolboxInstalled)
                    message += Environment.NewLine + "Note: Docker Toolbox will no longer work.";
                else if (IsVirtualBoxInstalled)
                    message += Environment.NewLine + "Note: VirtualBox will no longer work.";
            }
            return message;
        }

        private void ProcessError(IReadOnlyCollection<Feature> features, ITaskQueue taskQueue)
        {
            var descriptions = string.Join(" and ", features.Select(feature => feature.Description));

            var message = $"Docker failed to enable {descriptions}.{Environment.NewLine}Please enable ";
            message += features.Count == 1 ? "this feature" : "those features";
            message += " manually for Docker to be able to work properly.";

            if (features.Contains(Feature.HyperV))
            {
                message += $"{Environment.NewLine}To enable hyperV, follow the instructions here:{Environment.NewLine}{Urls.HyperVMicrosoftInstallDocumentation}";
            }

            QuitMessageBox.ShowOk(message);
            taskQueue.Queue(_quitAction);
        }

        public void AskToUser(BackendInstallFeatureException exception, ITaskQueue taskQueue)
        {
            var features = exception.Features;

            if (!BackendAskExceptionBox.Ask(FormatMessage(features))) return;

            taskQueue.Queue(() =>
            {
                var featuresFailedToInstall = _backend.InstallFeatures(features);
                if (featuresFailedToInstall.Length == 0)
                {
                    // TODO: Switch again or when it restarts
                    RestartComputer();
                }

                ProcessError(featuresFailedToInstall, taskQueue);
            });
        }

        private void RestartComputer()
        {
            var result = _cmd.Run("shutdown.exe", "/g /t 0", 0);
            if (result.ExitCode != 0)
                throw new DockerException($"Restarting the computer failed:{Environment.NewLine}{result.CombinedOutput}{Environment.NewLine}You will have to restart it manually.");
        }

        private static bool IsVirtualBoxInstalled
        {
            get
            {
                try
                {
                    using (var key = Registry.LocalMachine.OpenSubKey(VirtualBoxRegistryPath, false))
                    {
                        if (key == null)
                            return false;
                        var installDirectory = key.GetValue(VirtualBoxRegistryKey) as string;
                        return (installDirectory != null) && Directory.Exists(installDirectory);
                    }
                }
                catch
                {
                    return false;
                }
            }
        }
    }
}
