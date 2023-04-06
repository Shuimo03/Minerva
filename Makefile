out_file=events-exporter
out_path=bin

build-exporter:
	go build -o $(out_path)/$(out_file) cmd/main.go

pull-exporterImage:
