package read

func KeywordsRegexpsForWSFormat() []string {
	return []string{
		"select",
		"select[ \t\n\r\f]+top[ \t\n\r\f]+[0-9]",
		"from",
		"where",
		"group[ \t\n\r\f]+by",
		"order[ \t\n\r\f]+by",
		"[(]force[ \t\n\r\f]+order[)]",
		"update",
		"set",
		"union[ \t\n\r\f]+(all|)",
	}
}

func KeywordsRegexpsForWSReplace() map[string]string {
	return map[string]string{
		"(?i)select[ \t\n\r\f]+top[ \t\n\r\f]+": "select top ",
		"(?i)group[ \t\n\r\f]+by":               "group by",
		"(?i)order[ \t\n\r\f]+by":               "order by",
		"(?i)[(]force[ \t\n\r\f]+order[)]":      "force order",
		"(?i)union[ \t\n\r\f]+":                 "union ",
	}
}
