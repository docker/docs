using System;
using System.Runtime.Serialization.Formatters.Binary;

namespace Docker.Core.Pipe
{
    public class NamedPipeClient
    {
        private readonly Logger _logger;
        private readonly string _pipeName;

        public NamedPipeClient(string pipeName)
        {
            _logger = new Logger(GetType());
            _pipeName = pipeName;
        }

        public object Send(string action, params object[] parameters)
        {
            var retry = 5;
            while (true)
            {
                try
                {
                    return TrySend(action, parameters);
                }
                catch (TimeoutException ex)
                {
                    if (retry-- > 0)
                    {
                        continue;
                    }
                    throw new DockerException($"Failed to connect to {_pipeName}: time out", ex);
                }
            }
        }

        private object TrySend(string action, params object[] parameters)
        {
            _logger.Info($"Sending {action}({string.Join(", ", parameters)})...");

            using (var client = PipeHelper.NewImpersonatedClientStream(_pipeName))
            {
                client.Connect(1000);

                var bf = new BinaryFormatter();
                var request = new PipeRequest(action, parameters);
                bf.Serialize(client, request);
                client.Flush();
                client.WaitForPipeDrain();

                var response = (PipeResponse)bf.Deserialize(client);
                if (!response.Success)
                {
                    var error = (Exception)response.ReturnValue;
                    _logger.Error($"Unable to send {action}: {error.Message}");
                    throw error;
                }

                _logger.Info($"Received response for {action}");
                return response.ReturnValue;
            }
        }
    }
}