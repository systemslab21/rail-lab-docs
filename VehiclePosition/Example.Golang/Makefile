.PRECIOUS: %.proto

all: image

%.proto:
	cp ../$@ .
	sed -i '' 's/package esc;/package esc;\noption go_package = "main";/g' $@

%.pb.go: %.proto
	protoc --go_out=./ $<

image: vehicle_position.pb.go vehicle_attributes.pb.go
	docker build .

