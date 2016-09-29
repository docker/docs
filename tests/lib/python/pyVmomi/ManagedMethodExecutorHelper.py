# VMware vSphere Python SDK
# Copyright (c) 2008-2016 VMware, Inc. All Rights Reserved.
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

"""
This module provides convinent fns related to ManagedMethodExecutor
"""
__author__ = "VMware, Inc."

from pyVmomi import VmomiSupport, SoapAdapter, vmodl
from .SoapAdapter import SoapStubAdapterBase, SerializeToUnicode, Deserialize

## ManagedMethodExecutor soap stub adapter
#
class MMESoapStubAdapter(SoapStubAdapterBase):
   """ Managed method executor stub adapter  """

   ## Constructor
   #
   # The endpoint can be specified individually as either a host/port
   # combination, or with a URL (using a url= keyword).
   #
   # @param self self
   # @param mme  managed method executor
   def __init__(self, mme):
      stub = mme._stub
      SoapStubAdapterBase.__init__(self, version=stub.version)
      self.mme = mme

   ## Compute the version information for the specified namespace
   #
   # @param ns the namespace
   def ComputeVersionInfo(self, version):
      SoapStubAdapterBase.ComputeVersionInfo(self, version)
      self.versionId = self.versionId[1:-1]

   ## Invoke a managed method, with _ExecuteSoap. Wohooo!
   #
   # @param self self
   # @param mo the 'this'
   # @param info method info
   # @param args arguments
   def InvokeMethod(self, mo, info, args):
      # Serialize parameters to soap parameters
      methodArgs = None
      if info.params:
         methodArgs = vmodl.Reflect.ManagedMethodExecutor.SoapArgument.Array()
         for param, arg in zip(info.params, args):
            if arg is not None:
               # Serialize parameters to soap snippets
               soapVal = SerializeToUnicode(val=arg, info=param, version=self.version)

               # Insert argument
               soapArg = vmodl.Reflect.ManagedMethodExecutor.SoapArgument(
                                                  name=param.name, val=soapVal)
               methodArgs.append(soapArg)

      moid = mo._GetMoId()
      version = self.versionId
      methodName = VmomiSupport.GetVmodlName(info.type) + "." + info.name

      # Execute method
      result = self.mme.ExecuteSoap(moid=moid,
                                    version=version,
                                    method=methodName,
                                    argument=methodArgs)
      return self._DeserializeExecutorResult(result, info.result)

   ## Invoke a managed property accessor
   #
   # @param self self
   # @param mo the 'this'
   # @param info property info
   def InvokeAccessor(self, mo, info):
      moid = mo._GetMoId()
      version = self.versionId
      prop = info.name

      # Fetch property
      result = self.mme.FetchSoap(moid=moid, version=version, prop=prop)
      return self._DeserializeExecutorResult(result, info.type)

   ## Deserialize result from ExecuteSoap / FetchSoap
   #
   # @param self self
   # @param result result from ExecuteSoap / FetchSoap
   # @param resultType Expected result type
   def _DeserializeExecutorResult(self, result, resultType):
      obj = None
      if result:
         # Parse the return soap snippet. If fault, raise exception
         if result.response:
            # Deserialize back to result
            obj = Deserialize(result.response, resultType, stub=self)
         elif result.fault:
            # Deserialize back to fault (or vmomi fault)
            fault = Deserialize(result.fault.faultDetail,
                                object,
                                stub=self)
            # Silent pylint
            raise fault # pylint: disable-msg=E0702
         else:
            # Unexpected: result should have either response or fault
            msg = "Unexpected execute/fetchSoap error"
            reason = "execute/fetchSoap did not return response or fault"
            raise vmodl.Fault.SystemError(msg=msg, reason=reason)
      return obj
