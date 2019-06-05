# weather-checker
Simple Serverless project example using Go and AWS Lambda to return basic weather information for a given city name. The endpoint would be formatted like `/dev/weather/{city}`. For example: https://ktjlpq2ws7.execute-api.us-east-1.amazonaws.com/dev/weather/Rome,it

## Prerequisites

Install NPM and then Serverless.
```
npm install serverless -g
```

Register for an API key with [OpenWeatherMap](https://openweathermap.org/api) and then add it as an environment variable named `OWM_KEY` on your Lambda after its first deployment.

Create AWS Security Credentials. It doesn't *have* to be for an IAM user but I'd recommend it. 

## Setup

```
serverless config credentials --provider aws --key YOUR_AWS_KEY --secret YOUR_AWS_SECRET
```

## Deployment

```
make deploy
```

## Authors

* **Ted O'Connor** - [opub](https://github.com/opub)

## License

This project is licensed under the Apache License - see the [LICENSE](LICENSE) file for details.