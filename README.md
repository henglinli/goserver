
struct Boundary {
       unit64_t message_size_;
       unit64_t message_type_;
};

struct ProtbufMessage;

writer:
write Boundary;
write ProtbufMessage;

reader:
read Boundary;
read ProtbufMessage;