#if !DEBUG
using System.ServiceProcess;
#endif

namespace Docker.Service
{
    static class Program
    {
        static void Main()
        {
#if DEBUG
            var service = new DockerService();
            try
            {
                service.DoRun();
                do
                {
                    System.Threading.Thread.Sleep(1000);
                }
                while (System.Diagnostics.Process.GetProcessesByName("Docker for Windows").Length > 0);
            }
            finally
            {
                service.DoStop();
            }
#else
            ServiceBase.Run(new ServiceBase[]
            {
                new DockerService()
            });
#endif
        }
    }
}
