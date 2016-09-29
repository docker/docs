using System.Collections.Generic;
using Docker.Core;
using static Docker.Core.Cmd;

namespace Docker.Backend
{
    public interface IMobyCommand
    {
        ProcessExecutionInfo Run(string commandLine, IDictionary<string, string> envVariables = null, bool load = false);
    }

    public class MobyCommand : IMobyCommand
    {
        private readonly ICmd _cmd;

        public MobyCommand(ICmd cmd)
        {
            _cmd = cmd;
        }

        public ProcessExecutionInfo Run(string commandLine, IDictionary<string, string> envVariables = null, bool load = false)
        {
            if (load)
            {
                _cmd.Run(Paths.DockerExe, $"load -i \"{Paths.InResourcesDir("nsenter.tar")}\"");
            }

            var command = "run --rm --privileged --pid=host ";

            if (envVariables != null)
            {
                foreach (var envVariable in envVariables)
                {
                    command += $"-e {envVariable.Key}={Escape(envVariable.Value)} ";
                }
            }

            command += $"d4w/nsenter /bin/sh -c {Escape(commandLine)}";

            return _cmd.Run(Paths.DockerExe, command, 0);
        }

        private static string Escape(string value)
        {
            return $"\"{value.Replace("\"", "\\\"")}\"";
        }
    }
}