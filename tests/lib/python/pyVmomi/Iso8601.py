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
from __future__ import print_function
# TODO (hartsocks): Introduce logging to remove the need for print function.
"""
This module is for ISO 8601 parsing
"""
__author__ = 'VMware, Inc.'

from six import iteritems
import time
from datetime import datetime, timedelta, tzinfo
import re

""" Regular expression to parse a subset of ISO 8601 format """
_dtExpr = re.compile(
   # XMLSchema datetime. Mandatory to have - and :
   # See: http://www.w3.org/TR/xmlschema-2/#isoformats
   # Note: python datetime cannot handle the following:
   #       - leap second, ie. 0-60 seconds (not 0-59)
   #       - BC (negative years)
   # year [-]0000..9999
   r'(?P<year>-?\d{4})' \
   # month 01..12
   r'(-(?P<month>(0[1-9]|1[0-2]))' \
    # day 01..31
    r'(-(?P<day>(0[1-9]|[1-2]\d|3[01])))?)?' \
   # time separator 'T'
   r'(T' \
    # hour 00..24
    r'(?P<hour>([01]\d|2[0-4]))' \
    # minute 00..59
    r'((:(?P<minute>[0-5]\d))' \
     # seconds 00..60 (leap second ok)
     r'(:(?P<second>([0-5]\d|60))' \
      # microsecond. max 16 digits
      # - Should not allows trailing zeros. But python isoformat() put zeros
      #   after microseconds. Oh well, allows trailing zeros, quite harmless
      r'(\.(?P<microsecond>\d{1,16}))?)?)?' \
    # UTC 'Z', or...
    r'((?P<tzutc>Z)' \
    # tz [+-]00..13:0..59|14:00
    r'|((?P<tzhr>[+-](([0]\d)|(1[0-3])|(?P<tzlimit>)14))' \
      r'(:(?P<tzmin>(?(tzlimit)00|([0-5]\d))))?))?' \
   r')?$')

""" Default date time val. Key should match the tags in _dtExpr """
_dtExprKeyDefValMap = {'year' : None, 'month' : 1, 'day' : 1,
                       'hour' : 0, 'minute' : 0, 'second' : 0,
                       'microsecond' : 0}

class TZInfo(tzinfo):
   """ Timezone info class """

   timedelta0 = timedelta(hours=0)
   timedelta1 = timedelta(hours=1)

   def __init__(self, tzname='UTC', utcOffset=None, dst=None):
      self._tzname = tzname
      if not utcOffset:
         utcOffset = self.timedelta0
      self._utcOffset = utcOffset
      if not dst:
         dst = None
      self._dst = dst

   def utcoffset(self, dt):
      return self._utcOffset + self.dst(dt)

   def tzname(self, dt):
      return self._tzname

   def dst(self, dt):
      ret = self.timedelta0
      if self._dst:
         if self._dst[0] <= dt.replace(tzinfo=None) < self._dst[1]:
            ret = self.timedelta1
      return ret


class TZManager:
   """ Time zone manager """
   _tzInfos = {}

   @staticmethod
   def GetTZInfo(tzname='UTC', utcOffset=None, dst=None):
      """ Get / Add timezone info """
      key = (tzname, utcOffset, dst)
      tzInfo = TZManager._tzInfos.get(key)
      if not tzInfo:
         tzInfo = TZInfo(tzname, utcOffset, dst)
         TZManager._tzInfos[key] = tzInfo
      return tzInfo


def ParseISO8601(datetimeStr):
   """
   Parse ISO 8601 date time from string.
   Returns datetime if ok, None otherwise
   Note: Allows YYYY / YYYY-MM, but truncate YYYY -> YYYY-01-01,
                                             YYYY-MM -> YYYY-MM-01
         Truncate microsecond to most significant 6 digits
   """
   datetimeVal = None
   match = _dtExpr.match(datetimeStr)
   if match:
      try:
         dt = {}
         for key, defaultVal in iteritems(_dtExprKeyDefValMap):
            val = match.group(key)
            if val:
               if key == 'microsecond':
                  val = val[:6] + '0' * (6 - len(val))
               dt[key] = int(val)
            elif defaultVal:
               dt[key] = defaultVal

         # Orig. XMLSchema don't allow all zeros year. But newer draft is ok
         #if dt['year'] == 0:
         #   # Year cannot be all zeros
         #   raise Exception('Year cannot be all zeros')

         # 24 is a special case. It is actually represented as next day 00:00
         delta = None
         if dt.get('hour', 0) == 24:
            # Must be 24:00:00.0
            if dt.get('minute', 0) == 0 and dt.get('second', 0) == 0 and \
               dt.get('microsecond', 0) == 0:
               dt['hour'] = 23
               delta = timedelta(hours=1)
            else:
               return None

         # Set tzinfo
         # TODO: dst
         tzInfo = None
         val = match.group('tzutc')
         if val:
            tzInfo = TZManager.GetTZInfo()
         else:
            val = match.group('tzhr')
            if val:
               # tz hours offset
               tzhr = int(val)
               utcsign = val[0]

               # tz minutes offset
               tzmin = 0
               val = match.group('tzmin')
               if val:
                  tzmin = tzhr >= 0 and int(val) or -int(val)

               # Better tzname (map UTC +-00:00 to UTC)
               tzname = 'UTC'
               if tzhr != 0 or tzmin != 0:
                  tzname += ' %s%02d:%02d' % (utcsign, abs(tzhr), abs(tzmin))

               tzInfo = TZManager.GetTZInfo(tzname=tzname,
                                            utcOffset=timedelta(hours=tzhr,
                                                                minutes=tzmin))
         if tzInfo:
            dt['tzinfo'] = tzInfo

         datetimeVal = datetime(**dt)
         if delta:
            datetimeVal += delta
      except Exception as e:
         pass
   return datetimeVal


def ISO8601Format(dt):
   """
   Python datetime isoformat() has the following problems:
   - leave trailing 0 at the end of microseconds (violates XMLSchema rule)
   - tz print +00:00 instead of Z
   - Missing timezone offset for datetime without tzinfo
   """
   isoStr = dt.strftime('%Y-%m-%dT%H:%M:%S')
   if dt.microsecond:
      isoStr += ('.%06d' % dt.microsecond).rstrip('0')
   if dt.tzinfo:
      tz = dt.strftime('%z')
   else:
      if time.daylight and time.localtime().tm_isdst:
         utcOffset_minutes = -time.altzone / 60
      else:
         utcOffset_minutes = -time.timezone / 60
      tz = "%+.2d%.2d" % (utcOffset_minutes / 60, (abs(utcOffset_minutes) % 60))
   if tz == '+0000':
      return isoStr + 'Z'
   elif tz:
      return isoStr + tz[:3] + ':' + tz[3:]
   else:
      # Local offset is unknown
      return isoStr + '-00:00'


# Testing
if __name__ == '__main__':
   # Valid entries
   for testStr in [
         '1971', # 1971-01-01
         '1971-11', # 1971-11-01
         '1971-11-02',
         '1971-11-02T23',
         '1971-11-02T23Z',
         '1971-11-02T23:04',
         '1971-11-02T23:04Z',
         '1971-11-02T23:04:15',
         '1971-11-02T23:04:15Z',
         '1971-11-02T23:04:15.1',
         '1971-11-02T23:04:15.01',
         '1971-11-02T23:04:15.023456',
         '1971-11-02T23:04:15.103456Z',
         '1971-11-02T23:04:15.123456+11',
         '1971-11-02T23:04:15.123456-11',
         '1971-11-02T23:04:15.123456+11:30',
         '1971-11-02T23:04:15.123456-11:30',
         '1971-11-02T23:04:15.123456+00:00', # Same as Z
         '1971-11-02T23:04:15.123456-00:00', # Same as Z

         '1971-01-02T23:04:15+14',
         '1971-01-02T23:04:15+14:00',
         '1971-01-02T23:04:15-14',
         '1971-01-02T23:04:15-14:00',

         # Valid: Truncate microsec to 6 digits
         '1971-01-02T23:04:15.123456891+11',

         '1971-01-02T24', # 24 is valid. It should represent the 00:00 the
                          # next day
         '1971-01-02T24:00',
         '1971-01-02T24:00:00',
         '1971-01-02T24:00:00.0',

         # Should NOT be valid but python isoformat adding trailing zeros
         '1971-01-02T23:04:15.123430', # Microseconds ends in zero
         '1971-01-02T23:04:15.0', # Microseconds ends in zero

         # Should be valid but python datetime don't support it
         #'2005-12-31T23:59:60Z', # Leap second
         #'-0001', # BC 1
        ]:
      dt = ParseISO8601(testStr)
      if dt == None:
         print('Failed to parse ({0})'.format(testStr))
         assert(False)

      # Make sure we can translate back
      isoformat = ISO8601Format(dt)
      dt1 = ParseISO8601(isoformat)
      if dt.tzinfo is None:
         dt = dt.replace(tzinfo=dt1.tzinfo)
      if dt1 != dt:
         print('ParseISO8601 -> ISO8601Format -> ParseISO8601 failed ({0})'.format(testStr))
         assert(False)

      # Make sure we can parse python isoformat()
      dt2 = ParseISO8601(dt.isoformat())
      if dt2 == None:
         print('ParseISO8601("{0}".isoformat()) failed'.format(testStr))
         assert(False)

      print(testStr, '->', dt, isoformat)

   # Basic form
   for testStr in [
         '197111', # 1971-11-01
         '19711102',
         '19711102T23',
         '19711102T23Z',
         '19711102T2304',
         '19711102T2304Z',
         '19711102T230415',
         '19711102T230415Z',
         '19711102T230415.123456',
         '19711102T230415.123456Z',
         '19711102T230415.123456+11',
         '19711102T230415.123456-11',
         '19711102T230415.123456+1130',
         '19711102T230415.123456-1130',
        ]:
      # Reject for now
      dt = ParseISO8601(testStr)
      if dt != None:
         print('ParseISO8601 ({0}) should fail, but it did not'.format(testStr))
         assert(False)
      #print testStr, '->', dt
      #assert(dt != None)

   # Invalid entries
   for testStr in [
         # Xml schema reject year 0
         '0000', # 0 years are not allowed
         '+0001', # Leading + is not allowed

         '', # Empty datetime str
         '09', # Years must be at least 4 digits
         '1971-01-02T', # T not follow by time
         '1971-01-02TZ', # T not follow by time
         '1971-01-02T+10', # T not follow by time
         '1971-01-02T-10', # T not follow by time
         '1971-01-02T23:', # extra :
         '1971-01-02T23:04:', # extra :
         '1971-01-02T23:0d', # 0d
         '1971-01-02T23:04:15.', # Dot not follows by microsec
         '1971-01-02+12', # time without T
         '1971Z', # Z without T
         '1971-01-02T23:04:15.123456Z+11', # Z follows by +
         '1971-01-02T23:04:15.123456Z-11', # Z follows by -
         '1971-01-02T23:04:15.123456+:30', # extra :
         '1971-01-02T23:04:15.123456+30:', # extra :
         '1971-01-02T23:04:15.01234567890123456789', # Too many microseconds digits

         # Python isoformat leave trailing zeros in microseconds
         # Relax regular expression to accept it
         #'1971-01-02T23:04:15.123430', # Microseconds ends in zero
         #'1971-01-02T23:04:15.0', # Microseconds ends in zero

         # Timezone must be between +14 / -14
         '1971-01-02T23:04:15+15',
         '1971-01-02T23:04:15-15',
         '1971-01-02T23:04:15+14:01',
         '1971-01-02T23:04:15-14:01',

         # Mix basic form with extended format
         '197101-02T23:04:15.123456',
         '19710102T23:04:15.123456',
         '19710102T230415.123456+11:30',
         '1971-01-02T230415.123456',
         '1971-01-02T23:04:15.123456+1130',

         # Error captured by datetime class
         '1971-00-02', # Less than 1 month
         '1971-13-02', # Larger than 12 months
         '1971-01-00', # Less than 1 day
         '1971-11-32', # Larger than 30 days for Nov
         '1971-12-32', # Larger than 31 days
         '1971-01-02T24:01', # Larger than 23 hr
         '1971-01-02T23:61', # Larger than 60 min
         '1971-01-02T23:60:61', # Larger than 61 sec
        ]:
      dt = ParseISO8601(testStr)
      if dt != None:
         print('ParseISO8601 ({0}) should fail, but it did not'.format(testStr))
         assert(False)
