using System.Text.RegularExpressions;
using Docker.Core;

namespace Docker.Backend.Processes
{
    internal class SubProcessLogger
    {
        private static readonly Regex LogWithDate = new Regex("\\d\\d\\d\\d/\\d\\d/\\d\\d \\d\\d:\\d\\d:\\d\\d .*");

        private readonly ILogger _logger;

        public SubProcessLogger(ILogger logger)
        {
            _logger = logger;
        }

        public void Log(string line)
        {
            if (line == null) return;

            var lineNoDate = LogWithDate.Match(line).Success ? line.Substring(20) : line;

            if (lineNoDate.Contains(": [DEBUG] "))
            {
                _logger.Debug(lineNoDate.Replace(": [DEBUG] ", ": "));
            }
            else if (lineNoDate.Contains(": [WARNING] "))
            {
                _logger.Warning(lineNoDate.Replace(": [WARNING] ", ": "));
            }
            else if (lineNoDate.Contains(": [ERROR] "))
            {
                _logger.Error(lineNoDate.Replace(": [ERROR] ", ": "));
            }
            else if (lineNoDate.Contains(": [INFO] "))
            {
                _logger.Info(lineNoDate.Replace(": [INFO] ", ": "));
            }
            else
            {
                _logger.Info(lineNoDate);
            }
        }
    }
}
