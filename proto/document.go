package proto

import (
	"github.com/bryk-io/x/did"
)

func (d *DID_Document) Load(doc *did.Document) {
	d.Context = doc.Context
	d.Subject = doc.Subject
	d.Created = doc.Created
	d.Updated = doc.Updated
	d.Authentication = doc.Authentication
	d.PublicKey = []*DID_PublicKey{}
	for _, k := range doc.PublicKeys {
		dk := &DID_PublicKey{
			ID:          k.ID,
			Type:        k.Type,
			Controller:  k.Controller,
			ValueHex:    k.ValueHex,
			ValueBase58: k.ValueBase58,
			ValueBase64: k.ValueBase64,
			Private:     k.Private,
		}
		d.PublicKey = append(d.PublicKey, dk)
	}
	d.Service = []*DID_Service{}
	for _, s := range doc.Services {
		ds := &DID_Service{
			ID:              s.ID,
			Type:            s.Type,
			ServiceEndpoint: s.Endpoint,
		}
		d.Service = append(d.Service, ds)
	}
	if doc.Proof != nil {
		d.Proof = &DID_Proof{
			Context:    doc.Proof.Context,
			Type:       doc.Proof.Type,
			Creator:    doc.Proof.Creator,
			Created:    doc.Proof.Created,
			Domain:     doc.Proof.Domain,
			Nonce:      doc.Proof.Nonce,
			ProofValue: doc.Proof.Value,
		}
	}
}
