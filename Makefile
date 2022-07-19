UNICODE_DB = "https://www.unicode.org/Public/14.0.0/ucd/extracted/DerivedName.txt"

bin:
	go build -o bin/sluk main.go

download-db:
	curl -L $(UNICODE_DB) -o data/DerivedName.txt
