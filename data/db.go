package data

import _ "embed"

// DB is an embedded file extracted from https://www.unicode.org/Public/zipped/latest/UCD.zip
// extracted/DerivedName.txt
//go:embed DerivedName.txt
var DB string
