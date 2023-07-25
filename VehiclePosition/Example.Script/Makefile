all: image

%.proto:
	cp ../$@ .

image: vehicle_position.proto vehicle_attributes.proto
	docker build .
