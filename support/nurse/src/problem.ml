open Sexplib.Std

type t = {
  in_file: string;
  regexp: string;
  label: string;
  description: string;
  link_to_issues: (string * int) list;
  link_from_issues: (string * int) list;
} [@@deriving sexp]

let for_mac = "docker/for-mac"

let transfused_crash = {
  in_file = "docker-system.log";
  regexp = "Transport endpoint is not connected";
  label = "transfused-crash";
  description = "This error suggests that a component of the volume sharing system failed.";
  link_to_issues = [ for_mac, 5 ];
  link_from_issues = [];
}

let qcow2_corruption = {
  in_file = "docker-system.log";
  regexp = "Invalid_argument\\(\"Cstruct.sub";
  label = "qcow2-corruption";
  description = "This error suggests that the .qcow2 file is corrupt and will need to be recreated.";
  link_to_issues = [ for_mac, 11 ];
  link_from_issues = [];
}

let var_lib_docker_corruption1 = {
  in_file = "moby/var/log/docker.log";
  regexp = "invalid checksum digest format";
  label = "aufs-corruption";
  description = "The daemon cannot start because some files in /var/lib/docker are corrupt. \
  This can happen if the filesystem is unmounted uncleanly. The workaround is to reset to \
  factory defaults which will unfortunately delete /var/lib/docker: all containers and \
  images will need to be rebuilt.";
  link_to_issues = [ for_mac, 20 ];
  link_from_issues = [];
}

let var_lib_docker_corruption2 = { var_lib_docker_corruption1 with regexp = "invalid mount id value" }

let docker_didnt_start = {
  in_file = "moby/var/log/docker.log";
  regexp = "Error starting daemon:";
  label = "daemon-didnt-start";
  description = "The docker daemon inside the VM failed to start.";
  link_to_issues = [];
  link_from_issues = [];
}

let communication_with_networking = {
  in_file = "docker-system.log";
  regexp = "Communication with networking components failed";
  label = "vmnetd-comm-failure";
  description = "The UI failed to communicate with the networking helper process.";
  link_to_issues = [ for_mac, 61 ];
  link_from_issues = [];
}

let too_many_hyperkits = {
  in_file = "ps-ax.log";
  regexp = "hyperkit -A -m(.|\n)+hyperkit -A -m";
  label = "more-than-one-hyperkit";
  description = "More than one hyperkit process has started and is accessing the \
  qcow2 disk at the same time. This will lead to corruption and Docker must be \
  reset to factory defaults.";
  link_to_issues = [ for_mac, 71 ];
  link_from_issues = [];
}

let virtualbox = {
  in_file = "docker-system.log";
  regexp = "Please upgrade or uninstall Virtualbox";
  label = "uninstall-virtualbox";
  description = "Some parts of Virtualbox may still be installed. Please follow \
  https://www.virtualbox.org/manual/ch02.html#idm871 to completely uninstall \
  Virtualbox. Please note that dragging the application to the trash is not enough:\
  there are kernel modules installed as well.";
  link_to_issues = [ for_mac, 78 ];
  link_from_issues = [];
}

let moby_kernel = {
  in_file = "docker-system.log";
  regexp = "rcu_sched self-detected stall";
  label = "moby-kernel-froze";
  description = "The kernel inside the VM stalled so Docker must be restarted.";
  link_to_issues = [ for_mac, 87 ];
  link_from_issues = [];
}

let invariants = {
  in_file = "docker-system.log";
  regexp = "INVARIANT VIOLATED";
  label = "invariant-violated";
  description = "An invariant within the code of one of the components was \
  not held. Docker needs to be restarted.";
  link_to_issues = [ for_mac, 89 ];
  link_from_issues = [];
}

let compiled_in = [
  transfused_crash;
  qcow2_corruption;
  var_lib_docker_corruption1;
  var_lib_docker_corruption2;
  docker_didnt_start;
  communication_with_networking;
  too_many_hyperkits;
  virtualbox;
  moby_kernel;
  invariants;
]
