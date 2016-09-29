using System.Collections.Generic;
using System.Linq;

namespace Docker.Core
{
    public class CommandLineOptions
    {
        public bool DisableCheckForUpdates { get; }
        public bool DisableWelcomeWhale { get; }
        public bool DisableToolboxMigration { get; }
        public string Username { get; }
        public string Password { get; }

        public CommandLineOptions(IReadOnlyCollection<string> args)
        {
            DisableCheckForUpdates = args.Contains("-DisableCheckForUpdates");
            DisableWelcomeWhale = args.Contains("-DisableWelcomeWhale");
            DisableToolboxMigration = args.Contains("-DisableToolboxMigration");
            Username = ParseUsername(args);
            Password = ParsePassword(args);
        }

        private static string ParseUsername(IEnumerable<string> args)
        {
            return (from _ in args where _.StartsWith("-Username=") select FlagValue(_)).FirstOrDefault();
        }

        private static string ParsePassword(IEnumerable<string> args)
        {
            return (from _ in args where _.StartsWith("-Password=") select FlagValue(_)).FirstOrDefault();
        }

        private static string FlagValue(string flag)
        {
            return flag.Split(new [] {'='}, 2)[1];
        }
    }
}
