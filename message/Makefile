#
JAVA_OUT_PATH = java
GO_OUT_PATH = .
#
all: java_out_path  out.go out.java 
#
java_out_path:
	@mkdir -p $(JAVA_OUT_PATH) 
#
out.go: 	message.proto extensions.proto
	@protoc $^ --go_out=$(GO_OUT_PATH)
#
out.java:	message.proto extensions.proto
	@protoc $^ --java_out=$(JAVA_OUT_PATH)
#
clean:
	@rm -rf $(JAVA_OUT_PATH) *~
#
distclean:	clean
	@rm -rf *.pb.go
#
%.pb.go:	%.proto
	@protoc $< --go_out=$(GO_OUT_PATH)
#
%.java: 	%.proto
	@protoc $< --java_out=$(JAVA_OUT_PATH)
