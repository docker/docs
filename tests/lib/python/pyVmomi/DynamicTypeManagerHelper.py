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

"""
This module is a converter from dynamic type to pyVmomi type
"""
__author__ = "VMware, Inc."

from pyVmomi import VmomiSupport, vmodl
from Cache import Cache

## Dynamic type importer
#
class DynamicTypeImporter:
   """ Dynamic type importer  """

   ## Constructor
   #
   # @param  stub Server stub
   def __init__(self, stub, hostSystem=None):
      self.stub = stub
      self.hostSystem = hostSystem

   ## Get dynamic type manager
   #
   # @return moRef to dynamic type manager
   @Cache
   def GetTypeManager(self):
      """ Get dynamic type manager """
      dynTypeMgr = None
      if self.hostSystem:
         try:
            dynTypeMgr = self.hostSystem.RetrieveDynamicTypeManager()
         except vmodl.fault.MethodNotFound as err:
            pass

      if not dynTypeMgr:
         # Older host not support RetrieveDynamicTypeManager
         cmdlineTypesMoId = "ha-dynamic-type-manager"
         dynTypeMgr = vmodl.reflect.DynamicTypeManager(cmdlineTypesMoId,
                                                       self.stub)
      return dynTypeMgr

   ## Import dynamic types
   #
   # @param  prefix Only types with the specified prefix are imported
   # @return dynamic types
   @Cache
   def ImportTypes(self, prefix=''):
      """ Build dynamic types """
      # Use QueryTypeInfo to get all types
      dynTypeMgr = self.GetTypeManager()
      filterSpec = None
      if prefix != '':
         filterSpec = vmodl.reflect.DynamicTypeManager.TypeFilterSpec(
                                                              typeSubstr=prefix)
      allTypes = dynTypeMgr.QueryTypeInfo(filterSpec)

      ## Convert dynamic types to pyVmomi types
      #
      DynamicTypeConstructor().CreateTypes(allTypes)
      return allTypes


## Construct pyVmomi types from dynamic types definition
#
class DynamicTypeConstructor:
   """ Dynamic type constructor  """

   _mapFlags = { "optional": VmomiSupport.F_OPTIONAL,
                 "linkable": VmomiSupport.F_LINKABLE,
                 "link":     VmomiSupport.F_LINK,
                 "secret":   VmomiSupport.F_SECRET }

   ## Constructor
   #
   def __init__(self):
      """ Constructor """
      pass

   ## Create pyVmomi types from vmodl.reflect.DynamicTypeManager.AllTypeInfo
   #
   # @param  allTypes vmodl.reflect.DynamicTypeManager.AllTypeInfo
   def CreateTypes(self, allTypes):
      """
      Create pyVmomi types from vmodl.reflect.DynamicTypeManager.AllTypeInfo
      """
      enumTypes, dataTypes, managedTypes = self._ConvertAllTypes(allTypes)
      self._CreateAllTypes(enumTypes, dataTypes, managedTypes)

   ## Convert all dynamic types to pyVmomi type definitions
   #
   # @param  allTypes vmodl.reflect.DynamicTypeManager.AllTypeInfo
   # @return a tuple of (enumTypes, dataTypes, managedTypes)
   def _ConvertAllTypes(self, allTypes):
      """ Convert all dynamic types to pyVmomi type definitions """
      # Generate lists good for VmomiSupport.CreateXYZType
      enumTypes = self._Filter(self._ConvertEnumType, allTypes.enumTypeInfo)
      dataTypes = self._Filter(self._ConvertDataType, allTypes.dataTypeInfo)
      managedTypes = self._Filter(self._ConvertManagedType,
                                  allTypes.managedTypeInfo)
      retAllTypes = (enumTypes, dataTypes, managedTypes)
      return retAllTypes

   ## Create pyVmomi types from pyVmomi type definitions
   #
   # @param  enumTypes pyVmomi enum type definitions
   # @param  dataTypes pyVmomi data type definitions
   # @param  managedTypes pyVmomi managed type definitions
   def _CreateAllTypes(self, enumTypes, dataTypes, managedTypes):
      """ Create pyVmomi types from pyVmomi type definitions """

      # Create versions
      for typeInfo in managedTypes:
         name = typeInfo[0]
         version = typeInfo[3]
         VmomiSupport.AddVersion(version, '', '1.0', 0, name)
         VmomiSupport.AddVersionParent(version, 'vmodl.version.version0')
         VmomiSupport.AddVersionParent(version, 'vmodl.version.version1')
         VmomiSupport.AddVersionParent(version, version)

      # Create partial types
      for fn, infos in (VmomiSupport.CreateEnumType, enumTypes), \
                       (VmomiSupport.CreateDataType, dataTypes), \
                       (VmomiSupport.CreateManagedType, managedTypes):
         for typeInfo in infos:
            try:
               fn(*typeInfo)
            except Exception as err:
               #Ignore errors due to duplicate importing
               pass

   def _ConvertAnnotations(self, annotations):
      """ Convert annotations to pyVmomi flags """
      flags = 0
      if annotations:
         for annotation in annotations:
            flags |= self._mapFlags.get(annotation.name, 0)
      return flags

   @staticmethod
   def _Filter(fn, types):
      """ Call fn for each non null element in types. Similiar to filter """
      if types:
         return [fn(prop) for prop in types if prop is not None]
      else:
         return []

   def _ConvertParamType(self, paramType):
      """
      Convert vmodl.reflect.DynamicTypeManager.ParamTypeInfo to pyVmomi param
      definition
      """
      if paramType:
         name = paramType.name
         version = paramType.version
         aType = paramType.type
         flags = self._ConvertAnnotations(paramType.annotation)
         privId = paramType.privId
         param = (name, aType, version, flags, privId)
      else:
         param = None
      return param

   def _ConvertMethodType(self, methodType):
      """
      Convert vmodl.reflect.DynamicTypeManager.MethodTypeInfo to pyVmomi method
      definition
      """
      if methodType:
         name = methodType.name
         wsdlName = methodType.wsdlName
         version = methodType.version
         params = self._Filter(self._ConvertParamType, methodType.paramTypeInfo)
         privId = methodType.privId
         faults = methodType.fault

         # Figure out reture info
         if methodType.returnTypeInfo:
            returnTypeInfo = methodType.returnTypeInfo
            retFlags = self._ConvertAnnotations(returnTypeInfo.annotation)
            methodRetType = returnTypeInfo.type
         else:
            retFlags = 0
            methodRetType = "void"
         if wsdlName.endswith("_Task"):
            # TODO: Need a seperate task return type for task, instead of
            #       hardcode vim.Task as return type
            retType = "vim.Task"
         else:
            retType = methodRetType
         retInfo = (retFlags, retType, methodRetType)

         method = (name, wsdlName, version, params, retInfo, privId, faults)
      else:
         method = None
      return method

   def _ConvertManagedPropertyType(self, propType):
      """
      Convert vmodl.reflect.DynamicTypeManager.PropertyTypeInfo to pyVmomi
      managed property definition
      """
      if propType:
         name = propType.name
         version = propType.version
         aType = propType.type
         flags = self._ConvertAnnotations(propType.annotation)
         privId = propType.privId
         prop = (name, aType, version, flags, privId)
      else:
         prop = None
      return prop

   def _ConvertManagedType(self, managedType):
      """
      Convert vmodl.reflect.DynamicTypeManager.ManagedTypeInfo to pyVmomi
      managed type definition
      """
      if managedType:
         vmodlName = managedType.name
         wsdlName = managedType.wsdlName
         version = managedType.version
         parent = managedType.base[0]
         props = self._Filter(self._ConvertManagedPropertyType, managedType.property)
         methods = self._Filter(self._ConvertMethodType, managedType.method)
         moType = (vmodlName, wsdlName, parent, version, props, methods)
      else:
         moType = None
      return moType

   def _ConvertDataPropertyType(self, propType):
      """
      Convert vmodl.reflect.DynamicTypeManager.PropertyTypeInfo to pyVmomi
      data property definition
      """
      if propType:
         name = propType.name
         version = propType.version
         aType = propType.type
         flags = self._ConvertAnnotations(propType.annotation)
         prop = (name, aType, version, flags)
      else:
         prop = None
      return prop

   def _ConvertDataType(self, dataType):
      """
      Convert vmodl.reflect.DynamicTypeManager.DataTypeInfo to pyVmomi data
      type definition
      """
      if dataType:
         vmodlName = dataType.name
         wsdlName = dataType.wsdlName
         version = dataType.version
         parent = dataType.base[0]
         props = self._Filter(self._ConvertDataPropertyType, dataType.property)
         doType = (vmodlName, wsdlName, parent, version, props)
      else:
         doType = None
      return doType

   def _ConvertEnumType(self, enumType):
      """
      Convert vmodl.reflect.DynamicTypeManager.EnumTypeInfo to pyVmomi enum
      type definition
      """
      if enumType:
         vmodlName = enumType.name
         wsdlName = enumType.wsdlName
         version = enumType.version
         values = enumType.value
         enumType = (vmodlName, wsdlName, version, values)
      else:
         enumType = None
      return enumType

