# tabbycat

Package tabbycat is a wrapper around text/tabwriter which ignores the width of
any text matching a given regular expression.  This can be used to properly
tabbify text containing non-printing ANSI terminal control codes, for example.
Note that text/tabwriter's ability to filter HTML tags is always enabled since
this mechanism is exploited to achieve this package's purpose.
