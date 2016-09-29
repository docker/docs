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
This module implements the cache decorator
"""
__author__ = "VMware, Inc."

def Cache(fn):
   """ Function cache decorator """
   def fnCache(*args, **kwargs):
      """ Cache function """
      key = (args and tuple(args) or None,
             kwargs and frozenset(kwargs.items()) or None)
      if key not in fn.__cached__:
         fn.__cached__[key] = cache = fn(*args, **kwargs)
      else:
         cache = fn.__cached__[key]
      return cache

   def ResetCache():
      """ Reset cache """
      fn.__cached__ = {}

   setattr(fn, "__cached__", {})
   setattr(fn, "__resetcache__", ResetCache)
   fnCache.__name__ = fn.__name__
   fnCache.__doc__ = fn.__doc__
   fnCache.__dict__.update(fn.__dict__)
   return fnCache
