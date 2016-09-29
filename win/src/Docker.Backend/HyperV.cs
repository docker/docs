using System;
using System.Collections.Generic;
using System.IO;
using System.Management.Automation;
using System.Security.AccessControl;
using System.Security.Principal;
using Docker.Core;
using System.ServiceProcess;
using Docker.Core.backend;

namespace Docker.Backend
{
    public interface IHyperV
    {
        void CheckHyperVState();
        void Start();
        void Create(Settings settings);
        void Stop();
        void Destroy(bool keepVolume);
        string GetId();
        string DownloadLogs();
    }

    public class HyperV : IHyperV
    {
        private const string Script = "MobyLinux.ps1";

        private readonly Logger _logger;
        private readonly IPowerShell _powerShell;
        private readonly ICmd _cmd;

        public HyperV(IPowerShell powerShell, ICmd cmd)
        {
            _logger = new Logger(GetType());
            _powerShell = powerShell;
            _cmd = cmd;
        }

        public void CheckHyperVState()
        {
            try
            {
                var sc = new ServiceController { ServiceName = "vmms" };
                if (sc.Status != ServiceControllerStatus.Running)
                {
                    _logger.Info("Starting Hyper-V Virtual Machine Management service");

                    sc.Start();
                    sc.WaitForStatus(ServiceControllerStatus.Running, TimeSpan.FromSeconds(2));

                    if (sc.Status != ServiceControllerStatus.Running)
                        throw new DockerException("Failed to start Hyper-V Virtual Machine Management service");
                }
            }
            catch (Exception ex)
            {
                var message = ex.Message + Environment.NewLine;
                message += Environment.NewLine + "If you just enabled Hyper-V, please restart now.";
                throw new BackendQuitException(message);
            }

            var biosHypervisorPresent = true;
            try
            {
                var wmi = new Wmi.Wmi(@"root\CIMV2");
                var containers = wmi.GetOrFail("Win32_computersystem");
                if (containers.Get("HypervisorPresent") == "False")
                    biosHypervisorPresent = false;
            }
            catch
            {
                //ignored
            }

            if (!biosHypervisorPresent)
                throw new BackendQuitException("Hardware assisted virtualization and data execution protection must be enabled in the BIOS");

            _logger.Info("Hyper-V is running");
        }

        public void Start()
        {
            _logger.Info("Start");

            RunScript("start", new Dictionary<string, object>
            {
                {"Start", true}
            });
        }

        public void Create(Settings settings)
        {
            _logger.Info("Create");

            var isoPath = Paths.InResourcesDir("mobylinux.iso");

            RunScript("create", new Dictionary<string, object>
            {
                {"Create", true},
                {"SwitchSubnetAddress", settings.SubnetAddress},
                {"SwitchSubnetMaskSize", settings.SubnetMaskSize},
                {"CPUs", settings.VmCpus},
                {"Memory", settings.VmMemory},
                {"IsoFile", isoPath}
            });
        }

        public void Stop()
        {
            _logger.Info("Stop");

            RunScript("stop", new Dictionary<string, object>
            {
                {"Stop", true}
            });
        }

        public void Destroy(bool keepVolume)
        {
            _logger.Info("Destroy");

            if (keepVolume)
            {
                RunScript("destroy", new Dictionary<string, object>
                {
                    {"Destroy", true},
                    {"KeepVolume", true}
                });
            }
            else
            {
                RunScript("destroy", new Dictionary<string, object>
                {
                    {"Destroy", true}
                });
            }
        }

        public string GetId()
        {
            return _powerShell.Output("(Get-VM MobyLinuxVM).Id.Guid");
        }

        public string DownloadLogs()
        {
            string error;
            try
            {
                var tmpTar = Path.Combine(Path.GetTempPath(), Path.GetRandomFileName() + ".tar");

                var result = _cmd.RunAsAdministrator(Paths.MobyLogDownloaderExe, $"-vmid={GetId()} -o={tmpTar}", 0);
                if (result.ExitCode == 0)
                {
                    // Give access to this file to everybody
                    var accessControl = File.GetAccessControl(tmpTar);
                    accessControl.AddAccessRule(new FileSystemAccessRule(new SecurityIdentifier(WellKnownSidType.WorldSid, null), FileSystemRights.FullControl, AccessControlType.Allow));
                    File.SetAccessControl(tmpTar, accessControl);

                    return tmpTar;
                }

                error = result.ErrorOutput;
            }
            catch (Exception e)
            {
                error = e.Message;
            }

            _logger.Warning($"Unable to download logs: {error}");

            return "";
        }

        private void RunScript(string action, Dictionary<string, object> parameters)
        {
            try
            {
                var scriptPath = Paths.InResourcesDir(Script);
                var script = File.ReadAllText(scriptPath);

                _powerShell.Run(script, parameters, line => _logger.Info(line));
            }
            catch (RuntimeException e)
            {
                throw new RuntimeException($"Unable to {action}: {e.Message}{Environment.NewLine}{e.ErrorRecord.ScriptStackTrace}", e, e.ErrorRecord);
            }
            catch (Exception e)
            {
                throw new RuntimeException($"Unable to {action}", e);
            }
        }
    }
}