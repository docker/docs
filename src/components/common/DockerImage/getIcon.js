import {
  BittorrentSyncIcon,
  CentosIcon,
  CouchdbIcon,
  DatadogIcon,
  DebianIcon,
  DockerIcon,
  DrupalIcon,
  ElasticsearchIcon,
  FedoraIcon,
  GlassfishIcon,
  HaproxyIcon,
  InfluxdbIcon,
  JoomlaIcon,
  MariadbIcon,
  MemcachedIcon,
  MongodbIcon,
  MysqlIcon,
  NewrelicIcon,
  PostgresqlIcon,
  RabbitmqIcon,
  RedisIcon,
  RiakIcon,
  TomcatIcon,
  TutumIcon,
  UbuntuIcon,
  WikimediaIcon,
  WordpressIcon,
} from '../Icon';

const DOCKER_NS = 'library';
const TUTUM_NS = 'tutum';

function normalize(fullName) {
  let [namespace, nameAndTag] = fullName.split('/');
  if (!nameAndTag) {
    nameAndTag = namespace;
    namespace = DOCKER_NS;
  }

  // eslint-disable-next-line no-unused-vars
  const [name, tag] = nameAndTag.split(':');

  return `${namespace}/${name}`;
}

export default (fullName = '') => {
  const predicate = normalize(fullName);
  const d = DOCKER_NS;
  const t = TUTUM_NS;

  switch (predicate) {
    case `${d}/wordpress`:
    case `${t}/wordpress`:
    case `${t}/wordpress-stackable`:
      return WordpressIcon;
    case `${d}/centos`:
    case `${t}/centos`:
      return CentosIcon;
    case `${d}/tomcat`:
    case `${t}/tomcat`:
      return TomcatIcon;
    case `${t}/riak`:
      return RiakIcon;
    case `${d}/redis`:
    case `${t}/redis`:
      return RedisIcon;
    case `${d}/rabbitmq`:
    case `${t}/rabbitmq`:
      return RabbitmqIcon;
    case `${d}/postgres`:
    case `${t}/postgresql`:
      return PostgresqlIcon;
    case `${t}/newrelic-agent`:
      return NewrelicIcon;
    case `${d}/mysql`:
    case `${t}/myql`:
      return MysqlIcon;
    case `${d}/mongo`:
    case `${t}/mongodb`:
      return MongodbIcon;
    case `${d}/memcached`:
    case `${t}/memcached`:
      return MemcachedIcon;
    case `${d}/mariadb`:
    case `${t}/mariadb`:
      return MariadbIcon;
    case `${d}/joomla`:
    case `${t}/joomla`:
      return JoomlaIcon;
    case `${t}/influxdb`:
      return InfluxdbIcon;
    case `${d}/haproxy`:
    case `${t}/haproxy`:
      return HaproxyIcon;
    case `${d}/glassfish`:
    case `${t}/glassfish`:
      return GlassfishIcon;
    case `${d}/fedora`:
    case `${t}/fedora`:
      return FedoraIcon;
    case `${d}/elasticsearch`:
    case `${t}/elasticsearch`:
      return ElasticsearchIcon;
    case `${d}/drupal`:
    case `${t}/drupal`:
      return DrupalIcon;
    case `${d}/debian`:
    case `${t}/debian`:
      return DebianIcon;
    case `${t}/datadog-agent`:
      return DatadogIcon;
    case `${t}/couchdb`:
      return CouchdbIcon;
    case `${d}/cassandra`:
      return DockerIcon; // icon missing
    case `${t}/cassandra`:
      return TutumIcon; // icon missing
    case `${t}/btsync`:
      return BittorrentSyncIcon;
    case `${t}/wikimedia`:
      return WikimediaIcon;
    case `${d}/ubuntu`:
    case `${t}/ubuntu`:
      return UbuntuIcon;
    default:
      return predicate.split('/')[0] === t ?
      TutumIcon :
      DockerIcon;
  }
};
