
CompileDaemon -exclude-dir=.git -exclude=".#*" -include="*.html" -include="*.css" -build="go build cmd/forum.go" -command="forum ./conf.json"