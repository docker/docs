using System;
using System.Globalization;
using System.IO;
using System.Reflection;
using Docker.Core;
using Microsoft.Win32;
using Newtonsoft.Json;

namespace Docker.WPF
{
    public interface ISettingsLoader
    {
        Settings Load();
        void SaveChanges(Action<Settings> changes);
        bool IsFirstLaunchOrHasExecutableChanged();
    }

    public class SettingsLoader : ISettingsLoader
    {
        private const string Beta14DaemonDefault = "{\"debug\":true,\"labels\":[]}";
        private const string StartAtloginKey = @"SOFTWARE\Microsoft\Windows\CurrentVersion\Run";
        private const string DockerForWindowsKey = "Docker for Windows";

        private static string Filename => Path.Combine(Paths.LocalRoamingApplicationData, "settings.json");

        private readonly object _ioLock = new object();

        public Settings Load()
        {
            if (!File.Exists(Filename) || new FileInfo(Filename).Length == 0)
            {
                var settings = new Settings();
                Save(settings);
                return settings;
            }

            lock (_ioLock)
            {
                try
                {
                    using (var fileReader = File.OpenText(Filename))
                    {
                        return ParseJson(fileReader);
                    }
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"Failed to load settings: {ex.Message}");
                    var settings = new Settings();
                    Save(settings);
                    return settings;
                }
            }
        }

        internal static Settings ParseJson(TextReader reader)
        {
            var settings = (Settings)new JsonSerializer().Deserialize(reader, typeof(Settings));

            if (settings.DaemonOptions == Beta14DaemonDefault)
            {
                settings.DaemonOptions = Settings.DaemonDefault;
            }

            return settings;
        }

        public void SaveChanges(Action<Settings> changes)
        {
            var settings = Load();
            changes.Invoke(settings);
            Save(settings);
        }

        private void Save(Settings settings)
        {
            lock (_ioLock)
            {
                try
                {
                    using (var fileWriter = File.CreateText(Filename))
                    {
                        new JsonSerializer().Serialize(fileWriter, settings);
                    }
                }
                catch (Exception ex)
                {
                    throw new DockerException($"Could not save Settings to {Filename}: {ex.Message}", ex);
                }

                try
                {
                    using (var key = Registry.CurrentUser.OpenSubKey(StartAtloginKey, true))
                    {
                        if (settings.StartAtLogin)
                        {
                            key?.SetValue(DockerForWindowsKey, "\"C:\\Program Files\\Docker\\Docker\\Docker for Windows.exe\""); //TODO: Hardcoded path, why not use Paths.ExecutablePath ?
                        }
                        else
                        {
                            key?.DeleteValue(DockerForWindowsKey);
                        }
                    }
                }
                catch
                {
                    // Ignore
                }
            }
        }

        private static string CurrentExecutableDate
        {
            get
            {
                var location = Assembly.GetEntryAssembly().Location;
                var fi = new FileInfo(location);
                return fi.LastWriteTimeUtc.ToString(CultureInfo.InvariantCulture);
            }
        }

        public bool IsFirstLaunchOrHasExecutableChanged()
        {
            if (File.Exists(Filename) && (Load().ExecutableDate == CurrentExecutableDate))
            {
                return false;
            }

            var settings = Load();
            settings.ExecutableDate = CurrentExecutableDate;
            Save(settings);
            return true;
        }
    }
}