using System;
using System.Collections.Generic;

namespace Docker.Core
{
    [Serializable]
    public class Settings
    {
        public const string DaemonDefault = "{\"registry-mirrors\":[],\"insecure-registries\":[]}";

        public int VmCpus;
        public int VmMemory;
        public string SubnetAddress;
        public int SubnetMaskSize;
        public bool UseDnsForwarder;
        public string NameServer;
        public bool IsTracking;

        public readonly Dictionary<string, bool> SharedDrives;
        public string DaemonOptions;
        public string SysCtlConf;

        public string ExecutableDate;
        public bool AutoUpdateEnabled;
        public bool StartAtLogin;

        public string ProxyExclude;
        public string ProxyHttp;
        public string ProxyHttps;
        public bool UseHttpProxy;

        public Settings()
        {
            VmCpus = Math.Min(Environment.ProcessorCount, 2);
            VmMemory = 2048;
            SubnetAddress = "10.0.75.0";
            SubnetMaskSize = 24;
            UseDnsForwarder = true;
            NameServer = "8.8.8.8";
            IsTracking = true;

            SharedDrives = new Dictionary<string, bool>();
            DaemonOptions = DaemonDefault;
            SysCtlConf = "";

            ExecutableDate = "";
            AutoUpdateEnabled = true;
            StartAtLogin = true;

            ProxyExclude = "";
            ProxyHttp = "";
            ProxyHttps = "";
        }
    }
}