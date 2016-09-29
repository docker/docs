using System;
using System.IO;

namespace Docker.Core
{
    public interface IToolboxMigration
    {
        bool IsToolboxInstalled { get; }
        bool DefaultMachineExists { get; }
        bool HyperVDiskFolderExists { get; }
        void MigrateUser();
        string GetMachineVolumePath(string name);
    }

    public abstract class BaseToolboxMigration
    {
        public bool HyperVDiskFolderExists => new DirectoryInfo(Paths.HyperVDisksPath).Exists;

        public string GetMachineVolumePath(string name) => Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.UserProfile), ".docker", "machine", "machines", name, "disk.vmdk");
    }

    public class NopToolboxMigration : BaseToolboxMigration, IToolboxMigration
    {
        public bool DefaultMachineExists => false;

        public bool IsToolboxInstalled => false;

        public void MigrateUser()
        {
            // do nothing
        }
    }

    public class ToolboxMigration : BaseToolboxMigration, IToolboxMigration
    {
        public bool DefaultMachineExists => new FileInfo(GetMachineVolumePath("default")).Exists;

        public bool IsToolboxInstalled => new FileInfo(IDFilePath).Exists;

        private string IDFilePath => Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData), "DockerToolbox", "id.txt");

        public void MigrateUser()
        {
            try
            {
                File.Copy(IDFilePath, Path.Combine(Paths.LocalRoamingApplicationData, ".trackid"), true);
            }
            catch (Exception ex)
            {
                throw new Exception($"Failed to migrate user: {ex.Message}");
            }
        }
    }
}
