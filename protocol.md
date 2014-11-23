# Dorp Protocol

The dorp server and client communicate via tcp.

Before the communication begins, the client and server must both be configured with:

1. The pre-shared 32 byte private key.
2. The port to communicate on

Upon opening a connection, the server sends a 64 byte message to the client.
The first 40 bytes are a nacl.secretbox containing the nonce the server expect, and the remaining 24 are the nonce the server used to encrypt the message.

The client, when ready to send an update, should do the following:

1. Convert the states into bytes, where dorp.Positive is 1 and dorp.Negative is 0.
2. Stick them in an array as [door, light].
3. Encrypt the byte array using secretbox with the preshared key and the last nonce recieved from the server.
4. Send the encrypted array to the server.
5. Listen for the server to reply with a new nonce to expect. It's sent in the same format as the inital contact.
