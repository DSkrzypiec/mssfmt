package read

func keywordsRegexpsForWSFormat() map[string]string {
	return map[string]string{
		"(?i)[ \t\n\r\f]*select[ \t\n\r\f]+":                                  "select",
		"(?i)[ \t\n\r\f]*select[ \t\n\r\f]+top[ \t\n\r\f]+[0-9]+[ \t\n\r\f]+": "select top",
		"(?i)[ \t\n\r\f]+from[ \t\n\r\f]+":                                    "from",
		"(?i)[ \t\n\r\f]+where[ \t\n\r\f]+":                                   "where",
		"(?i)[ \t\n\r\f]+group[ \t\n\r\f]+by[ \t\n\r\f]+":                     "group by",
		"(?i)[ \t\n\r\f]+order[ \t\n\r\f]+by[ \t\n\r\f]+":                     "order by",
		"(?i)[ \t\n\r\f]*update[ \t\n\r\f]+":                                  "update",
		"(?i)[ \t\n\r\f]*set[ \t\n\r\f]+":                                     "set",
		"(?i)[ \t\n\r\f]+union[ \t\n\r\f]+(all|)[ \t\n\r\f]+":                 "union",
	}
}
