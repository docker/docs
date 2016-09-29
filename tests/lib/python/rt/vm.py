from abc import ABCMeta, abstractmethod
import atexit
import argparse
import uuid
import os
import platform
from subprocess import call

from pyVmomi import vim, vmodl
from pyVim.connect import SmartConnect, Disconnect

import rt.vmconfig as config


def main():
    # top level
    parser = argparse.ArgumentParser(
        description="Manage virtual machines for regression tests",
        prog="rt-vm"
        )
    parser.add_argument('--file',
                        help="location of configuration file")
    parser.add_argument('--debug', dest='debug',
                        action='store_true', help="debugging on")
    parser.set_defaults(debug=False)
    subparsers = parser.add_subparsers()

    # create the parser for the "init" command
    parser_init = subparsers.add_parser(
        'init',
        help="Initialize the configuration file"
    )

    parser_init.set_defaults(func=init)
    # create the parser for the "create" command
    parser_create = subparsers.add_parser('create', help="Create a VM")
    parser_create.add_argument('os', type=str)
    parser_create.add_argument('version', type=str)
    parser_create.set_defaults(func=create)
    # create the parser for the "ip" command
    parser_ip = subparsers.add_parser('ip', help="Get the IP of a VM")
    parser_ip.add_argument('name', type=str)
    parser_ip.set_defaults(func=ip)
    # create the parser for the "start" command
    parser_start = subparsers.add_parser('start', help="Start a VM")
    parser_start.add_argument('name', type=str)
    parser_start.set_defaults(func=start)
    # create the parser for the "stop" command
    parser_stop = subparsers.add_parser('stop', help="Stop a VM")
    parser_stop.add_argument('name', type=str)
    parser_stop.set_defaults(func=stop)
    # create the parser for the "destroy" command
    parser_destroy = subparsers.add_parser('destroy', help="Destroy a VM")
    parser_destroy.add_argument('name', type=str)
    parser_destroy.set_defaults(func=destroy)
    # create the parser for the "revert" command
    parser_revert = subparsers.add_parser('revert',
                                          help="Revert a VM to a snapshot")
    parser_revert.add_argument('name', type=str)
    parser_revert.add_argument('snapshot', type=str)
    parser_revert.set_defaults(func=revert)
    # create the parser for the "upload" command
    parser_revert = subparsers.add_parser(
        'upload',
        help="Upload a template to the server"
    )
    parser_revert.add_argument('path', type=str)
    parser_revert.add_argument('datastore', type=str)
    parser_revert.set_defaults(func=upload)

    args = parser.parse_args()
    return args.func(args)


def get_driver(cfg):
    if cfg["driver"] == "esx":
        return EsxDriver(cfg)
    elif cfg["driver"] == "vmrun":
        return VmrunDriver(cfg)
    else:
        raise Exception("Unsupported driver")


def init(args):
    config.init_vm_config(args.file)
    return 0


def create(args):
    try:
        cfg = config.parse_vm_config(args.file)
        driver = get_driver(cfg)
        driver.create(args.os, args.version)
    except Exception as e:
        print("Error: {}".format(e))
        if args.debug:
            import traceback
            traceback.print_exc()
        return 1
    return 0


def ip(args):
    try:
        cfg = config.parse_vm_config(args.file)
        driver = get_driver(cfg)
        driver.ip(args.name)
    except Exception as e:
        print("Error: {}".format(e))
        if args.debug:
            import traceback
            traceback.print_exc()
        return 1
    return 0


def start(args):
    try:
        cfg = config.parse_vm_config(args.file)
        driver = get_driver(cfg)
        driver.start(args.name)
    except Exception as e:
        print("Error: {}".format(e))
        if args.debug:
            import traceback
            traceback.print_exc()
        return 1
    return 0


def stop(args):
    try:
        cfg = config.parse_vm_config(args.file)
        driver = get_driver(cfg)
        driver.stop(args.name)
    except Exception as e:
        print("Error: {}".format(e))
        if args.debug:
            import traceback
            traceback.print_exc()
        return 1
    return 0


def destroy(args):
    try:
        cfg = config.parse_vm_config(args.file)
        driver = get_driver(cfg)
        driver.destroy(args.name)
    except Exception as e:
        print("Error: {}".format(e))
        if args.debug:
            import traceback
            traceback.print_exc()
        return 1
    return 0


def revert(args):
    try:
        cfg = config.parse_vm_config(args.file)
        driver = get_driver(cfg)
        driver.revert(args.name, args.snapshot)
    except Exception as e:
        print("Error: {}".format(e))
        if args.debug:
            import traceback
            traceback.print_exc()
        return 1
    return 0


def upload(args):
    try:
        cfg = config.parse_vm_config(args.file)
        driver = get_driver(cfg)
        driver.upload(args.path, args.datastore)
    except Exception as e:
        print ("Error: {}".format(e))
        if args.debug:
            import traceback
            traceback.print_exc()
        return 1
    return 0


def get_vm_name(os, version):
    s = str(uuid.uuid4())[-12:]
    return "test-{0}-{1}-{2}".format(os, version, s)


def get_template_name(os, version):
    return "test-{0}-{1}-template".format(os, version)


class BaseDriver(object):
    __metaclass__ = ABCMeta

    def __init__(self, config):
        self.config = config

    @abstractmethod
    def create(self, os, version):
        return

    @abstractmethod
    def ip(self, name):
        return

    @abstractmethod
    def start(self, name):
        return

    @abstractmethod
    def stop(self, name):
        return

    @abstractmethod
    def destroy(self, name):
        return

    @abstractmethod
    def revert(self, name, snapshot):
        return

    @abstractmethod
    def upload(self, name, path, datastore):
        return


class EsxDriver(BaseDriver):

    def __init__(self, config):
        super(EsxDriver, self).__init__(config)
        self.validate()

    def validate(self):
        """
        Validate the configuration
        """
        if (
            not self.config["esx"] or
            not self.config["esx"]["username"] or
            not self.config["esx"]["password"] or
            not self.config["esx"]["host"] or
            not self.config["esx"]["datacenter"] or
            not self.config["esx"]["cluster"]
           ):
            raise Exception("Invalid ESX Config")

        if "port" not in self.config["esx"].keys():
            self.config["esx"]["port"] = 443

    def create(self, os, version):
        """
        Create a VM
        This will clone a template and create a new VM for testing
        """
        self.connect()
        template_name = get_template_name(os, version)
        vm_name = get_vm_name(os, version)

        template = None
        template = self.get_obj([vim.VirtualMachine], template_name)
        if not template:
            raise Exception("Template {} not found".format(template_name))

        datacenter = self.get_obj(
            [vim.Datacenter],
            self.config["esx"]["datacenter"]
        )
        destfolder = datacenter.vmFolder
        cluster = self.get_obj(
            [vim.ClusterComputeResource],
            self.config["esx"]["cluster"]
        )
        resource_pool = cluster.resourcePool

        relospec = vim.vm.RelocateSpec(
            pool=resource_pool,
            host=template.runtime.host,
            diskMoveType='createNewChildDiskBacking'
        )

        clonespec = vim.vm.CloneSpec(
            location=relospec,
            powerOn=True,
            template=False,
            snapshot=template.snapshot.rootSnapshotList[0].snapshot
        )

        task = template.Clone(folder=destfolder, name=vm_name, spec=clonespec)
        self.wait_for_task(task)
        print(vm_name)

    def ip(self, name):
        """
        Get a VM's IP Address
        """
        self.connect()
        vm = self.get_obj([vim.VirtualMachine], name)
        if vm is not None:
            summary = vm.summary
        else:
            raise Exception("VM {} does not exist".format(name))
        if summary.guest is not None:
            ip = summary.guest.ipAddress
            if ip is not None:
                print(ip)
            else:
                print("")
        else:
            raise Exception("Cannot query VM tools for VM {}".format(name))

    def start(self, name):
        """
        Start a VM
        """
        self.connect()
        vm = self.get_obj([vim.VirtualMachine], name)
        if format(vm.runtime.powerState) == "poweredOff":
            if vm is not None:
                task = vm.PowerOnVM_Task()
                self.wait_for_task(task)
            else:
                raise Exception("VM {} does not exist".format(name))

    def stop(self, name):
        """
        Stop a VM
        """
        self.connect()
        vm = self.get_obj([vim.VirtualMachine], name)
        if format(vm.runtime.powerState) == "poweredOn":
            if vm is not None:
                task = vm.PowerOffVM_Task()
                self.wait_for_task(task)
            else:
                raise Exception("VM {} does not exist".format(name))

    def destroy(self, name):
        """
        Destroy a VM
        """
        self.connect()
        vm = self.get_obj([vim.VirtualMachine], name)
        if vm is not None:
            task = vm.Destroy_Task()
            self.wait_for_task(task)
        else:
            raise Exception("VM {} does not exist".format(name))

    def revert(self, name, snapshot_name):
        """
        Revert a VM to a snapshot
        """
        self.connect()
        vm = self.get_obj([vim.VirtualMachine], name)
        if vm is not None:
            snapshots = vm.snapshot.rootSnapshotList
            for snapshot in snapshots:
                if snapshot.name == snapshot_name:
                    s = snapshot.snapshot
                    task = s.RevertToSnapshot_Task()
                    self.wait_for_task(task)
                    return
            raise Exception("Snapshot not found")
        else:
            raise Exception("VM {} does not exist".format(name))

    def upload(self, path, datastore):
        self.connect()
        # 1) Delete existing template
        _, name = os.path.split(path)
        old_template = self.get_obj([vim.VirtualMachine], name)

        if old_template is not None:
            self.destroy(name)

        # 2) ovftool upload
        vi_string = "vi://%s:%s@%s/%s/host/%s" % (
            self.config["esx"]["username"],
            self.config["esx"]["password"],
            self.config["esx"]["host"],
            self.config["esx"]["datacenter"],
            self.config["esx"]["cluster"]
        )

        upload_cmd = [
            "ovftool",
            "--name=%s" % (name),
            "--datastore=%s" % (datastore),
            "--diskMode=thin",
            "--noSSLVerify",
            "--acceptAllEulas",
            "--allowExtraConfig",
            "--extraConfig:smc.present=TRUE",
            "%s" % (os.path.join(path, "%s.vmx" % (name))),
            vi_string
        ]
        rc = call(upload_cmd)
        if rc != 0:
            raise Exception("ovftool failed with exit code %d" % (rc))

        # 3) Set nested virt

        # Initial upload took a while, disconnect and reconnect
        Disconnect(self.si)
        self.connect(
                )
        template = self.get_obj([vim.VirtualMachine], name)
        spec = vim.vm.ConfigSpec()
        hardware_virt = vim.option.OptionValue(key="featMask.vm.hv.capable",
                                               value="Min:1")
        spec.nestedHVEnabled = True
        spec.extraConfig = [hardware_virt]
        task = template.ReconfigVM_Task(spec)
        self.wait_for_task(task)

        # 4) Take snapshot
        task = template.CreateSnapshot_Task(name=name,
                                            description="Initial Snapshot",
                                            memory=False,
                                            quiesce=False)
        self.wait_for_task(task)

        # 4) Mark as template
        template.MarkAsTemplate()

    def get_obj(self, vimtype, name):
        """
        Return an object by name, if name is None the
        first found object is returned
        """
        obj = None
        content = self.si.RetrieveContent()
        container = content.viewManager.CreateContainerView(
            content.rootFolder, vimtype, True)
        for c in container.view:
            if c.name == name:
                obj = c
                break
        return obj

    def wait_for_task(self, task):
        """ wait for a vCenter task to finish """
        content = self.si.RetrieveContent()
        property_collector = content.propertyCollector
        obj_specs = [vmodl.query.PropertyCollector.ObjectSpec(obj=task)]
        property_spec = vmodl.query.PropertyCollector.PropertySpec(
            type=vim.Task,
            pathSet=[],
            all=True
        )
        filter_spec = vmodl.query.PropertyCollector.FilterSpec()
        filter_spec.objectSet = obj_specs
        filter_spec.propSet = [property_spec]
        pcfilter = property_collector.CreateFilter(filter_spec, True)

        try:
            version, state = None, None
            completed = False
            # Loop looking for updates till the state moves
            # to a completed state.
            while not completed:
                update = property_collector.WaitForUpdates(version)
                for filter_set in update.filterSet:
                    for obj_set in filter_set.objectSet:
                        t = obj_set.obj
                        for change in obj_set.changeSet:
                            if change.name == 'info':
                                state = change.val.state
                            elif change.name == 'info.state':
                                state = change.val
                            else:
                                continue

                            if not str(t) == str(task):
                                continue

                            if state == vim.TaskInfo.State.success:
                                # Remove task from taskList
                                completed = True
                            elif state == vim.TaskInfo.State.error:
                                raise task.info.error
                # Move to next version
                version = update.version
        finally:
            if pcfilter:
                pcfilter.Destroy()

    def connect(self):
        import ssl
        if hasattr(ssl, '_create_unverified_context'):
            ssl_context = ssl._create_unverified_context()
        else:
            ssl_context = None

        self.si = SmartConnect(
            host=self.config["esx"]["host"],
            user=self.config["esx"]["username"],
            pwd=self.config["esx"]["password"],
            port=self.config["esx"]["port"],
            sslContext=ssl_context)
        atexit.register(Disconnect, self.si)


class VmrunDriver(BaseDriver):

    def __init__(self, config):
        super(VmrunDriver, self).__init__(config)
        if not self.config["vmrun"]["vm-folder"]:
            raise Exception("Invalid Config")

        # try to find vmrun, first in the path then in the default location
        if platform.system() == "Windows":
            vmrun = "vmrun.exe"
            # XXX is this the correct path?
            default_path = "C:\Program Files\VMware\VMware VIX"
        else:
            vmrun = "vmrun"
            default_path = "/Applications/VMware Fusion.app/Contents/Library"

        paths = os.environ["PATH"].split(os.pathsep)
        paths.append(default_path)
        self.vmrunpath = None
        for path in paths:
            path = path.strip('"')
            tmp = os.path.join(path, vmrun)
            if os.path.isfile(tmp) and os.access(tmp, os.X_OK):
                self.vmrunpath = tmp
                break
        if not self.vmrunpath:
            raise Exception("could not find vmrun in PATH nor in %s" %
                            default_path)

    def create(self, os, version):
        template_name = get_template_name(os, version)
        vm_name = get_vm_name(os, version)

        template_vmx = self.get_vmx(template_name)
        vm_vmx = self.get_vmx(vm_name)

        self.vmrun('clone "{0}" "{1}" linked'.format(template_vmx, vm_vmx))
        print(vm_name)

    def ip(self, name):
        vmx = self.get_vmx(name)
        self.vmrun('readVariable "{}" guestVar ip'.format(vmx))

    def start(self, name):
        vmx = self.get_vmx(name)
        self.vmrun('start "{}"'.format(vmx))

    def stop(self, name):
        vmx = self.get_vmx(name)
        self.vmrun('stop "{}" soft'.format(vmx))

    def destroy(self, name):
        vmx = self.get_vmx(name)
        self.vmrun('deleteVM "{}"'.format(vmx))

    def revert(self, name, snapshot):
        vmx = self.get_vmx(name)
        self.vmrun('revertToSnapshot "{0}" "{1}"'.format(vmx, snapshot))

    def upload(self, path, datastore):
        raise Exception("Operation not supported in the vmrun driver")

    def vmrun(self, args):
        if platform.system() == "Windows":
            t = "-T ws"
        else:
            t = "-T fusion"
        command = '"{}" {} {}'.format(self.vmrunpath, t, args)
        code = call(command, shell=True)
        if code != 0:
            raise Exception("vmrun {} : failed with exit code {}".format(args,
                                                                         code))

    def get_vmx(self, name):
        return os.path.join(self.config["vmrun"]["vm-folder"],
                            name,
                            "{}.vmx".format(name))
