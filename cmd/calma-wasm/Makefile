serve:
	python3 -m http.server 8080

servedoc:
	cd ../../docs && python3 -m http.server 8080

# tinygo の version up したら
cpwasmexec:
	cp $(shell tinygo env TINYGOROOT)/targets/wasm_exec.js .

# tinygo でダメだったら
cpgowasmexec:
	cp $(shell go env GOROOT)/lib/wasm/wasm_exec.js .

tinybuild:
	tinygo build -o main.wasm -target wasm -tags tinywasm main.go
gobuild:
	GOOS=js GOARCH=wasm go build -o main.wasm main.go

deploy: gobuild
	cleancss -o ../../docs/toast.css toast.css
	uglifyjs main.js --compress drop_console=true -o ../../docs/main.js
	uglifyjs toast.js --compress drop_console=true -o ../../docs/toast.js
	uglifyjs wasm_exec.js --compress drop_console=true -o ../../docs/wasm_exec.js
	cp index.html ../../docs
	mv main.wasm ../../docs