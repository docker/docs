using System;
using System.Net;

namespace Docker.Core
{
    public class IpHelper
    {
        public string VmIp(string subnetAddress, int subnetMaskSize)
        {
            return IpAtIndex(subnetAddress, 2, subnetMaskSize);
        }

        public string SwitchIp(string subnetAddress, int subnetMaskSize)
        {
            return IpAtIndex(subnetAddress, 1, subnetMaskSize);
        }

        private string IpAtIndex(string subnetAddress, int index, int subnetMaskSize)
        {
            CheckMask(subnetAddress, subnetMaskSize);
            var addr = ConvertAddressToByteArray(subnetAddress);

            return $"{addr[0]}.{addr[1]}.{addr[2]}.{addr[3] + index}";
        }

        public string SubnetMask(string subnetAddress, int subnetMaskSize)
        {
            var mask = CheckMask(subnetAddress, subnetMaskSize);
            var addr = ConvertAddressToByteArray(mask);

            return $"{addr[0]}.{addr[1]}.{addr[2]}.{addr[3]}";
        }

        private string CheckMask(string subnetAddress, int subnetMaskSize)
        {
            var mask = GetMaskFromSize(subnetMaskSize);
            if (!IsAddressMatchWithMask(subnetAddress, mask))
            {
                throw new DockerException($"Subnet address [{subnetAddress}] is incompatible with mask [{mask}]");
            }

            return mask;
        }

        private static int GetBitSize(byte value)
        {
            var bitSize = 0;
            byte bitMask = 128;
            while ((value & bitMask) != 0)
            {
                bitSize++;
                value ^= bitMask;
                bitMask /= 2;
            }
            if (value != 0) return -1;
            return bitSize;
        }

        public int GetMaskSize(string mask)
        {
            var ipParts = mask.Split('.');
            if (ipParts.Length != 4) return -1;

            byte ip0;
            if (!byte.TryParse(ipParts[0], out ip0)) return -1;
            var maskSize0 = GetBitSize(ip0);
            if (maskSize0 < 0) return -1;
            byte ip1;
            if (!byte.TryParse(ipParts[1], out ip1)) return -1;
            var maskSize1 = GetBitSize(ip1);
            if ((maskSize1 == -1) || (maskSize1 > 0 && maskSize0 != 8)) return -1;
            byte ip2;
            if (!byte.TryParse(ipParts[2], out ip2)) return -1;
            var maskSize2 = GetBitSize(ip2);
            if ((maskSize2 == -1) || (maskSize2 > 0 && maskSize1 != 8)) return -1;
            byte ip3;
            if (!byte.TryParse(ipParts[3], out ip3)) return -1;
            var maskSize3 = GetBitSize(ip3);
            if ((maskSize3 == -1) || (maskSize3 > 0 && maskSize2 != 8)) return -1;
            return maskSize0 + maskSize1 + maskSize2 + maskSize3;
        }

        public string GetMaskFromSize(int maskSize)
        {
            checked
            {
                var ip3 = 256 - (1 << (8 - Math.Max(0, Math.Min(maskSize - 24, 8))));
                var ip2 = 256 - (1 << (8 - Math.Max(0, Math.Min(maskSize - 16, 8))));
                var ip1 = 256 - (1 << (8 - Math.Max(0, Math.Min(maskSize - 8, 8))));
                var ip0 = 256 - (1 << (8 - Math.Max(0, Math.Min(maskSize, 8))));

                return $"{ip0}.{ip1}.{ip2}.{ip3}";
            }
        }

        public bool IsValidAdress(string address)
        {
            if (address == null)
                throw new ArgumentNullException(nameof(address));

            IPAddress addr;
            return IPAddress.TryParse(address, out addr) && address.Split('.').Length == 4;
        }

        private byte[] ConvertAddressToByteArray(string address)
        {
            if (!IsValidAdress(address)) throw new DockerException($"Invalid address: {address}");
            try
            {
                var ip = address.Split('.');
                return new[] { byte.Parse(ip[0]), byte.Parse(ip[1]), byte.Parse(ip[2]), byte.Parse(ip[3]) };
            }
            catch
            {
                return new byte[] { 0, 0, 0, 0 };
            }
        }

        public bool IsAddressMatchWithMask(string address, string mask)
        {
            var a = ConvertAddressToByteArray(address);
            var m = ConvertAddressToByteArray(mask);
            return (a[0] | m[0]) == m[0] && (a[1] | m[1]) == m[1] && (a[2] | m[2]) == m[2] && (a[3] | m[3]) == m[3];
        }
    }
}
