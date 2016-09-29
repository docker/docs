using System;
using System.Diagnostics;
using System.IO;
using Docker.Core;

namespace Docker.Backend
{
    public interface IDockerMachineImport
    {
        void MigrateVolume(string inputVolumePath, string outputVolumePath);
    }

    public class DockerMachineImport : IDockerMachineImport
    {
        private readonly Logger _logger;
  
        public DockerMachineImport()
        {
            _logger = new Logger(GetType());
        }

        public void MigrateVolume(string inputVolumePath, string outputVolumePath)
        {
            if (!new FileInfo(inputVolumePath).Exists)
                throw new DockerException($"Volume not found: {inputVolumePath}");

            try
            {
                var qemuImg = Path.Combine(Paths.ResourcesPath, "qemu-img", "qemu-img.exe");
                var psi = new ProcessStartInfo(qemuImg, ImportArguments(inputVolumePath, outputVolumePath))
                {
                    UseShellExecute = false,
                    RedirectStandardOutput = true,
                    RedirectStandardError = true,
                    CreateNoWindow = true
                };

                var process = Process.Start(psi);
                if (process == null)
                    throw new DockerException($"Failed to migrate volume: {inputVolumePath}");

                process.OutputDataReceived += (sender, args) => _logger.Info(args.Data);
                process.ErrorDataReceived += (sender, args) => _logger.Error(args.Data);
                process.BeginOutputReadLine();
                process.BeginErrorReadLine();
                process.WaitForExit();
                if (process.ExitCode != 0)
                    throw new DockerException($"Failed to migrate volume: {inputVolumePath}");
            }
            catch (DockerException)
            {
                throw;
            }
            catch (Exception ex)
            {
                throw new DockerException($"Failed to migrate volume: {ex.Message}", ex);
            }
        }

        internal static string ImportArguments(string inputVolumePath, string outputVolumePath)
        {
            return $"convert \"{inputVolumePath}\" -O vhdx -o subformat=dynamic -p \"{outputVolumePath}\"";
        }
    }
}
