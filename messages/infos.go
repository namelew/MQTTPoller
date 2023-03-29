package messages

type Info struct {
	MemoryDisplay bool
	CpuDisplay    bool
	DiscDisplay   bool
}

type InfoDisplay struct{
	Cpu 	string
	Ram		uint64
	Disk 	uint64 
}