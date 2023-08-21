# Vehicle Position Interface Specification

The distributed rail lab allows testing scenarios including landside and trackside equipment spanning multiple physical locations and laboratories.
An interface exists to stream the live position data of a train on a test track into the distributed lab environment for analysis and experiment control purposes.

<!-- ## Overview

The following figure illustrates the protocol stack for transmission of the train position data.
 -->

## Example Code

Two client implementations for the vehicle position interface exist:

 - [Example.Golang](Example.Golang)
 - [Example.Script](Example.Script)

## Train Position Sensor

A sensing device records the train position. This can include onboard odometry sensors, a GPS sensor or fiberoptic sensing equipment.
In the following, we will use the term 'train position sensor' to identify the entire system that records and forwards sensor data to the lab.


## UDP / DTLS

On the application level, the communication is encrypted and authenticated using X.509 client/server certificates.
The sensor data endpoint of the lab infrastructure is available on UDP port 5684.

## CoAP

The constrained application protocol is used to transmit new sensor data values to the lab infrastructure.

Upon startup, the train position sensor must send an attribute telegram to the lab:
- CoAP Path: /api/v1/attributes
- Message Type: Confirmable

For each new sensor reading, exactly one position telegram must be sent to the lab. Sensor readings do not have to be retransmitted in case of an intermittent packet / network loss (non-confirmable).
- CoAP Path: /api/v1/telemetry
- Message Type: Non-Confirmable

## Protobuf Message: Attribute Telegram

- NID_ENGINE (string)
  - The NID_ENGINE of the active ETCS On Board equipment of the test train (if applicable)

**Protobuf Schema**

```proto
syntax = "proto3";

package esc;

message VehicleAttributes {
    string nid_engine = 1;
}
```

## Protobuf Message: Position Telegram

- Timestamp
  - The absolute timestamp of the sensor reading, in milliseconds since 1970-01-01
  - Time is be coordinated within the lab
  - The train position sensor upon startup must perform a time synchronization using the NTP time servers of the Physikalisch Technische Bundesanstalt (PTB)
    - ptbtime1.ptb.de
    - ptbtime2.ptb.de
    - ptbtime3.ptb.de
    - ptbtime4.ptb.de
- Position
  - Upon startup, the train position sensor defines a starting point and the forward/reverse direction of the train.
  - The determined train position used for further analysis in the lab, in centimeters from the starting point, in positive or negative direction.
  - Depending of the driving direction, negative values are possible
  - The position measurement must refer to a fixed but arbitrary position of the test vehicle.
  - The position must increase or decrease monotonically unless the driving direction of the vehicle has changed.
  - Upon change of driving direction, the reported offset may not jump.
  - Upon change of driving direction, the position counts in the opposite direction.
- Speed
  - The current speed of the vehicle in m/s.
  - This value can be used by the lab to interpolate and predict the current location of the train.
- Session ID
  - Opaque identifier string that can be freely chosen
  - Remains constant unless a train-side reset was performed
  - Can be changed without resetting the distance
  - Should be reset when the reference value of the sensor is reset

**Protobuf Schema**

```proto
syntax = "proto3";

package esc;

message VehiclePosition {
    // in ms, number of milliseconds elapsed since the 01/01/1970
    uint64 timestamp = 1;
    // in cm, current position of the train
    int64 position = 2;
    // in m/s, current speed of the train
    double speed = 3;
    // opaque session id
    string session_id = 4;
}
```
