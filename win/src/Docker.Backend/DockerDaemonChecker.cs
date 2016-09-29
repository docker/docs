using Docker.Core;

namespace Docker.Backend
{
    public interface IDockerDaemonChecker
    {
        void Check();
    }

    public class DockerDaemonChecker : IDockerDaemonChecker
    {
        private readonly ILogger _logger;
        private readonly ICmd _cmd;

        public DockerDaemonChecker(ICmd cmd)
        {
            _logger = new Logger(GetType());
            _cmd = cmd;
        }

        public void Check()
        {
            var error = "Docker daemon is not running";

            for (var i = 0; i < 10; i++)
            {
                var result = _cmd.Run(Paths.DockerExe, "ps", 0);
                if (result.ExitCode == 0)
                {
                    _logger.Info("Docker daemon is running");
                    return;
                }

                if (result.CombinedOutput?.Length > 0)
                {
                    error = result.CombinedOutput;
                }
            }

            throw new DockerException(error);
        }
    }
}