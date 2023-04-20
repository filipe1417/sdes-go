# sdes-go
Implementation of a simplified Data Encryption Standard (DES) in GO
## Usage: 
### Ecrypt:
go run sdes.go '8-bit plain text' '10-bit key'

Example: go run sdes.go 11100101 1011001010
       
### Decrypt (add 'd' flag):
go run sdes.go '8 bit encrypted text' '10 bit key' d

Example: go run sdes.go 11011111 1011110111 d

## TO-DO
- Implement Cipher Block Chaining (CBC) system
- Translate variable names
- Add code comments
- Search about the possibility of padding the plain text entry
         
