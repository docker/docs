import json
import os.path

default_vm_config = """{
  "driver": "",
  "vmrun": {
    "vm-folder": ""
  },
  "esx": {
    "username": "",
    "password": "",
    "host": "",
    "datacenter": "",
    "cluster": ""
  }
}
"""


def init_vm_config(fname):
    if not fname:
        fname = "vm.json"
    if os.path.isfile(fname):
        return
    f = open(fname, "w")
    f.write(default_vm_config)
    f.close()
    return


def parse_vm_config(fname):
    if not fname:
        fname = "vm.json"
    if not os.path.isfile(fname):
        raise Exception("Config file not found")
    f = open(fname)
    cfg = f.read()
    data = json.loads(cfg)
    f.close()
    return data
