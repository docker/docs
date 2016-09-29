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
from pyVmomi.VmomiSupport import GetVmodlType

class StubAdapterAccessorMixin:
   def __init__(self):
      self._pc = None
      self._pcType = GetVmodlType("vmodl.query.PropertyCollector")
      self._siType = GetVmodlType("vim.ServiceInstance")

   ## Retrieve a managed property
   #
   # @param self self
   # @param mo managed object
   # @param info property info
   def InvokeAccessor(self, mo, info):
      filterSpec = self._pcType.FilterSpec(
         objectSet=[self._pcType.ObjectSpec(obj=mo, skip=False)],
         propSet=[self._pcType.PropertySpec(all=False, type=mo.__class__,
                                                 pathSet=[info.name])],
         )
      ## Cache the property collector if it isn't already
      #  No need to lock _pc since multiple instances of PropertyCollector on
      #  the client will talk to the same instance on the server.
      if not self._pc:
         si = self._siType("ServiceInstance", self)
         self._pc = si.RetrieveContent().propertyCollector
      result = self._pc.RetrievePropertiesEx(specSet=[filterSpec],
                                             options=self._pcType.RetrieveOptions(maxObjects=1))
      objectContent = result.objects[0]
      if len(objectContent.propSet) > 0:
         return objectContent.propSet[0].val
      if len(objectContent.missingSet) > 0 and objectContent.missingSet[0].fault:
         raise objectContent.missingSet[0].fault
      return None
