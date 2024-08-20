use std::fmt;
use std::sync::Arc;

use bitcoin_hashes::sha256;
use bitcoin_hashes::Hash;
use ecies::decrypt;
use ecies::encrypt;
use libsecp256k1::sign;
use libsecp256k1::verify;
use libsecp256k1::Message;
use libsecp256k1::PublicKey;
use libsecp256k1::SecretKey;
use libsecp256k1::Signature;

#[derive(Clone, Copy, Debug)]
pub enum CryptoError {
    Secp256k1Error(libsecp256k1::Error),
}

#[derive(Clone)]
pub struct KeyPair {
    private_key: Vec<u8>,
    public_key: Vec<u8>,
}

impl KeyPair {
    pub fn get_public_key(&self) -> Vec<u8> {
        self.public_key.clone()
    }

    pub fn get_private_key(&self) -> Vec<u8> {
        self.private_key.clone()
    }
}

impl fmt::Display for CryptoError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            Self::Secp256k1Error(err) => write!(f, "Secp256k1 error {}", err),
        }
    }
}

pub fn sign_ecdsa(msg: Vec<u8>, private_key_bytes: Vec<u8>) -> Result<Vec<u8>, CryptoError> {
    let sk = SecretKey::parse_slice(&private_key_bytes).map_err(CryptoError::Secp256k1Error)?;
    let hashed_message = sha256::Hash::hash(&msg);
    let msg = Message::parse_slice(hashed_message.as_byte_array())
        .map_err(CryptoError::Secp256k1Error)?;
    let (signature, _) = sign(&msg, &sk);
    Ok(signature.serialize_der().as_ref().to_vec())
}

pub fn verify_ecdsa(
    msg: Vec<u8>,
    signature_bytes: Vec<u8>,
    public_key_bytes: Vec<u8>,
) -> Result<bool, CryptoError> {
    let compressed = public_key_bytes.len() == 33;
    let pk = match compressed {
        true => PublicKey::parse_slice(
            &public_key_bytes,
            Some(libsecp256k1::PublicKeyFormat::Compressed),
        )
        .map_err(CryptoError::Secp256k1Error)?,
        false => {
            PublicKey::parse_slice(&public_key_bytes, Some(libsecp256k1::PublicKeyFormat::Full))
                .map_err(CryptoError::Secp256k1Error)?
        }
    };
    let hashed_message = sha256::Hash::hash(&msg);
    let msg = Message::parse_slice(hashed_message.as_byte_array())
        .map_err(CryptoError::Secp256k1Error)?;
    let sig = Signature::parse_der(&signature_bytes).map_err(CryptoError::Secp256k1Error)?;
    Ok(verify(&msg, &sig, &pk))
}

pub fn encrypt_ecies(msg: Vec<u8>, public_key_bytes: Vec<u8>) -> Result<Vec<u8>, CryptoError> {
    encrypt(&public_key_bytes, &msg).map_err(CryptoError::Secp256k1Error)
}

pub fn decrypt_ecies(
    cipher_text: Vec<u8>,
    private_key_bytes: Vec<u8>,
) -> Result<Vec<u8>, CryptoError> {
    decrypt(&private_key_bytes, &cipher_text).map_err(CryptoError::Secp256k1Error)
}

pub fn generate_keypair() -> Result<Arc<KeyPair>, CryptoError> {
    let (sk, pk) = ecies::utils::generate_keypair();
    let keypair = KeyPair {
        private_key: sk.serialize().to_vec(),
        public_key: pk.serialize().to_vec(),
    };
    Ok(keypair.into())
}

#[cfg(test)]
mod tests {
    use ecies::utils::generate_keypair;

    use super::*;

    #[test]
    fn test_ecdsa() {
        let (sk, pk) = generate_keypair();
        let msg = b"hello world";
        let signature = sign_ecdsa(msg.to_vec(), sk.serialize().to_vec()).unwrap();
        let result = verify_ecdsa(msg.to_vec(), signature, pk.serialize().to_vec()).unwrap();
        assert_eq!(result, true);
    }

    #[test]
    fn test_ecies() {
        let (sk, pk) = generate_keypair();
        let msg = b"hello world";
        let cipher_text = encrypt_ecies(msg.to_vec(), pk.serialize().to_vec()).unwrap();
        let plain_text = decrypt_ecies(cipher_text, sk.serialize().to_vec()).unwrap();
        assert_eq!(plain_text, msg.to_vec());
    }
}
