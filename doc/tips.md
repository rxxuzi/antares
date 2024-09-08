# tips

## Search Tips

* **Basic Search:** Simply type your search term (e.g., "document")
* **OR Search:** Use "||" or "OR" between terms (e.g., "jpg || png")
* **AND Search:** Use "&&" or "AND" between terms (e.g., "document && pdf")
* **NOT Search:** Use "NOT" or "!!" before a term (e.g., "image NOT png")
* **Case Sensitive:** Toggle to make search case-sensitive
* **Regex:** Toggle to use regular expressions in your search

### Regular Expression Quick Guide

*   `.` - Matches any single character
*   `*` - Matches zero or more of the preceding character
*   `+` - Matches one or more of the preceding character
*   `?` - Matches zero or one of the preceding character
*   `^` - Matches the start of the string
*   `$` - Matches the end of the string
*   `[abc]` - Matches any one character in the brackets
*   `[^abc]` - Matches any character not in the brackets
*   `\d` - Matches any digit (0-9)
*   `\w` - Matches any word character (a-z, A-Z, 0-9, \_)

#### Examples:

*   `\.pdf$` - Matches filenames ending with ".pdf"
*   `^img_\d+\.jpg$` - Matches filenames like "img\_123.jpg"
*   `[aeiou]` - Matches any vowel in the filename