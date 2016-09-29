using System;
using Docker.Core;

namespace Docker.WPF.Crash
{
    public interface ICrashReport
    {
        string Send(Exception exception);
        string SendDiagnostic();
    }

    public class CrashReport : ICrashReport
    {
        private readonly Logger _logger;
        private readonly DebugInfo _debugInfo;
        private readonly BugSnag _bugSnag;
        private readonly S3 _s3;

        public CrashReport(DebugInfo debugInfo, BugSnag bugSnag, S3 s3)
        {
            _logger = new Logger(GetType());
            _debugInfo = debugInfo;
            _bugSnag = bugSnag;
            _s3 = s3;
        }

        public string Send(Exception exception)
        {
            _debugInfo.DumpDebugInfo();
            _bugSnag.Notify(exception);
            return _s3.UploadReport();
        }

        public string SendDiagnostic()
        {
            return Send(new DockerException("Diagnostic"));
        }

        public void SendOrFailSilently(Exception ex)
        {
            try
            {
                Send(ex);
            }
            catch
            {
                _logger.Error($"Unable to send crash report: {ex.Message}");
            }
        }
    }
}