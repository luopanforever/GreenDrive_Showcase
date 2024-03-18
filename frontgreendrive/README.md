# Getting Started with Create React App
You can learn more in the [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started).





#### OPEN SSL 报错

```json
{
  opensslErrorStack: [ 'error:03000086:digital envelope routines::initialization error' ],
  library: 'digital envelope routines',
  reason: 'unsupported',
  code: 'ERR_OSSL_EVP_UNSUPPORTED'
}
```
`export NODE_OPTIONS=--openssl-legacy-provider`
