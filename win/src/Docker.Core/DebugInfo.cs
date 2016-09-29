using System;
using System.IO;
using Docker.Core.backend;

namespace Docker.Core
{
    public class DebugInfo
    {
        private readonly Logger _logger;
        private readonly ICmd _cmd;
        private readonly IBackend _backend;

        private string _debugInfo;
        private bool _fileHasBeenDumped;

        public DebugInfo(IBackend backend, ICmd cmd)
        {
            _logger = new Logger(GetType());
            _backend = backend;
            _cmd = cmd;
        }

        public string GetDebugInfo()
        {
            if (_debugInfo != null)
                return _debugInfo;

            try
            {
                _debugInfo = _backend.GetDebugInfo();
            }
            catch (Exception ex)
            {
                try
                {
                    var result = _cmd.Run("powershell.exe",
                        $"-ExecutionPolicy UNRESTRICTED -NoProfile -NonInteractive -command \"& {{& '{Paths.InResourcesDir("DockerDebugInfo.ps1")}'}}\"", 0);
                    _debugInfo =  $"Error: {ex.Message}, running DebugInfo script as user{Environment.NewLine}{result.CombinedOutput}";
                }
                catch (Exception ex2)
                {
                    _logger.Error(ex2.Message);
                    _debugInfo = ex2.Message;
                }
            }

            return _debugInfo;
        }

        public void DumpDebugInfo()
        {
            if (_fileHasBeenDumped)
                return;

            var filename = Path.Combine(Paths.LocalApplicationData, "DebugInfo-" + DateTime.Now.ToString("yyyy-MM-dd_HH-mm-ss") + ".txt");
            try
            {
                File.WriteAllText(filename, GetDebugInfo());
                _fileHasBeenDumped = true;
            }
            catch (Exception ex)
            {
                _logger.Error(ex.Message);
            }
        }
    }
}
