using System;
using System.Collections.Generic;
using System.IO;
using System.IO.Compression;
using Docker.Core;
using System.Net.Http;
using Docker.Core.backend;
using Docker.Core.Tracking;

namespace Docker.WPF.Crash
{
    public class S3
    {
        private readonly DebugInfo _debugInfo;
        private readonly IBackend _backend;
        private readonly string _bucket;
        private readonly Tracking _tracking;

        public S3(DebugInfo debugInfo, IBackend backend, Tracking tracking)
        {
            _debugInfo = debugInfo;
            _backend = backend;
            _tracking = tracking;
            _bucket = "docker-pinata-support";
        }

        public string UploadReport()
        {
            var timeStamp = DateTime.Now.ToString("yyyy-MM-dd_HH-mm-ss");

            Upload(timeStamp, new Dictionary<string, object>
            {
                 { "Log", Logger.AllLogs() },
                 { "DebugInfo", _debugInfo.GetDebugInfo() }
            });

            return $"{_tracking.Id}/{timeStamp}";
        }

        internal void Upload(string timeStamp, Dictionary<string, object> content)
        {
            var tmpDir = Path.Combine(Path.GetTempPath(), Path.GetRandomFileName());
            Directory.CreateDirectory(tmpDir);

            foreach (var key in content.Keys)
            {
                File.WriteAllText(Path.Combine(tmpDir, $"{key}-{timeStamp}.txt"), content[key].ToString());
            }

            var logTar = _backend.DownloadVmLogs();
            if (!string.IsNullOrEmpty(logTar))
            {
                File.Copy(logTar, Path.Combine(tmpDir, "mobyLogs.tar"));
                File.Delete(logTar);
            }

            var zipFile = Path.Combine(Path.GetTempPath(), Path.GetRandomFileName());
            ZipFile.CreateFromDirectory(tmpDir, zipFile);

            UploadToS3(timeStamp, new FileInfo(zipFile));
        }

        private void UploadToS3(string timeStamp, FileInfo zipFile)
        {
            try
            {
                using (var client = new HttpClient(new HttpClientHandler {AllowAutoRedirect = true}))
                {
                    client.DefaultRequestHeaders.Add("x-amz-acl", "bucket-owner-full-control");

                    using (var fs = zipFile.Open(FileMode.Open, FileAccess.Read))
                    {
                        var content = new StreamContent(fs);
                        content.Headers.ContentType = new System.Net.Http.Headers.MediaTypeHeaderValue("application/zip");
                        var response =
                            client.PutAsync(
                                $"https://{_bucket}.s3.amazonaws.com/incoming/w1/{_tracking.Id}/{timeStamp}/report.zip",
                                content).Result;
                        if (!response.IsSuccessStatusCode)
                        {
                            throw new DockerException("Unable to upload to S3");
                        }
                    }

                }
            }
            catch (DockerException)
            {
                throw;
            }
            catch (Exception ex)
            {
                throw new DockerException("Unable to upload to S3", ex);
            }
        }
    }
}
