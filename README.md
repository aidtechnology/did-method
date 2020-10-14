# DID Method

[![Build Status](https://github.com/bryk-io/did-method/workflows/ci/badge.svg?branch=master)](https://github.com/bryk-io/did-method/actions)
[![Version](https://img.shields.io/github/tag/bryk-io/did-method.svg)](https://github.com/bryk-io/did-method/releases)
[![Software License](https://img.shields.io/badge/license-BSD3-red.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bryk-io/did-method?style=flat)](https://goreportcard.com/report/github.com/bryk-io/did-method)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0-ff69b4.svg)](.github/CODE_OF_CONDUCT.md)

The present document describes the __"bryk"__ DID Method specification. The
definitions, conventions and technical details included intend to provide a 
solid base for further developments while maintaining compliance with the work,
still in progress, on the [W3C Credentials Community Group](https://w3c-ccg.github.io/did-spec/).

For more information about the origin and purpose of Decentralized Identifiers please
refer to the original [DID Primer.](https://github.com/WebOfTrustInfo/rwot5-boston/blob/master/topics-and-advance-readings/did-primer.md)

To facilitate adoption and testing, and promote open discussions about the subjects
treated, this repository also includes an open source reference implementation for a
CLI client and network agent. You can directly download the binary from the
[published releases](https://github.com/bryk-io/did-method/releases). 

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
build a new __(3P)__ digital identity: __Private, Permanent__ and
__Portable__.

## 2. Access Considerations
In order to be considered open, the system must be publicly available. Any user
should be able to freely register, publish and update as many identifiers as
desired without the express authorization of any third party. This characteristic
of the model permits us to classify it as __censorship resistant.__

At the same time, this level of openness makes the model vulnerable to malicious
intentions and abuse. In such a way that a bad actor may prevent legitimate access
to the system by consuming the available resources. This kind of cyber-attack is
known as a [DoS (Denial-of-Service) attack](https://en.wikipedia.org/wiki/Denial-of-service_attack).

> In computing, a denial-of-service attack (DoS attack) is a cyber-attack in which
  the perpetrator seeks to make a machine or network resource unavailable to its
  intended users by temporarily or indefinitely disrupting services of a host
  connected to the Internet. Denial of service is typically accomplished by
  flooding the targeted machine or resource with superfluous requests in an attempt
  to overload systems and prevent some or all legitimate requests from being
  fulfilled.

The "bryk" DID Method specification includes a __"Request Ticket"__ security
mechanism designed to mitigate risks of abuse while ensuring open access and
censorship resistance.

## 3. DID Method Specification
The method specification provides all the technical considerations, guidelines and
recommendations produced for the design and deployment of the DID method
implementation. The document is organized in 3 main sections.

1. __DID Schema.__ Definitions and conventions used to generate valid identifier
   instances.
2. __DID Document.__ Considerations on how to generate and use the DID document
   associated with a given identifier instance.
3. __Agent Protocol.__ Technical specifications detailing how to perform basic
  network operations, and the risk mitigation mechanisms in place, for tasks such as:
    - Publish a new identifier instance.
    - Update an existing identifier instance.
    - Resolve an existing identifier and retrieve the latest published version of
      its DID Document.

### 3.1 DID Schema

A Decentralized Identifier is defined as a [RFC3986](https://tools.ietf.org/html/rfc3986)
Uniform Resource Identifier, with a format based on the generic DID schema. Fore more
information you can refer to the
[original documentation](https://w3c.github.io/did-core/#generic-did-syntax).

```abnf
did                = "did:" method-name ":" method-specific-id
method-name        = 1*method-char
method-char        = %x61-7A / DIGIT
method-specific-id = *idchar *( ":" *idchar )
idchar             = ALPHA / DIGIT / "." / "-" / "_"
did-url            = did *( ";" param ) path-abempty [ "?" query ]
                     [ "#" fragment ]
param              = param-name [ "=" param-value ]
param-name         = 1*param-char
param-value        = *param-char
param-char         = ALPHA / DIGIT / "." / "-" / "_" / ":" /
                     pct-encoded
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
The id string should be a randomly generated 32 bytes [SHA3-256](https://goo.gl/Wx8pTY)
hash value, encoded in hexadecimal format as a lower-case string of 64 characters.
The formal schema for the `specific-idstring` field on this mode is the following.

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
contains, among other relevant details, cryptographic material to support
authentication of the DID subject.

The document is a Linked Data structure that ensures a high degree of flexibility
while facilitating the process of acquiring, parsing and using the contained
information. For the moment, the suggested encoding format for the document is
[JSON-LD](https://www.w3.org/TR/json-ld/). Other formats could be used in the future.

> The term Linked Data is used to describe a recommended best practice for exposing
  sharing, and connecting information on the Web using standards, such as URLs,
  to identify things and their properties. When information is presented as Linked
  Data, other related information can be easily discovered and new information can be
  easily linked to it. Linked Data is extensible in a decentralized way, greatly
  reducing barriers to large scale integration. 

At the very least, the document must include the DID subject it's referring to under
the `id` key.

```json
{
  "@context": "https://www.w3.org/ns/did/v1",
  "id": "did:bryk:c137:b616fca9-ad86-4be5-bc9c-0e3f8e27dc8d"
}
```

As it stands, this document is not very useful in itself. Other relevant details that
are often included in a DID Document are:

- [Created:](https://w3c-ccg.github.io/did-spec/#created-optional)
  Timestamp of the original creation.
- [Updated:](https://w3c-ccg.github.io/did-spec/#updated-optional)
  Timestamp of the most recent change.
- [Public Keys:](https://w3c-ccg.github.io/did-spec/#public-keys)
  Public keys are used for digital signatures, encryption and other cryptographic
  operations, which in turn are the basis for purposes such as authentication, secure
  communication, etc.
- [Authentication:](https://w3c-ccg.github.io/did-spec/#authentication)
  List the enabled mechanisms by which the DID subject can cryptographically prove
  that they are, in fact, associated with a DID Document.
- [Services:](https://w3c-ccg.github.io/did-spec/#service-endpoints)
  In addition to publication of authentication and authorization mechanisms, the
  other primary purpose of a DID Document is to enable discovery of service endpoints
  for the subject. A service endpoint may represent any type of service the subject
  wishes to advertise, including decentralized identity management services for
  further discovery, authentication, authorization, or interaction. 

Additionally, the DID Document may include any other fields deemed relevant for the
particular use case or implementation.

Example of a more complete, and useful, DID Document.
```json
{
  "@context": [
    "https://www.w3.org/ns/did/v1",
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
  ]
}
```

Is important to note that the official specifications around service endpoints are
still in a very early stage at this point. Where appropriate or required the present
Method specification builds on it and introduces new considerations. 

#### 3.2.1 Method Requirements
Building upon the base requirements and recommendations from the original
specification, the "bryk" DID method introduces the following additional guidelines.

- The fields `created` and `updated` are required for all generated
  DID Documents.
- All service endpoints included in the DID Document may include an additional `data`
  field. Is recommended to include all extra parameters required for the particular
  service under this field.
- Supported public keys and signature formats
  - [Ed25519](https://w3c-ccg.github.io/ld-cryptosuite-registry/#ed25519signature2018)
  - [RSA](https://w3c-ccg.github.io/ld-cryptosuite-registry/#rsasignature2018)
    (with a minimum length of 4096 bits).
  - [secp256k1](https://w3c-ccg.github.io/ld-cryptosuite-registry/#eddsasasignaturesecp256k1)

More information on the official keys and signatures formats is available at
[LD Cryptographic Suite Registry](https://w3c-ccg.github.io/ld-cryptosuite-registry/).

#### 3.2.2 Proofs

[proof:](https://w3c-ccg.github.io/did-spec/#proof-optional) Cryptographic proof
of the integrity of the DID Document according its subject. Recently it was removed
from the DID core document. This method still generates valid proofs for all mutations
performed on the DID documents and returns it under the `proof` element of all
resolved identifiers.

```json
{
  "document": "...",
  "proof": {
    "@context": [
      "https://w3id.org/security/v1"
    ],
    "type": "Ed25519Signature2018",
    "created": "2020-08-08T03:12:53Z",
    "domain": "did.bryk.io",
    "nonce": "3ec84acf8b301f3d7e0bba25a24b438a",
    "proofPurpose": "authentication",
    "verificationMethod": "did:bryk:46389176-6109-4de7-bdb4-67e4fcf0230d#master",
    "proofValue": "QvVkJxTWHf6BQO5A/RzgqDoz6neKaagHWspwSeWqztWnjnt7Rlc73KKiHRs9++C2tdV3pZQtPiKDk6C7Q7nFAQ=="
  }
}
```

> More information about this change is [available here](https://github.com/w3c/did-core/issues/293).

### 3.3 Agent Protocol
The method implementation introduces the concept of a __network agent__. A network 
agent is responsible for handling incoming client requests. It's very important to
note that the agent itself adheres to an operational protocol. The protocol is
independent of the data storage and message delivery mechanisms used. The method
protocol can be implemented using a __Distributed Ledger Platform__, as well as any
other infrastructure suitable for the particular use case.

There are two main groups of operations available, __read__ and __write__. Write
operations are required when a user wishes to publish a new identifier record to
the network, or update the available information for an existing one. Read
operations enable resolution and retrieval of DID Documents and other relevant
assets published in the network.

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
  int64  timestamp
  int64  nonce
  string key_id
  bytes  document
  bytes  proof
  bytes  signature
}
```

The client generates a ticket for the request using the following algorithm.

1. Let the __"bel"__ function be a method to produce a deterministic binary-encoded
   representation of a given input value using little endian byte order.
2. Let the __"hex"__ function be a method to produce a deterministic hexadecimal
   binary-encoded representation of a given input value.
3. __"timestamp"__ is set to the current UNIX time at the moment of creating the
   ticket.
4. __"nonce"__ is a randomly generated integer of 64bit precision.
5. __"key_id"__ is set to the identifier from the cryptographic key used to
   generate the ticket signature, MUST be enabled as an authentication key for the
   DID instance.
6. __"document"__ is set to the JSON-encoded DID Document to process.
7. __"proof"__ is set to the JSON-encoded valid proof for the DID Document to process.
8. A HashCash round is initiated for the ticket. The hash mechanism used MUST be
   SHA3-256 and the content submitted for each iteration of the round is a byte
   concatenation of the form:
   `"bel(timestamp) | bel(nonce) | hex(key_id) | document | proof"`.
9. The __"nonce"__ value of the ticket is atomically increased by one for each
   iteration of the round.
10. The ticket's __"challenge"__ is implicitly set to the produced hash from the
   HashCash round.
11. The __"signature"__ for the ticket is generated using the selected key of the DID
    and the obtained challenge value: `did.keys["key_id"].sign(challenge)`

Upon receiving a new write request the network agent validates the request ticket
using the following procedure.

1. Verify the ticket's `challenge` is valid by performing a HashCash
   verification.
2. Validate `document` are a properly encoded DID Document.
3. Validate `proof` is valid for the DID Document included in the ticket.
4. DID instance `method` value is properly set and supported by the agent.
5. Ensure `document` don’t include any private key. For security reasons no
   private keys should ever be published on the network.
6. Verify `signature` is valid.
    - For operations submitting a new entry, the key contents are obtained directly
      from the ticket contents. This ensures the user submitting the new DID instance
      is the one in control of the corresponding private key.
    - For operations updating an existing entry, the key contents are obtained from
      the previously stored record. This ensures the user submitting the request is
      the one in control of the original private key.
7. If the request is valid, the entry will be created or updated accordingly.

A sample implementation of the described __Request Ticket__ mechanism is available
[here](https://github.com/bryk-io/did-method/blob/master/proto/v1/ticket.go).

#### 3.3.2 DID Resolution

The simplest mechanism to resolve a particular DID instance to the latest published
version of its corresponding DID Document is using the agent's HTTP interface.

The resolution and data retrieval is done by performing a __GET__ request of the form:

`https://did.bryk.io/v1/retrieve/{{method}}/{{subject}}`

For example:

```bash
curl -v https://did.bryk.io/v1/retrieve/bryk/4d81bd52-2edb-4703-b8fc-b26d514a9c56
```

If the subject is valid, and information has been published to the network, the
response will include the latest version available of its corresponding DID Document
encoded in JSON-LD with a __200__ status code. If no information is available the
response will be a JSON encoded error message with a __404__ status code.

```json
{
  "document": "...",
  "proof": "..."
}
```

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

To publish a new identifier instance or to update an existing one you can also use
the agent's HTTP interface or the provided SDK and clients.

When using HTTP the operation should be a __POST__ request with a properly
constructed and JSON-encoded request as the request's data. Binary data should be
encoded in standard [Base64](https://en.wikipedia.org/wiki/Base64) when transmitted
using JSON.

You can also publish and update a DID identifier instance using the provided SDK and
RPC interface. For example, using the Go client.

```go
// Error handling omitted for brevity
res, _ := client.Process(context.TODO(), request)
if res.Ok {
	// ...
}
```

## 4. Client Operations

> To enable the full functionality of DIDs and DID Documents on a particular
  distributed ledger or network (called the target system), a DID method
  specification MUST specify how each of the following CRUD operations is performed
  by a client. Each operation MUST be specified to the level of detail necessary to
  build and test interoperable client implementations with the target system.

The following sections provide detailed descriptions and examples of all required
CRUD base operations and some more advanced use cases. As described earlier, all
supported operations can be accessed using either the agent's HTTP interface or the
provided SDK and CLI client tool.

For brevity the following examples use the provided CLI client tool.

### 4.1 CRUD Operations

Basic operations enabling the users to create, read, update and delete identifier
instances.

#### 4.1.1 Create (Register)

To locally create a new DID instance.

```
didctl create [reference name]
```

The value provided for `reference name` is an easy-to-remember alias you choose for
the new identifier instance, __it won't have any use in the network context__.
The CLI also performs the following tasks for the newly generated identifier.

- Create a new `master` Ed25519 private key for the identifier
- Set the `master` key as an authentication mechanism for the identifier
- Generates a cryptographic integrity proof for the identifier using the `master` key

If required, the `master` key can be recovered using the selected `recovery-mode`,
for more information inspect the options available for the `create` command.

```
Creates a new DID locally

Usage:
  didctl register [flags]

Aliases:
  register, create, new

Examples:
didctl register [DID reference name]

Flags:
  -h, --help                    help for register
      --recovery-mode string    choose a recovery mechanism for your primary key, 'passphrase' or 'secret-sharing' (default "secret-sharing")
      --secret-sharing string   specify the number of shares and threshold value in the following format: shares,threshold (default "3,2")
      --tag string              specify a tag value for the identifier instance
``` 

#### 4.1.2 Read (Verify)

You can retrieve a list of all your existing identifiers using the following command.

```
didctl list
```

The output produced will be something like this.

```
Reference Name    Recovery Mode     DID
dev               passphrase        did:bryk:4d81bd52-2edb-4703-b8fc-b26d514a9c56
sample            secret-sharing    did:bryk:99dc4a30-7434-42e5-ac75-5f330be0ea0a
```

To inspect the DID Document of your local identifiers.

```
didctl info [reference name]
```

The generated document will be something similar for the following example.

```json
{
  "@context": [
    "https://www.w3.org/ns/did/v1",
    "https://w3id.org/security/v1"
  ],
  "id": "did:bryk:99dc4a30-7434-42e5-ac75-5f330be0ea0a",
  "created": "2019-03-14T12:02:33-04:00",
  "updated": "2019-03-14T12:02:33-04:00",
  "publicKey": [
    {
      "id": "did:bryk:99dc4a30-7434-42e5-ac75-5f330be0ea0a#master",
      "type": "Ed25519VerificationKey2018",
      "controller": "did:bryk:99dc4a30-7434-42e5-ac75-5f330be0ea0a",
      "publicKeyHex": "e5271fa5208eedf6c95611320ed8c4300dcd04ab57207364a0909fc64c5e30d7"
    }
  ],
  "authentication": [
    "did:bryk:99dc4a30-7434-42e5-ac75-5f330be0ea0a#master"
  ],
  "proof": {
    "@context": [
      "https://w3id.org/security/v1"
    ],
    "type": "Ed25519Signature2018",
    "creator": "did:bryk:99dc4a30-7434-42e5-ac75-5f330be0ea0a#master",
    "created": "2019-03-14T16:02:34Z",
    "domain": "did.bryk.io",
    "nonce": "7e0fda2827eec4418df4513d2d6874c5",
    "proofValue": "cKDKBgNS6itQF1zDaOUd6bDo+5CIKoSN+lOb8PkZqvT+K3c2wvDUVqMYN8mKA0Om+B8wYM1qDz9mI0iWva0qBg=="
  }
}
```

If a certain DID identifier has previously been published to the network, you can
resolve it and retrieve the latest version of its corresponding DID Document using
the `get` command. To run a verification of the cryptographic integrity proof
contained in the document you can add the `--verify` option.

```
didctl get --verify did:bryk:4d81bd52-2edb-4703-b8fc-b26d514a9c56
```

The command will perform the required network operations and verifications.

```
[Mar 14 12:15:23.674]  INFO establishing connection to the network with node: rpc-did.bryk.io:80
[Mar 14 12:15:23.713] DEBUG retrieving record
[Mar 14 12:15:23.721] DEBUG decoding contents
[Mar 14 12:15:23.721]  INFO verifying the received DID document
[Mar 14 12:15:24.277]  INFO integrity proof is valid
{
  "@context": [
    "https://www.w3.org/ns/did/v1",
    "https://w3id.org/security/v1"
  ],
  "id": "did:bryk:4d81bd52-2edb-4703-b8fc-b26d514a9c56",
  "created": "2019-03-10T13:42:34-04:00",
  "updated": "2019-03-12T10:07:55-04:00",
  "publicKey": [
    {
      "id": "did:bryk:4d81bd52-2edb-4703-b8fc-b26d514a9c56#master",
      "type": "Ed25519VerificationKey2018",
      "controller": "did:bryk:4d81bd52-2edb-4703-b8fc-b26d514a9c56",
      "publicKeyHex": "be4db03c2f809aa79ea3055a2da8ddfd807fecd073356e337561cd0640251d9f"
    },
    {
      "id": "did:bryk:4d81bd52-2edb-4703-b8fc-b26d514a9c56#code-sign",
      "type": "Ed25519VerificationKey2018",
      "controller": "did:bryk:4d81bd52-2edb-4703-b8fc-b26d514a9c56",
      "publicKeyHex": "e7cc93d399e467a39fca74e32795b1ab1110a7dc94e8623830cd069c1cac72b8"
    }
  ],
  "authentication": [
    "did:bryk:4d81bd52-2edb-4703-b8fc-b26d514a9c56#master"
  ],
  "proof": {
    "@context": [
      "https://w3id.org/security/v1"
    ],
    "type": "Ed25519Signature2018",
    "creator": "did:bryk:4d81bd52-2edb-4703-b8fc-b26d514a9c56#master",
    "created": "2019-03-12T14:07:56Z",
    "domain": "did.bryk.io",
    "nonce": "09206a2a195cd14f5a6cac70279bba35",
    "proofValue": "YdY1+GxDNwlc55alKZIJ0if55FwQsE2Gan91l+fuv+UF1UAnI10l/DelGPyBSOO2OUiTNzXC6x/jojOum/RNDg=="
  }
}
```

#### 4.1.3 Update (Publish)

Whenever you wish to make one of your identifiers, in its current state, accessible
to the world, you can publish it to the network.

```
didctl sync sample
```

The CLI tool will generate the __Request Ticket__, submit the operation for
processing to the network and present the final result.

```
[Mar 14 12:26:30.435] DEBUG key selected for the operation: did:bryk:99dc4a30-7434-42e5-ac75-5f330be0ea0a#master
[Mar 14 12:26:30.435]  INFO updating record proof
[Mar 14 12:26:31.075]  INFO publishing: sample
[Mar 14 12:26:31.075]  INFO generating request ticket
[Mar 14 12:28:00.230] DEBUG ticket obtained: 00000042186234bab7d3e39207a9fcde7c8e71c2b4e84cf528f0328b3d6e8a32
[Mar 14 12:28:00.230] DEBUG time: 1m29.154355394s (rounds completed 21723044)
[Mar 14 12:28:00.231]  INFO establishing connection to the network with node: rpc-did.bryk.io:80
[Mar 14 12:28:00.232]  INFO submitting request to the network
[Mar 14 12:28:00.234] DEBUG request status: true
```

Once an identifier is published any user can retrieve and validate your DID document.
If you make local changes to your identifier, like adding a new cryptographic key or
service endpoint, and you wish these adjustments to be accessible to the rest of the
users, you'll need to publish it again.

#### 4.1.4 Delete (Deactivate)

If at some point you wish to prevent other users to resolve one of yours previously
published identifiers you may submit a __deactivation__ request by adding the
`--deactivate` option to the sync command.

```
didctl sync sample --deactivate
```

__No information is destroyed or lost with these operations__, the identifier and all
its related data is safely stored on your local machine. This will only prevent other
users from retrieving your DID Document from the network.

### 4.2 DID Instance Management

The CLI client also facilitates some tasks required to manage a DID instance.

#### 4.2.1 Key Management

A DID Document list all public keys in use for the referenced DID instance. Public
keys are used for digital signatures, encryption and other cryptographic operations,
which in turn are the basis for purposes such as authentication, secure communication,
etc.

```
Manage cryptographic keys associated with the DID

Usage:
  didctl edit key [command]

Available Commands:
  add         Add a new cryptographic key for the DID
  recover     Recover a previously generated Ed25519 cryptographic key
  remove      Remove an existing cryptographic key for the DID
  sign        Produce a linked digital signature

Flags:
  -h, --help   help for key

Global Flags:
      --config string   config file ($HOME/.didctl/config.yaml)
      --home string     home directory ($HOME/.didctl)

Use "didctl edit key [command] --help" for more information about a command.
```

To add a new cryptographic key to one of your identifiers you can use the `did key add`
command.

```
Add a new cryptographic key for the DID

Usage:
  didctl edit key add [flags]

Examples:
didctl edit key add [DID reference name] --name my-new-key --type ed --authentication

Flags:
      --authentication   enable this key for authentication purposes
  -h, --help             help for add
      --name string      name to be assigned to the newly added key (default "key-#")
      --type string      type of cryptographic key, either RSA (rsa) or Ed25519 (ed) (default "ed")
```

It will produce and properly add a public key entry. The cryptographic
integrity proof on the DID Document will also be updated accordingly.

```json
{
  "id": "did:bryk:4d81bd52-2edb-4703-b8fc-b26d514a9c56#code-sign",
  "type": "Ed25519VerificationKey2018",
  "controller": "did:bryk:4d81bd52-2edb-4703-b8fc-b26d514a9c56",
  "publicKeyHex": "e7cc93d399e467a39fca74e32795b1ab1110a7dc94e8623830cd069c1cac72b8"
}
```

You can also safely remove an existing key from your identifier using the
`edit key remove` command.

```
edit key remove [DID reference name] [key name]
```

#### 4.2.2 Linked Data Signatures

The CLI client also facilitates the process of generating and validating [Linked Data
Signatures](https://w3c-dvcg.github.io/ld-signatures/). For example, to create a new
signature document from an existing file you can run the following command.

```
cat file_to_sign | didctl sign dev
```

The output produced will be a valid JSON-LD document containing the signature details.

```json
{
  "@context": [
    "https://w3id.org/security/v1"
  ],
  "type": "Ed25519Signature2018",
  "creator": "did:bryk:4d81bd52-2edb-4703-b8fc-b26d514a9c56#master",
  "created": "2019-03-15T14:05:54Z",
  "domain": "did.bryk.io",
  "nonce": "f14d4619a39f7deb5a382bf32b220726",
  "signatureValue": "khqsBcnCViYm/3QFjgAQX2iOGDbNpsD5rPWsokWNLsBxhtRf79A+qV1f+9sphjVCxNP02jesOOni3t9zMCZbBw=="
}
```

You can save and share the produced JSON output. Other users will be able to verify the
integrity and authenticity of the signature using the `verify` command.

```
cat file_to_sign | didctl verify signature.json
```

The CLI will inspect the signature file, retrieve the DID Document for the creator
and use the public key to verify the integrity and authenticity of the signature.

```
[Mar 15 10:10:22.286]  INFO verifying LD signature
[Mar 15 10:10:22.286] DEBUG load signature file
[Mar 15 10:10:22.286] DEBUG decoding contents
[Mar 15 10:10:22.286] DEBUG validating signature creator
[Mar 15 10:10:22.287]  INFO establishing connection to the network with node: rpc-did.bryk.io:80
[Mar 15 10:10:22.458] DEBUG retrieving record
[Mar 15 10:10:22.471] DEBUG decoding contents
[Mar 15 10:10:22.973]  INFO signature is valid
```

#### 4.2.3 Service Management

As mentioned in earlier sections, one of the more relevant aspects of a DID Document
is its capability to list interaction mechanisms available for a particular subject.
This is done by including information of __Service Endpoints__ in the document. Using
the CLI client you can manage the services enabled for any of your identifiers.

```
Manage services enabled for the identifier

Usage:
  didctl edit service [command]

Available Commands:
  add         Register a new service entry for the DID
  remove      Remove an existing service entry for the DID

Flags:
  -h, --help   help for service

Global Flags:
      --config string   config file ($HOME/.didctl/config.yaml)
      --home string     home directory ($HOME/.didctl)

Use "didctl edit service [command] --help" for more information about a command.
```

To add a new service you can use the `did service add` command.

```
Register a new service entry for the DID

Usage:
  didctl edit service add [flags]

Examples:
didctl edit service add [DID reference name] --name "service name" --endpoint https://www.agency.com/user_id

Flags:
      --endpoint string   main URL to access the service
  -h, --help              help for add
      --name string       service's reference name (default "external-service-#")
      --type string       type identifier for the service handler (default "identity.bryk.io.ExternalService")
```

It will produce and properly add a service endpoint entry. The cryptographic
integrity proof on the DID Document will also be updated accordingly.

```json
{
  "id": "did:bryk:99dc4a30-7434-42e5-ac75-5f330be0ea0a;iadb-bonds",
  "type": "identity.bryk.io.ExternalService",
  "serviceEndpoint": "https://www.iadb.org/bonds"
}
```

You can also safely remove a service from your identifier using the
`edit service remove` command.

```
edit service remove [DID reference name] [service name]
```
