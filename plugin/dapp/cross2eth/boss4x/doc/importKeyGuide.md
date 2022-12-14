###  向硬件签名PCIe卡中导入secp256k1私钥实施步骤


***
#### １.生成sm2密钥对
```
[1] 查询密钥
[2] 删除密钥
[3] 生成密钥
[0] 退出
Please input the item number:3
请输入要生成的密钥类型，0-SM4, 1-SM1，2-SM2，3-ECC_SECP_256R1, 4-RSA, 5-AES, 6-DES，7-SM7, 8-ECC_SECP_256K1, 9-HMAC: 2
请输入要存储的密钥索引号(有效范围[1,64])：11
Function [TassCtlGenerateKey] run success
sk_keyLen = 32, sk_key = A1DC40596A7D45FCFDB3569222AB9AA0BF80E87CEE21DE06EBB34D1DB6891DAB
pk_kcvLen = 64, pk_kcv = 3744522891E6216AB8DAFF91598F361987FFEB907CCDD040815471F89AE179E2BBECDBB9E5FD1FAE8F8601F58A73E4D91213D6576F72B5A69D68F352ADFD00AF
```



***
#### 2.生成sm2密钥对(在本地,在导出卡中密钥时，需要使用)

```
$./ebcli_A sm2 create
sm2 public  key = df88444bb03100ae594bf09857e2f9183cf9d2aa6b5287282b17fa2f88d3d0bc6ecc3b6074e09b65c876257d22581bd6e68c4628b9d4edc6479a8ab733d0bbc4
sm2 private key = 1cca3f93369cf987637d856169e13f0f623be4637e2dbf46040b14d15a8a4ada

```


***
#### 3.导出sm2密钥对
```
[1] 导出非对称密钥
[2] 导入非对称密钥
[0] 退出
Please input the item number:1
请输入SM2的保护公钥: df88444bb03100ae594bf09857e2f9183cf9d2aa6b5287282b17fa2f88d3d0bc6ecc3b6074e09b65c876257d22581bd6e68c4628b9d4edc6479a8ab733d0bbc4
请输入导出密钥的类型：2-SM2，3-ECC_SECP_256R1, 4-RSA, 8-ECC_SECP_256K1: 2
请输入要导出密钥的索引：11
Function [TassCtlExportKey] run success
导出非对称密钥成功！
随机对称密钥：symmKeyCipherLen = 112, symmKeyCipher: D67AE0CD8D6F38D0CEB6F37A6A32DA9F44A825BB3D760DD7DA0AD7C31DC708FE9ECBEBDF8E5897BC55A6D3D80CC0EDF7B13BF21838F268D99E0D6160E61D776863EA2DFFE8E56F7DDA937F164A7E6B5FFC83E98029E06D6692B33B29BDCC18E2B60509B374B4692420D9C275E6009B35

导出的公钥：exportedPkLen = 64, exportedPk: 3744522891E6216AB8DAFF91598F361987FFEB907CCDD040815471F89AE179E2BBECDBB9E5FD1FAE8F8601F58A73E4D91213D6576F72B5A69D68F352ADFD00AF
导出的私钥：exportedKeyCipherByPkLen = 32, exportedKeyCipherByPk: DB15564A8A9CB01FB4A9A2998554C548CF06E0BC2634CF859EB330CA752AE2A0
```


***
#### 4.解密导出的sm2私钥
```
./ebcli_A sm2 decipher -c DB15564A8A9CB01FB4A9A2998554C548CF06E0BC2634CF859EB330CA752AE2A0 -k 1cca3f93369cf987637d856169e13f0f623be4637e2dbf46040b14d15a8a4ada -s D67AE0CD8D6F38D0CEB6F37A6A32DA9F44A825BB3D760DD7DA0AD7C31DC708FE9ECBEBDF8E5897BC55A6D3D80CC0EDF7B13BF21838F268D99E0D6160E61D776863EA2DFFE8E56F7DDA937F164A7E6B5FFC83E98029E06D6692B33B29BDCC18E2B60509B374B4692420D9C275E6009B35
0x546563941bdd7639179c5f272bfbb708e53dc92245791c9f1ec1a9fb8674a4eb
```


***
#### 5.在本地加密需要导入的secp256k1私钥
```
./ebcli_A sm2 encipher -k 0xcc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944 -t 0x546563941bdd7639179c5f272bfbb708e53dc92245791c9f1ec1a9fb8674a4eb -f 0x37d244cae2278d575dc93ca70a0b6f6c
随机对称密钥:74462058d38d033f09617d790159d4083108b9088deed07d760dd1ea56f16e7228b7f3612bf5b30904f07a0f1314f0ee17cd4a8921d8b17d20ebe77484aee53b75fd1b37a91582cfb7f77524c97cb748fda5944851067a51e61d4cdd4d74c61ec496458b8284a0ce1ac558745a6b805a len: 112
需要导入的公钥:b8eff229e59682ddd087158a2a44179b9083900b6f30c7ea7aac6aa5856f43c3506a16bc46f313c0cb47107ebf6657e0f4418fe984d67c372a7a4cf8d07b21bc len: 64
需要导入的私钥:c9779be48615a0fb817e6675089fe272b47a47e5bf9e1bb6f4fbb4255259ff0e len: 32
```




***
#### ６.导入secp256k1私钥
```
[1] 应用密钥管理（增加、查询、删除）
[2] 导入外部明文密钥到设备
[3] 修改私钥使用权限口令
[4] 导入导出非对称密钥
[0] 退出
Please input the item number:4
[1] 导出非对称密钥
[2] 导入非对称密钥
[0] 退出
Please input the item number:2
请输入SM2的保护私钥(密文,从加密卡中导出,其公钥用于sm4加密): A1DC40596A7D45FCFDB3569222AB9AA0BF80E87CEE21DE06EBB34D1DB6891DAB
请输入随机对称密钥: 50b8f08687690f7ee123abf4993741ad06afb34cfc20fcf6a8dd9c82993ad0cdeeeec063fa956e24b99027351a596a9f706fcdc445903fb05616a1af052a0886a57cdd0588a4009760b9ff2061632343aa3ab9309d1aa39c9853d1c88df79d005c35c6aea130ff27229d71110c09832f
请输入要导入的公钥: 504fa1c28caaf1d5a20fefb87c50a49724ff401043420cb3ba271997eb5a43879f2f7508d37165db9d9721b819ccaa3ef08a20bbab986c18b79d44a7e4201b8e
请输入要导入的私钥(密文): c497234d2eab47d5cc96fa174e8d9ef7417bfc045d4ad8a87697c48a16282ff9
请输入要导入密钥的类型：2-SM2，3-ECC_SECP_256R1, 4-RSA, 8-ECC_SECP_256K1: 8
请输入要导入密钥的索引：8
Function [TassCtlImportKey] run success
导入非对称密钥成功！
```
