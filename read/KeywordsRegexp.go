package read

func keywordsRegexpsForWSFormat() []string {
	return []string{
		"(?i)[ \t\n\r\f]*select[ \t\n\r\f]+",
		"(?i)[ \t\n\r\f]*select[ \t\n\r\f]+top[ \t\n\r\f]+[0-9]+[ \t\n\r\f]+",
		"(?i)[ \t\n\r\f]+from[ \t\n\r\f]+",
		"(?i)[ \t\n\r\f]+where[ \t\n\r\f]+",
		"(?i)[ \t\n\r\f]+group[ \t\n\r\f]+by[ \t\n\r\f]+",
		"(?i)[ \t\n\r\f]+order[ \t\n\r\f]+by[ \t\n\r\f]+",
		"(?i)[ \t\n\r\f]*update[ \t\n\r\f]+",
		"(?i)[ \t\n\r\f]*set[ \t\n\r\f]+",
		"(?i)[ \t\n\r\f]+union[ \t\n\r\f]+(all|)[ \t\n\r\f]+",
	}
}

func keywordsRegexpsForWSReplace() map[string]string {
	return map[string]string{
		"select[ \t\n\r\f]+top[ \t\n\r\f]+": "select top ",
		"group[ \t\n\r\f]+by":               "group by",
		"order[ \t\n\r\f]+by":               "order by",
		"[(]force[ \t\n\r\f]+order[)]":      "force order",
		"union[ \t\n\r\f]+":                 "union ",
	}
}
