using System;
using System.IO;
using System.Reflection;

namespace Docker.Core
{
    public static class Paths
    {
        // Folders
        public static string ProgramFiles => EnsureDirectoryIsCreated(Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.ProgramFiles), "Docker"));
        public static string CommonApplicationData => EnsureDirectoryIsCreated(Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.CommonApplicationData), "Docker"));
        public static string LocalApplicationData => EnsureDirectoryIsCreated(Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData), "Docker"));
        public static string LocalRoamingApplicationData => EnsureDirectoryIsCreated(Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData), "Docker"));
        public static string HyperVDisksPath => @"C:\Users\Public\Documents\Hyper-V\Virtual hard disks";

        // Files
        public static string KitematicPath => FullPath(ProgramFiles, "Kitematic", "Kitematic.exe");
        public static string DockerExe => FullPath(ResourcesPath, "bin", "docker.exe");
        public static string WindowsDockerDaemonExe => FullPath(ResourcesPath, "bin", "dockerd.exe");
        public static string DockerProxyExe => FullPath(ResourcesPath, "com.docker.proxy.exe");
        public static string MobyLogDownloaderExe => FullPath(ResourcesPath, "moby-diag-dl.exe");
        public static string DaemonSocketPath => FullPath(CommonApplicationData, "daemonsocket");
        public static string TrackId = FullPath(LocalRoamingApplicationData, ".trackid");
        public static string MobyDiskPath => FullPath(HyperVDisksPath, "MobyLinuxVM.vhdx");
        public static string TmpVolumeMigrationPath => FullPath(HyperVDisksPath, "MigratedVolume.vhdx");
        public static string DockerPidPath => FullPath(Environment.GetFolderPath(Environment.SpecialFolder.CommonApplicationData), "docker.pid");

        // Pipes
        public static string DockerDaemonNamedPipe = "docker_engine_windows";

        private static string ExecutablePath => Path.GetDirectoryName(Assembly.GetExecutingAssembly().Location);

        public static string ResourcesPath
        {
            get
            {
                var envVar = Environment.GetEnvironmentVariable("D4W_RESOURCES_PATH");
                if (!string.IsNullOrEmpty(envVar))
                {
                    return envVar;
                }
#if DEBUG
                return FullPath(ExecutablePath, "..", "..", "src", "Resources");
#else
                return FullPath(ExecutablePath, "Resources");
#endif
            }
        }

        private static string _logFilename;

        public static string LogFilename => _logFilename ?? (_logFilename = Path.Combine(LocalApplicationData, "log.txt"));

        public static string EnsureDirectoryIsCreated(string path)
        {
            new DirectoryInfo(path).Create();
            return path;
        }

        public static string InResourcesDir(string script)
        {
            return Path.Combine(ResourcesPath, script);
        }

        private static string FullPath(params string[] args)
        {
            return Path.GetFullPath(Path.Combine(args));
        }
    }
}