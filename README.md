### InterPlanetary Secure Extended Shell

Secure shell session leveraging IPFS's libp2p to enable more secure handshake and extended capabilities over hosts spanning a solar system and beyond.

#### Status: verrry pre-alpha. Proof of concept complete. Some time before it's a "daily driver".

#### Motivation: 

A frequently encountered prompt while using conventional ssh client is the following:
```
The authenticity of host 'example.com (xxx.xxx.xxx.xxx)' can't be established.
ECDSA key fingerprint is SHA256:KI0kMBUx4KAV4TIIhNLdiw1qEU27+7oOa+2M2KAtL+o.
Are you sure you want to continue connecting (yes/no/[fingerprint])? 
```
The power user knows the extra steps needed to prevent this message, but the typical user may choose `yes` as an answer to this prompt. This can leave this session open to a "man in the middle" attack so that all data communicated over the channel defined by this session can be collected by an adversary.

Our goal is a more secure handshake "out of the box" using less configuration than the conventional solution. We have leveraged some of the security features of IPFS's libp2p networking stack to design an updated session handshake that has fast resolution and requires very little configuration compared to former solutions. See the full spec of the handshake below. A happy surprise in this new design is we found it securely enables a shell session multiplexed over a family of hosts, perhaps one on each planet...


// TODO write handshake spec
