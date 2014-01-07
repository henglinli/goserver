
	struct Boundary {
       	       unit32_t message_size_;
	};

	struct ProtbufMessage;

	writer:
	write Boundary;
	write ProtbufMessage;

	reader:
	read Boundary;
	read ProtbufMessage;