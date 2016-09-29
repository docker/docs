using System;
using System.Diagnostics;
using System.IO;
using Docker.Core;

namespace Docker.Backend.Processes
{
    public interface IWindowsDockerDaemon
    {
        void Start();
        void Stop();
    }

    public class WindowsDockerDaemon : IWindowsDockerDaemon
    {
        private readonly ExternalProcess _process;

        public WindowsDockerDaemon()
        {
            _process = new ExternalProcess(GetType(), "dockerd.exe", "dockerd");
        }

        public void Start()
        {
            _process.Start(() => $"-H npipe:////./pipe/{Paths.DockerDaemonNamedPipe}");
        }

        public void Stop()
        {
            _process.Stop();
        }
    }
}
