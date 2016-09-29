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

## Diff any two objects
import six
from six.moves import zip

from pyVmomi import VmomiSupport, types
import logging
from VmomiSupport import GetWsdlName, Type

__Log__ = logging.getLogger('ObjDiffer')

def LogIf(condition, message):
   """Log a message if the condition is met"""
   if condition:
      __Log__.debug(message)

def IsPrimitiveType(obj):
   """See if the passed in type is a Primitive Type"""
   return (isinstance(obj, types.bool) or isinstance(obj, types.byte) or
      isinstance(obj, types.short) or isinstance(obj, six.integer_types) or
      isinstance(obj, types.double) or isinstance(obj, types.float) or
      isinstance(obj, six.string_types) or
      isinstance(obj, types.PropertyPath) or
      isinstance(obj, types.ManagedMethod) or
      isinstance(obj, types.datetime) or
      isinstance(obj, types.URI) or isinstance(obj, type))


class Differ:
   """Class for comparing two Objects"""
   def __init__(self, looseMatch=False, ignoreArrayOrder=True):
      self._looseMatch = looseMatch
      self._ignoreArrayOrder = ignoreArrayOrder

   def DiffAnyObjects(self, oldObj, newObj, isObjLink=False):
      """Diff any two Objects"""
      if oldObj == newObj:
         return True
      if not oldObj or not newObj:
         __Log__.debug('DiffAnyObjects: One of the objects is unset.')
         return self._looseMatch
      oldObjInstance = oldObj
      newObjInstance = newObj
      if isinstance(oldObj, list):
         oldObjInstance = oldObj[0]
      if isinstance(newObj, list):
         newObjInstance = newObj[0]
      # Need to see if it is a primitive type first since type information
      #   will not be available for them.
      if (IsPrimitiveType(oldObj) and IsPrimitiveType(newObj)
         and oldObj.__class__.__name__ == newObj.__class__.__name__):
         if oldObj == newObj:
            return True
         elif oldObj == None or newObj == None:
            __Log__.debug('DiffAnyObjects: One of the objects in None')
         return False
      oldType = Type(oldObjInstance)
      newType = Type(newObjInstance)
      if oldType != newType:
         __Log__.debug('DiffAnyObjects: Types do not match %s != %s' %
            (repr(GetWsdlName(oldObjInstance.__class__)),
             repr(GetWsdlName(newObjInstance.__class__))))
         return False
      elif isinstance(oldObj, list):
         return self.DiffArrayObjects(oldObj, newObj, isObjLink)
      elif isinstance(oldObjInstance, types.ManagedObject):
         return (not oldObj and not newObj) or (oldObj and newObj
            and oldObj._moId == newObj._moId)
      elif isinstance(oldObjInstance, types.DataObject):
         if isObjLink:
            bMatch = oldObj.GetKey() == newObj.GetKey()
            LogIf(not bMatch, 'DiffAnyObjects: Keys do not match %s != %s'
               % (oldObj.GetKey(), newObj.GetKey()))
            return bMatch
         return self.DiffDataObjects(oldObj, newObj)

      else:
         raise TypeError("Unknown type: "+repr(GetWsdlName(oldObj.__class__)))

   def DiffDoArrays(self, oldObj, newObj, isElementLinks):
      """Diff two DataObject arrays"""
      if len(oldObj) != len(newObj):
         __Log__.debug('DiffDoArrays: Array lengths do not match %d != %d'
            % (len(oldObj), len(newObj)))
         return False
      for i, j in zip(oldObj, newObj):
         if isElementLinks:
            if i.GetKey() != j.GetKey():
               __Log__.debug('DiffDoArrays: Keys do not match %s != %s'
                  % (i.GetKey(), j.GetKey()))
               return False
         else:
            if not self.DiffDataObjects(i, j):
               __Log__.debug(
                  'DiffDoArrays: one of the elements do not match')
               return False
      return True

   def DiffAnyArrays(self, oldObj, newObj, isElementLinks):
      """Diff two arrays which contain Any objects"""
      if len(oldObj) != len(newObj):
         __Log__.debug('DiffAnyArrays: Array lengths do not match. %d != %d'
            % (len(oldObj), len(newObj)))
         return False
      for i, j in zip(oldObj, newObj):
         if not self.DiffAnyObjects(i, j, isElementLinks):
            __Log__.debug('DiffAnyArrays: One of the elements do not match.')
            return False
      return True

   def DiffPrimitiveArrays(self, oldObj, newObj):
      """Diff two primitive arrays"""
      if len(oldObj) != len(newObj):
         __Log__.debug('DiffDoArrays: Array lengths do not match %d != %d'
            % (len(oldObj), len(newObj)))
         return False
      match = True
      if self._ignoreArrayOrder:
         oldSet = oldObj and frozenset(oldObj) or frozenset()
         newSet = newObj and frozenset(newObj) or frozenset()
         match = (oldSet == newSet)
      else:
         for i, j in zip(oldObj, newObj):
            if i != j:
               match = False
               break
      if not match:
         __Log__.debug(
            'DiffPrimitiveArrays: One of the elements do not match.')
         return False
      return True


   def DiffArrayObjects(self, oldObj, newObj, isElementLinks=False):
      """Method which deligates the diffing of arrays based on the type"""
      if oldObj == newObj:
         return True
      if not oldObj or not newObj:
         return False
      if len(oldObj) != len(newObj):
         __Log__.debug('DiffArrayObjects: Array lengths do not match %d != %d'
            % (len(oldObj), len(newObj)))
         return False
      firstObj = oldObj[0]
      if IsPrimitiveType(firstObj):
         return self.DiffPrimitiveArrays(oldObj, newObj)
      elif isinstance(firstObj, types.ManagedObject):
         return self.DiffAnyArrays(oldObj, newObj, isElementLinks)
      elif isinstance(firstObj, types.DataObject):
         return self.DiffDoArrays(oldObj, newObj, isElementLinks)
      else:
         raise TypeError("Unknown type: %s" % oldObj.__class__)


   def DiffDataObjects(self, oldObj, newObj):
      """Diff Data Objects"""
      if oldObj == newObj:
         return True
      if not oldObj or not newObj:
         __Log__.debug('DiffDataObjects: One of the objects in None')
         return False
      oldType = Type(oldObj)
      newType = Type(newObj)
      if oldType != newType:
         __Log__.debug(
            'DiffDataObjects: Types do not match for dataobjects. %s != %s'
            % (oldObj._wsdlName, newObj._wsdlName))
         return False
      for prop in oldObj._GetPropertyList():
         oldProp = getattr(oldObj, prop.name)
         newProp = getattr(newObj, prop.name)
         propType = oldObj._GetPropertyInfo(prop.name).type
         if not oldProp and not newProp:
            continue
         elif ((prop.flags & VmomiSupport.F_OPTIONAL) and
               self._looseMatch and (not newProp or not oldProp)):
            continue
         elif not oldProp or not newProp:
            __Log__.debug(
               'DiffDataObjects: One of the objects has property %s unset'
               % prop.name)
            return False

         bMatch = True
         if IsPrimitiveType(oldProp):
            bMatch = oldProp == newProp
         elif isinstance(oldProp, types.ManagedObject):
            bMatch = self.DiffAnyObjects(oldProp, newProp, prop.flags
               & VmomiSupport.F_LINK)
         elif isinstance(oldProp, types.DataObject):
            if prop.flags & VmomiSupport.F_LINK:
               bMatch = oldObj.GetKey() == newObj.GetKey()
               LogIf(not bMatch, 'DiffDataObjects: Key match failed %s != %s'
                  % (oldObj.GetKey(), newObj.GetKey()))
            else:
               bMatch = self.DiffAnyObjects(oldProp, newProp, prop.flags
                  & VmomiSupport.F_LINK)
         elif isinstance(oldProp, list):
            bMatch = self.DiffArrayObjects(oldProp, newProp, prop.flags
               & VmomiSupport.F_LINK)
         else:
            raise TypeError("Unknown type: "+repr(propType))

         if not bMatch:
            __Log__.debug('DiffDataObjects: Objects differ in property %s'
               % prop.name)
            return False
      return True


def DiffAnys(obj1, obj2, looseMatch=False, ignoreArrayOrder=True):
   """Diff any two objects. Objects can either be primitive type
      or DataObjects"""
   differ = Differ(looseMatch = looseMatch, ignoreArrayOrder = ignoreArrayOrder)
   return differ.DiffAnyObjects(obj1, obj2)
