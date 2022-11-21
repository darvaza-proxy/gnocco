package version

//go:generate ./version.sh

import _ "embed"

//Version contains the git hashtag injected by make
//go:embed version.txt
var Version string

//BuildDate contains the build timestamp injected by make
//go:embed builddate.txt
var BuildDate string
