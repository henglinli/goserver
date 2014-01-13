#
JAVA_OUT = java
#
all: java_out  out.go out.java 
#
java_out:
	@mkdir -p $(JAVA_OUT) 
#
out.go: message.pb.go
#
out.java:	message.proto
	@protoc $< --java_out=$(JAVA_OUT)
#
clean:
	@rm -rf $(JAVA_OUT) *~
#
%.pb.go:	%.proto
	@protoc $< --go_out=.
#
%.java: 	%.proto
	@protoc $< --java_out=java
