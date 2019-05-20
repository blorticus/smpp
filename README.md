# Golang SMPP v3.4 Library

## Overview

This is a golang library for [SMPP v3.4](http://opensmpp.org/specs/SMPP_v3_4_Issue1_2.pdf) message generation and parsing.

## Usage

The library is in the package **smpp**.  There are two basic object types: **smpp.Parameter** and **smpp.PDU**.  One may *Encode* a **PDU** object to an octet stream in network byte order, and may *Decode* an octet stream in network byte order into a **PDU**.

An SMPP v3.4 PDU has a header and one or more **Parameter**s.  Parameters are Mandatory or Optional.  Mandatory Parameters are at fixed offsets in a PDU, based on the PDU type, and are either of fixed size, or are null terminated (or it is a _short_message_, which is treated specially, as described below).  Optional Parameters are TLV (Tag/Length/Value) encoded.  There are different constructors for different **Paramter** types:

```golang
p := smpp.NewFLParameter(value interface{})
p := smpp.NewCOctetStringParameter(value string)
p := smpp.NewOctetStringFromString(value string)
p := smpp.NewTLVParameter(tag uint16, value interface{})
```

`NewFLParameter` is for fixed length parameters, and the length and encoding is inferred from the `value` type, which may be `uint8`, `uint16` or `uint32`.  `NetCOctetStringParameter` is for C-Octet-Strings (which are null--that is, byte with a value of 0--terminated).  `OctetStringFromString` produces an Octet-String (which is not null terminated) from a string.  The only Parameter type that uses this is _short_messsage_, which must be preceded by an _sm_length_ Parameter that provides the _short_message_ length.  `NewTLVParameter` generates an Optional Parameter.  The Length of the TLV is inferred from type of `value`, which may be `uint8`, `uint16`, `uint32`, `string` or `[]byte`.

None of the **Parameter** constructors return an error, so that they can be used inline with a PDU constructor.  They will, however, return `nil` if something goes wrong.

To create a **PDU**:

```golang
pdu := smpp.NewPDU(id smpp.CommandIDType, status uint32, sequence uint32, mandatoryParams []*smpp.Parameter, optionalParams []*smpp.Parameter)
```

To decode an incoming byte stream (which must be a complete PDU):

```golang
pdu, err := smpp.DecodePDU(stream []byte)
```

## Examples

There are examples in the *examples/* directory.

## Status

In **pdu.go**, the global variable `pduTypeDefinition` maps the Mandatory Parameters for each PDU type.  Particularly, neither SubmitMulti nor SubmitMultiResp are implemented, and messages of this type will neither encode nor decode properly.  It requires a bit of extra logic in the code to support these.  Additionally, there are no unit tests for a subset of the message types, so their encode/decode methods are not thoroughly tested.

## TODO

# Implement a maximum length parameter to prevent rapid buffer memory consumption in Peer object
