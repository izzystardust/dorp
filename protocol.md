# Dorp Protocol

Dorp communicates via pushing encrypted data from the client to /set on the server.

Messages sent to the server are built as follows:
1. Create a JSON payload containing door and light state. It should look like
    ```json
    {
        "DoorState": "open|closed",
        "LightState": "on|off",
    }
    ```

2. Create a random nonce and encrypt the JSON payload using a NaCl secretbox.

3. Encode the message and nonce as base64 strings. Join the ciphertext and nonce together with the delimeter character `;`. It should look like `base64CipherText;base64Nonce`.

4. Send that message to the server! Woo!
