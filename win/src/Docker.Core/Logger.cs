using System;
using System.IO;
using NLog;
using NLog.Config;
using NLog.Targets;

namespace Docker.Core
{
    public enum ChannelType
    {
        Info,
        Warning,
        Error,
        Debug
    }

    [Serializable]
    public class LogEntry
    {
        private readonly DateTime _date;
        private readonly string _category;
        private readonly ChannelType _channel;
        private readonly string _message;

        public LogEntry(DateTime date, string category, ChannelType channel, string message)
        {
            _date = date;
            _category = category;
            _channel = channel;
            _message = message;
        }

        public string ToLine()
        {
            return $"[{_date:HH:mm:ss.fff}][{_category,-15}][{_channel,-7}] {_message}";
        }
    }

    public interface ILogEntryListener
    {
        void Log(LogEntry entry);
    }

    public interface ILogger
    {
        void Info(string message);
        void Warning(string message);
        void Debug(string message);
        void Error(string message);
    }

    public class Logger : ILogger
    {
        private static ILogEntryListener _listener;
        private static string _logFile;

        public static void Initialize(string filename)
        {
            _logFile = filename;

            var config = new LoggingConfiguration();

            var fileTarget = new FileTarget("file")
            {
                FileName = filename,
                MaxArchiveFiles = 7,
                ArchiveOldFileOnStartup = true,
                ArchiveNumbering = ArchiveNumberingMode.DateAndSequence,
                ArchiveEvery = FileArchivePeriod.Day,
                Layout = "${message}"
            };
            config.AddTarget(fileTarget);
            config.AddRuleForAllLevels(fileTarget);

#if DEBUG
            var debugTarget = new DebuggerTarget("debug")
            {
                Layout = "${message}"
            };
            config.AddTarget(debugTarget);
            config.AddRuleForAllLevels(debugTarget);
#endif

            LogManager.Configuration = config;
        }

        public static void SetListener(ILogEntryListener listener)
        {
            _listener = listener;
        }

        private readonly NLog.Logger _logger;
        private readonly string _category;

        public Logger(Type type) : this(type.Name)
        {
        }

        public Logger(string name)
        {
            _category = name;
            _logger = LogManager.GetLogger(name);
        }

        // Should be internal
        public void WriteLog(LogEntry entry)
        {
            string line;
            try
            {
                line = entry.ToLine();
            }
            catch
            {
                return; // ignored
            }

            _logger.Info(line);
        }

        private void Log(ChannelType channel, string message)
        {
            var entry = new LogEntry(DateTime.Now, _category, channel, message);
            WriteLog(entry);
            _listener?.Log(entry);
        }

        public void Info(string message)
        {
            Log(ChannelType.Info, message);
        }

        public void Warning(string message)
        {
            Log(ChannelType.Warning, message);
        }

        public void Error(string message)
        {
            Log(ChannelType.Error, message);
        }

        public void Debug(string message)
        {
            Log(ChannelType.Debug, message);
        }

        public static string AllLogs()
        {
            return File.ReadAllText(_logFile);
        }
    }
}