# VMware vSphere Python SDK
# Copyright (c) 2008-2015 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# ******* WARNING - AUTO GENERATED CODE - DO NOT EDIT *******
from __future__ import absolute_import
from pyVmomi.VmomiSupport import CreateDataType, CreateManagedType, CreateEnumType, AddVersion, AddVersionParent, F_LINK, F_LINKABLE, F_OPTIONAL, F_SECRET
from pyVmomi.VmomiSupport import newestVersions, currentVersions, stableVersions, matureVersions, publicVersions, oldestVersions

AddVersion("vmodl.version.version2", "", "", 0, "vim25")
AddVersion("vmodl.version.version1", "", "", 0, "vim25")
AddVersion("vmodl.version.version0", "", "", 0, "vim25")
AddVersionParent("vmodl.version.version2", "vmodl.version.version2")
AddVersionParent("vmodl.version.version2", "vmodl.version.version1")
AddVersionParent("vmodl.version.version2", "vmodl.version.version0")
AddVersionParent("vmodl.version.version1", "vmodl.version.version1")
AddVersionParent("vmodl.version.version1", "vmodl.version.version0")
AddVersionParent("vmodl.version.version0", "vmodl.version.version0")

newestVersions.Add("vmodl.version.version2")
currentVersions.Add("vmodl.version.version2")
stableVersions.Add("vmodl.version.version2")
matureVersions.Add("vmodl.version.version2")
publicVersions.Add("vmodl.version.version2")
oldestVersions.Add("vmodl.version.version0")

CreateDataType("vmodl.DynamicArray", "DynamicArray", "vmodl.DataObject", "vmodl.version.version0", [("dynamicType", "string", "vmodl.version.version0", F_OPTIONAL), ("val", "anyType[]", "vmodl.version.version0", 0)])
CreateDataType("vmodl.DynamicData", "DynamicData", "vmodl.DataObject", "vmodl.version.version0", [("dynamicType", "string", "vmodl.version.version0", F_OPTIONAL), ("dynamicProperty", "vmodl.DynamicProperty[]", "vmodl.version.version0", F_OPTIONAL)])
CreateDataType("vmodl.DynamicProperty", "DynamicProperty", "vmodl.DataObject", "vmodl.version.version0", [("name", "vmodl.PropertyPath", "vmodl.version.version0", 0), ("val", "anyType", "vmodl.version.version0", 0)])
CreateDataType("vmodl.KeyAnyValue", "KeyAnyValue", "vmodl.DynamicData", "vmodl.version.version1", [("key", "string", "vmodl.version.version1", 0), ("value", "anyType", "vmodl.version.version1", 0)])
CreateDataType("vmodl.LocalizableMessage", "LocalizableMessage", "vmodl.DynamicData", "vmodl.version.version1", [("key", "string", "vmodl.version.version1", 0), ("arg", "vmodl.KeyAnyValue[]", "vmodl.version.version1", F_OPTIONAL), ("message", "string", "vmodl.version.version1", F_OPTIONAL)])
CreateDataType("vmodl.fault.HostCommunication", "HostCommunication", "vmodl.RuntimeFault", "vmodl.version.version0", None)
CreateDataType("vmodl.fault.HostNotConnected", "HostNotConnected", "vmodl.fault.HostCommunication", "vmodl.version.version0", None)
CreateDataType("vmodl.fault.HostNotReachable", "HostNotReachable", "vmodl.fault.HostCommunication", "vmodl.version.version0", None)
CreateDataType("vmodl.fault.InvalidArgument", "InvalidArgument", "vmodl.RuntimeFault", "vmodl.version.version0", [("invalidProperty", "vmodl.PropertyPath", "vmodl.version.version0", F_OPTIONAL)])
CreateDataType("vmodl.fault.InvalidRequest", "InvalidRequest", "vmodl.RuntimeFault", "vmodl.version.version0", None)
CreateDataType("vmodl.fault.InvalidType", "InvalidType", "vmodl.fault.InvalidRequest", "vmodl.version.version0", [("argument", "vmodl.PropertyPath", "vmodl.version.version0", F_OPTIONAL)])
CreateDataType("vmodl.fault.ManagedObjectNotFound", "ManagedObjectNotFound", "vmodl.RuntimeFault", "vmodl.version.version0", [("obj", "vmodl.ManagedObject", "vmodl.version.version0", 0)])
CreateDataType("vmodl.fault.MethodNotFound", "MethodNotFound", "vmodl.fault.InvalidRequest", "vmodl.version.version0", [("receiver", "vmodl.ManagedObject", "vmodl.version.version0", 0), ("method", "string", "vmodl.version.version0", 0)])
CreateDataType("vmodl.fault.NotEnoughLicenses", "NotEnoughLicenses", "vmodl.RuntimeFault", "vmodl.version.version0", None)
CreateDataType("vmodl.fault.NotImplemented", "NotImplemented", "vmodl.RuntimeFault", "vmodl.version.version0", None)
CreateDataType("vmodl.fault.NotSupported", "NotSupported", "vmodl.RuntimeFault", "vmodl.version.version0", None)
CreateDataType("vmodl.fault.RequestCanceled", "RequestCanceled", "vmodl.RuntimeFault", "vmodl.version.version0", None)
CreateDataType("vmodl.fault.SecurityError", "SecurityError", "vmodl.RuntimeFault", "vmodl.version.version0", None)
CreateDataType("vmodl.fault.SystemError", "SystemError", "vmodl.RuntimeFault", "vmodl.version.version0", [("reason", "string", "vmodl.version.version0", 0)])
CreateDataType("vmodl.fault.UnexpectedFault", "UnexpectedFault", "vmodl.RuntimeFault", "vmodl.version.version0", [("faultName", "vmodl.TypeName", "vmodl.version.version0", 0), ("fault", "vmodl.MethodFault", "vmodl.version.version0", F_OPTIONAL)])
