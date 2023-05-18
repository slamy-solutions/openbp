# System vault service
The `system_vault` service is a crucial component responsible for robust security management within the OpenBP system. Its primary function is to establish communication with the Hardware Security Module (***HSM***) through ***PKCS*** interfaces. It offers a user-friendly API that enables seamless integration with other services.

By leveraging the capabilities of the ***HSM*** and implementing ***PKCS*** interfaces, the `system_vault` service ensures a high level of security for managing critical information within the OpenBP system. Its user-friendly API further facilitates seamless interaction with other services, making it an essential component for secure and efficient operations.

## Seal / Unseal
Upon system startup, the vault enters a `sealed` state, which means that communication with the service is not possible. To make effective use of the service, the vault needs to be manually `unsealed` using a password. It is crucial to keep this password secure as losing it would render the service inaccessible and result in permanent loss of all stored secrets.

!!! warning
    By default, the service comes with a predefined unsealing password, which it attempts to use during startup. However, it is highly recommended to change this password to enhance security. This can be achieved through appropriate API request within the service.

## Usage examples

!!! info
    Those examples are only for external developers. If you are contributor to the OpenBP project, remember to change imports so they are pointing from the inside of the OpenBP repository.

??? example "Encrypt/Decrypt arbitrary data"
    In this example we will encrypt and decrypt string using vaults HSM.

    === "GO"
        ```golang
        package main

        import (
            "bytes"
            "context"
            "fmt"

            openbp "github.com/slamy-solutions/openbp-go"
            "github.com/slamy-solutions/openbp-go/modules"
            "github.com/slamy-solutions/openbp-go/modules/system"
            "github.com/slamy-solutions/openbp-go/modules/system/proto/vault"
        )

        func main() {
            ctx := context.TODO()

            obp, err := openbp.ConnectToModules(ctx, modules.NewStubConfig().
                WithSystemModule(system.NewSystemStubConfig().WithVault()),
            )
            if err != nil {
                panic(err)
            }

            data := []byte("Hello OpenBP")

            // Encrypt data
            encryptClient, _ := obp.System.Vault.EncryptStream(ctx)
            encryptClient.Send(&vault.EncryptStreamRequest{
                PlainData: data,
            })
            encryptClient.CloseSend()

            encryptResponse, _ := encryptClient.Recv()
            encryptedData := encryptResponse.EncryptedData

            // Decrypt data
            decryptClient, _ := obp.System.Vault.DecryptStream(ctx)
            decryptClient.Send(&vault.DecryptStreamRequest{
                EncryptedData: encryptedData,
            })
            decryptClient.CloseSend()

            decryptResponse, _ := decryptClient.Recv()

            // Compare to original data
            fmt.Printf("Data restored correctly: %t", bytes.Equal(data, decryptResponse.PlainData))
        }
        ```

??? example "Generate signature (HMAC)"
    In this example, we will generate HMAC using vaults HSM

    === "GO"
        ```golang
        package main

        import (
            "context"
            "encoding/base64"
            "fmt"

            openbp "github.com/slamy-solutions/openbp-go"
            "github.com/slamy-solutions/openbp-go/modules"
            "github.com/slamy-solutions/openbp-go/modules/system"
            "github.com/slamy-solutions/openbp-go/modules/system/proto/vault"
        )

        func main() {
            ctx := context.TODO()

            obp, err := openbp.ConnectToModules(ctx, modules.NewStubConfig().
                WithSystemModule(system.NewSystemStubConfig().WithVault()),
            )
            if err != nil {
                panic(err)
            }

            data := []byte("Hello OpenBP")

            response, _ := obp.System.Vault.HMACSign(ctx, &vault.HMACSignRequest{
                Data: data,
            })

            signature := base64.URLEncoding.EncodeToString(response.Signature)
            fmt.Println("Signature: " + signature)
        }
        ```

??? example "Generate RSA and sign message with it"
    In this example, we will sign message with RSA key-pair using vaults HSM

    === "GO"
        ```golang
        package main

        import (
            "context"
            "encoding/base64"
            "fmt"

            openbp "github.com/slamy-solutions/openbp-go"
            "github.com/slamy-solutions/openbp-go/modules"
            "github.com/slamy-solutions/openbp-go/modules/system"
            "github.com/slamy-solutions/openbp-go/modules/system/proto/vault"
        )

        func main() {
            ctx := context.TODO()

            obp, err := openbp.ConnectToModules(ctx, modules.NewStubConfig().
                WithSystemModule(system.NewSystemStubConfig().WithVault()),
            )
            if err != nil {
                panic(err)
            }

            keyName := "my-super-rsa-key-pair"
            data := []byte("Hello OpenBP")

            obp.System.Vault.EnsureRSAKeyPair(ctx, &vault.EnsureRSAKeyPairRequest{
                KeyName: keyName,
            })

            response, _ := obp.System.Vault.RSASign(ctx, &vault.RSASignRequest{
                KeyName: keyName,
                Data:    data,
            })

            signature := base64.URLEncoding.EncodeToString(response.Signature)
            fmt.Println("Signature: " + signature)
        }
        ```

## Communication API
The system offers communication capabilities with the service through the gRPC interface. The interface definitions, specified in the proto file, are provided by the `system` module. This allows clients to interact with the `system_vault` service seamlessly using gRPC-based communication protocols.

### General

??? example "rpc Seal(SealRequest) returns (SealResponse) {};"
    Close and encrypt vault. After sealing, most of the operations will not be accessible.

    === "Request"
        The request doesn't have parameters         
    === "OK"
        The response will always be positive, even when the vault is already sealed
??? example "rpc Unseal(UnsealRequest) returns (UnsealResponse) {};"
    Decrypt and open vault. Must be done before most of the operations with vault secrets.

    === "Request"
        | Parameter | Type   | Description                   |
        | --------- | ------ | ----------------------------- |
        | secret    | string | Password to decrypt the vault |
    === "OK"
        | Response value | Type | Description                            |
        | -------------- | ---- | -------------------------------------- |
        | success        | bool | Indicates if vault was unsealed or not |
??? example "rpc UpdateSealSecret(UpdateSealSecretRequest) returns (UpdateSealSecretResponse) {};"
    Set up new seal secret and reincrypt vault. The vault must be unsealed before this operation. You don't need to unseal vault after this operation.

    === "Request"
        This operation requires you to have administrator access to the HSM. Check PKCS11 spec. If you are using emulated HSM (by default) this will be the same as the seal/unseal secret by default. Change it.

        | Parameter            | Type   | Description                     | Default    |
        | -------------------- | ------ | ------------------------------- | ---------- |
        | currentAdminPassword | string | Current administrator password. | "12345678" |
        | newPassword          | string | New administrator password.     | -          |
        | newSecret            | string | New seal/unseal secret.         | -          |
    === "OK"
        The seal/unseal secret was successfully updated. The administrator password was updated.
    === "NOT_AUTHORIZED"
        Bad administrator password
??? example "rpc GetStatus(GetStatusRequest) returns (GetStatusResponse) {};"
    Get current status of the vault.

    === "Request"
        The request doesn't have parameters
    === "OK"
        | Response value | Type | Description            |
        | -------------- | ---- | ---------------------- |
        | sealed         | bool | Is vault sealed or not |
    
### RSA
The `system_vault` service plays a crucial role in securely managing RSA key-pairs. It employs a highly secure approach where the private key is never exposed outside of the Hardware Security Module (HSM) and cannot be extracted from it. On the other hand, the public key, which is intended to be shared and accessible to all, can be retrieved.

!!! info
    To ensure the utmost security, all operations involving the private key are performed exclusively within the HSM. This means that any cryptographic operations, such as signing or decryption, utilize the private key directly within the HSM's protected environment. This approach significantly reduces the risk of private key compromise or unauthorized access.

By securely managing RSA key-pairs within the `system_vault` service and relying on the protection provided by the HSM, the confidentiality and integrity of the private key are maintained while allowing the public key to be widely accessible for encryption, verification, or other public operations.

??? example "rpc EnsureRSAKeyPair(EnsureRSAKeyPairRequest) returns (EnsureRSAKeyPairResponse) {};"
    Create RSA key pair if it doesnt exist. The private key never leaves the HSM.

    === "Request"
        | Parameter | Type   | Description                 |
        | --------- | ------ | --------------------------- |
        | keyName   | string | Unique name of the key-pair |
    === "OK"
        The key-pair was created or already existed before this operation.
    === "FAILED_PRECONDITION"
        The vault is sealed.

??? example "rpc GetRSAPublicKey(GetRSAPublicKeyRequest) returns (GetRSAPublicKeyResponse) {};"
    Get public key of the RSA key-pair.

    === "Request"
        | Parameter | Type   | Description                 |
        | --------- | ------ | --------------------------- |
        | keyName   | string | Unique name of the key-pair |
    === "OK"
        | Response value | Type  | Description                                 |
        | -------------- | ----- | ------------------------------------------- |
        | publicKey      | bytes | Public key in the PKCS #1, ASN.1 DER format |
    === "NOT_FOUND"
        Key-pair with provided name doesn't exist
    === "FAILED_PRECONDITION"
        The vault is sealed.

??? example "rpc RSASignStream(stream RSASignStreamRequest) returns (RSASignStreamResponse) {};"
    Sign message stream with RSA private key. It will use SHA512_RSA_PKCS (RS512) algorithm to sign the message.

    !!! tip
        If your data is short (maximum several kilobytes) - use `RSASign` instead. It will be faster and use less resources.

    === "Request"
        | Parameter | Type   | Description                 |
        | --------- | ------ | --------------------------- |
        | keyName   | string | Unique name of the key-pair |
        | data      | bytes  | Data chunk to sign          |
    === "OK"
        | Response value | Type  | Description                    |
        | -------------- | ----- | ------------------------------ |
        | signature      | bytes | Signature of the provided data |
    === "NOT_FOUND"
        Key-pair with provided name doesn't exist
    === "FAILED_PRECONDITION"
        The vault is sealed.

??? example "rpc RSAVerifyStream(stream RSAVerifyStreamRequest) returns (RSAVerifyStreamResponse) {};"
    Validate signature of the message stream using RSA key-pairs public key. It will use SHA512_RSA_PKCS (RS512) algorithm to verify the message.

    !!! tip
        If your data is short (maximum several kilobytes) - use `RSAVerify` instead. It will be faster and use less resources.

    === "Request"
        | Parameter | Type   | Description                 |
        | --------- | ------ | --------------------------- |
        | keyName   | string | Unique name of the key-pair |
        | data      | bytes  | Data chunk to validate      |
        | signature | bytes  | Signature to validate       |
    === "OK"
        | Response value | Type | Description                                                                   |
        | -------------- | ---- | ----------------------------------------------------------------------------- |
        | valid          | bool | True if and only if provided data and its signature matches provided key-pair |
    === "NOT_FOUND"
        Key-pair with provided name doesn't exist
    === "FAILED_PRECONDITION"
        The vault is sealed.

??? example "rpc RSASign(RSASignRequest) returns (RSASignResponse) {};"
    Sign message with RSA. It will use SHA512_RSA_PKCS (RS512) algorithm to sign the message.

    !!! warning
        The data must be short (max several kilobytes). If it is longer - use `RSASignStream` instead.

    === "Request"
        | Parameter | Type   | Description                 |
        | --------- | ------ | --------------------------- |
        | keyName   | string | Unique name of the key-pair |
        | data      | bytes  | Data to sign                |
    === "OK"
        | Response value | Type  | Description                    |
        | -------------- | ----- | ------------------------------ |
        | signature      | bytes | Signature of the provided data |
    === "NOT_FOUND"
        Key-pair with provided name doesn't exist
    === "FAILED_PRECONDITION"
        The vault is sealed.

??? example "rpc RSAVerify(RSAVerifyRequest) returns (RSAVerifyResponse) {};"
    Validate signature of the message using RSA key-pairs public key. It will use SHA512_RSA_PKCS (RS512) algorithm to verify the message.

    !!! warning
        The data must be short (max several kilobytes). If it is longer - use `RSAVerifyStream` instead.

    === "Request"
        | Parameter | Type   | Description                 |
        | --------- | ------ | --------------------------- |
        | keyName   | string | Unique name of the key-pair |
        | data      | bytes  | Data to validate            |
        | signature | bytes  | Signature to validate       |
    === "OK"
        | Response value | Type | Description                                                                   |
        | -------------- | ---- | ----------------------------------------------------------------------------- |
        | valid          | bool | True if and only if provided data and its signature matches provided key-pair |
    === "NOT_FOUND"
        Key-pair with provided name doesn't exist
    === "FAILED_PRECONDITION"
        The vault is sealed.

### HMAC

The `system_vault` service offers the capability to generate HMAC (Hash-based Message Authentication Code) for data. One of the key security measures in this process is that the HMAC secret, which is utilized in generating the HMAC, never leaves the HSM.

!!! info
    To ensure the utmost security, all operations involving HMAC are performed exclusively within the HSM.

??? example "rpc HMACSignStream(stream HMACSignStreamRequest) returns (HMACSignStreamResponse) {};"
    Calculate HMAC signature for input data stream. HMAC secret never leaves the HSM (hardware security module). It automatically uses the best available HMAC algorithm for currently used HSM.

    !!! tip
        If your data is short (maximum several kilobytes) - use `HMACSign` instead. It will be faster and use less resources.

    === "Request"
        | Parameter | Type  | Description        |
        | --------- | ----- | ------------------ |
        | data      | bytes | Data chunk to sign |
    === "OK"
        | Response value | Type  | Description                    |
        | -------------- | ----- | ------------------------------ |
        | signature      | bytes | Signature of the provided data |
    === "FAILED_PRECONDITION"
        The vault is sealed.

??? example "rpc HMACVerifyStream(stream HMACVerifyStreamRequest) returns (HMACVerifyStreamResponse) {};"
    Verify HMAC signature of the data stream.

    !!! tip
        If your data is short (maximum several kilobytes) - use `HMACVerify` instead. It will be faster and use less resources.

    === "Request"
        | Parameter | Type  | Description            |
        | --------- | ----- | ---------------------- |
        | data      | bytes | Data chunk to validate |
        | signature | bytes | Signature to validate  |
    === "OK"
        | Response value | Type | Description                                                 |
        | -------------- | ---- | ----------------------------------------------------------- |
        | valid          | bool | True if and only if provided data and its signature matches |
    === "FAILED_PRECONDITION"
        The vault is sealed.

??? example "rpc HMACSign(HMACSignRequest) returns (HMACSignResponse) {};"
    Calculate HMAC signature for input data. HMAC secret never leaves the HSM (hardware security module). It automatically uses the best available HMAC algorithm for currently used HSM.

    !!! warning
        The data must be short (max several kilobytes). If it is longer - use `HMACSignStream` instead.

    === "Request"
        | Parameter | Type  | Description  |
        | --------- | ----- | ------------ |
        | data      | bytes | Data to sign |
    === "OK"
        | Response value | Type  | Description                    |
        | -------------- | ----- | ------------------------------ |
        | signature      | bytes | Signature of the provided data |
    === "FAILED_PRECONDITION"
        The vault is sealed.

??? example "rpc HMACVerify(HMACVerifyRequest) returns (HMACVerifyResponse) {};"
    Verifiy HMAC signature of the data. 
    
    !!! warning
        The data must be short (max several kilobytes). If it is longer - use `HMACVerifyStream` instead.

    === "Request"
        | Parameter | Type  | Description           |
        | --------- | ----- | --------------------- |
        | data      | bytes | Data to validate      |
        | signature | bytes | Signature to validate |
    === "OK"
        | Response value | Type | Description                                                 |
        | -------------- | ---- | ----------------------------------------------------------- |
        | valid          | bool | True if and only if provided data and its signature matches |
    === "FAILED_PRECONDITION"
        The vault is sealed.

### Encryption

The `system_vault` service includes the capability to encrypt data, adding an extra layer of security to sensitive information. The encryption key used for this process is securely managed and protected within the HSM, safeguarding it from potential threats.

??? example "rpc HMACSignStream(stream HMACSignStreamRequest) returns (HMACSignStreamResponse) {};"
    Encrypt data stream. Encryption secret never leaves the HSM (hardware security module). It automatically uses the best available encryption algorithm for currently used HSM. It will only select the algorithm that is capable of decrypting from the middle of the whole data (partial decryption).

    !!! tip
        If your data is short (maximum several kilobytes) - use `HMACSign` instead. It will be faster and use less resources.

    === "Request"
        | Parameter | Type  | Description           |
        | --------- | ----- | --------------------- |
        | plainData | bytes | Data chunk to encrypt |
    === "OK"
        | Response value | Type  | Description             |
        | -------------- | ----- | ----------------------- |
        | encryptedData  | bytes | Encrypted chunk of data |
    === "FAILED_PRECONDITION"
        The vault is sealed.

??? example "rpc HMACVerifyStream(stream HMACVerifyStreamRequest) returns (HMACVerifyStreamResponse) {};"
    Decrypt data stream. If you want to decrypt part of the information from the middle of the whole data - ensure, that the first chunk of the data, that you are sending is padded by 2048 bit.

    !!! tip
        If your data is short (maximum several kilobytes) - use `HMACVerify` instead. It will be faster and use less resources.

    === "Request"
        | Parameter     | Type  | Description                 |
        | ------------- | ----- | --------------------------- |
        | encryptedData | bytes | Encrypted chunk of the data |
    === "OK"
        | Response value | Type  | Description                 |
        | -------------- | ----- | --------------------------- |
        | plainData      | bytes | Decrypted chunk of the data |
    === "FAILED_PRECONDITION"
        The vault is sealed.


## HSM
The system_vault service offers the flexibility to work with multiple Hardware Security Module (HSM) providers, allowing users to configure their preferred provider before the initial startup of the system. This capability enables seamless integration with different HSM technologies based on specific requirements or preferences.

To select a particular HSM provider, you can set the `HSM_PROVIDER` environment variable. The following providers are currently supported:

- ***`softhsm2`***: This is the default provider that emulates an HSM using the SoftHSM2 library.
- ***`dynamic`***: This provider offers a dynamic interface, allowing potential connections with various PKCS11 client libraries. It provides flexibility for utilizing different HSM implementations.

!!! warning
    It is worth mentioning that due to the nature of making secrets non-extractable for security reasons, switching HSM providers at runtime may not be possible. Therefore, it is recommended to configure the desired provider before the initial startup of the system.

Furthermore, the system is designed to accommodate additional HSM providers in the future. Contributions to expand the range of supported providers are highly appreciated, ensuring ongoing enhancement and versatility in choosing the most suitable HSM solution.

### SoftHSM2
The SoftHSM2 provider does not require any additional configuration. It comes preconfigured with optimal settings that adhere to best practices, alleviating any concerns or need for manual adjustments on your part.

### Dynamic HSM provider
To ensure the proper functioning of the dynamic provider, you need to follow a few steps. First, run OpenBP and it will create a volume for the `system_vault` service at `./data/system/vault/`pkcs (depending on what setup you are using). In this folder, copy your PKCS11 library and name it `pkcs.so`.

Additionally, you need to define the slot by setting the `DYNAMIC_PKCS11_SLOT` environment variable. This variable specifies the slot that the service should attempt to open when loading the PKCS11 library.

By configuring the environment variable and placing the PKCS11 library in the designated folder, the `system_vault` service will automatically load the library and try to open the provided slot during its startup process.

Start OpenBP and check the logs. You will receive a message if everything is ok.