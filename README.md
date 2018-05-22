# smpp-go

SMPP Libraries and Simulators in golang

---

SmppPdu
    CommandLength uint32
    CommandId uint32
    CommandStatus uint32
    SequenceNumber uin32
    [] MandatoryParameters
    [] OptionalParameters

    encode()
    decode()

SmppParameter
    Name string
    Type enum
    Value {}Interface
