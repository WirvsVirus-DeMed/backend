package node

type PayloadType uint32

const (
	Hello               PayloadType = 0
	Goodbye             PayloadType = 1
	MessageBroadcast    PayloadType = 2
	RequestResource     PayloadType = 3
	PeerRequest         PayloadType = 4
	PeerResponse        PayloadType = 5
	MedicineOffer       PayloadType = 6
	FullMedicineRequest PayloadType = 7
)
