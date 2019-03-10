# DID Method
[![Build Status](https://travis-ci.org/bryk-io/did-method.svg?branch=master)](https://travis-ci.org/bryk-io/did-method)
[![Version](https://img.shields.io/github/tag/bryk-io/did-method.svg)](https://github.com/bryk-io/did-method/releases)
[![Software License](https://img.shields.io/badge/license-BSD3-red.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bryk-io/did-method?style=flat)](https://goreportcard.com/report/github.com/bryk-io/did-method)

The present document describes the __"bryk"__ DID Method specification. The definitions,
conventions and technical details included intend to provide a solid base for further
developments while maintaining compliance with the work, still in progress, on the 
[W3C Credentials Community Group](https://w3c-ccg.github.io/did-spec/).

To facilitate adoption and testing, and promote open discussions about the subjects
treated, this repository also includes an open source reference implementation for a
CLI client and network agent.

Team members for the project.
- Sandra Murcia / [smurcia@iadb.org](mailto:smurcia@iadb.org)
- Marcos Allende / [marcosal@iadb.org](mailto:marcosal@iadb.org)
- Flavia Munhoso / [flaviamu@iadb.org](mailto:flaviamu@iadb.org)
- Ruben Cessa / [rcessa@iadb.org](mailto:rcessa@iadb.org)

## 1. Decentralized Identifiers
__In order to access online, i.e. digital, services, we need to be electronically
identifiable.__ It means we need an electronic profile that, with a certain level
of assurance, the service provider (either another person or an entity) can trust
it corresponds to our real identity. 

__Conventional identity management systems are based on centralized authorities.__
These authorities establish a process by which to entitle a user temporary access
to a given identifier element. Nevertheless, the true ownership of the identifier
remains on the assigner side and thus, can be removed, revoked and reassigned if
deemed adequate. This creates and intrinsically asymmetric power relationship between
the authority entity and the user. Some examples of this kind of identifiers include:

- Domain names
- Email addresses
- Phone numbers
- IP addresses
- User names

Additionally, from the standpoint of cryptographic trust verification, each of these
centralized authorities serves as its own
[Root of Trust](https://csrc.nist.gov/Projects/Hardware-Roots-of-Trust).

__An alternative model to manage digital identifiers must be open and user-centric__.
It should be considered as such by satisfying at least the following considerations:

- Anyone must have access to freely register, publish and update as many identifiers
  as considered necessary.
- There should be no centralized authority required for the generation and assignment
  of identifiers.
- The end user must have true ownership of the assigned identifiers, i.e. no one but
  the user should be able to remove, revoke and/or reassign the user's identifiers.

This model is commonly referred to as __Decentralized Identifiers__, and allows us to
build a new __(3P)__ digital identity, one that is: __Private, Permanent__ and
__Portable__.

## 2. Access Considerations
In order to be considered open, the system must be publicly available. Any user should
be able to freely register, publish and update as many identifiers as desired without
the express authorization of any third party. This characteristic of the model permits
us to classify it as __censorship resistant.__

At the same time, this level of openness makes the model vulnerable to malicious intentions
and abuse. In such a way that a bad actor may prevent legitimate access to the system by
consuming the available resources. This kind of cyber-attack is known as a 
[DoS (Denial-of-Service) attack](https://en.wikipedia.org/wiki/Denial-of-service_attack).

> In computing, a denial-of-service attack (DoS attack) is a cyber-attack in which the 
  perpetrator seeks to make a machine or network resource unavailable to its intended users
  by temporarily or indefinitely disrupting services of a host connected to the Internet.
  Denial of service is typically accomplished by flooding the targeted machine or resource
  with superfluous requests in an attempt to overload systems and prevent some or all
  legitimate requests from being fulfilled.

The "bryk" DID Method specification includes a __"Request Ticket"__ security mechanism
designed to mitigate risks of abuse while ensuring open access and censorship resistance.

## 3. DID Method Specification
The method specification provides all the technical considerations, guidelines and
recommendations produced for the design and deployment of the DID method implementation.
The document is organized in 3 main sections.

1. __DID Schema.__ Definitions and conventions used to generate valid identifier instances.
2. __DID Document.__ Considerations on how to generate and use the DID document associated
   with a given identifier instance.
3. __Network Operations.__ Technical specifications detailing how to perform basic
  network operations, and the risk mitigation mechanisms in place, for tasks such as:
    - Publish a new identifier instance.
    - Update an existing identifier instance.
    - Resolve an existing identifier and retrieve the latest published version of its DID
    Document.

### 3.1 DID Schema

A Decentralized Identifier is defined as a [RFC3986](https://tools.ietf.org/html/rfc3986)
Uniform Resource Identifier, with a format based on the generic DID schema. Fore more
information you can refer to the
[original documentation](https://w3c-ccg.github.io/did-spec/#decentralized-identifiers-dids).

```abnf
did-reference      = did [ "/" did-path ] [ "#" did-fragment ]
did                = "did:" method ":" specific-idstring
method             = 1*methodchar
methodchar         = %x61-7A / DIGIT
specific-idstring  = idstring *( ":" idstring )
idstring           = 1*idchar
idchar             = ALPHA / DIGIT / "." / "-"
```

Example of a simple Decentralized Identifier (DID).

```
did:example:123456789abcdefghi
```

Expanding on the previous definitions the bryk DID Method specification use the
following format.

```abnf
did                = "did:bryk:" [tag ":"] specific-idstring
tag                = 1*tagchar
specific-idstring  = depends on the particular use case
tagchar            = ALPHA / DIGIT / "." / "-"
```

The optional `tag` element provides a flexible namespace mechanism that can be used
to classify identifier instances into logical groups of arbitrary complexity.

The `specific-idstring` field does not impose any format requirements to ensure the
maximum level of flexibility to end users and implementers. The official implementation
however, proposes and recommends two formal modes for id strings.

#### 3.1.1 Mode UUID
The id string should be a randomly generated lower-case UUID v4 instance as defined by
[RFC4122](https://tools.ietf.org/html/rfc4122). The formal schema for the
`specific-idstring` field on this mode is the following.

```abnf
specific-idstring      = time-low "-" time-mid "-"
                         time-high-and-version "-"
                         clock-seq-and-reserved
                         clock-seq-low "-" node
time-low               = 4hexOctet
time-mid               = 2hexOctet
time-high-and-version  = 2hexOctet
clock-seq-and-reserved = hexOctet
clock-seq-low          = hexOctet
node                   = 6hexOctet
hexOctet               = hexDigit hexDigit
hexDigit               = "0" / "1" / "2" / "3" / "4" / "5" / "6" / "7" /
                         "8" / "9" / "a" / "b" / "c" / "d" / "e" / "f"
```

Example of a DID instance of mode UUID with a `tag` value of `c137`.

```abnf
did:bryk:c137:02825c9d-6660-4f17-92db-2bd22c4ed902
```

#### 3.1.2 Mode Hash
The id string should be a randomly generated 32 byte [SHA3-256](https://goo.gl/Wx8pTY)
hash value, encoded in hexadecimal format as a lower-case string of 64 characters. The formal
schema for the `specific-idstring` field on this mode is the following.

```abnf
specific-idstring = 32hexOctet
hexOctet          = hexDigit hexDigit
hexDigit          = "0" / "1" / "2" / "3" / "4" / "5" / "6" / "7" /
                    "8" / "9" / "a" / "b" / "c" / "d" / "e" / "f"
```

Example of a DID instance of mode hash with a `tag` value of `c137`.

```
did:bryk:c137:85d48aebe67da2fdd273d03071de663d4fdd470cff2f5f3b8b41839f8b07075c
```

### 3.2 DID Document
A Decentralized Identifier, regardless of its particular method, can be resolved
to a standard resource describing the subject. This resource is called a
[DID Document](https://w3c-ccg.github.io/did-spec/#did-documents), and typically
contains, among other relevant details, cryptographic material that enables
authentication of the DID subject.

The document is a Linked Data structure that ensures a high degree of flexibility
while facilitating the process of acquiring, parsing and using the contained information.
For the moment, the suggested encoding format for the document is
[JSON-LD](https://www.w3.org/TR/json-ld/). Other formats could be used in the future.

> The term Linked Data is used to describe a recommended best practice for exposing
  sharing, and connecting information on the Web using standards, such as URLs,
  to identify things and their properties. When information is presented as Linked
  Data, other related information can be easily discovered and new information can be
  easily linked to it. Linked Data is extensible in a decentralized way, greatly
  reducing barriers to large scale integration. 

At the very least, the document must include the DID subject it's referring to under the `id` key.

```json
{
  "@context": "https://w3id.org/did/v1",
  "id": "did:bryk:c137:b616fca9-ad86-4be5-bc9c-0e3f8e27dc8d"
}
```

As it stands, this document is not very useful in itself. Other relevant details that
are often included in a DID Document are:

- [Created](https://w3c-ccg.github.io/did-spec/#created-optional):
  Timestamp of the original creation.
- [Updated](https://w3c-ccg.github.io/did-spec/#updated-optional):
  Timestamp of the most recent change.
- [Public Keys](https://w3c-ccg.github.io/did-spec/#public-keys):
  Public keys are used for digital signatures, encryption and other cryptographic
  operations, which in turn are the basis for purposes such as authentication, secure
  communication, etc.
- [Authentication](https://w3c-ccg.github.io/did-spec/#authentication):
  List the enabled mechanisms by which the DID subject can cryptographically prove
  that they are, in fact, associated with a DID Document.
- [Services](https://w3c-ccg.github.io/did-spec/#service-endpoints):
  In addition to publication of authentication and authorization mechanisms, the
  other primary purpose of a DID Document is to enable discovery of service endpoints
  for the subject. A service endpoint may represent any type of service the subject
  wishes to advertise, including decentralized identity management services for
  further discovery, authentication, authorization, or interaction. 
- [Proof](https://w3c-ccg.github.io/did-spec/#proof-optional):
  Cryptographic proof of the integrity of the DID Document according its subject.

Is important to note that the official specifications around service endpoints are
still in a very early stage at this point. Where appropriate or required the present
Method specification builds on it and introduces new considerations.
    
Additionally, the DID Document may include any other fields deemed relevant for the
particular use case or implementation.

Example of a more complete, and useful, DID Document.
```json
{
  "@context": [
    "https://w3id.org/did/v1",
    "https://w3id.org/security/v1"
  ],
  "id": "did:bryk:c137:eeb0c865-ce21-4ad6-baf8-5ba287ba8683",
  "created": "2019-03-09T15:44:15Z",
  "updated": "2019-03-09T15:44:15Z",
  "publicKey": [
    {
      "id": "did:bryk:c137:eeb0c865-ce21-4ad6-baf8-5ba287ba8683#master",
      "type": "Ed25519VerificationKey2018",
      "controller": "did:bryk:c137:eeb0c865-ce21-4ad6-baf8-5ba287ba8683",
      "publicKeyBase58": "CmTmF8kiepYmsPBXgyjdPW5dpvpW9J2RfVHGUcJJHhMg"
    },
    {
      "id": "did:bryk:c137:eeb0c865-ce21-4ad6-baf8-5ba287ba8683#backup",
      "type": "Ed25519VerificationKey2018",
      "controller": "did:bryk:c137:eeb0c865-ce21-4ad6-baf8-5ba287ba8683",
      "publicKeyHex": "5d0b7e4efb804fdb967890ab66cd5b793db47e7a1d9692e3edd8893cc41a8d63"
    }
  ],
  "authentication": [
    "did:bryk:c137:eeb0c865-ce21-4ad6-baf8-5ba287ba8683#master",
    "did:bryk:c137:eeb0c865-ce21-4ad6-baf8-5ba287ba8683#backup"
  ],
  "service": [
    {
      "id": "did:bryk:c137:eeb0c865-ce21-4ad6-baf8-5ba287ba8683;portal.service",
      "type": "identity.bryk.io.ExternalService",
      "serviceEndpoint": "http://c137.com/path/to/user",
      "data": {
        "individual": "rick"
      }
    }
  ],
  "proof": {
    "@context": [
      "https://w3id.org/security/v1"
    ],
    "type": "Ed25519Signature2018",
    "creator": "did:bryk:c137:eeb0c865-ce21-4ad6-baf8-5ba287ba8683#master",
    "created": "2019-03-09T15:44:16Z",
    "domain": "identity.bryk.io",
    "nonce": "c3dcd2ec89e439f18cea8767abf379c7",
    "proofValue": "QAYz9GlVsVhf4KaZdnu5KMGCKTPK026CZg3fXxQYU7EDZ/0URlgYwBdHIOzAG8ZBIUCGdDQEk7nmlj3DwTJaDg=="
  }
}
```

#### 3.2.1 Method Requirements
Building upon the base requirements and recommendations from the original specification,
the "bryk" DID method introduces the following additional guidelines.

- The fields `created`, `updated`, and `proof` are required for all generated
  DID Documents.
- All service endpoints included in the DID Document may include an additional `data`
  field. Is recommended to include all extra parameters required for the particular
  service under these field.
- Supported public keys and signature formats
  - [Ed25519](https://w3c-ccg.github.io/ld-cryptosuite-registry/#ed25519signature2018)
  - [RSA](https://w3c-ccg.github.io/ld-cryptosuite-registry/#rsasignature2018),
    with a minimum length of 4096 bits.
  - [secp256k1](https://w3c-ccg.github.io/ld-cryptosuite-registry/#eddsasasignaturesecp256k1),
    is __not__ supported.

More information on the official keys and signatures formats is available at
[LD Cryptographic Suite Registry](https://w3c-ccg.github.io/ld-cryptosuite-registry/).

### 3.3 Network Operations
The method implementation introduces the concept of a __network agent__. A network 
agent is responsible for handling incoming client requests. It's very important to
note that the agent itself adheres to an operational protocol. The protocol is
independent of the data storage and message delivery mechanisms used. The method
protocol can be implemented using a __Distributed Ledger Platform__, as well as any
other infrastructure applicable for the particular use case.

There are two main groups of operations available, __read__ and __write__. Write operations
are required when a user wishes to publish a new identifier record to the network, or
update the available information on an existing one. Read operations enable resolution
and retrieval of DID Documents, and other relevant assets, published in the network.

#### 3.3.1 Request Ticket
As described earlier, a security mechanism is required to prevent malicious and
abusive activities. For these purposes, we introduce a __ticket__ requirement for all
write network operations. The mechanism is based on the original
[HashCash](http://www.hashcash.org/hashcash.pdf) algorithm and aims to mitigate
the following problems.

- __Discourage [DoS Attacks](https://en.wikipedia.org/wiki/Denial-of-service_attack)__.
  By making the user cover the “costs” of submitting a request for processing.
- __Prevent [Replay Attacks](https://en.wikipedia.org/wiki/Replay_attack)__.
  Validating the ticket was specifically generated for the request being processed.
- __Requests Authentication__.
  Ensuring the user submitting the ticket is the owner of the DID, by incorporating
  a digital signature requirement that covers both the ticket details and the
  DID instance.

A request ticket has the following structure.

```
ticket {
  int64 timestamp
  int64 nonce
  string key_id
  bytes content
  bytes signature
}
```

The client generates a ticket for the request using the following algorithm.

1. Let the __"bel"__ function be a method to produce a binary-encoded representation
   of a given input value using little endian byte order.
2. Let the __"hex"__ function be a method to produce a hexadecimal binary-encoded
   representation of a given input value.
3. __"timestamp"__ is set to the current UNIX time at the moment of creating the ticket.
4. __"nonce"__ is a randomly generated integer of 64 bit precision.
5. __"key-id"__ is set to the identifier from the cryptographic key used to generate the
   ticket signature, MUST be enabled as an authentication key for the DID instance.
6. __"content"__ is set to a deterministic binary encoding of the DID Document to
  process.
7. A HashCash round is initiated for the ticket. The hash mechanism used MUST be SHA3-256
   and the content submitted for each iteration of the round is a byte concatenation of
   the form: `"bel(timestamp) | bel(nonce) | hex(key-id) | content"`.
8. The __"nonce"__ value of the ticket is atomically increased by one for each iteration
   of the round.
9. The ticket's __"challenge"__ is implicitly set to the produced hash from the
   HashCash round.
10. The __"signature"__ for the ticket is generated using the selected key of the DID
   and the obtained challenge value: `did.keys["key-id"].sign(challenge)`

Upon receiving a new write request the network agent validates the request ticket using
the following procedure.

1. Verify the ticket's __"challenge"__ is valid by forming a HashCash verification.
2. Validate __“contents”__ are a properly encoded DID instance.
3. DID instance’s __“method”__ value is properly set, in this case to “bryk”
4. Ensure __“contents”__ don’t include any private key. For security reasons no private
   keys should ever be published on the network.
5. Verify __“signature”__ is valid.
    - For operations submitting a new entry, the key contents are obtained directly
      from the ticket contents. This ensures the user submitting the new DID instance
      is the one in control of the corresponding private key.
    - For operations updating an existing entry, the key contents are obtained from
      the previously stored record. This ensures the user submitting the request is
      the one in control of the original private key.
6. If the request is valid, the entry will be created or updated accordingly.

#### 3.3.2 DID Resolution

The simplest mechanism to resolve a particular DID instance to the latest published
version of its corresponding DID Document is using the agent's HTTP interface.

The resolution and data retrieval is done by performing a __GET__ request of the form:

`https://did.bryk.io/v1/retrieve?subjet={DID subject}`

For example:

```bash
curl -v https://did.bryk.io/v1/retrieve?subjet=c137:eeb0c865-ce21-4ad6-baf8-5ba287ba8683
```

If the subject is valid, and information has been published to the network the response
will be latest version available of its corresponding DID Document encoded in JSON-LD
with a __200__ status code. If no information is available the response will be a JSON
encoded error message with a __400__ status code.

You can also retrieve an existing subject using the provided SDK and RPC interface.
For example, using the Go client.

```go
// Error handling omitted for brevity
sub := "c137:eeb0c865-ce21-4ad6-baf8-5ba287ba8683"
response, _ := client.Retrieve(context.TODO(), proto.Request{Subject:sub})
if response.Ok {
	id := new(did.Identifier)
	id.Decode(response.Contents)
}
```

#### 3.3.3 DID Publishing and Update

To publish a new identifier instance or to update an existing one you can also use the
agent's HTTP interface or the provided SDK and clients.

When using HTTP the operation should be a __POST__ request with a properly constructed
and JSON encoded request ticket as the request's data. Binary data should be encoded in
standard [Base64](https://en.wikipedia.org/wiki/Base64) when transmitted using JSON.

Example of publish operation.

```bash
# Binary contents redacted for brevity
curl -v \
--header "Content-Type: application/json" \
--request POST \
--data \
'{"timestamp":"1552226666","nonce":"36219","keyId":"master","content":"...","signature":"..."}' \
https://did.bryk.io/v1/process
``` 

You can also publish and update a DID identifier instance using the provided SDK and
RPC interface. For example, using the Go client.

```go
// Error handling omitted for brevity
res, _ := client.Process(context.TODO(), ticket)
if res.Ok {
	// ...
}
```
