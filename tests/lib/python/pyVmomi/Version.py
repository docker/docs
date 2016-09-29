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
from __future__ import absolute_import
from pyVmomi.VmomiSupport import nsMap, versionMap, versionIdMap, serviceNsMap, parentMap

## Add an API version
def AddVersion(version, ns, versionId='', isLegacy=0, serviceNs=''):
  if not ns:
     ns = serviceNs
  if not (version in parentMap):
      nsMap[version] = ns
      if len(versionId) > 0:
         versionMap[ns + '/' + versionId] = version
      if isLegacy or ns is "":
         versionMap[ns] = version
      versionIdMap[version] = versionId
      if not serviceNs:
         serviceNs = ns
      serviceNsMap[version] = serviceNs
      parentMap[version] = {}

## Check if a version is a child of another
def IsChildVersion(child, parent):
   return child == parent or parent in parentMap[child]
