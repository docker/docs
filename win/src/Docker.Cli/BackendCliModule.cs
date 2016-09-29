using System.Collections.Generic;
using Docker.Core;
using Docker.Core.backend;
using Docker.Core.Backend;
using Docker.Core.di;
using Docker.WPF;
using Docker.WPF.Credentials;

namespace Docker.Cli
{
    internal class BackendCliModule : Module
    {
        private readonly IReadOnlyCollection<string> _args;

        public BackendCliModule(IReadOnlyCollection<string> args)
        {
            _args = args;
        }

        protected override void Configure()
        {
            var commandLineOptions = new CommandLineOptions(_args);
            Bind(new Logger("BackendCli"));
            Bind(commandLineOptions);
            Bind(Channel.Load());
            Bind<INotifications, CliNotifications>();
            Bind<IBackend, BackendClient>();
            Bind<ISharedDrives, SettingsBasedSharedDrives>();
            Bind<IWelcomeShower, NopWelcomeShower>();
            Bind<IToolboxMigration, NopToolboxMigration>();
            Bind<ICredentialAsker>(new FixedCredentialAsker(commandLineOptions.Username, commandLineOptions.Password));
        }
    }
}