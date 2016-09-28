import {
  AmazonIcon,
  DigitaloceanIcon,
  AzureIcon,
  PacketIcon,
  SoftlayerIcon,
  PrivateIcon,
} from '../Icon';

export default (name) => {
  switch (name) {
    case 'aws':
      return AmazonIcon;
    case 'digitalocean':
      return DigitaloceanIcon;
    case 'azure':
      return AzureIcon;
    case 'packet':
      return PacketIcon;
    case 'softlayer':
      return SoftlayerIcon;
    default:
      return PrivateIcon;
  }
};
