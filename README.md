## Overview

This is a golang library for (SMPP v3.4)[http://opensmpp.org/specs/SMPP_v3_4_Issue1_2.pdf] message generation and parsing.

## Usage

The library is in the package **smpp**.  There are two basic object types: **smpp.Parameter** and **smpp.PDU**.  One may *Encode* a **PDU** object to an octet stream in network byte order, and may *Decode* an octet stream in network byte order into a **PDU**.

An SMPP v3.4 PDU has a header and one or more **Parameter**s.  Parameters are Mandatory or Optional.  Mandatory Parameters are at fixed offsets in a PDU, based on the PDU type, and are either of fixed size, or are null terminated.  Optional Parameters are TLV (Tag/Length/Value) encoded.  There are different constructors for different **Paramter** types:

```golang
p := NewFLParameter(value interface{})
p := NewCOctetStringParameter(value string)
p := NewOctetStringFromString(value string)
p := NewTLVParameter(tag uint16, value interface{})
```

`NewFLParameter` is for fixed length parameters, and the length and encoding is inferred from the `value` type, which may be `uint8`, `uint16` or `uint32`.  `NetCOctetStringParameter` is for C-Octet-Strings (which are null--that is, byte with value--terminated).  `OctetStringFromString` produces an Octet-String (which are not null terminated) from a string.  The only Parameter type that uses this is _short_messsage_, which must be accompanied by an _sm_length_ Parameter that provides the _short_message_ length.  `NewTLVParameter` generates an Optional Parameter.  The Length of the TLV is inferred from type of `value`, which may be `uint8`, `uint16`, `uint32`, `string` or `[]byte`.

To create a **PDU**:

```golang
pdu := NewPDU(id CommandIDType, status uint32, sequence uint32, mandatoryParams []*Parameter, optionalParams []*Parameter)
```

## Status

In **pdu.go**, the global variable `pduTypeDefinition` maps the Mandatory Parameters for each PDU type.  It is incomplete; not all types are currently defined.  Those that are do have a limited set of unit tests.  If missing types are defined, unit tests should be created for those types.
