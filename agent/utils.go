package agent

import (
	didpb "github.com/bryk-io/did-method/proto"
	"go.bryk.io/x/did"
)

func populateDocument(id *did.Identifier) *didpb.Document {
	src := id.Document()
	dp := &didpb.Document{
		Context:        src.Context,
		Subject:        src.Subject,
		Created:        src.Created,
		Updated:        src.Updated,
		Authentication: src.Authentication,
		Keys:           []*didpb.PublicKey{},
		Services:       []*didpb.ServiceEndpoint{},
	}
	for _, k := range src.PublicKeys {
		nk := &didpb.PublicKey{
			Type:        k.Type.String(),
			Controller:  k.Controller,
			Id:          k.ID,
			ValueHex:    k.ValueHex,
			ValueBase58: k.ValueBase58,
			ValueBase64: k.ValueBase64,
		}
		dp.Keys = append(dp.Keys, nk)
	}
	for _, s := range src.Services {
		ns := &didpb.ServiceEndpoint{
			Id:       s.ID,
			Type:     s.Type,
			Endpoint: s.Endpoint,
		}
		dp.Services = append(dp.Services, ns)
	}
	if src.Proof != nil {
		dp.Proof = &didpb.Proof{
			Context: src.Proof.Context,
			Type:    src.Proof.Type,
			Created: src.Proof.Created,
			Value:   src.Proof.Value,
			Creator: src.Proof.Creator,
			Domain:  src.Proof.Domain,
			Nonce:   src.Proof.Nonce,
		}
	}
	return dp
}
