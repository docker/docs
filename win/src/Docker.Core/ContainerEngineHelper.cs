using System;
using System.Diagnostics;
using System.IO;

namespace Docker.Core
{
    public interface IContainerEngineHelper
    {
        void Switch();
        bool UseWindowsContainerEngine { get; }
        bool UseLinuxContainerEngine { get; }
        bool CanUseWindowsContainers();
        void ForceKillLingeringDaemon();
    }

    public class ContainerEngineHelper : IContainerEngineHelper
    {
        private readonly ILogger _logger;

        public ContainerEngineHelper()
        {
            _logger = new Logger(GetType());
        }

        public bool UseWindowsContainerEngine => File.Exists(Paths.DaemonSocketPath);
        public bool UseLinuxContainerEngine => !UseWindowsContainerEngine;

        public bool CanUseWindowsContainers()
        {
            int buildNumber;

            return int.TryParse(Env.Os.BuildNumber, out buildNumber) && buildNumber >= 14372;
        }

        public void Switch()
        {
            _logger.Info("Switching daemon...");

            if (UseWindowsContainerEngine)
            {
                _logger.Info("Switching to Linux containers");

                try
                {
                    File.Delete(Paths.DaemonSocketPath);
                }
                catch (Exception e)
                {
                    throw new DockerException($"Unable to switch to Linux containers: {e.Message}", e);
                }
            }
            else
            {
                _logger.Info("Switching to Windows containers");

                try
                {
                    File.WriteAllText(Paths.DaemonSocketPath, Paths.DockerDaemonNamedPipe);
                }
                catch (Exception e)
                {
                    throw new DockerException($"Unable to switch to Windows containers: {e.Message}", e);
                }
            }
        }

        public void ForceKillLingeringDaemon()
        {
            if (!File.Exists(Paths.DockerPidPath)) return;

            try
            {
                var pid = File.ReadAllText(Paths.DockerPidPath).Trim();
                Process.GetProcessById(int.Parse(pid)).Kill();
            }
            catch (Exception e)
            {
                _logger.Error($"Might have failed to kill a running dockerd process: {e.Message}");
            }

            try
            {
                File.Delete(Paths.DockerPidPath);
            }
            catch (Exception e)
            {
                _logger.Warning($"Can't delete dockerd pid1: {e.Message}");
            }
        }

    }
}