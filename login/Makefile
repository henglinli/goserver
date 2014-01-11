include $(GOPATH)/src/code.google.com/p/goprotobuf/Make.protobuf
#
JAVA_OUT = java
#
all: java_out login.pb.go Login.java
#

java_out:
	mkdir -p $(JAVA_OUT) 
#
Login.java : login.proto
	protoc $< --java_out=./java

clean:
	rm -rf $(JAVA_OUT)
