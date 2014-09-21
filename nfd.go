package ndn

import (
	"github.com/taylorchu/tlv"
	"time"
)

type ControlPacket struct {
	Name      Command   `tlv:"7"`
	Selectors Selectors `tlv:"9?"`
	Nonce     []byte    `tlv:"10"`
	Scope     uint64    `tlv:"11?"`
	LifeTime  uint64    `tlv:"12?"`
}

// see http://redmine.named-data.net/projects/nfd/wiki/Management
type Command struct {
	Localhost      string                  `tlv:"8"`
	Nfd            string                  `tlv:"8"`
	Module         string                  `tlv:"8"`
	Command        string                  `tlv:"8"`
	Parameters     parametersComponent     `tlv:"8"`
	Timestamp      uint64                  `tlv:"8"`
	Nonce          []byte                  `tlv:"8"`
	SignatureInfo  signatureInfoComponent  `tlv:"8"`
	SignatureValue signatureValueComponent `tlv:"8*"`
}

// WriteTo writes control interest packet to tlv.Writer after it signs the name automatically
//
// Everything except Module, Command and Parameters will be populated.
func (this *ControlPacket) WriteTo(w tlv.Writer) (err error) {
	this.Name.Localhost = "localhost"
	this.Name.Nfd = "nfd"
	this.Name.Timestamp = uint64(time.Now().UnixNano() / 1000000)
	this.Name.Nonce = newNonce()
	this.Name.SignatureInfo.SignatureInfo.SignatureType = SignKey.SignatureType()
	this.Name.SignatureInfo.SignatureInfo.KeyLocator.Name = SignKey.Name.CertificateName()

	digest, err := newSha256(this.Name)
	if err != nil {
		return
	}
	this.Name.SignatureValue.SignatureValue, err = SignKey.sign(digest)
	if err != nil {
		return
	}

	if this.LifeTime == 0 {
		this.LifeTime = 4000
	}
	this.Nonce = newNonce()
	err = tlv.Marshal(w, this, 5)
	return
}

func (this *ControlPacket) ReadFrom(r tlv.PeekReader) error {
	return tlv.Unmarshal(r, this, 5)
}

type parametersComponent struct {
	Parameters Parameters `tlv:"104"`
}

type signatureInfoComponent struct {
	SignatureInfo SignatureInfo `tlv:"22"`
}

type signatureValueComponent struct {
	SignatureValue []byte `tlv:"23"`
}

type Parameters struct {
	Name                Name     `tlv:"7?"`
	FaceId              uint64   `tlv:"105?"`
	Uri                 string   `tlv:"114?"`
	LocalControlFeature uint64   `tlv:"110?"`
	Origin              uint64   `tlv:"111?"`
	Cost                uint64   `tlv:"106?"`
	Flags               uint64   `tlv:"108?"`
	Strategy            Strategy `tlv:"107?"`
	ExpirationPeriod    uint64   `tlv:"109?"`
}

type Strategy struct {
	Name Name `tlv:"7"`
}

type ControlResponse struct {
	StatusCode uint64     `tlv:"102"`
	StatusText string     `tlv:"103"`
	Parameters Parameters `tlv:"104?"`
}

type NextHopRecord struct {
	FaceId uint64 `tlv:"105"`
	Cost   uint64 `tlv:"106"`
}

type FibEntry struct {
	Name     Name            `tlv:"7"`
	NextHops []NextHopRecord `tlv:"129"`
}

type FaceEntry struct {
	FaceId      uint64 `tlv:"105"`
	Uri         string `tlv:"114"`
	LocalUri    string `tlv:"129"`
	FaceFlag    uint64 `tlv:"194"`
	InInterest  uint64 `tlv:"144"`
	InData      uint64 `tlv:"145"`
	OutInterest uint64 `tlv:"146"`
	OutData     uint64 `tlv:"147"`
}

type ForwarderStatus struct {
	NfdVersion       uint64 `tlv:"128"`
	StartTimestamp   uint64 `tlv:"129"`
	CurrentTimestamp uint64 `tlv:"130"`
	NameTreeEntry    uint64 `tlv:"131"`
	FibEntry         uint64 `tlv:"132"`
	PitEntry         uint64 `tlv:"133"`
	MeasurementEntry uint64 `tlv:"134"`
	CsEntry          uint64 `tlv:"135"`
	InInterest       uint64 `tlv:"144"`
	InData           uint64 `tlv:"145"`
	OutInterest      uint64 `tlv:"146"`
	OutData          uint64 `tlv:"147"`
}
