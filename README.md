# URL SHORTNER

A REST-API application in golang, that provide allows users to shorten long URLs into short, shareable links. The app is capable of generating unique shortened URLs, mapping them to their corresponding original URLs, and redirecting users when they access the shortened link.


## API Reference

#### Handling the URL

```http
  POST /short
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `fullurlname` | `string` | **Required**. URL that required to be short. |

#### Get item

 - Get the redirected URL in `{URL}/shorten/{shortkey}`


## Deployment

To deploy this project run
