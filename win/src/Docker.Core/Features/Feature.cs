using System;

namespace Docker.Core.Features
{
    [Serializable]
    public class Feature
    {
        public const string HyperVFeatureName = "Microsoft-Hyper-V";
        public static readonly Feature HyperV = new Feature(HyperVFeatureName, "Hyper-V");
        public static readonly Feature Containers = new Feature("Containers", "Containers");

        public string Name { get; }
        public string Description { get; }

        public Feature(string name, string description)
        {
            Name = name;
            Description = description;
        }

        private bool Equals(Feature other)
        {
            return string.Equals(Name, other.Name);
        }

        public override bool Equals(object obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != GetType()) return false;
            return Equals((Feature)obj);
        }

        public override int GetHashCode()
        {
            return Name.GetHashCode();
        }

        public override string ToString()
        {
            return $"{Name} ({Description})";
        }
    }
}
