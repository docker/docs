using System.Runtime.Serialization.Formatters.Binary;

namespace Docker.Core.Pipe
{
    public class LogPipeClient
    {
        private readonly string _pipeName;
        private readonly Logger _logger;

        public LogPipeClient(string pipeName, Logger logger)
        {
            _pipeName = pipeName;
            _logger = logger;
        }

        public void Run()
        {
            while (true)
            {
                var bf = new BinaryFormatter();

                try
                {
                    using (var client = PipeHelper.NewClientStream(_pipeName))
                    {
                        client.Connect();

                        while (true)
                        {
                            var logEntry = (LogEntry)bf.Deserialize(client);
                            _logger.WriteLog(logEntry);
                        }
                    }
                }
                catch
                {
                    // Ignore
                }
            }
            // ReSharper disable once FunctionNeverReturns
        }
    }
}