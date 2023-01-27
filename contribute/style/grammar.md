---
title: Grammar and style
description: Grammar and style guidelines for technical documentation
keywords: grammar, style, contribute
toc_max: 2
--- 

Docker documentation should always be written in US English with US grammar. 

## Acronyms and initialisms

An acronym is an abbreviation you would speak as a word, for example, ROM (for read only memory). Other examples include radar and scuba, which started out as acronyms but are now considered words in their own right.

An initialism is a type of acronym that comprises a group of initial letters used as an abbreviation for a name or expression. If you were using the acronym in a spoken conversation, you would enunciate each letter: H-T-M-L for Hypertext Markup Language.

### Best practice:

- Spell out lesser-known acronyms or initialisms on first use, then follow with the acronym or initialism in parentheses. After this, throughout the rest of your page or document, use the acronym or initialism alone.
> ‘You can use single sign-on (SSO) to sign in to Notion. You may need to ask your administrator to enable SSO.’
- Where the acronym or initialism is more commonly used than the full phrase, for example, URL, HTML, you do not need to follow this spell-it-out rule.
- Use all caps for acronyms of file types (a JPEG file).
- Don't use apostrophes for plural acronyms. ✅URLs ❌URL’S
- Avoid using an acronym for the first time in a title or heading. If the first use of the acronym is in a title or heading, introduce the acronym (in parentheses, following the spelled-out term) in the first body text that follows.

## Bolds and italics

Unless you're referring to Ul text or user-defined text, you should not add emphasis to text. Clear, front-loaded wording makes the subject of a sentence clear.

### Best practice:

- Don't use bold to refer to a feature name.
- Use italics sparingly, as this type of formatting can be difficult to read in digital experiences.
Notable exceptions are titles of articles, blog posts, or specifcation documents.

## Capitalization

Use sentence case for just about everything. Sentence case means capitalizing only the first word, as you would in a standard sentence.

The following content elements should use sentence case:

- Titles of webinars and events
- Headings and subheadings in all content types
- Calls to action
- Headers in boxed text
- Column and row headers in tables
- Links
- Sentences (of course)
- Anything in the Ul including navigation labels, buttons, headings

### Best practice

- As a general rule, it's best to avoid the use of ALL CAPITALS in most content types. They are more difficult to scan and take up more space. While all caps can convey emphasis, they can also give the impression of shouting.
- If a company name is all lowercase or all uppercase letters, follow their style, even at the beginning of sentences: DISH and bluecrux. When in doubt, check the company's website.
- Use title case for Docker solutions : Docker Extensions, Docker Hub.
- Capitalize a job title if it immediately precedes a name (Chief Executive Officer Scott Johnston).
- Do not capitalize a job title that follows a name or is a generic reference (Scott Johnston, chief executive officer of Docker).
- Capitalize department names when you refer to the name of a department, but use lower case if you are talking about the field of work and not the actual department.
- When referring to specific user interface text, like a button label or menu item, use the same capitalization that’s displayed in the user interface. 

## Contractions

A contraction results from letters being left out from the original word or phrase, such as you're for you are or it's for it is.

Following our conversational approach, it's acceptable to use contractions in almost all content types. Just don't get carried away. Too many contractions in a sentence can make it difficult to read.

### Best practice:

- Stay consistent - don't switch between contractions and their spelled-out equivalents in body copy or Ul text.
- Avoid negative contractions (can't, don't, wouldn't, and shouldn't) whenever possible. Try to rewrite your sentence to align with our helpful approach that puts the focus on solutions, not problems.
- Never contract a noun with is, does, has, or was as this can make it look like the noun is possessive. Your container is ready. Your container’s ready.

## Dates

You should use the U.S. format of month day, year format for dates: June 26, 2021.

When possible, use the month's full name (October 1, 2022). If there are space constraints, use 3-letter abbreviations followed by a period (Oct. 1. 2022).

## Decimals and fractions

In all UI content and technical documentation, use decimals instead of fractions.

### Best practice:

- Always carry decimals to at least the hundredth place (33.76).

- In tables, align decimals on the decimal point.

- Add a zero before the decimal point for decimal fractions less than one (0.5 cm).

## Lists

Lists are a great way to break down complex ideas and make them easier to read and scan.

### Best practice:

- Bulleted lists should contain relatively few words or short phrases. For most content types, limit the number of items in a list to five.
- Don’t add commas (,) or semicolons (;) to the ends of list items.
- Some content types may use bulleted lists that contain 10 items, but it's preferable to break longer lists into several lists, each with its own subheading or introduction.
- Never create a bulleted list with only one bullet, and never use a dash to indicate a bulleted list.
- If your list items are fragments, capitalize the list items for ease of scanning but do not use terminal punctuation. 
Example:
    
    I went to the shops to buy:
    
    - Milk
    - Flour
    - Eggs
- Make sure all your list items are parallel. This means you should structure each list item in the same way. They should all be fragments, or they should all be complete sentences. If you start one list item with a verb, then start every list item with a verb.
- Every item in your list should start with a capital letter unless they’re parameters or commands.
- Nested sequential lists are labelled with lowercase letters. For example:
    1. Top level
    2. Top level
        1. Child step
        2. Child step


## Numbers

When you work with numbers in content, the best practices include:

- Spell out the numbers one to nine, except in units of measure such as 4 MB.
- Represent numbers with two or more digits as numerals (10, 625, 773,925).
- Recast a sentence, rather than begin it with a number (unless it's a year).
- When you cite numbers in an example, write them out as they appear in any accompanying screenshots.
- Write numbers out as they appear on the platform when you cite them in an example.
- To refer to large numbers in abstract (such as thousands, millions, and billions), use a combination of words and numbers. Do not abbreviate numeric signifiers.
- Avoid using commas in numbers because they can represent decimals in different cultures. For numbers that are five digits or more, use a space to separate.
    
    
    | Correct | Incorrect |
    | --- | --- |
    | 1000 | 1,000 |
    | 14 586 | 14,586 |
    | 1 390 680 | 1,390,680 |

### Versions

When writing version numbers for release notes, follow the example below:

Version 4.8.2

## Punctuation

### Colons and semicolons

- Use colons to: introduce a list inline in a sentence; introduce a bulleted or numbered list; and provide an explanation.
- Do not use semicolons. Use two sentences instead.

### Commas

- Use the serial or Oxford comma - a comma before the coordinating conjunction (and, or) in a list of three or more things.
- If a series contains more than three items or the items are long, consider a bulleted list to improve readability.

### Dashes and hyphens

- Use the em dash (-) sparingly when you want the reader to pause, to create parenthetical statements, or to emphasize specific words or phrases. Always put a space on either side of the em dash.
- Use an en dash (-) to indicate spans of numbers, dates, or time.
- Use a hyphen to join two or more words to form compound adjectives that precede a noun. The purpose of joining words to form a compound adjective is to differentiate the meaning from the adjectives used separately (for example, ‘up-to-date documentation’ ‘lump-sum payment’, and ‘well-stocked cupboard’. You can also use a hyphen to:
    - Avoid awkward doubling of vowels. For example ‘semi-independence*’,* or ‘re-elect’.
    - Prevent misreading of certain words. For example ‘Re-collect’ means to collect again; without a hyphen the word recollect has a different meaning.

### Parentheses

Don't use parentheses in technical documentation. They can reduce the readability of a sentence.  