using System;
using System.Globalization;
using System.IO;
using System.Threading;
using Docker.Core;
using Docker.Core.di;

namespace Docker.Cli
{
    internal class MainBackendCli
    {
        public static void Main(string[] args)
        {
            Environment.Exit(Run(args) ? 0 : 1);
        }

        private static bool Run(string[] args)
        {
            if (args.Length == 0)
            {
                BackendCli.Usage();
                return false;
            }

            Logger logger = null;
            try
            {
                Thread.CurrentThread.CurrentCulture = new CultureInfo("en-US");

                Logger.Initialize(Path.Combine(Paths.LocalApplicationData, "cli.txt"));

                logger = new Logger("BackendCli");
                logger.Info("Starting on: " + DateTime.Now);

                using (var singletons = new Singletons(new BackendCliModule(args)))
                {
                    var backendCli = singletons.Get<BackendCli>();
                    backendCli.Run(args);
                }

                return true;
            }
            catch (Exception ex)
            {
                logger?.Error(ex.Message);
                return false;
            }
        }
    }
}