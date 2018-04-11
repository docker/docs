---
title: Promotion policies templates
description: Learn how to use templates when setting your promotion policies to rename your images
keywords: registry, promotion, mirror
---

When defining promotion policies you can use templates to dynamically name the
tag that is going to be created.

You can use these template keywords to define your new tag:

| Template | Description                     | Example result    |
|:---------|:--------------------------------|:------------------|
| %n       | The tag to promote              | 1, 4.5, latest    |
| %A       | Day of the week                 | Sunday, Monday    |
| %a       | Day of the week, abbreviated    | Sun, Mon, Tue     |
| %w       | Day of the week, as a number    | 0, 1, 6           |
| %d       | Number for the day of the month | 01, 15, 31        |
| %B       | Month                           | January, December |
| %b       | Month, abbreviated              | Jan, Jun, Dec     |
| %m       | Month, as a number              | 01, 06, 12        |
| %Y       | Year                            | 1999, 2015, 2048  |
| %y       | Year, two digits                | 99, 15, 48        |
| %H       | Hour, in 24 hour format         | 00, 12, 23        |
| %I       | Hour, in 12 hour format         | 01, 10, 10        |
| %p       | Period of the day               | AM, PM            |
| %M       | Minute                          | 00, 10, 59        |
| %S       | Second                          | 00, 10, 59        |
| %f       | Microsecond                     | 000000, 999999    |
| %Z       | Name for the timezone           | UTC, PST, EST     |
| %j       | Day of the year                 | 001, 200, 366     |
| %W       | Week of the year                | 00, 10, 53        |
