using System;
using System.Collections.Concurrent;
using System.Runtime.Serialization.Formatters.Binary;

namespace Docker.Core.Pipe
{
    public class LogPipeServer : ILogEntryListener, IDisposable
    {
        private readonly string _pipeName;
        private readonly BlockingCollection<LogEntry> _logEntries;

        public LogPipeServer(string pipeName)
        {
            _pipeName = pipeName;
            _logEntries = new BlockingCollection<LogEntry>();
        }

        public void Log(LogEntry logEntry)
        {
            _logEntries.Add(logEntry);
        }

        public void Run()
        {
            while (true)
            {
                var bf = new BinaryFormatter();

                try
                {
                    using (var pipe = PipeHelper.NewServerStream(_pipeName))
                    {
                        pipe.WaitForConnection();

                        foreach (var logEntry in _logEntries.GetConsumingEnumerable())
                        {
                            bf.Serialize(pipe, logEntry);
                            pipe.Flush();
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

        //https://msdn.microsoft.com/library/ms244737.aspx
        ~LogPipeServer()
        {
            Dispose(false);
        }

        public void Dispose()
        {
            Dispose(true);
            GC.SuppressFinalize(this);
        }

        protected virtual void Dispose(bool disposing)
        {
            if (disposing)
                _logEntries.Dispose();
        }
    }
}