using System;
using System.Collections.Generic;
using System.Linq;
using System.Management.Automation;
using System.Management.Automation.Runspaces;

namespace Docker.Core
{
    public interface IPowerShell
    {
        void Run(string script, Dictionary<string, object> parameters, Action<string> lineReceived);
        string Output(string script);
    }

    public class PowerShell : IPowerShell
    {
        private readonly Logger _logger;
        private readonly ISupplier<RunspacePool> _runspacePool;

        public PowerShell()
        {
            _logger = new Logger(GetType());

            _runspacePool = Suppliers<RunspacePool>.Memoize(() =>
            {
                var initialState = InitialSessionState.CreateDefault();

                var runspacePool = RunspaceFactory.CreateRunspacePool(initialState);
                runspacePool.Open();
                return runspacePool;
            });
        }

        public void Run(string script, Dictionary<string, object> parameters, Action<string> lineReceived)
        {
            _logger.Info($"Run script with parameters:{ToString(parameters)}...");

            using (var ps = CreatePowerShell())
            {
                ps.AddScript(script);

                if (parameters != null)
                {
                    foreach (var parameter in parameters)
                    {
                        ps.AddParameter(parameter.Key, parameter.Value);
                    }
                }

                var output = new PSDataCollection<PSObject>();
                var result = ps.BeginInvoke<int, PSObject>(null, output);

                foreach (var item in output)
                {
                    lineReceived?.Invoke(item.ToString());
                }

                ps.EndInvoke(result);
            }
        }

        public string Output(string script)
        {
            _logger.Info(script.Length < 40 ? "Run script..." : $"Run script '{script}'...");

            using (var ps = CreatePowerShell())
            {
                ps.AddScript(script);

                var results = ps.Invoke();

                return string.Join(Environment.NewLine, results);
            }
        }

        private System.Management.Automation.PowerShell CreatePowerShell()
        {
            var ps = System.Management.Automation.PowerShell.Create();
            ps.RunspacePool = _runspacePool.Get();
            return ps;
        }

        private static string ToString(Dictionary<string, object> parameters)
        {
            if (parameters == null) return "";

            return parameters.Aggregate("", (current, parameter) => current + $" -{parameter.Key} {parameter.Value}");
        }
    }
}