using System;
using System.Collections.Generic;
using System.IO;
using System.IO.Pipes;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;

namespace SerialWatcher
{
    class Program
    {
        private static readonly string _namedPipeHostname = ".";
        private static readonly string _namedPipePortname = "MobyLinuxVM-com1";
        static void Main(string[] args)
        {
            StreamWriter sw = null;
            new Thread(() =>
                    {
                        try
                        {
                            while (true)
                            {
                                using (NamedPipeClientStream pipeStream = new NamedPipeClientStream(_namedPipeHostname, _namedPipePortname, PipeDirection.InOut))
                                {
                                    try { pipeStream.Connect(2000); } catch { }
                                    if (pipeStream.IsConnected)
                                    {
                                        Console.WriteLine("********************************************************************************");
                                        Console.WriteLine("---------------------     N E W  C O N N E C T I O N    ------------------------");
                                        Console.WriteLine("********************************************************************************");
                                        StreamReader sr = new StreamReader(pipeStream);
                                        sw = new StreamWriter(pipeStream);
                                        try
                                        {
                                            sw.AutoFlush = true;
                                            while (pipeStream.IsConnected)
                                            {
                                                string line;
                                                while ((line = sr.ReadLine()) != null)
                                                {
                                                    Console.WriteLine(line);
                                                }
                                            }
                                        }
                                        catch (Exception ex)
                                        {
                                            Console.WriteLine("Failed reading from serial: {0}", ex.Message);
                                        }
                                        finally
                                        {
                                            if (pipeStream.IsConnected)
                                            {
                                                sr.Dispose();
                                                sw.Dispose();
                                            }
                                            sr = null;
                                            sw = null;
                                        }
                                    }
                                }
                                Console.WriteLine("********************************************************************************");
                                Console.WriteLine("*********************     C O N N E C T I O N  L O S T    **********************");
                                Console.WriteLine("********************************************************************************");
                            }
                        }
                        // Catch the IOException that is raised if the pipe is
                        // broken or disconnected.
                        catch (IOException e)
                        {
                            Console.WriteLine("Client error: {0}", e.Message);
                        }
                    }).Start();
            string message;
            Console.WriteLine("Docker Serial Watcher v{0}.{1}", System.Reflection.Assembly.GetExecutingAssembly().GetName().Version.Major, System.Reflection.Assembly.GetExecutingAssembly().GetName().Version.Minor);
            Console.WriteLine("Please type a message and press [ENTER], or type 'quit' to exit the program");
            while ((message = Console.ReadLine()) != null)
            {
                if (message == "quit")
                    break;

                if (sw != null)
                    sw.WriteLine(message);
            }
        }
    }
}
